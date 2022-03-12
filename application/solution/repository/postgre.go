package repository

import (
	"liokoredu/application/solution"

	"github.com/jackc/pgx/v4/pgxpool"
)

type SolutionDatabase struct {
	pool *pgxpool.Pool
}

func NewSolutionDatabase(conn *pgxpool.Pool) solution.Repository {
	return &SolutionDatabase{pool: conn}
}
