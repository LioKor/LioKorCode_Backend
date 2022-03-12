package main

import (
	"liokoredu/application/server"
)

func main() {
	s := server.NewServer()
	s.ListenAndServe()
}
