package solution

import "liokoredu/application/models"

type Repository interface {
	InsertSolution(taskId uint64, testsTotal int) (uint64, error)
	UpdateSolution(id uint64, code int, tests int) error
	GetSolutions(taskId uint64) (models.SolutionsSQL, error)
}
