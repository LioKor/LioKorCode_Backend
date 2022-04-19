package task

import "liokoredu/application/models"

type UseCase interface {
	GetTask(id uint64, uid uint64, forCheck bool) (*models.Task, error)
	GetTasks(uid uint64, page int) (models.ShortTasks, error)
	IsCleared(taskId uint64, uid uint64) (bool, error)
	GetUserTasks(uid uint64, page int) (models.ShortTasks, error)
	CreateTask(t *models.TaskNew) (uint64, error)
	DeleteTask(id uint64, uid uint64) error
	UpdateTask(id uint64, t *models.TaskNew) error
	MarkTaskDone(id uint64, uid uint64) error
}
