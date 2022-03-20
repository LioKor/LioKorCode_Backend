package usecase

import (
	"liokoredu/application/models"
	"liokoredu/application/solution"
	"time"
)

type SolutionUseCase struct {
	repo solution.Repository
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

func (sd *SolutionUseCase) UpdateSolution(id uint64, code int, tests int) error {
	return sd.repo.UpdateSolution(id, code, tests)
}

func (s *SolutionUseCase) InsertSolution(taskId uint64, uid uint64, code string, testsTotal int) (uint64, error) {
	location, _ := time.LoadLocation("Europe/London")

	received := time.Now().In(location)
	return s.repo.InsertSolution(taskId, uid, code, testsTotal, received)
}

func NewSolutionUseCase(s solution.Repository) solution.UseCase {
	return &SolutionUseCase{repo: s}
}
