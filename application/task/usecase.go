package task

import "liokoredu/application/models"

type UseCase interface {
	GetTask(id uint64) (*models.Task, error)
}
