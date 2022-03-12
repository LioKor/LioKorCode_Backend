package http

import (
	"liokoredu/application/task"
	"liokoredu/pkg/constants"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
	"github.com/labstack/gommon/log"
	"github.com/mailru/easyjson"
)

type TaskHandler struct {
	uc task.UseCase
}

func CreateTaskHandler(e *echo.Echo,
	uc task.UseCase) {
	taskHandler := TaskHandler{
		uc: uc,
	}
	e.GET("/api/v1/task/:id", taskHandler.getTask)

}

func (th *TaskHandler) getTask(c echo.Context) error {
	defer c.Request().Body.Close()

	id := c.Param(constants.IdKey)
	n, _ := strconv.ParseUint(string(id), 10, 64)

	t, err := th.uc.GetTask(n)

	if err != nil {
		return err
	}

	if _, err = easyjson.MarshalToWriter(t, c.Response().Writer); err != nil {
		log.Error(c, err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return nil
}
