package repository

import (
	"context"
	"errors"
	"liokoredu/application/models"
	"liokoredu/application/solution"
	"log"
	"net/http"
	"time"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/labstack/echo"
)

type SolutionDatabase struct {
	pool *pgxpool.Pool
}

// GetSolutions implements solution.Repository
func (sd *SolutionDatabase) GetSolutions(taskId uint64) (models.SolutionsSQL, error) {
	var sln models.SolutionsSQL
	err := pgxscan.Select(context.Background(), sd.pool, &sln,
		`SELECT * FROM solutions WHERE task_id = $1`, taskId)
	log.Println(err)
	if errors.As(err, &pgx.ErrNoRows) || len(sln) == 0 {
		return models.SolutionsSQL{}, echo.NewHTTPError(http.StatusNotFound, errors.New("not found"))
	}

	if err != nil {
		return models.SolutionsSQL{}, err
	}

	slns := sln

	return slns, nil
}

// UpdateSolution implements solution.Repository
func (sd *SolutionDatabase) UpdateSolution(id uint64, code int, tests int) error {
	_, err := sd.pool.Exec(context.Background(),
		`UPDATE solutions SET "check_result" = $1, "tests_passed" = $2 WHERE id = $3`,
		code, tests, id)

	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (sd *SolutionDatabase) InsertSolution(taskId uint64, code string,
	testsTotal int, receivedTime time.Time) (uint64, error) {
	var id uint64
	err := sd.pool.QueryRow(context.Background(),
		`INSERT INTO solutions (task_id, check_result, tests_passed, tests_total, 
			received_date_time, source_code) 
		VALUES ($1, 1, 0, $2, $3, $4) RETURNING id`,
		taskId, testsTotal, receivedTime, code).Scan(&id)
	if err != nil {
		log.Println(err)
		return 0, echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return id, nil
}

func NewSolutionDatabase(conn *pgxpool.Pool) solution.Repository {
	return &SolutionDatabase{pool: conn}
}
