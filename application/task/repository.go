package task

import "liokoredu/application/models"

type Repository interface {
	GetTask(id uint64) (*models.TaskSQL, error)
}
