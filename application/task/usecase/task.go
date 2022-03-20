package usecase

import (
	"liokoredu/application/models"
	"liokoredu/application/task"
)

type TaskUseCase struct {
	repo task.Repository
}

// DeleteTask implements task.UseCase
func (tuc *TaskUseCase) DeleteTask(id uint64, uid uint64) error {
	return tuc.repo.DeleteTask(id, uid)
}

func (uc *TaskUseCase) GetTasks(page int) (*models.ShortTasks, error) {
	tsks, err := uc.repo.GetTasks(page)
	if err != nil {
		return &models.ShortTasks{}, err
	}

	return tsks, nil
}

// CreateTask implements task.UseCase
func (uc *TaskUseCase) CreateTask(t *models.TaskNew) (uint64, error) {
	return uc.repo.CreateTask(t.ConvertNewTaskToTaskSQL())
}

func NewTaskUseCase(t task.Repository) task.UseCase {
	return &TaskUseCase{repo: t}
}

func (uc TaskUseCase) GetTask(id uint64) (*models.Task, error) {
	t, err := uc.repo.GetTask(id)
	if err != nil {
		return &models.Task{}, err
	}

	tsk := t.ConvertToTask()
	return tsk, nil
}
