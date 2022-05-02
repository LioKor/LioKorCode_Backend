package http

import (
	"encoding/json"
	"log"

	"github.com/gorilla/websocket"
	"github.com/nitrous-io/ot.go/ot/session"
)

type Event struct {
	Name string      `json:"e"`
	Data interface{} `json:"d,omitempty"`
}

type Connection struct {
	ID      string
	Session *Session
	Ws      *websocket.Conn
}

type ConnEvent struct {
	Conn *Connection
	*Event
}

func NewConnection(session *Session, ws *websocket.Conn) *Connection {
	return &Connection{Session: session, Ws: ws}
}

func (c *Connection) Handle() error {
	s := c.Session
	log.Println("got session")

	var err error
	var filename string
	var filesession *session.Session

	for filename, filesession = range s.FileSessions {
		err = c.Send(&Event{"doc", map[string]interface{}{
			"document": filesession.Document,
			"revision": len(filesession.Operations),
			"clients":  filesession.Clients,
			"filename": filename,
		}})
		break
	}
	if err != nil {
		log.Println(err)
		return err
	}
	log.Println("sent")

	s.RegisterConnection(c, filename)
	log.Println("registered connection")

	for {
		e, err := c.ReadEvent()
		if err != nil {
			break
		}

		s.EventChan <- ConnEvent{c, e}
	}

	// we have to get current file, where client is
	for filename, filesession = range s.FileSessions {
		client := filesession.Clients[c.ID]
		if client != nil {
			break
		}
	}

	log.Println("quit")
	s.UnregisterConnection(c, filename)
	if s.FileSessions[filename].Clients[c.ID] == nil {
		return nil
	}
	c.Broadcast(&Event{"quit", map[string]interface{}{
		"client_id": c.ID,
		"username":  s.FileSessions[filename].Clients[c.ID].Name,
	}})

	return nil
}

func (c *Connection) ReadEvent() (*Event, error) {
	_, msg, err := c.Ws.ReadMessage()
	log.Println(err)
	if err != nil {
		return nil, err
	}
	log.Println("msg", string(msg[:]))
	m := &Event{}
	if err = json.Unmarshal(msg, &m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *Connection) Send(msg *Event) error {
	log.Println(msg)
	j, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	if err = c.Ws.WriteMessage(websocket.TextMessage, j); err != nil {
		return err
	}
	return nil
}

func (c *Connection) Broadcast(msg *Event) {
	for conn := range c.Session.Connections {
		if conn != c {
			conn.Send(msg)
		}
	}
}
