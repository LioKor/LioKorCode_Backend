package usecase

import (
	"liokoredu/application/models"
	"liokoredu/application/task"
)

type TaskUseCase struct {
	repo task.Repository
}

// FindTasksFull implements task.UseCase
func (tuc *TaskUseCase) FindTasksFull(str string, useSolved bool, solved bool, useMine bool, mine bool, uid uint64, page int, count int) (models.ShortTasks, int, error) {
	tsks, num, err := tuc.repo.FindTasksFull(str, useSolved, solved, useMine, mine, uid, page, count)
	if err != nil {
		return models.ShortTasks{}, num, err
	}
	if uid == 0 {
		return *tsks, num, nil
	}
	tsksArr := models.ShortTasks{}
	for _, tsk := range *tsks {
		isDone, err := tuc.repo.IsCleared(tsk.Id, uid)
		if err != nil {
			return models.ShortTasks{}, num, err
		}
		tsk.IsCleared = isDone
		tsksArr = append(tsksArr, tsk)
	}

	return tsksArr, num, nil
}

func (tuc *TaskUseCase) IsCleared(taskId uint64, uid uint64) (bool, error) {
	return tuc.repo.IsCleared(taskId, uid)
}

func (tuc *TaskUseCase) MarkTaskDone(id uint64, uid uint64) error {
	res, err := tuc.repo.IsCleared(id, uid)
	if err != nil {
		return err
	}
	if res {
		return nil
	}
	return tuc.repo.MarkTaskDone(id, uid)
}

func (tuc *TaskUseCase) UpdateTask(id uint64, t *models.TaskNew) error {
	tsk := t.ConvertNewTaskToTaskSQL()
	tsk.Id = id
	return tuc.repo.UpdateTask(tsk)
}

func (tuc *TaskUseCase) DeleteTask(id uint64, uid uint64) error {
	return tuc.repo.DeleteTask(id, uid)
}

func (tuc *TaskUseCase) FindTasks(str string, uid uint64, page int, count int) (models.ShortTasks, error) {
	tsks, err := tuc.repo.FindTasks(str, page, count)
	if err != nil {
		return models.ShortTasks{}, err
	}
	if uid == 0 {
		return *tsks, nil
	}
	tsksArr := models.ShortTasks{}
	for _, tsk := range *tsks {
		isDone, err := tuc.repo.IsCleared(tsk.Id, uid)
		if err != nil {
			return models.ShortTasks{}, err
		}
		tsk.IsCleared = isDone
		tsksArr = append(tsksArr, tsk)
	}

	return tsksArr, nil
}

func (tuc *TaskUseCase) GetTasks(uid uint64, page int, count int) (models.ShortTasks, error) {
	tsks, err := tuc.repo.GetTasks(page, count)
	if err != nil {
		return models.ShortTasks{}, err
	}
	if uid == 0 {
		return *tsks, nil
	}
	tsksArr := models.ShortTasks{}
	for _, tsk := range *tsks {
		isDone, err := tuc.repo.IsCleared(tsk.Id, uid)
		if err != nil {
			return models.ShortTasks{}, err
		}
		tsk.IsCleared = isDone
		tsksArr = append(tsksArr, tsk)
	}

	return tsksArr, nil
}

func (tuc *TaskUseCase) GetSolvedTasks(uid uint64, page int, count int) (models.ShortTasks, error) {
	tsks, err := tuc.repo.GetSolvedTasks(uid, page, count)
	if err != nil {
		return models.ShortTasks{}, err
	}
	if uid == 0 {
		return *tsks, nil
	}
	tsksArr := models.ShortTasks{}
	for _, tsk := range *tsks {
		tsk.IsCleared = true
		tsksArr = append(tsksArr, tsk)
	}

	return tsksArr, nil
}

func (tuc *TaskUseCase) GetUnsolvedTasks(uid uint64, page int, count int) (models.ShortTasks, error) {
	tsks, err := tuc.repo.GetUnsolvedTasks(uid, page, count)
	if err != nil {
		return models.ShortTasks{}, err
	}
	if uid == 0 {
		return *tsks, nil
	}
	tsksArr := models.ShortTasks{}
	for _, tsk := range *tsks {
		tsk.IsCleared = false
		tsksArr = append(tsksArr, tsk)
	}

	return tsksArr, nil
}

func (uc *TaskUseCase) GetUserTasks(uid uint64, page int, count int) (models.ShortTasks, error) {
	tsks, err := uc.repo.GetUserTasks(uid, page, count)
	if err != nil {
		return models.ShortTasks{}, err
	}

	return *tsks, nil
}

func (uc *TaskUseCase) CreateTask(t *models.TaskNew) (uint64, error) {
	return uc.repo.CreateTask(t.ConvertNewTaskToTaskSQL())
}

func NewTaskUseCase(t task.Repository) task.UseCase {
	return &TaskUseCase{repo: t}
}

func (uc TaskUseCase) GetTask(id uint64, uid uint64, forCheck bool) (*models.Task, error) {
	t, err := uc.repo.GetTask(id)
	if err != nil {
		return &models.Task{}, err
	}

	isCreator := false
	if t.Creator == uid || forCheck {
		isCreator = true
	}

	isCleared, err := uc.repo.IsCleared(id, uid)
	if err != nil {
		return &models.Task{}, err
	}

	tsk := t.ConvertToTask(isCreator, isCleared)
	return tsk, nil
}

func (uc TaskUseCase) GetPages(count int) (int, error) {
	n, err := uc.repo.GetPages()
	if err != nil {
		return 0, err
	}
	if n%count == 0 {
		return n / count, nil
	}
	return n/count + 1, nil
}
