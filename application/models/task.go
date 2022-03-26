package models

import (
	"database/sql"
	json "encoding/json"
	"log"
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

type ShortTask struct {
	Id          uint64 `json:"id"`
	Title       string `json:"name"`
	Description string `json:"description"`
	TestAmount  int    `json:"testsAmount"`
}

//easyjson:json
type ShortTasks []ShortTask

type TaskNew struct {
	Title       string     `json:"name"`
	Description string     `json:"description"`
	Input       string     `json:"stdinDescription"`
	Output      string     `json:"stdoutDescription"`
	Hints       string     `json:"hints"`
	Tests       InputTests `json:"tests"`
	Creator     uint64     `json:"creator"`
	IsPrivate   bool       `json:"is_private"`
	Code        string     `json:"code"`
}

//easyjson:json
type Tasks []Task

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

//easyjson:json
type TasksSQL []TaskSQL

func (tsql TaskSQL) ConvertToTask(isCreator bool) *Task {
	t := &Task{}
	t.Id = tsql.Id
	t.Title = tsql.Title
	t.Description = tsql.Description
	t.Hints = tsql.Hints.String
	t.Input = tsql.Input
	t.Output = tsql.Output
	t.TestsAmount = tsql.TestAmount
	err := json.Unmarshal([]byte(tsql.Tests), &t.Tests)
	if err != nil {
		log.Println("error converting tests: ", err)
	}

	if !isCreator {
		if len(t.Tests) >= 2 {
			t.Tests = t.Tests[:2]
		}
	}

	return t
}

func (tsksSQL TasksSQL) ConvertToTasks() *Tasks {
	tsks := Tasks{}
	for _, elem := range tsksSQL {
		tsks = append(tsks, *elem.ConvertToTask(false))
	}
	return &tsks
}

func (tn TaskNew) ConvertNewTaskToTaskSQL() *TaskSQL {
	t := &TaskSQL{}
	t.Title = tn.Title
	t.Description = tn.Description
	t.Hints = NewNullString(tn.Hints)
	t.Input = tn.Input
	t.Output = tn.Output
	t.TestAmount = len(tn.Tests)

	bts, err := json.Marshal(tn.Tests)
	if err != nil {
		log.Println("error converting tests: ", err)
	}

	str := string(bts[:])
	t.Tests = str
	t.Creator = tn.Creator
	t.IsPrivate = tn.IsPrivate
	t.Code = NewNullString(tn.Code)

	location, _ := time.LoadLocation("Europe/London")
	received := time.Now().In(location)

	t.Date = received

	return t
}

func (tn TaskNew) Validate() bool {
	if len(tn.Title) == 0 {
		return false
	}
	if len(tn.Description) == 0 {
		return false
	}
	if len(tn.Hints) == 0 {
		return false
	}
	if len(tn.Input) == 0 {
		return false
	}
	if len(tn.Output) == 0 {
		return false
	}
	if len(tn.Tests) < 2 {
		return false
	}

	return true
}

func NewNullString(s string) sql.NullString {
	if len(s) == 0 {
		return sql.NullString{}
	}
	return sql.NullString{
		String: s,
		Valid:  true,
	}
}
