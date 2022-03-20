package task

import "liokoredu/application/models"

type UseCase interface {
	GetTask(id uint64) (*models.Task, error)
	GetTasks(page int) (*models.ShortTasks, error)
	CreateTask(t *models.TaskNew) (uint64, error)
	DeleteTask(id uint64, uid uint64) error
	UpdateTask(id uint64, t *models.TaskNew) error
}
