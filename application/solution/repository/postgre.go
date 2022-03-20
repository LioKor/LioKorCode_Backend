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

func (sd *SolutionDatabase) GetSolution(id uint64, taskId uint64, uid uint64) (models.SolutionSQL, error) {
	var sln models.SolutionsSQL
	err := pgxscan.Select(context.Background(), sd.pool, &sln,
		`SELECT * FROM solutions WHERE id = $1 AND task_id = $2 AND uid = $3`, id, taskId, uid)
	if errors.As(err, &pgx.ErrNoRows) && len(sln) == 0 {
		log.Println("solution repo: GetSolution: error getting solution: no solution")
		return models.SolutionSQL{}, nil
	}
	if err != nil {
		log.Println("solution repo: GetSolutions: error getting solutions:", err)
		return models.SolutionSQL{}, echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return sln[0], nil
}

func (sd *SolutionDatabase) DeleteSolution(id uint64, uid uint64) error {
	resp, err := sd.pool.Exec(context.Background(),
		`DELETE from solutions WHERE id = $1 AND uid = $2`,
		id, uid)

	if err != nil {
		log.Println("solution repo: DeleteSolution: error deleting solution:", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	if resp.RowsAffected() == 0 {
		log.Println("solution repo: DeleteSolution: error deleting solution: no solution to delete")
		return echo.NewHTTPError(http.StatusNotFound, "solution with solutionId from this user not found")
	}

	return nil
}

// GetSolutions implements solution.Repository
func (sd *SolutionDatabase) GetSolutions(taskId uint64, uid uint64) (models.SolutionsSQL, error) {
	var sln models.SolutionsSQL
	err := pgxscan.Select(context.Background(), sd.pool, &sln,
		`SELECT * FROM solutions WHERE task_id = $1 AND uid = $2`, taskId, uid)
	if errors.As(err, &pgx.ErrNoRows) && len(sln) == 0 {
		log.Println("solution repo: GetSolutions: error getting solution: no solutions")
		return models.SolutionsSQL{}, nil
	}
	if err != nil {
		log.Println("solution repo: GetSolutions: error getting solutions:", err)
		return models.SolutionsSQL{}, echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return sln, nil
}

// UpdateSolution implements solution.Repository
func (sd *SolutionDatabase) UpdateSolution(id uint64, code int, tests int) error {
	_, err := sd.pool.Exec(context.Background(),
		`UPDATE solutions SET "check_result" = $1, "tests_passed" = $2, "task_updated" = false WHERE id = $3`,
		code, tests, id)

	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (sd *SolutionDatabase) InsertSolution(taskId uint64, uid uint64, code string, makefile string,
	testsTotal int, receivedTime time.Time) (uint64, error) {
	var id uint64
	err := sd.pool.QueryRow(context.Background(),
		`INSERT INTO solutions (task_id, check_result, tests_passed, tests_total, 
			received_date_time, source_code, uid, makefile) 
		VALUES ($1, 1, 0, $2, $3, $4, $5, $6) RETURNING id`,
		taskId, testsTotal, receivedTime, code, uid, makefile).Scan(&id)
	if err != nil {
		log.Println(err)
		return 0, echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return id, nil
}

func NewSolutionDatabase(conn *pgxpool.Pool) solution.Repository {
	return &SolutionDatabase{pool: conn}
}
