package http

import (
	"liokoredu/application/models"
	"liokoredu/application/user"
	"liokoredu/pkg/constants"
	"liokoredu/pkg/generators"
	"log"
	"net/http"
	"time"

	"github.com/labstack/echo"
	"github.com/mailru/easyjson"
)

type UserHandler struct {
	uc user.UseCase
}

func CreateUserHandler(e *echo.Echo, uc user.UseCase) {
	userHandler := UserHandler{
		uc: uc,
	}

	e.POST("/api/v1/users", userHandler.createUser)
	e.POST("/api/v1/user/auth", userHandler.login)
	e.GET("/api/v1/user", userHandler.getUserProfile)
	e.DELETE("/api/v1/user/session", userHandler.logout)
}

func (uh *UserHandler) getUserProfile(c echo.Context) error {
	defer c.Request().Body.Close()

	cookie, err := c.Cookie(constants.SessionCookieName)
	if err != nil && cookie != nil {
		log.Println("user handler: getUserProfile: error getting cookie")
		return echo.NewHTTPError(http.StatusBadRequest, "error getting cookie")
	}

	if cookie == nil {
		log.Println("user handler: getUserProfile: no cookie")
		return echo.NewHTTPError(http.StatusUnauthorized, "Not authenticated")
	}

	uid, err := uh.uc.CheckSession(cookie.Value)
	if err != nil {
		return err
	}

	usr, err := uh.uc.GetUserByUid(uid)
	if err != nil {
		return err
	}

	if _, err = easyjson.MarshalToWriter(usr, c.Response().Writer); err != nil {
		log.Println(c, err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return nil
}

func (uh *UserHandler) createUser(c echo.Context) error {
	defer c.Request().Body.Close()

	cookie, err := c.Cookie(constants.SessionCookieName)
	if err != nil && cookie != nil {
		log.Println("user handler: createUser: error getting cookie")
		return echo.NewHTTPError(http.StatusBadRequest, "error getting cookie")
	}

	if cookie != nil {
		uid, err := uh.uc.CheckSession(cookie.Value)
		if err != nil {
			return err
		}

		if uid != 0 {
			return echo.NewHTTPError(http.StatusBadRequest, "user is already logged in")
		}
	}

	newUser := &models.User{}

	err = easyjson.UnmarshalFromReader(c.Request().Body, newUser)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	uid, err := uh.uc.CreateUser(*newUser)
	if err != nil {
		return err
	}

	token, err := uh.uc.StoreSession(uid)
	if err != nil {
		return err
	}

	cookie = generators.CreateCookieWithValue(token)
	c.SetCookie(cookie)
	return nil
}

func (uh *UserHandler) login(c echo.Context) error {
	defer c.Request().Body.Close()

	cookie, err := c.Cookie(constants.SessionCookieName)
	if err != nil && cookie != nil {
		log.Println("user handler: login: error getting cookie")
		return echo.NewHTTPError(http.StatusBadRequest, "error getting cookie")
	}

	if cookie != nil {
		uid, err := uh.uc.CheckSession(cookie.Value)
		if err != nil {
			return err
		}

		if uid != 0 {
			return echo.NewHTTPError(http.StatusBadRequest, "user is already logged in")
		}
	}

	usr := &models.UserAuth{}

	err = easyjson.UnmarshalFromReader(c.Request().Body, usr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	uid, err := uh.uc.LoginUser(*usr)
	if err != nil {
		return err
	}

	token, err := uh.uc.StoreSession(uid)
	if err != nil {
		return err
	}

	cookie = generators.CreateCookieWithValue(token)
	c.SetCookie(cookie)
	return nil
}

func (uh *UserHandler) logout(c echo.Context) error {
	defer c.Request().Body.Close()

	cookie, err := c.Cookie(constants.SessionCookieName)
	if err != nil && cookie != nil {
		log.Println("user handler: logout: error getting cookie")
		return echo.NewHTTPError(http.StatusBadRequest, "error getting cookie")
	}

	if cookie == nil {
		log.Println("user handler: logout: no cookie")
		return echo.NewHTTPError(http.StatusUnauthorized, "Not authenticated")
	}

	err = uh.uc.DeleteSession(cookie.Value)
	if err != nil {
		return err
	}

	cookie.Expires = time.Now().AddDate(0, 0, -1)
	c.SetCookie(cookie)

	return nil
}
