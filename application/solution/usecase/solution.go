package usecase

import "liokoredu/application/solution"

type SolutionUseCase struct {
	repo solution.Repository
}

func NewSolutionUseCase(s solution.Repository) solution.UseCase {
	return &SolutionUseCase{repo: s}
}
