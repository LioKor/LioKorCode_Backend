package server

import (
	"context"
	"liokoredu/application/microservices/redactor/proto"
)

type RedactorServer struct {
}

func NewRedactorServer() *RedactorServer {
	return &RedactorServer{}
}

func (rs *RedactorServer) CreateConnection(c context.Context, uid *proto.IdValue) (*proto.CreateAnswer, error) {
	str := ""
	return &proto.CreateAnswer{ConnectionId: str, Flag: false, Msg: ""}, nil
}

func (rs *RedactorServer) Connect(c context.Context, id *proto.IdValue) (*proto.ConnectAnswer, error) {
	{
		/*
			flag, userId, err := a.usecase.Check(s.Value)
			if err != nil {
				return nil, status.Error(codes.InvalidArgument, err.Error())
			}
		*/

		return &proto.ConnectAnswer{}, nil
	}
}
