package usecase

import "liokoredu/application/user"

type UserUseCase struct {
	repo user.Repository
}

func NewUserUseCase(u user.Repository) user.UseCase {
	return &UserUseCase{repo: u}
}
