package models

import "github.com/jackc/pgx/v5/pgxpool"

type Models struct {
	Users UserModelInterface
}

func NewModels(db *pgxpool.Pool) Models {
	return Models{
		Users: UserModel{DB: db},
	}
}
