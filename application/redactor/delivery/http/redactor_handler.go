package http

import (
	"liokoredu/application/microservices/redactor/client"
	"liokoredu/application/models"
	"liokoredu/application/server/middleware"
	"liokoredu/pkg/constants"
	"log"
	"net/http"

	"github.com/labstack/echo"
	"github.com/mailru/easyjson"
)

type RedactorHandler struct {
	rpcRedactor *client.RedactorClient
}

func CreateRedactorHandler(e *echo.Echo, rpcR *client.RedactorClient, a middleware.Auth) {

	redactorHandler := RedactorHandler{rpcRedactor: rpcR}

	e.POST("/api/v1/redactor", redactorHandler.CreateConnection, a.GetSession)
}

func (rh RedactorHandler) CreateConnection(c echo.Context) error {
	defer c.Request().Body.Close()

	uid := c.Get(constants.UserIdKey).(uint64)

	id, err, code := rh.rpcRedactor.CreateConnection(uid)
	if err != nil {
		log.Println(id, err, code)
		return err
	}

	if _, err = easyjson.MarshalToWriter(&models.IdValue{Id: id}, c.Response().Writer); err != nil {
		log.Println(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return nil
}
