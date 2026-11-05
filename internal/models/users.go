package models

import (
	"errors"
	"time"

	"github.com/antonisgkamitsios/verese/internal/validator"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type UserModel struct {
	DB *pgxpool.Pool
}

type User struct {
	ID        int64
	Name      string
	Email     string
	Password  password
	Activated bool
	Version   int32
	CreatedAt time.Time
	UpdatedAt time.Time

	validator.Validator
}

type password struct {
	plaintext *string
	hash      []byte
}

type UserModelInterface interface{}

var AnonymousUser = &User{}

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

func (u *User) Validate() {
	u.Check(u.Name != "", "name", "must be provided")
	u.Check(len(u.Name) <= 500, "name", "must be less than 500 bytes long")

	u.Check(u.Email != "", "email", "must be provided")
	u.Check(validator.ValidMail(u.Email), "email", "must be a valid email address")

	if u.Password.plaintext != nil {
		u.Check(*u.Password.plaintext != "", "password", "must be provided")
		u.Check(len(*u.Password.plaintext) >= 8, "password", "must be at least 8 bytes long")
		u.Check(len(*u.Password.plaintext) <= 72, "password", "must be less than 72 bytes long")
	}

	if u.Password.hash == nil {
		panic("missing password hash for user")
	}
}
