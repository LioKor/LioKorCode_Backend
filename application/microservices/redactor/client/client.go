package client

import (
	"context"
	"liokoredu/application/microservices/redactor/proto"
	"liokoredu/pkg/constants"
	"log"
	"net/http"

	"github.com/labstack/echo"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type RedactorClient struct {
	client proto.RedactorClient
	gConn  *grpc.ClientConn
}

func NewRedactorClient(port string) (*RedactorClient, error) {
	gConn, err := grpc.Dial(
		constants.Localhost+port,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}

	return &RedactorClient{client: proto.NewRedactorClient(gConn), gConn: gConn}, nil
}

func (rc *RedactorClient) CreateConnection(uid uint64) (string, error, int) {
	idValue := proto.IdValue{Id: uid}
	answer, err := rc.client.CreateConnection(context.Background(), &idValue)
	if err != nil {
		return "", err, http.StatusInternalServerError
	}
	if answer.Flag {
		return "", echo.NewHTTPError(http.StatusBadRequest, answer.Msg), http.StatusBadRequest
	}

	return answer.ConnectionId, nil, http.StatusOK
}

func (rc *RedactorClient) Connect(value string) (bool, uint64, error, int) {
	/*
		sessionValue := &proto.Session{Value: value}

		answer, err := a.client.Check(context.Background(), sessionValue)
		if err != nil {
			return false, 0, err, http.StatusInternalServerError
		}
	*/

	return true, 0, nil, http.StatusOK
}

func (rc *RedactorClient) Close() {
	if err := rc.gConn.Close(); err != nil {
		log.Println(err)
	}
}
