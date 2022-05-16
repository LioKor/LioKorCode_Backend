package repository

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/labstack/echo"

	"liokoredu/application/models"
	"liokoredu/application/task"
)

type TaskDatabase struct {
	pool *pgxpool.Pool
}

func (td *TaskDatabase) FindTasksFull(str string, useSolved bool, solved bool, useMine bool, mine bool, uid uint64, page int, count int) (*models.ShortTasks, int, error) {
	t := models.ShortTasks{}
	switch {
	case !useSolved && !useMine:
		err := pgxscan.Select(context.Background(), td.pool, &t,
			`SELECT t.id, t.title, t.description, t.test_amount, t.creator as creator_id, u.username as creator
			FROM tasks t
			JOIN users u ON u.id = t.creator
			WHERE is_private = false and (LOWER(title) LIKE '%' || $1 || '%'
			OR LOWER(description) LIKE '%' || $1 || '%' OR to_char(t.id, '999') LIKE '%' || $1 || '%')
			ORDER BY id DESC 
			LIMIT $2
			OFFSET $3`,
			strings.ToLower(str), count, (page-1)*count)

		if err != nil {
			log.Println("task repository: findTasksFull: error getting tasks", err)
			return &models.ShortTasks{}, 0, err
		}

		n := []int{}
		err = pgxscan.Select(context.Background(), td.pool, &n,
			`select count (*) 
			from tasks t
			JOIN users u ON u.id = t.creator
			WHERE is_private = false and (LOWER(title) LIKE '%' || $1 || '%'
			OR LOWER(description) LIKE '%' || $1 || '%' OR to_char(t.id, '999') LIKE '%' || $1 || '%')
			LIMIT $2
			OFFSET $3;`, strings.ToLower(str), count, (page-1)*count)

		if err != nil {
			log.Println("task repository: FindTasksFull: error getting num:", err)
			return &models.ShortTasks{}, 0, err
		}

		if len(n) == 0 {
			return &models.ShortTasks{}, 0, err
		}

		return &t, n[0], nil
	case useSolved && !useMine:
		s := "="
		if !solved {
			s = "!="
		}

		err := pgxscan.Select(context.Background(), td.pool, &t,
			`SELECT t.id, t.title, t.description, t.test_amount, t.creator as creator_id, u.username as creator
		FROM tasks t
		JOIN users u ON u.id = t.creator
		JOIN tasks_done td ON td.uid = $1 and td.task_id`+s+` t.id
		WHERE is_private = false and (LOWER(title) LIKE '%' || $1 || '%'
		OR LOWER(description) LIKE '%' || $2 || '%' OR to_char(t.id, '999') LIKE '%' || $2 || '%')
		ORDER BY id DESC 
		LIMIT $3
		OFFSET $4;`,
			uid, str, count, (page-1)*count)

		if err != nil {
			log.Println("task repository: findTasksFull: error getting tasks", err)
			return &models.ShortTasks{}, 0, err
		}

		n := []int{}
		err = pgxscan.Select(context.Background(), td.pool, &n,
			`select count (*)
			FROM tasks t
			JOIN users u ON u.id = t.creator
			JOIN tasks_done td ON td.uid = $1 and td.task_id `+s+` t.id
			WHERE is_private = false and (LOWER(title) LIKE '%' || $2 || '%'
			OR LOWER(description) LIKE '%' || $2 || '%' OR to_char(t.id, '999') LIKE '%' || $2 || '%')
			LIMIT $3
			OFFSET $4;`,
			uid, str, count, (page-1)*count)

		if err != nil {
			log.Println("task repository: FindTasksFull: error getting num:", err)
			return &models.ShortTasks{}, 0, err
		}

		if len(n) == 0 {
			return &models.ShortTasks{}, 0, err
		}

		return &t, n[0], nil

	case !useSolved && useMine:
		err := pgxscan.Select(context.Background(), td.pool, &t,
			`SELECT t.id, t.title, t.description, t.test_amount, t.creator as creator_id, u.username AS creator
			FROM tasks t
 			JOIN users u ON u.id = t.creator
			WHERE is_private = false AND t.creator = $1 and (LOWER(title) LIKE '%' || $2 || '%'
			OR LOWER(description) LIKE '%' || $2 || '%' OR to_char(t.id, '999') LIKE '%' || $2 || '%')
			LIMIT $3 
			OFFSET $4`,
			uid, str, count, (page-1)*count)
		if err != nil {
			log.Println("task repository: getTasks: error getting tasks", err)
			return &models.ShortTasks{}, 0, err
		}

		n := []int{}
		err = pgxscan.Select(context.Background(), td.pool, &n,
			`select count (*) 
				FROM tasks t
				JOIN users u ON u.id = t.creator
				WHERE is_private = false AND t.creator = $1 and (LOWER(title) LIKE '%' || $2 || '%'
				OR LOWER(description) LIKE '%' || $2 || '%' OR to_char(t.id, '999') LIKE '%' || $2 || '%')
				LIMIT $3 
				OFFSET $4`,
			uid, str, count, (page-1)*count)

		if err != nil {
			log.Println("task repository: FindTasksFull: error getting num:", err)
			return &models.ShortTasks{}, 0, err
		}

		if len(n) == 0 {
			return &models.ShortTasks{}, 0, err
		}

		return &t, n[0], nil
	case useSolved && useMine:
		s := "="
		if !solved {
			s = "!="
		}

		err := pgxscan.Select(context.Background(), td.pool, &t,
			`SELECT t.id, t.title, t.description, t.test_amount, t.creator as creator_id, u.username as creator
		FROM tasks t
		JOIN users u ON u.id = t.creator
		JOIN tasks_done td ON td.uid = $1 and td.task_id`+s+` t.id
		WHERE is_private = false and t.creator = $2 (LOWER(title) LIKE '%' || $3 || '%'
		OR LOWER(description) LIKE '%' || $3 || '%' OR to_char(t.id, '999') LIKE '%' || $3 || '%')
		LIMIT $4
		OFFSET $5;`,
			uid, uid, str, count, (page-1)*count)

		if err != nil {
			log.Println("task repository: findTasksFull: error getting tasks", err)
			return &models.ShortTasks{}, 0, err
		}

		n := []int{}
		err = pgxscan.Select(context.Background(), td.pool, &n,
			`select count (*) 
			FROM tasks t
			JOIN users u ON u.id = t.creator
			JOIN tasks_done td ON td.uid = $1 and td.task_id `+s+` t.id
			WHERE is_private = false and t.creator = $2 (LOWER(title) LIKE '%' || $3 || '%'
			OR LOWER(description) LIKE '%' || $3 || '%' OR to_char(t.id, '999') LIKE '%' || $3 || '%')
			LIMIT $4
			OFFSET $5;`,
			uid, uid, str, count, (page-1)*count)

		if err != nil {
			log.Println("task repository: FindTasksFull: error getting num:", err)
			return &models.ShortTasks{}, 0, err
		}

		if len(n) == 0 {
			return &models.ShortTasks{}, 0, err
		}

		return &t, n[0], nil
	}

	return &models.ShortTasks{}, 0, nil
}

