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
	/*
		sessionValue := usr.Value
		if len(sessionValue) != 0 {
			flag, _, err := a.usecase.Check(sessionValue)
			if err != nil {
				return &proto.LoginAnswer{Value: sessionValue, Flag: false}, err
			}
			if flag {
				return &proto.LoginAnswer{Value: sessionValue,
					Flag: true,
					Msg:  "user is already logged in"}, nil
			}
		}

		sessionValue, flag, err := a.usecase.Login(usr.Login, usr.Password)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		if flag {
			return &proto.LoginAnswer{Value: sessionValue,
				Flag: true,
				Msg:  "incorrect data"}, nil
		}
	*/

	return &proto.CreateAnswer{}, nil
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
