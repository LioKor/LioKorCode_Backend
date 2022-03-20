package task

import "liokoredu/application/models"

type Repository interface {
	GetTask(id uint64) (*models.TaskSQL, error)
	GetTasks(page int) (*models.ShortTasks, error)
	CreateTask(t *models.TaskSQL) (uint64, error)
	DeleteTask(id uint64, uid uint64) error
	UpdateTask(t *models.TaskSQL) error
}