func (td *TaskDatabase) GetPages() (int, error) {
	n := []int{}
	err := pgxscan.Select(context.Background(), td.pool, &n,
		`select count (*) from tasks;`)

	if err != nil {
		log.Println("task repository: GetPages: error getting num:", err)
		return 0, echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	if len(n) == 0 {
		return 0, nil
	}

	return n[0], nil
}

func (td *TaskDatabase) FindTasks(str string, page int, count int) (*models.ShortTasks, error) {
	t := models.ShortTasks{}
	err := pgxscan.Select(context.Background(), td.pool, &t,
		`SELECT t.id, t.title, t.description, t.test_amount, t.creator as creator_id, u.username as creator
			FROM tasks t
			JOIN users u ON u.id = t.creator
			WHERE is_private = false and (LOWER(title) LIKE '%' || $1 || '%'
			OR LOWER(description) LIKE '%' || $1 || '%' OR to_char(t.id, '999') LIKE '%' || $1 || '%')
			ORDER BY id DESC 
			LIMIT $2
			OFFSET $3`,
		strings.ToLower(str), count, (page-1)*count)
	if err != nil {
		log.Println("task repository: findTasks: error getting tasks", err)
		return &models.ShortTasks{}, err
	}

	return &t, nil
}

func (td *TaskDatabase) GetSolvedTasks(uid uint64, page int, count int) (*models.ShortTasks, error) {
	t := models.ShortTasks{}
	err := pgxscan.Select(context.Background(), td.pool, &t,
		`SELECT t.id, t.title, t.description, t.test_amount, t.creator as creator_id, u.username as creator
		FROM tasks t
		JOIN users u ON t.creator = u.id
		JOIN tasks_done td ON td.uid = $1 and td.task_id = t.id
		WHERE is_private = false
		ORDER BY t.id DESC LIMIT $2 OFFSET $3`,
		uid, count, (page-1)*count)
	if err != nil {
		log.Println("task repository: getSolvedTasks: error getting tasks", err)
		return &models.ShortTasks{}, err
	}

	return &t, nil
}

func (td *TaskDatabase) GetUnsolvedTasks(uid uint64, page int, count int) (*models.ShortTasks, error) {
	t := models.ShortTasks{}
	err := pgxscan.Select(context.Background(), td.pool, &t,
		`SELECT t.id, t.title, t.description, t.test_amount, t.creator as creator_id, u.username as creator
		FROM tasks t
		JOIN users u ON t.creator = u.id
		JOIN tasks_done td ON td.uid = $1 and td.task_id != t.id
		WHERE is_private = false
		ORDER BY t.id DESC LIMIT $2 OFFSET $3`,
		uid, count, (page-1)*count)
	if err != nil {
		log.Println("task repository: getSolvedTasks: error getting tasks", err)
		return &models.ShortTasks{}, err
	}

	return &t, nil
}

func (td *TaskDatabase) IsCleared(taskId uint64, uid uint64) (bool, error) {
	var id []uint64
	err := pgxscan.Select(context.Background(), td.pool, &id,
		`SELECT uid FROM tasks_done WHERE uid = $1 and task_id = $2`,
		uid, taskId)
	if err != nil {
		log.Println("task repository: IsCleared: error checking task tasks", err)
		return false, err
	}

	if len(id) == 0 {
		return false, nil
	}

	return true, nil
}

func (td *TaskDatabase) MarkTaskDone(id uint64, uid uint64) error {
	_, err := td.pool.Exec(context.Background(),
		`INSERT INTO tasks_done (uid, task_id) VALUES ($1, $2);`, uid, id)

	if err != nil {
		log.Println("task repository: MarkTaskDone: error marking task_done:", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return nil
}

func (td *TaskDatabase) UpdateTask(t *models.TaskSQL) error {
	resp, err := td.pool.Exec(context.Background(),
		`UPDATE tasks set title = $1, description = $2, hints = $3, 
		input = $4, output = $5, test_amount = $6, tests = $7 WHERE creator = $8 AND id = $9;`,
		t.Title, t.Description, t.Hints, t.Input, t.Output, t.TestAmount, t.Tests, t.Creator,
		t.Id)

	if err != nil {
		log.Println("task repository: UpdateTask: error updating task:", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	if resp.RowsAffected() == 0 {
		return echo.NewHTTPError(http.StatusNotFound, "no task from user")
	}

	return nil
}

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

func (td *TaskDatabase) GetTasks(page int, count int) (*models.ShortTasks, error) {
	t := models.ShortTasks{}
	err := pgxscan.Select(context.Background(), td.pool, &t,
		`SELECT t.id, t.title, t.description, t.test_amount, t.creator as creator_id, u.username as creator
			FROM tasks t
			JOIN users u ON u.id = t.creator
			WHERE is_private = false 
			ORDER BY id DESC 
			LIMIT $1 
			OFFSET $2`,
		count, (page-1)*count)
	if err != nil {
		log.Println("task repository: getTasks: error getting tasks", err)
		return &models.ShortTasks{}, err
	}

	return &t, nil
}

func (td *TaskDatabase) GetUserTasks(uid uint64, page int, count int) (*models.ShortTasks, error) {
	t := models.ShortTasks{}
	err := pgxscan.Select(context.Background(), td.pool, &t,
		`SELECT t.id, t.title, t.description, t.test_amount, t.creator as creator_id, u.username AS creator
			FROM tasks t
 			JOIN users u ON u.id = t.creator
			WHERE is_private = false AND t.creator = $1 
			ORDER BY id DESC 
			LIMIT $2 
			OFFSET $3`,
		uid, count, (page-1)*count)
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
		return &models.TaskSQL{}, echo.NewHTTPError(http.StatusNotFound, "Task with id "+fmt.Sprint(id)+" not found")
	}

	return &t[0], nil
}
