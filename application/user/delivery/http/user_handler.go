package http

import (
	"liokoredu/application/user"

	"github.com/labstack/echo"
)

type UserHandler struct {
	uc user.UseCase
}

func CreateUserHandler(e *echo.Echo, uc user.UseCase) {
	userHandler := UserHandler{
		uc: uc,
	}

	e.GET("/api/v1/users/get", userHandler.getUser)
	e.POST("/api/v1/users/store", userHandler.storeSession)

}

func (uh *UserHandler) getUser(c echo.Context) error {
	defer c.Request().Body.Close()

	return nil
}

func (uh *UserHandler) storeSession(c echo.Context) error {
	defer c.Request().Body.Close()

	return nil
}
