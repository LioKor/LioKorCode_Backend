package http

import (
	"liokoredu/application/models"
	"liokoredu/application/server/middleware"
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
	uc task.UseCase, uuc user.UseCase, a middleware.Auth) {
	taskHandler := TaskHandler{
		uc:  uc,
		uuc: uuc,
	}
	e.GET("/api/v1/tasks/:id", taskHandler.getTask)

	e.POST("/api/v1/tasks", taskHandler.createTask, a.GetSession)
	e.GET("/api/v1/tasks", taskHandler.getTasks)
	e.GET("/api/v1/tasks/search", taskHandler.findTasks)
	e.GET("/api/v1/tasks/solved", taskHandler.getSolvedTasks, a.GetSession)
	e.GET("/api/v1/tasks/unsolved", taskHandler.getUnsolvedTasks, a.GetSession)
	e.GET("/api/v1/tasks/user", taskHandler.getUserTasks, a.GetSession)
	e.DELETE("/api/v1/tasks/:id", taskHandler.deleteTask, a.GetSession)
	e.PUT("/api/v1/tasks/:id", taskHandler.updateTask, a.GetSession)
}

func (th *TaskHandler) getTask(c echo.Context) error {
	defer c.Request().Body.Close()
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)

	id := c.Param(constants.IdKey)
	n, _ := strconv.ParseUint(string(id), 10, 64)

	cookie, err := c.Cookie(constants.SessionCookieName)
	if err != nil && cookie != nil {
		log.Println("user handler: createTask: error getting cookie")
		return echo.NewHTTPError(http.StatusBadRequest, "error getting cookie")
	}

	var uid uint64

	if cookie == nil {
		uid = 0
	} else {
		uid, err = th.uuc.CheckSession(cookie.Value)
		if err != nil {
			return err
		}
	}

	t, err := th.uc.GetTask(n, uid, false)
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
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)

	page := c.QueryParams().Get(constants.PageKey)
	p, _ := strconv.Atoi(string(page))
	if p == 0 {
		p = 1
	}

	count := c.QueryParams().Get(constants.CountKey)
	cc, _ := strconv.Atoi(string(count))
	if cc == 0 {
		cc = constants.TasksPerPage
	}

	uid := uint64(0)

	cookie, err := c.Cookie(constants.SessionCookieName)
	if err != nil && cookie != nil {
		log.Println("user handler: getTasks: error getting cookie")
		return echo.NewHTTPError(http.StatusBadRequest, "error getting cookie")
	}

	if cookie == nil {
		uid = 0
	} else {
		uid, err = th.uuc.CheckSession(cookie.Value)
		if err != nil {
			return err
		}
	}

	tsks, err := th.uc.GetTasks(uid, p, cc)
	if err != nil {
		return err
	}

	if _, err = easyjson.MarshalToWriter(tsks, c.Response().Writer); err != nil {
		log.Println("task handler: getTasks: error marshaling answer to writer", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return nil
}

func (th *TaskHandler) findTasks(c echo.Context) error {
	defer c.Request().Body.Close()
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)

	str := c.QueryParam("find")
	if str == "" {
		return c.JSON(200, models.ShortTasks{})
	}
	page := c.QueryParams().Get(constants.PageKey)
	p, _ := strconv.Atoi(string(page))
	if p == 0 {
		p = 1
	}

	count := c.QueryParams().Get(constants.CountKey)
	cc, _ := strconv.Atoi(string(count))
	if cc == 0 {
		cc = constants.TasksPerPage
	}

	uid := uint64(0)

	cookie, err := c.Cookie(constants.SessionCookieName)
	if err != nil && cookie != nil {
		log.Println("user handler: getTasks: error getting cookie")
		return echo.NewHTTPError(http.StatusBadRequest, "error getting cookie")
	}

	if cookie == nil {
		uid = 0
	} else {
		uid, err = th.uuc.CheckSession(cookie.Value)
		if err != nil {
			return err
		}
	}

	tsks, err := th.uc.FindTasks(str, uid, p, cc)
	if err != nil {
		return err
	}

	if _, err = easyjson.MarshalToWriter(tsks, c.Response().Writer); err != nil {
		log.Println("task handler: findTasks: error marshaling answer to writer", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return nil
}

func (th *TaskHandler) getSolvedTasks(c echo.Context) error {
	defer c.Request().Body.Close()
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)

	page := c.QueryParams().Get(constants.PageKey)
	p, _ := strconv.Atoi(string(page))
	if p == 0 {
		p = 1
	}

	count := c.QueryParams().Get(constants.CountKey)
	cc, _ := strconv.Atoi(string(count))
	if cc == 0 {
		cc = constants.TasksPerPage
	}

	uid := c.Get(constants.UserIdKey).(uint64)

	tsks, err := th.uc.GetSolvedTasks(uid, p, cc)
	if err != nil {
		return err
	}

	if _, err = easyjson.MarshalToWriter(tsks, c.Response().Writer); err != nil {
		log.Println("task handler: getSolvedTasks: error marshaling answer to writer", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return nil
}

func (th *TaskHandler) getUnsolvedTasks(c echo.Context) error {
	defer c.Request().Body.Close()
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)

	page := c.QueryParams().Get(constants.PageKey)
	p, _ := strconv.Atoi(string(page))
	if p == 0 {
		p = 1
	}

	count := c.QueryParams().Get(constants.CountKey)
	cc, _ := strconv.Atoi(string(count))
	if cc == 0 {
		cc = constants.TasksPerPage
	}

	uid := c.Get(constants.UserIdKey).(uint64)

	tsks, err := th.uc.GetUnsolvedTasks(uid, p, cc)
	if err != nil {
		return err
	}

	if _, err = easyjson.MarshalToWriter(tsks, c.Response().Writer); err != nil {
		log.Println("task handler: getSolvedTasks: error marshaling answer to writer", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return nil
}

func (th *TaskHandler) getUserTasks(c echo.Context) error {
	defer c.Request().Body.Close()
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)

	page := c.QueryParams().Get(constants.PageKey)
	p, _ := strconv.Atoi(string(page))
	if p == 0 {
		p = 1
	}

	count := c.QueryParams().Get(constants.CountKey)
	cc, _ := strconv.Atoi(string(count))
	if cc == 0 {
		cc = constants.TasksPerPage
	}

	uid := c.Get(constants.UserIdKey).(uint64)

	tsks, err := th.uc.GetUserTasks(uid, p, cc)
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
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)

	uid := c.Get(constants.UserIdKey).(uint64)

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
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)

	uid := c.Get(constants.UserIdKey).(uint64)

	id := c.Param(constants.IdKey)
	iid, _ := strconv.ParseUint(string(id), 10, 64)

	err := th.uc.DeleteTask(iid, uid)
	if err != nil {
		return err
	}

	return nil
}

func (th *TaskHandler) updateTask(c echo.Context) error {
	defer c.Request().Body.Close()
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)

	uid := c.Get(constants.UserIdKey).(uint64)

	tn := &models.TaskNew{}

	if err := easyjson.UnmarshalFromReader(c.Request().Body, tn); err != nil {
		log.Println("task handler: updateTask: error unmarshaling task from reader", err)
		return echo.NewHTTPError(http.StatusTeapot, err.Error())
	}

	tn.Creator = uid
	id := c.Param(constants.IdKey)
	iid, _ := strconv.ParseUint(string(id), 10, 64)

	err := th.uc.UpdateTask(iid, tn)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return nil
}
