package server

import (
	"log"
	"net"

	"liokoredu/application/microservices/redactor/proto"

	"google.golang.org/grpc"
)

type Server struct {
	port     string
	redactor *RedactorServer
}

func NewServer(port string) *Server {
	/*
		pool, err := pgxpool.Connect(context.Background(), constants.DBConnect)
		if err != nil {
			logger.Fatal(err)
		}
		err = pool.Ping(context.Background())
		if err != nil {
			logger.Fatal(err)
		}


		conn, err := tarantool.Connect(constants.TarantoolAddress, tarantool.Opts{
			User: constants.TarantoolUser,
			Pass: constants.TarantoolPassword,
		})
		if err != nil {
			logger.Fatal(err)
		}

		_, err = conn.Ping()
		if err != nil {
			logger.Fatal(err)
		}
	*/

	return &Server{
		port:     port,
		redactor: NewRedactorServer(),
	}
}

func (s *Server) ListenAndServe() error {
	gServer := grpc.NewServer()

	listener, err := net.Listen("tcp", s.port)
	defer listener.Close()
	if err != nil {
		log.Println(err)
		return err
	}
	proto.RegisterRedactorServer(gServer, s.redactor)
	log.Println("starting server at " + s.port)
	err = gServer.Serve(listener)

	if err != nil {
		return nil
	}

	return nil
}
