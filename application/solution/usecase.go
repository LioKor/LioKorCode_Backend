package solution

import "liokoredu/application/models"

type UseCase interface {
	InsertSolution(taskId uint64, uid uint64, code map[string]interface{}, testsTotal int) (uint64, error)
	UpdateSolution(id uint64, upd models.SolutionUpdate) error
	DeleteSolution(id uint64, uid uint64) error
	GetSolutions(taskId uint64, uid uint64) (models.Solutions, error)
	GetSolution(solId uint64, taskId uint64, uid uint64) (models.SolutionFull, error)
}
