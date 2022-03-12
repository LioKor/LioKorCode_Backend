package models

import (
	"database/sql"
	json "encoding/json"
	"time"
)

type Task struct {
	Id          uint64     `json:"id"`
	Title       string     `json:"name"`
	Description string     `json:"description"`
	Input       string     `json:"stdinDescription"`
	Output      string     `json:"stdoutDescription"`
	Hints       string     `json:"hints"`
	TestsAmount int        `json:"testsAmount"`
	Tests       InputTests `json:"tests"`
}

type TaskSQL struct {
	Id          uint64         `sql:"id"`
	Title       string         `sql:"title"`
	Description string         `sql:"description"`
	Hints       sql.NullString `sql:"hints"`
	Input       string         `sql:"input"`
	Output      string         `sql:"output"`
	TestAmount  int            `sql:"test_amount"`
	Tests       string         `sql:"tests"`
	Creator     uint64         `sql:"creator"`
	IsPrivate   bool           `sql:"is_private"`
	Code        sql.NullString `sql:"code"`
	Date        time.Time      `sql:"date"`
}

func (tsql TaskSQL) ConvertToTask() *Task {
	t := &Task{}
	t.Id = tsql.Id
	t.Title = tsql.Title
	t.Description = tsql.Description
	t.Hints = tsql.Hints.String
	t.Input = tsql.Input
	t.Output = tsql.Output
	t.TestsAmount = tsql.TestAmount
	_ = json.Unmarshal([]byte(tsql.Tests), &t.Tests)

	return t
}
