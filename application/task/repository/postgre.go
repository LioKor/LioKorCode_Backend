package repository

import (
	"context"
	"errors"
	"fmt"
	"liokoredu/application/models"
	"liokoredu/application/task"
	"liokoredu/pkg/constants"
	"log"
	"net/http"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/labstack/echo"
)

type TaskDatabase struct {
	pool *pgxpool.Pool
}

// DeleteTask implements task.Repository
func (td *TaskDatabase) DeleteTask(id uint64, uid uint64) error {
	resp, err := td.pool.Exec(context.Background(),
		`DELETE from tasks WHERE id = $1 AND creator = $2`,
		id, uid)

	if err != nil {
		log.Println("task repo: DeleteTask: error deleting task:", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	if resp.RowsAffected() == 0 {
		log.Println("task repo: DeleteTask: error deleting task: no task to delete")
		return echo.NewHTTPError(http.StatusNotFound, "task with taskId from this user not found")
	}

	return nil
}

// GetTasks implements task.Repository
func (td *TaskDatabase) GetTasks(page int) (*models.ShortTasks, error) {
	var t models.ShortTasks
	err := pgxscan.Select(context.Background(), td.pool, &t,
		`SELECT id, title, description, test_amount FROM tasks WHERE 
		is_private = false ORDER BY id DESC LIMIT $1 OFFSET $2`,
		constants.TasksPerPage, (page-1)*constants.TasksPerPage)
	if err != nil {
		log.Println("task repository: getTasks: error getting tasks", err)
		return &models.ShortTasks{}, err
	}

	return &t, nil
}

func (td *TaskDatabase) CreateTask(t *models.TaskSQL) (uint64, error) {
	var id uint64
	err := td.pool.QueryRow(context.Background(),
		`INSERT INTO tasks (title, description, hints, input, output, test_amount, tests, creator,
				is_private, code, date) 
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) RETURNING id`,
		t.Title, t.Description, t.Hints, t.Input, t.Output, t.TestAmount, t.Tests, t.Creator,
		t.IsPrivate, t.Code, t.Date).Scan(&id)

	if err != nil {
		log.Println("task repository: createTask: error creating task:", err)
		return 0, echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return id, nil
}

func NewTaskDatabase(conn *pgxpool.Pool) task.Repository {
	return &TaskDatabase{pool: conn}
}

func (td TaskDatabase) GetTask(id uint64) (*models.TaskSQL, error) {
	var t []models.TaskSQL
	err := pgxscan.Select(context.Background(), td.pool, &t,
		`SELECT * FROM tasks WHERE id = $1`, id)
	if err != nil {
		log.Println("task repository: getTask: error getting task", err)
		return &models.TaskSQL{}, err
	}

	if len(t) == 0 {
		return &models.TaskSQL{}, echo.NewHTTPError(http.StatusNotFound, errors.New("Task with id "+fmt.Sprint(id)+" not found"))
	}

	return &t[0], nil
}
