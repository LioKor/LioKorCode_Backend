package models

type Solution struct {
	SourceCode string `json:"sourceCode"`
}

type SolutionSend struct {
	Id         uint64     `json:"id"`
	SourceCode string     `json:"sourceCode"`
	Tests      InputTests `json:"tests"`
}

//easyjson:json
type InputTests [][]string

type SolutionUpdate struct {
	Code         int    `json:"checkResult"`
	CheckMessage string `json:"checkMessage"`
	Passed       int    `json:"testsPassed"`
	TestsTotal   int    `json:"testsTotal"`
}

type SolutionSQL struct {
	Id          uint64
	TaskId      uint64
	CheckResult int
	TestsPassed int
	TestsTotal  int
}

//easyjson:json
type Solutions []SolutionSQL

type SolutionId struct {
	Id uint64 `json:"id"`
}
