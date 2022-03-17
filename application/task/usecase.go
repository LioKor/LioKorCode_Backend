package task

import "liokoredu/application/models"

type UseCase interface {
	GetTask(id uint64) (*models.Task, error)
	GetTasks(page int) (*models.Tasks, error)
	CreateTask(t *models.TaskNew) (uint64, error)
}
