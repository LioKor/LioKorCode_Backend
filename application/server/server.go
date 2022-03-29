package server

import (
	"context"
	"log"

	"github.com/gomodule/redigo/redis"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/labstack/echo"

	"liokoredu/application/server/middleware"
	slhttp "liokoredu/application/solution/delivery/http"
	slrep "liokoredu/application/solution/repository"
	sluc "liokoredu/application/solution/usecase"
	thttp "liokoredu/application/task/delivery/http"
	trep "liokoredu/application/task/repository"
	tuc "liokoredu/application/task/usecase"
	uhttp "liokoredu/application/user/delivery/http"
	urep "liokoredu/application/user/repository"
	uuc "liokoredu/application/user/usecase"
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

	redisPool := &redis.Pool{
		MaxIdle:   80,
		MaxActive: 12000,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", ":6379")
			if err != nil {
				panic(err.Error())
			}
			return c, err
		},
	}

	userRep := urep.NewUserDatabase(redisPool, pool)
	solutionRep := slrep.NewSolutionDatabase(pool)
	taskRep := trep.NewTaskDatabase(pool)

	userUC := uuc.NewUserUseCase(userRep)

	taskUC := tuc.NewTaskUseCase(taskRep)
	solutionUC := sluc.NewSolutionUseCase(solutionRep, taskUC)

	a := middleware.NewAuth(userUC)

	uhttp.CreateUserHandler(e, userUC, a)
	slhttp.CreateSolutionHandler(e, solutionUC, taskUC, userUC)
	thttp.CreateTaskHandler(e, taskUC, userUC)

	server.e = e
	return &server
}

func (s Server) ListenAndServe() {
	s.e.Logger.Fatal(s.e.Start(":1323"))
}
