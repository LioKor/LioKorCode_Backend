package models

import (
	"encoding/json"
	"time"
)

type Solution struct {
	SourceCode map[string]interface{} `json:"sourceCode"`
}

type Files struct {
	Files []string `json:"files"`
}

type FileName struct {
	Filename string `json:"filename"`
}

type SolutionFile struct {
	Filename string `json:"filename"`
	Text     string `json:"text"`
}

//easyjson:json
type SolutionFiles []SolutionFile

type RedactorFile struct {
	Text string `json:"text"`
}

type SolutionSend struct {
	Id         uint64                 `json:"id"`
	SourceCode map[string]interface{} `json:"sourceCode"`
	Tests      InputTests             `json:"tests"`
}

//easyjson:json
type InputTests [][]string

type SolutionUpdate struct {
	Code            int       `json:"checkResult"`
	CheckMessage    string    `json:"checkMessage"`
	CheckedDateTime time.Time `json:"checkedDatetime"`
	CheckTime       float32   `json:"checkTime"`
	CompileTime     float32   `json:"buildTime"`
	Passed          int       `json:"testsPassed"`
	TestsTotal      int       `json:"testsTotal"`
}

type SolutionSQL struct {
	Id               uint64
	ReceivedDateTime time.Time
	CheckedDateTime  time.Time
	SourceCode       string
	TaskId           uint64
	CheckResult      int
	CheckTime        float32
	CompileTime      float32
	CheckMessage     string
	TestsPassed      int
	TestsTotal       int
	Uid              uint64
}

type SolutionOne struct {
	Id               uint64    `json:"id"`
	SourceCode       string    `json:"sourceCode"`
	ReceivedDateTime time.Time `json:"receivedDatetime"`
	CheckedDateTime  time.Time `json:"checkedDatetime"`
	CheckResult      int       `json:"checkResult"`
	CheckTime        float32   `json:"checkTime"`
	CheckMessage     string    `json:"checkMessage"`
	CompileTime      float32   `json:"compileTime"`
	TestsPassed      int       `json:"testsPassed"`
	TestsTotal       int       `json:"testsTotal"`
}

type SolutionFull struct {
	Id               uint64                 `json:"id"`
	SourceCode       map[string]interface{} `json:"sourceCode"`
	ReceivedDateTime time.Time              `json:"receivedDatetime"`
	CheckedDateTime  time.Time              `json:"checkedDatetime"`
	CheckResult      int                    `json:"checkResult"`
	CheckMessage     string                 `json:"checkMessage"`
	CheckTime        float32                `json:"checkTime"`
	CompileTime      float32                `json:"compileTime"`
	Tests            TestResults            `json:"tests"`
	TestsPassed      int                    `json:"testsPassed"`
	TestsTotal       int                    `json:"testsTotal"`
}

//easyjson:json
type TestResults []TestResult

type TestResult struct {
	Stdin  string `json:"stdin"`
	Stdout string `json:"stdout"`
	Passed bool   `json:"passed"`
}

//easyjson:json
type Solutions []SolutionOne

//easyjson:json
type SolutionsSQL []SolutionSQL

type ReturnId struct {
	Id uint64 `json:"id"`
}

func (slnsSQL SolutionsSQL) ConvertToJson() Solutions {
	res := Solutions{}
	for _, elem := range slnsSQL {
		res = append(res, elem.ConvertToShort())
	}
	return res
}

func (slnSQL SolutionSQL) ConvertToShort() SolutionOne {
	newElem := SolutionOne{}

	newElem.Id = slnSQL.Id

	newElem.ReceivedDateTime = slnSQL.ReceivedDateTime
	newElem.CheckedDateTime = slnSQL.CheckedDateTime

	newElem.CheckResult = slnSQL.CheckResult
	newElem.CheckMessage = slnSQL.CheckMessage

	newElem.CheckTime = slnSQL.CheckTime
	newElem.CompileTime = slnSQL.CompileTime

	newElem.TestsPassed = slnSQL.TestsPassed
	newElem.TestsTotal = slnSQL.TestsTotal
	return newElem
}

func (slnSQL SolutionSQL) ConvertToFull(tsk *Task) SolutionFull {
	newElem := SolutionFull{}

	newElem.Id = slnSQL.Id

	newElem.ReceivedDateTime = slnSQL.ReceivedDateTime
	newElem.CheckedDateTime = slnSQL.CheckedDateTime

	newElem.CheckResult = slnSQL.CheckResult
	newElem.CheckMessage = slnSQL.CheckMessage

	newElem.CheckTime = slnSQL.CheckTime
	newElem.CompileTime = slnSQL.CompileTime

	newElem.TestsPassed = slnSQL.TestsPassed
	newElem.TestsTotal = slnSQL.TestsTotal

	i := 0
	tests := TestResults{}
	for ; i < slnSQL.TestsPassed; i++ {
		test := TestResult{}
		test.Stdin = tsk.Tests[i][0]
		test.Stdout = tsk.Tests[i][1]
		test.Passed = true
		tests = append(tests, test)
	}

	if slnSQL.TestsPassed < slnSQL.TestsTotal {
		test := TestResult{}

		test.Stdin = tsk.Tests[i][0]
		test.Stdout = tsk.Tests[i][1]
		test.Passed = false
		tests = append(tests, test)
	}

	newElem.Tests = tests

	var code map[string]interface{}
	_ = json.Unmarshal([]byte(slnSQL.SourceCode), &code)
	newElem.SourceCode = code

	return newElem
}
