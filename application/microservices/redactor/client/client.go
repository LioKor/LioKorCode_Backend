package client

import (
	"liokoredu/application/microservices/redactor/proto"
	"liokoredu/pkg/constants"
	"log"
	"net/http"

	"google.golang.org/grpc"
)

type RedactorClient struct {
	client proto.RedactorClient
	gConn  *grpc.ClientConn
}

func NewRedactorClient(port string) (*RedactorClient, error) {
	gConn, err := grpc.Dial(
		constants.Localhost+port,
		grpc.WithInsecure(),
	)
	if err != nil {
		return nil, err
	}

	return &RedactorClient{client: proto.NewRedactorClient(gConn), gConn: gConn}, nil
}

func (rc *RedactorClient) CreateConnection(login string, password string, value string) (uint64, string, error, int) {

	/*
		answer, err := a.client.Login(context.Background(), usr)
		if err != nil {
			return 0, "", err, http.StatusInternalServerError
		}
		if answer.Flag {
			return 0, "", echo.NewHTTPError(http.StatusBadRequest, answer.Msg), http.StatusBadRequest
		}
	*/

	return 0, " ", nil, http.StatusOK
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
