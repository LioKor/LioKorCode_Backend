package repository

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"liokoredu/application/models"
	"liokoredu/application/solution"
	"liokoredu/pkg/constants"
	"liokoredu/pkg/generators"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4"
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
		return models.SolutionSQL{}, echo.NewHTTPError(http.StatusNotFound, "solution for task from user not found")
	}
	if err != nil {
		log.Println("solution repo: GetSolution: error getting solutions:", err)
		return models.SolutionSQL{}, echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	if len(sln) == 0 {
		log.Println("solution repo: GetSolution: error getting solution: no solution")
		return models.SolutionSQL{}, echo.NewHTTPError(http.StatusNotFound, "solution for task from user not found")
	}

	base, _ := os.Getwd()
	sourceCode, _ := ioutil.ReadFile(filepath.ToSlash(base + constants.SolutionsDir + sln[0].SourceCode))
	sln[0].SourceCode = string(sourceCode)

	return sln[0], nil
}

func (sd *SolutionDatabase) DeleteSolution(id uint64, uid uint64) error {
	var filename []string

	err := pgxscan.Select(context.Background(), sd.pool, &filename,
		`SELECT source_code from solutions WHERE id = $1 AND uid = $2`,
		id, uid)

	if err != nil {
		log.Println("solution repo: DeleteSolution: error getting filename:", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	if len(filename) == 0 {
		log.Println("solution repo: DeleteSolution: error deleting solution: no file")
		return echo.NewHTTPError(http.StatusNotFound, "solution with solutionId from this user not found")
	}

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

	base, _ := os.Getwd()
	_ = os.Remove(filepath.ToSlash(base + constants.SolutionsDir + filename[0]))

	return nil
}

// GetSolutions implements solution.Repository
func (sd *SolutionDatabase) GetSolutions(taskId uint64, uid uint64) (models.SolutionsSQL, error) {
	var sln models.SolutionsSQL
	err := pgxscan.Select(context.Background(), sd.pool, &sln,
		`SELECT * FROM solutions WHERE task_id = $1 AND uid = $2`, taskId, uid)

	if err != nil {
		log.Println("solution repo: GetSolutions: error getting solutions:", err)
		return models.SolutionsSQL{}, echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	if len(sln) == 0 {
		log.Println("solution repo: GetSolutions: no solutions")
		return models.SolutionsSQL{}, nil
	}

	base, _ := os.Getwd()
	for _, elem := range sln {
		sourceCode, _ := ioutil.ReadFile(filepath.ToSlash(base + constants.SolutionsDir + elem.SourceCode))
		elem.SourceCode = string(sourceCode)
	}

	return sln, nil
}

func (sd *SolutionDatabase) UpdateSolution(id uint64, upd *models.SolutionUpdate) error {
	// PostgreSQL unable to store \x00 in string field
	upd.CheckMessage = strings.Replace(upd.CheckMessage, "\x00", "", -1)

	_, err := sd.pool.Exec(context.Background(),
		`UPDATE solutions SET check_result = $1, tests_passed = $2, check_message = $3,
		check_time = $4, compile_time = $5, checked_date_time = $6 WHERE id = $7`,
		upd.Code, upd.Passed, upd.CheckMessage, upd.CheckTime, upd.CompileTime, upd.CheckedDateTime, id)

	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (sd *SolutionDatabase) InsertSolution(taskId uint64, uid uint64, code map[string]interface{},
	testsTotal int, receivedTime time.Time) (uint64, error) {
	var id uint64

	filename := time.Now().Format("2006-01-02T15-04-05") +
		generators.RandStringRunes(constants.PrivateLength)

	b, _ := json.Marshal(code)
	data := []byte(b)

	base, _ := os.Getwd()
	file, err := os.Create(filepath.ToSlash(base + constants.SolutionsDir + filename))
	if err != nil {
		log.Println(err)
	}
	defer file.Close()

	_, err = file.Write(data)
	if err != nil {
		log.Println("solution repo: InsertSolution: error creating file:", err)
		return 0, echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	err = sd.pool.QueryRow(context.Background(),
		`INSERT INTO solutions (task_id, check_result, tests_passed, tests_total, 
			received_date_time, source_code, uid) 
		VALUES ($1, 1, 0, $2, $3, $4, $5) RETURNING id`,
		taskId, testsTotal, receivedTime, filename, uid).Scan(&id)
	if err != nil {
		log.Println("solution repo: InsertSolution: error inserting solution", err)
		return 0, echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return id, nil
}

func NewSolutionDatabase(conn *pgxpool.Pool) solution.Repository {
	return &SolutionDatabase{pool: conn}
}
