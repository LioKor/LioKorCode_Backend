package http

import (
	"liokoredu/application/models"
	"liokoredu/application/task"
	"liokoredu/pkg/constants"
	"log"
	"net/http"
	"strconv"

	"github.com/labstack/echo"

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
	e.GET("/api/v1/tasks/:id", taskHandler.getTask)
	e.POST("/api/v1/tasks/create", taskHandler.createTask)
	e.GET("/api/v1/tasks", taskHandler.getTasks)

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
		log.Println(c, err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return nil
}

func (th *TaskHandler) getTasks(c echo.Context) error {
	defer c.Request().Body.Close()

	page := c.QueryParams().Get(constants.PageKey)
	p, _ := strconv.Atoi(string(page))
	if p == 0 {
		p = 1
	}

	tsks, err := th.uc.GetTasks(p)
	if err != nil {
		return err
	}

	if _, err = easyjson.MarshalToWriter(tsks, c.Response().Writer); err != nil {
		log.Println("task handler: getTasks: error marshaling answer to writer", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return nil
}

func (th *TaskHandler) createTask(c echo.Context) error {
	defer c.Request().Body.Close()

	tn := &models.TaskNew{}

	if err := easyjson.UnmarshalFromReader(c.Request().Body, tn); err != nil {
		log.Println("task handler: createTask: error unmarshaling task from reader", err)
		return echo.NewHTTPError(http.StatusTeapot, err.Error())
	}

	tid, err := th.uc.CreateTask(tn)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if _, err = easyjson.MarshalToWriter(&models.ReturnId{Id: tid}, c.Response().Writer); err != nil {
		log.Println("task handler: createTask: error marshaling answer to writer", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return nil
}
