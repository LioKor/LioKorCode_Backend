package solution

import (
	"liokoredu/application/models"
	"time"
)

type Repository interface {
	InsertSolution(taskId uint64, uid uint64, code string, testsTotal int, receivedTime time.Time) (uint64, error)
	UpdateSolution(id uint64, code int, tests int) error
	DeleteSolution(id uint64, uid uint64) error
	GetSolutions(taskId uint64, uid uint64) (models.SolutionsSQL, error)
	GetSolution(id uint64, taskId uint64, uid uint64) (models.SolutionSQL, error)
}
