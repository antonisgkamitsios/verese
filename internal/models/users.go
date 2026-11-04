package models

import "github.com/jackc/pgx/v5/pgxpool"

type UserModel struct {
	DB *pgxpool.Pool
}

type UserModelInterface interface{}
