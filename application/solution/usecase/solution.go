package usecase

import (
	"liokoredu/application/models"
	"liokoredu/application/solution"
	"liokoredu/application/task"
	"time"
)

type SolutionUseCase struct {
	repo   solution.Repository
	ucTask task.UseCase
}

// GetSolution implements solution.UseCase
func (suc *SolutionUseCase) GetSolution(solId uint64, taskId uint64, uid uint64) (models.SolutionFull, error) {
	sln, err := suc.repo.GetSolution(solId, taskId, uid)
	if err != nil {
		return models.SolutionFull{}, err
	}
	tsk, err := suc.ucTask.GetTask(taskId, uid, true)
	if err != nil {
		return models.SolutionFull{}, err
	}

	return sln.ConvertToFull(tsk), nil
}

// DeleteSolution implements solution.UseCase
func (suc *SolutionUseCase) DeleteSolution(id uint64, uid uint64) error {
	return suc.repo.DeleteSolution(id, uid)
}

func (sd *SolutionUseCase) GetSolutions(taskId uint64, uid uint64) (models.Solutions, error) {
	slnsSQL, err := sd.repo.GetSolutions(taskId, uid)
	if err != nil {
		return models.Solutions{}, err
	}

	return slnsSQL.ConvertToJson(), nil
}

func (sd *SolutionUseCase) UpdateSolution(id uint64, upd models.SolutionUpdate) error {
	location, _ := time.LoadLocation("Europe/London")

	checked := time.Now().In(location)
	upd.CheckedDateTime = checked
	return sd.repo.UpdateSolution(id, &upd)
}

func (s *SolutionUseCase) InsertSolution(taskId uint64, uid uint64, code map[string]interface{}, testsTotal int) (uint64, error) {
	location, _ := time.LoadLocation("Europe/London")

	received := time.Now().In(location)
	return s.repo.InsertSolution(taskId, uid, code, testsTotal, received)
}

func NewSolutionUseCase(s solution.Repository, t task.UseCase) solution.UseCase {
	return &SolutionUseCase{repo: s, ucTask: t}
}
