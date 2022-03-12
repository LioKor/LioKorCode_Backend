package http

import (
	"liokoredu/application/solution"

	"github.com/labstack/echo"
)

type SolutionHandler struct {
	UseCase solution.UseCase
}

func CreateSolutionHandler(e *echo.Echo,
	uc solution.UseCase) {
	//solutionHandler := SolutionHandler{
	//	UseCase: uc,
	//}
}

func (sh SolutionHandler) PostSolution() {

}
