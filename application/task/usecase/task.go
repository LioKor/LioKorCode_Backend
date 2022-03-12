package usecase

import (
	"liokoredu/application/models"
	"liokoredu/application/task"
)

type TaskUseCase struct {
	repo task.Repository
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
