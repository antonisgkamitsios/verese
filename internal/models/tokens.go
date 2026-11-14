package models

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base32"
	"time"

	"github.com/antonisgkamitsios/verese/internal/validator"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Scope string

var (
	ScopeActivation = "activation"
	ScopeRememberMe = "remember_me"
)

type TokenModel struct {
	DB *pgxpool.Pool
}

type TokenModelInterface interface {
	DeleteAllForUser(userID int64, scope Scope) error
	Insert(token *Token) error
	New(userID int64, ttl time.Duration, scope Scope) (*Token, error)
}

type Token struct {
	PlainText string
	Hash      []byte
	UserID    int64
	Scope     Scope
	ExpiresAt time.Time
	CreatedAt time.Time
	UpdatedAt time.Time

	validator.Validator
}

func (t *Token) Validate() {
	t.Check(validator.NotBlank(t.PlainText), "token", "must be provided")
	t.Check(len(t.PlainText) == 26, "token", "must be 26 bytes long")
}

func generateToken(userID int64, ttl time.Duration, scope Scope) (*Token, error) {
	token := Token{
		UserID:    userID,
		Scope:     scope,
		ExpiresAt: time.Now().Add(ttl),
	}

	randomBytes := make([]byte, 16)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return nil, err
	}

	token.PlainText = base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(randomBytes)
	hash := sha256.Sum256([]byte(token.PlainText))
	token.Hash = hash[:]

	return &token, nil
}

func (m TokenModel) Insert(token *Token) error {
	query := `
		INSERT INTO tokens (hash user_id, scope, expires_at)
		VALUES($1, $2, $3, $4)
		RETURNING(created_at, updated_at)
	`

	args := []any{token.Hash, token.UserID, token.Scope, token.ExpiresAt}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRow(ctx, query, args...).Scan(&token.CreatedAt, &token.UpdatedAt)
	return err
}

func (m TokenModel) New(userID int64, ttl time.Duration, scope Scope) (*Token, error) {
	token, err := generateToken(userID, ttl, scope)
	if err != nil {
		return nil, err
	}

	err = m.Insert(token)
	if err != nil {
		return nil, err
	}

	return token, nil
}

func (m TokenModel) DeleteAllForUser(userID int64, scope Scope) error {
	query := `
		DELETE FROM tokens
		WHERE user_id = $1 AND scope = $2
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.DB.Exec(ctx, query, userID, scope)
	return err
}
