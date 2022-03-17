package task

import "liokoredu/application/models"

type Repository interface {
	GetTask(id uint64) (*models.TaskSQL, error)
	GetTasks(page int) (*models.TasksSQL, error)
	CreateTask(t *models.TaskSQL) (uint64, error)
}
