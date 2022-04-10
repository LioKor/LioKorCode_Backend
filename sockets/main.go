package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/nitrous-io/ot.go/ot"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var defaultSession = NewSession(`package main
import "fmt"
func main() {
	fmt.Println("Hello, playground")
}`)

// это стор, где хранятся сессии
var subscriptions = make(map[string]*Session)

func main() {
	ot.TextEncoding = ot.TextEncodingTypeUTF16
	r := mux.NewRouter()
	r.HandleFunc("/ws/{room}", serveWs)
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("public")))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8090"
	}

	go defaultSession.HandleEvents()

	fmt.Printf("Listening on port %s\n", port)
	err := http.ListenAndServe(fmt.Sprintf(":%s", port), r)
	if err != nil {
		log.Fatal("Error: ", err)
	}
	/*

		go h.run()

		router := gin.New()

		router.GET("/ws/:roomId", func(c *gin.Context) {
			roomId := c.Param("roomId")
			serveWs(c.Writer, c.Request, roomId)
		})

		router.Run("0.0.0.0:8090")
	*/

}

func createSession(room string, code string) {
	subscriptions[room] = NewSession(code)
}

func serveWs(w http.ResponseWriter, r *http.Request) {
	var err error
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		if _, ok := err.(websocket.HandshakeError); !ok {
			log.Println(err)
		}
		return
	}
	room := mux.Vars(r)["room"]

	//c := &connection{send: make(chan []byte, 256), ws: conn}
	//s := subscription{c, roomId}
	//h.register <- s
	//go s.writePump()
	//go s.readPump()

	var session *Session
	if session = subscriptions[room]; session == nil {
		session = NewSession("code")
		subscriptions[room] = session
		log.Println("created")
	}
	log.Println(session.nextConnID)
	log.Println(len(session.Connections))

	NewConnection(session, conn).Handle()
}
