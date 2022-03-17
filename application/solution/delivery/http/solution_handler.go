package http

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"liokoredu/application/models"
	"liokoredu/application/solution"
	"liokoredu/application/task"
	"liokoredu/pkg/constants"
	"log"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
	"github.com/mailru/easyjson"
)

type SolutionHandler struct {
	UseCase  solution.UseCase
	TUseCase task.UseCase
}

func CreateSolutionHandler(e *echo.Echo,
	uc solution.UseCase, tuc task.UseCase) {
	solutionHandler := SolutionHandler{
		UseCase:  uc,
		TUseCase: tuc,
	}
	e.POST("/api/v1/tasks/:id/solutions", solutionHandler.PostSolution)
	e.POST("/api/v1/solutions/update/:id", solutionHandler.UpdateSolution)
	e.GET("/api/v1/tasks/:id/solutions", solutionHandler.GetSolutions)
	e.GET("/api/v1/tasks/:taskId/solutions/:solutionId", solutionHandler.GetSolutions)
}

func (sh SolutionHandler) PostSolution(c echo.Context) error {
	defer c.Request().Body.Close()

	sln := &models.Solution{}
	id := c.Param(constants.IdKey)
	uid, _ := strconv.ParseUint(string(id), 10, 64)
	log.Println(uid)

	if err := easyjson.UnmarshalFromReader(c.Request().Body, sln); err != nil {
		log.Println(err)
		return echo.NewHTTPError(http.StatusTeapot, err.Error())
	}

	task, err := sh.TUseCase.GetTask(uid)
	if err != nil {
		return err
	}
	log.Println(sln.SourceCode)
	testAmount := task.TestsAmount

	solId, err := sh.UseCase.InsertSolution(uid, sln.SourceCode, testAmount)

	ss := models.SolutionSend{
		Id:         solId,
		SourceCode: sln.SourceCode,
		Tests:      models.InputTests(task.Tests),
	}

	reqBody, err := json.Marshal(ss)
	resp, err := http.Post("http://167.172.51.136:7070/check_task/long",
		"application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		print(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		print(err)
	}

	log.Println(string(body))
	c.Response().Write(body)

	update := &models.SolutionUpdate{}

	_ = json.Unmarshal(body, update)
	err = sh.UseCase.UpdateSolution(solId, update.Code, update.Passed)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	//ans := &models.ReturnId{Id: solId}
	//if _, err = easyjson.MarshalToWriter(ans, c.Response().Writer); err != nil {
	//	log.Println(c, err)
	//	return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	//}

	return nil
}

func (sh SolutionHandler) UpdateSolution(c echo.Context) error {
	defer c.Request().Body.Close()

	id := c.Param(constants.IdKey)
	uid, _ := strconv.ParseUint(string(id), 10, 64)

	info := &models.SolutionUpdate{}
	if err := easyjson.UnmarshalFromReader(c.Request().Body, info); err != nil {
		log.Println(err)
		return echo.NewHTTPError(http.StatusTeapot, err.Error())
	}

	err := sh.UseCase.UpdateSolution(uid, info.Code, info.Passed)
	if err != nil {
		log.Println(err)
		return echo.NewHTTPError(http.StatusBadGateway, err.Error())
	}

	return nil
}

func (sh SolutionHandler) GetSolutions(c echo.Context) error {
	defer c.Request().Body.Close()

	id := c.Param(constants.IdKey)
	uid, _ := strconv.ParseUint(string(id), 10, 64)

	slns, err := sh.UseCase.GetSolutions(uid)
	if err != nil {
		log.Println(err)
		return echo.NewHTTPError(http.StatusBadGateway, err.Error())
	}

	if _, err = easyjson.MarshalToWriter(slns, c.Response().Writer); err != nil {
		log.Println(c, err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return nil
}
