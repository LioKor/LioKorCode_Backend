package middleware

import (
	"liokoredu/application/user"
	"liokoredu/pkg/constants"
	"log"
	"net/http"

	"github.com/labstack/echo"
)

type Auth struct {
	uuc user.UseCase
}

func NewAuth(uuc user.UseCase) Auth {
	return Auth{uuc: uuc}
}

func (a Auth) GetSession(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		cookie, err := ctx.Cookie(constants.SessionCookieName)
		if err != nil && cookie != nil {
			log.Println("middleware: GetSession: error getting cookie", err.Error())
			return echo.NewHTTPError(http.StatusBadRequest, "Error getting cookie")
		}

		if cookie == nil {
			log.Println("middleware: GetSession: not authenticated")
			return echo.NewHTTPError(http.StatusUnauthorized, "Not authenticated")
		}

		uid, err := a.uuc.CheckSession(cookie.Value)
		if err != nil {
			return err
		}

		if uid == 0 {
			log.Println("middleware: GetSession: uid 0")
			return echo.NewHTTPError(http.StatusUnauthorized, "Not authenticated")
		}

		ctx.Set(constants.SessionCookieName, cookie.Value)
		ctx.Set(constants.UserIdKey, uid)
		return next(ctx)
	}
}
