package usecase

import (
	"liokoredu/application/user"
	"liokoredu/pkg/constants"
	"liokoredu/pkg/generators"
)

type UserUseCase struct {
	repo user.Repository
}

func (uuc *UserUseCase) StoreSession(uid uint64) (string, error) {
	token := generators.CreateCookieValue(constants.CookieLength)
	return token, uuc.repo.StoreSession(token, uid)
}

func (uuc *UserUseCase) GetId(token string) (uint64, error) {
	return uuc.repo.GetId(token)
}

func NewUserUseCase(u user.Repository) user.UseCase {
	return &UserUseCase{repo: u}
}
