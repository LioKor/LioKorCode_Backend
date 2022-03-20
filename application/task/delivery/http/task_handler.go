package http

import (
	"liokoredu/application/models"
	"liokoredu/application/task"
	"liokoredu/application/user"
	"liokoredu/pkg/constants"
	"log"
	"net/http"
	"strconv"

	"github.com/labstack/echo"

	"github.com/mailru/easyjson"
)

type TaskHandler struct {
	uc  task.UseCase
	uuc user.UseCase
}

func CreateTaskHandler(e *echo.Echo,
	uc task.UseCase, uuc user.UseCase) {
	taskHandler := TaskHandler{
		uc:  uc,
		uuc: uuc,
	}
	e.GET("/api/v1/tasks/:id", taskHandler.getTask)
	e.POST("/api/v1/tasks", taskHandler.createTask)
	e.GET("/api/v1/tasks", taskHandler.getTasks)
	e.DELETE("/api/v1/tasks/:id", taskHandler.deleteTask)
	e.PUT("/api/v1/tasks/:id", taskHandler.updateTask)
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

	cookie, err := c.Cookie(constants.SessionCookieName)
	if err != nil && cookie != nil {
		log.Println("user handler: createTask: error getting cookie")
		return echo.NewHTTPError(http.StatusBadRequest, "error getting cookie")
	}

	if cookie == nil {
		log.Println("user handler: createTask: no cookie")
		return echo.NewHTTPError(http.StatusUnauthorized, "Not authenticated")
	}

	uid, err := th.uuc.CheckSession(cookie.Value)
	if err != nil {
		return err
	}

	if uid == 0 {
		log.Println("user handler: createTask: uid 0")
		return echo.NewHTTPError(http.StatusUnauthorized, "Not authenticated")
	}

	tn := &models.TaskNew{}

	if err := easyjson.UnmarshalFromReader(c.Request().Body, tn); err != nil {
		log.Println("task handler: createTask: error unmarshaling task from reader", err)
		return echo.NewHTTPError(http.StatusTeapot, err.Error())
	}

	tn.Creator = uid

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

func (th *TaskHandler) deleteTask(c echo.Context) error {
	defer c.Request().Body.Close()

	cookie, err := c.Cookie(constants.SessionCookieName)
	if err != nil && cookie != nil {
		log.Println("user handler: deleteTask: error getting cookie")
		return echo.NewHTTPError(http.StatusBadRequest, "error getting cookie")
	}

	if cookie == nil {
		log.Println("user handler: deleteTask: no cookie")
		return echo.NewHTTPError(http.StatusUnauthorized, "Not authenticated")
	}

	uid, err := th.uuc.CheckSession(cookie.Value)
	if err != nil {
		return err
	}

	if uid == 0 {
		log.Println("user handler: deleteTask: uid 0")
		return echo.NewHTTPError(http.StatusUnauthorized, "Not authenticated")
	}

	id := c.Param(constants.IdKey)
	iid, _ := strconv.ParseUint(string(id), 10, 64)

	err = th.uc.DeleteTask(iid, uid)
	if err != nil {
		return err
	}

	return nil
}

func (th *TaskHandler) updateTask(c echo.Context) error {
	defer c.Request().Body.Close()

	cookie, err := c.Cookie(constants.SessionCookieName)
	if err != nil && cookie != nil {
		log.Println("user handler: updateTask: error getting cookie")
		return echo.NewHTTPError(http.StatusBadRequest, "error getting cookie")
	}

	if cookie == nil {
		log.Println("user handler: updateTask: no cookie")
		return echo.NewHTTPError(http.StatusUnauthorized, "Not authenticated")
	}

	uid, err := th.uuc.CheckSession(cookie.Value)
	if err != nil {
		return err
	}

	if uid == 0 {
		log.Println("user handler: updateTask: uid 0")
		return echo.NewHTTPError(http.StatusUnauthorized, "Not authenticated")
	}

	tn := &models.TaskNew{}

	if err := easyjson.UnmarshalFromReader(c.Request().Body, tn); err != nil {
		log.Println("task handler: updateTask: error unmarshaling task from reader", err)
		return echo.NewHTTPError(http.StatusTeapot, err.Error())
	}

	tn.Creator = uid
	id := c.Param(constants.IdKey)
	iid, _ := strconv.ParseUint(string(id), 10, 64)

	err = th.uc.UpdateTask(iid, tn)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return nil
}
