package repository

import (
	"context"
	"errors"
	"fmt"
	"liokoredu/application/models"
	"liokoredu/application/task"
	"log"
	"net/http"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/labstack/echo"
)

type TaskDatabase struct {
	pool *pgxpool.Pool
}

func NewTaskDatabase(conn *pgxpool.Pool) task.Repository {
	return &TaskDatabase{pool: conn}
}

func (td TaskDatabase) GetTask(id uint64) (*models.TaskSQL, error) {
	var t []models.TaskSQL
	log.Println("got")
	log.Println(id)
	err := pgxscan.Select(context.Background(), td.pool, &t,
		`SELECT * FROM tasks WHERE id = $1`, id)
	log.Println(err)

	if errors.As(err, &pgx.ErrNoRows) || len(t) == 0 {
		return &models.TaskSQL{}, echo.NewHTTPError(http.StatusNotFound, errors.New("Task with id "+fmt.Sprint(id)+" not found"))
	}

	if err != nil {
		log.Println(err)
		return &models.TaskSQL{}, err
	}

	return &t[0], nil
}
