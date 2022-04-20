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

// UpdateUserAvatar implements user.UseCase
func (uuc *UserUseCase) UpdateUserAvatar(uid uint64, avt *models.Avatar) error {
	return uuc.repo.UpdateUserAvatar(uid, avt)
}

// UpdatePassword implements user.UseCase
func (uuc *UserUseCase) UpdatePassword(uid uint64, data models.PasswordNew) error {
	usr, err := uuc.GetUserByUid(uid)
	if err != nil {
		return err
	}

	if !generators.CheckHashedPassword(usr.Password, data.Old) {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid password data provided")
	}

	return uuc.repo.UpdatePassword(uid, generators.HashPassword(data.New))
}

// UpdateUser implements user.UseCase
func (uuc *UserUseCase) UpdateUser(uid uint64, usr models.UserUpdate) error {
	if !usr.Validate() {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid user data provided")
	}
	usrs, err := uuc.repo.GetUserByEmailSubmitted(usr.Email)
	if err != nil {
		return err
	}

	if len(*usrs) != 0 && ((*usrs)[0].Id != uid) {
		return echo.NewHTTPError(http.StatusBadRequest, "Email has already been taken and verified")
	}

	return uuc.repo.UpdateUser(uid, usr)
}

// GetUserByUid implements user.UseCase
func (uuc *UserUseCase) GetUserByUid(uid uint64) (*models.User, error) {
	usr, err := uuc.repo.GetUserByUid(uid)
	if err != nil {
		return &models.User{}, err
	}
	(*usr).JWT = generators.CreateToken(usr)

	return usr, nil
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
		return 0, echo.NewHTTPError(500, err.Error())
	}

	if u != nil {
		return 0, echo.NewHTTPError(409, "user with this usermame or email already exists")
	}

	location, _ := time.LoadLocation("Europe/London")
	usr.JoinedDate = time.Now().In(location)

	usr.Password = generators.HashPassword(usr.Password)

	uid, err := uuc.repo.InsertUser(usr)
	if err != nil {
		return 0, echo.NewHTTPError(400, err.Error())
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

/*
func (uuc *UserUseCase) GetId(token string) (uint64, error) {
	return uuc.repo.GetId(token)
}
*/

func NewUserUseCase(u user.Repository) user.UseCase {
	return &UserUseCase{repo: u}
}
