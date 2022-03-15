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
	Id          uint64 `json:"checkResult"`
	TaskId      uint64
	CheckResult int
	TestsPassed int
	TestsTotal  int
}

type SolutionOne struct {
	Id          uint64 `json:"id"`
	CheckResult int    `json:"checkResult"`
	TestsPassed int    `json:"testsPassed"`
	TestsTotal  int    `json:"testsTotal"`
}

//easyjson:json
type Solutions []SolutionOne

//easyjson:json
type SolutionsSQL []SolutionSQL

type SolutionId struct {
	Id uint64 `json:"id"`
}

func (slnsSQL SolutionsSQL) ConvertToJson() Solutions {
	res := Solutions{}
	for _, elem := range slnsSQL {
		newElem := SolutionOne{}
		newElem.Id = elem.Id
		newElem.CheckResult = elem.CheckResult
		newElem.TestsPassed = elem.TestsPassed
		newElem.TestsTotal = elem.TestsTotal
		res = append(res, newElem)
	}
	return res
}