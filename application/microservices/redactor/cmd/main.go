package main

import (
	"liokoredu/application/microservices/redactor/server"
	"liokoredu/pkg/constants"
)

func main() {
	s := server.NewServer(constants.RedactorServicePort)
	s.ListenAndServe()
}


