package http

import (
	"liokoredu/application/user"

	"github.com/labstack/echo"
)

type UserHandler struct {
	uc user.UseCase
}

func CreateUserHandler(e *echo.Echo, uc user.UseCase) {
	_ = UserHandler{
		uc: uc,
	}

}
