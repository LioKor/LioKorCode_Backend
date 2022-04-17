package main

import (
	"encoding/json"
	"log"

	"github.com/gorilla/websocket"
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

	err := c.Send(&Event{"doc", map[string]interface{}{
		"document": s.Document,
		"revision": len(c.Session.Operations),
		"clients":  s.Clients,
	}})
	if err != nil {
		log.Println(err)
		return err
	}
	log.Println("sent")

	s.RegisterConnection(c)
	log.Println("registered connection")

	for {
		log.Println(1)
		e, err := c.ReadEvent()
		if err != nil {
			break
		}
		log.Println("read event")

		s.EventChan <- ConnEvent{c, e}
	}

	log.Println(2)
	s.UnregisterConnection(c)
	c.Broadcast(&Event{"quit", c.ID})

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
