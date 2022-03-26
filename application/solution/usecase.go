package solution

import "liokoredu/application/models"

type UseCase interface {
	InsertSolution(taskId uint64, uid uint64, code string, makefile string, testsTotal int) (uint64, error)
	UpdateSolution(id uint64, code int, tests int) error
	DeleteSolution(id uint64, uid uint64) error
	GetSolutions(taskId uint64, uid uint64) (models.Solutions, error)
	GetSolution(solId uint64, taskId uint64, uid uint64) (models.SolutionFull, error)
}
