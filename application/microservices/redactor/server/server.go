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
