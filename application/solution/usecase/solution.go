package usecase

import (
	"liokoredu/application/models"
	"liokoredu/application/solution"
)

type SolutionUseCase struct {
	repo solution.Repository
}

func (sd *SolutionUseCase) GetSolutions(taskId uint64) (models.Solutions, error) {
	return sd.repo.GetSolutions(taskId)
}

func (sd *SolutionUseCase) UpdateSolution(id uint64, code int, tests int) error {
	return sd.repo.UpdateSolution(id, code, tests)
}

func (s *SolutionUseCase) InsertSolution(taskId uint64, testsTotal int) (uint64, error) {
	return s.repo.InsertSolution(taskId, testsTotal)
}

func NewSolutionUseCase(s solution.Repository) solution.UseCase {
	return &SolutionUseCase{repo: s}
}
