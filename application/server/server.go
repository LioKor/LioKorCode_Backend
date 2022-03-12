package server

import (
	"context"
	"log"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/labstack/echo"

	slhttp "liokoredu/application/solution/delivery/http"
	slrep "liokoredu/application/solution/repository"
	sluc "liokoredu/application/solution/usecase"
	thttp "liokoredu/application/task/delivery/http"
	trep "liokoredu/application/task/repository"
	tuc "liokoredu/application/task/usecase"
	"liokoredu/pkg/constants"
)

type Server struct {
	e *echo.Echo
}

func NewServer() *Server {
	var server Server

	e := echo.New()

	pool, err := pgxpool.Connect(context.Background(),
		"user=lk"+
			" password=liokor"+constants.DBConnect)
	if err != nil {
		log.Fatal(err)
	}
	err = pool.Ping(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	solutionRep := slrep.NewSolutionDatabase(pool)
	taskRep := trep.NewTaskDatabase(pool)

	solutionUC := sluc.NewSolutionUseCase(solutionRep)
	taskUC := tuc.NewTaskUseCase(taskRep)

	slhttp.CreateSolutionHandler(e, solutionUC)
	thttp.CreateTaskHandler(e, taskUC)

	server.e = e
	return &server
}

func (s Server) ListenAndServe() {
	s.e.Logger.Fatal(s.e.Start(":1323"))
}
