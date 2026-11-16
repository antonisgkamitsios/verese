package models

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/antonisgkamitsios/verese/internal/validator"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

var ErrDuplicateEmail = errors.New("duplicate email")

var AnonymousUser = &User{}

type UserModel struct {
	DB *pgxpool.Pool
}

type User struct {
	ID        int64     `db:"id"`
	Name      string    `db:"name"`
	Email     string    `db:"email"`
	Password  password  `db:"password_hash"`
	Activated bool      `db:"activated"`
	Version   int32     `db:"version"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`

	validator.Validator `db:"-"`
}

type password struct {
	plaintext *string
	hash      []byte
}

type UserModelInterface interface {
	GetById(id int64) (*User, error)
	GetByEmail(email string) (*User, error)
	FindByToken(scope Scope, tokenPlainText string) (*User, error)
	Insert(user *User) error
	Update(user *User) error
}

func (u *User) IsAnonymous() bool {
	return u == AnonymousUser
}

func (p *password) Set(plaintextPassword string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plaintextPassword), 12)
	if err != nil {
		return err
	}

	p.plaintext = &plaintextPassword
	p.hash = hash
	return nil
}

func (p *password) Matches(plaintextPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(p.hash, []byte(plaintextPassword))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}
	}

	return true, nil
}

func (p *password) Scan(dest any) error {
	switch v := dest.(type) {
	case []byte:
		h := make([]byte, len(v))
		copy(h, v)
		p.hash = h
		return nil
	case nil:
		p.hash = nil
		return nil
	default:
		return fmt.Errorf("password: cannot scan %T", dest)
	}
}

func (u *User) Validate() {
	u.Check(validator.NotBlank(u.Name), "name", "must be provided")
	u.Check(len(u.Name) <= 500, "name", "must be less than 500 bytes long")

	u.Check(validator.NotBlank(u.Email), "email", "must be provided")
	u.Check(validator.ValidMail(u.Email), "email", "must be a valid email address")

	if u.Password.plaintext != nil {
		u.Check(validator.NotBlank(*u.Password.plaintext), "password", "must be provided")
		u.Check(len(*u.Password.plaintext) >= 8, "password", "must be at least 8 bytes long")
		u.Check(len(*u.Password.plaintext) <= 72, "password", "must be less than 72 bytes long")
	}

	if u.Password.hash == nil {
		panic("missing password hash for user")
	}
}

func (m UserModel) GetById(id int64) (*User, error) {
	query := `
		SELECT id, name, email,password_hash, activated, version, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.Query(ctx, query, id)
	if err != nil {
		return nil, err
	}

	user, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[User])
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &user, nil
}

func (m UserModel) GetByEmail(email string) (*User, error) {
	query := `
		SELECT id, name, email, password_hash, activated, version, created_at, updated_at
		FROM users
		WHERE email = $1
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.Query(ctx, query, email)
	if err != nil {
		return nil, err
	}

	user, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[User])
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &user, nil
}

func (m UserModel) FindByToken(scope Scope, tokenPlainText string) (*User, error) {
	tokenHash := sha256.Sum256([]byte(tokenPlainText))

	query := `
		SELECT users.id, users.name, users.email, users.password_hash, users.activated, users.version, users.created_at, users.updated_at	
		FROM users
		INNER JOIN tokens
		ON tokens.user_id = users.id
		WHERE tokens.hash = $1 AND tokens.scope = $2 AND tokens.expires_at > $3
	`

	args := []any{tokenHash[:], scope, time.Now()}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	user, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[User])
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &user, nil
}

func (m UserModel) Insert(user *User) error {
	query := `
		INSERT INTO users (name, email, password_hash, activated)
		VALUES($1, $2, $3, $4)
		RETURNING id, created_at, updated_at, version
	`

	args := []any{user.Name, user.Email, user.Password.hash, user.Activated}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRow(ctx, query, args...).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt, &user.Version)
	if err != nil {
		switch {
		case strings.Contains(err.Error(), `"users_email_key" (SQLSTATE 23505)`):
			return ErrDuplicateEmail
		default:
			return err
		}
	}

	return nil
}

func (m UserModel) Update(user *User) error {
	query := `
		UPDATE users
		SET name = $1, email = $2, password_hash = $3, activated = $4, version = version + 1, updated_at = $5
		WHERE id = $6 AND version = $7
		RETURNING version, updated_at
	`

	args := []any{user.Name, user.Email, user.Password.hash, user.Activated, time.Now(), user.ID, user.Version}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRow(ctx, query, args...).Scan(&user.Version, &user.UpdatedAt)
	if err != nil {
		switch {
		case strings.Contains(err.Error(), `"users_email_key" (SQLSTATE 23505)`):
			return ErrDuplicateEmail
		case errors.Is(err, pgx.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}

	return nil
}

type LoginForm struct {
	Email    string `form:"email"`
	Password string `form:"password"`

	validator.Validator `form:"-"`
}

func (f *LoginForm) Validate() {
	f.Check(validator.NotBlank(f.Email), "email", "must be provided")
	f.Check(validator.ValidMail(f.Email), "email", "must be a valid email address")

	f.Check(validator.NotBlank(f.Password), "password", "must be provided")
}
