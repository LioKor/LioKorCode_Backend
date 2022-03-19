package usecase

import (
	"liokoredu/application/models"
	"liokoredu/application/user"
	"liokoredu/pkg/constants"
	"liokoredu/pkg/generators"
	"net/http"
	"time"

	"github.com/labstack/echo"
)

type UserUseCase struct {
	repo user.Repository
}

// GetUserByUid implements user.UseCase
func (uuc *UserUseCase) GetUserByUid(uid uint64) (*models.User, error) {
	return uuc.repo.GetUserByUid(uid)
}

// DeleteSession implements user.UseCase
func (uuc *UserUseCase) DeleteSession(token string) error {
	return uuc.repo.DeleteSession(token)
}

// LoginUser implements user.UseCase
func (uuc *UserUseCase) LoginUser(usr models.UserAuth) (uint64, error) {
	u, err := uuc.repo.CheckUser(usr)
	if err != nil {
		return 0, err
	}
	if u == nil {
		return 0, echo.NewHTTPError(http.StatusForbidden, "invalid login or password")
	}
	if !(generators.CheckHashedPassword(u.Password, usr.Password)) {
		return 0, echo.NewHTTPError(http.StatusForbidden, "invalid login or password")
	}

	return u.Id, nil
}

// CreateUser implements user.UseCase
func (uuc *UserUseCase) CreateUser(usr models.User) (uint64, error) {
	u, err := uuc.GetUserByUsernameOrEmail(usr.Username, usr.Email)
	if err != nil {
		return 0, echo.NewHTTPError(500, err)
	}

	if u != nil {
		return 0, echo.NewHTTPError(409, err)
	}

	location, _ := time.LoadLocation("Europe/London")
	usr.JoinedDate = time.Now().In(location)

	usr.Password = generators.HashPassword(usr.Password)

	uid, err := uuc.repo.InsertUser(usr)
	if err != nil {
		return 0, echo.NewHTTPError(400, err)
	}

	return uid, nil
}

// GetUserByUsernameOrEmail implements user.UseCase
func (uuc *UserUseCase) GetUserByUsernameOrEmail(username string, email string) (*models.User, error) {
	return uuc.repo.GetUserByUsernameOrEmail(username, email)
}

func (uuc *UserUseCase) StoreSession(uid uint64) (string, error) {
	token := generators.CreateCookieValue(constants.CookieLength)
	return token, uuc.repo.StoreSession(token, uid)
}

// CheckSession implements user.UseCase
func (uuc *UserUseCase) CheckSession(token string) (uint64, error) {
	value, err := uuc.repo.CheckSession(token)
	if err != nil {
		return 0, err
	}

	if value != nil {
		return *value, nil
	}

	return 0, nil
}

func (uuc *UserUseCase) GetId(token string) (uint64, error) {
	return uuc.repo.GetId(token)
}

func NewUserUseCase(u user.Repository) user.UseCase {
	return &UserUseCase{repo: u}
}
