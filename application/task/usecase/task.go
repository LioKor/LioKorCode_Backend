package usecase

import (
	"liokoredu/application/models"
	"liokoredu/application/task"
)

type TaskUseCase struct {
	repo task.Repository
}

// UpdateTask implements task.UseCase
func (tuc *TaskUseCase) UpdateTask(id uint64, t *models.TaskNew) error {
	tsk := t.ConvertNewTaskToTaskSQL()
	tsk.Id = id
	return tuc.repo.UpdateTask(tsk)
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

func (uc *TaskUseCase) GetUserTasks(uid uint64, page int) (*models.ShortTasks, error) {
	tsks, err := uc.repo.GetUserTasks(uid, page)
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

func (uc TaskUseCase) GetTask(id uint64, uid uint64) (*models.Task, error) {
	t, err := uc.repo.GetTask(id)
	if err != nil {
		return &models.Task{}, err
	}

	isCreator := false
	if t.Creator == uid {
		isCreator = true
	}

	tsk := t.ConvertToTask(isCreator)
	return tsk, nil
}
