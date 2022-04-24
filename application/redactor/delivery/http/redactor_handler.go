package http

import (
	"liokoredu/application/microservices/redactor/client"
	"liokoredu/application/models"
	"liokoredu/application/server/middleware"
	"liokoredu/pkg/constants"
	"liokoredu/pkg/generators"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo"
	"github.com/mailru/easyjson"
)

type RedactorHandler struct {
	rpcRedactor *client.RedactorClient
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// это стор, где хранятся сессии
var subscriptions = make(map[string]*Session)

func CreateRedactorHandler(e *echo.Echo, a middleware.Auth) {

	redactorHandler := RedactorHandler{}

	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	e.POST("/api/v1/redactor", redactorHandler.CreateConnection)
	e.GET("/api/v1/ws/redactor/:id", redactorHandler.ConnectToRoom)
	e.GET("/api/v1/redactor/:id/:filepath", redactorHandler.GetFileText)

}

func (rh *RedactorHandler) CreateConnection(c echo.Context) error {
	defer c.Request().Body.Close()
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)

	sln := &models.SolutionFiles{}

	if err := easyjson.UnmarshalFromReader(c.Request().Body, sln); err != nil {
		log.Println(err)
		return echo.NewHTTPError(http.StatusTeapot, err.Error())
	}
	log.Println(sln)

	roomId, _ := createRoom(*sln)
	//serveWs(c, session)

	/*
		id, err, code := rh.rpcRedactor.CreateConnection(uid)
		if err != nil {
			log.Println(id, err, code)
			return err
		}
	*/

	//if _, err := easyjson.MarshalToWriter(&models.IdValue{Id: roomId}, c.Response().Writer); err != nil {
	//	log.Println(err)
	//	return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	//}

	return c.JSON(http.StatusOK, &models.IdValue{Id: roomId})
}

func (rh *RedactorHandler) ConnectToRoom(c echo.Context) error {
	defer c.Request().Body.Close()
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)

	id := c.Param(constants.IdKey)
	log.Println(id)
	s := getRoom(id)
	if s == nil {
		return c.JSON(http.StatusNotFound, nil)
	}
	serveWs(c, s)

	//uid := c.Get(constants.UserIdKey).(uint64)
	/*

		id, err, code := rh.rpcRedactor.CreateConnection(uid)
		if err != nil {
			log.Println(id, err, code)
			return err
		}

		if _, err = easyjson.MarshalToWriter(&models.IdValue{Id: id}, c.Response().Writer); err != nil {
			log.Println(err)
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	*/

	return c.JSON(http.StatusOK, nil)
}

func (rh *RedactorHandler) GetFileText(c echo.Context) error {
	defer c.Request().Body.Close()
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)

	id := c.Param(constants.IdKey)
	log.Println(id)
	filename := c.Param("filepath")

	s := getRoom(id)
	if s == nil {
		return c.JSON(http.StatusNotFound, nil)
	}
	text := &models.RedactorFile{Text: s.GetDocument(filename)}
	if _, err := easyjson.MarshalToWriter(text, c.Response().Writer); err != nil {
		log.Println(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	//uid := c.Get(constants.UserIdKey).(uint64)
	/*

		id, err, code := rh.rpcRedactor.CreateConnection(uid)
		if err != nil {
			log.Println(id, err, code)
			return err
		}

		if _, err = easyjson.MarshalToWriter(&models.IdValue{Id: id}, c.Response().Writer); err != nil {
			log.Println(err)
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	*/

	return c.JSON(http.StatusOK, nil)
}

func createRoom(code models.SolutionFiles) (string, *Session) {
	roomId := generators.RandStringRunes(constants.WSLength)
	log.Println("created room:", roomId)
	session := NewSession(code)
	go session.HandleEvents()
	subscriptions[roomId] = session
	return roomId, session
}

func getRoom(id string) *Session {
	log.Println("got room:", id)
	return subscriptions[id]
}

func serveWs(c echo.Context, s *Session) {
	conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		log.Println(err)
		if _, ok := err.(websocket.HandshakeError); !ok {
			log.Println(err)
		}
		return
	}

	NewConnection(s, conn).Handle()
}
