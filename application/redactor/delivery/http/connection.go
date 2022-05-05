package http

import (
	"encoding/json"
	"liokoredu/pkg/constants"
	"log"
	"time"

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

func (c *Connection) write(mt int, payload []byte) error {
	c.Ws.SetWriteDeadline(time.Now().Add(constants.WriteWait))
	return c.Ws.WriteMessage(mt, payload)
}

func (c *Connection) Handle() error {
	s := c.Session

	err := c.Send(&Event{"doc", map[string]interface{}{
		"document": s.Document,
		"revision": len(c.Session.Operations),
		"clients":  s.Clients,
	}})
	if err != nil {
		log.Println(err)
		return err
	}

	s.RegisterConnection(c)
	go c.pingPong()

	for {
		e, err := c.ReadEvent()
		if err != nil {
			break
		}
		s.EventChan <- ConnEvent{c, e}
	}

	c.Broadcast(&Event{"quit", map[string]interface{}{
		"client_id": c.ID,
		"username":  s.Clients[c.ID].Name,
	}})
	s.UnregisterConnection(c)
	return nil
}

func (c *Connection) pingPong() {
	ticker := time.NewTicker(constants.PingPeriod)
	defer func() {
		ticker.Stop()
	}()

	for {
		select {
		case <-ticker.C:
			if err := c.write(websocket.PingMessage, []byte{}); err != nil {
				c.Broadcast(&Event{"quit", map[string]interface{}{
					"client_id": c.ID,
					"username":  c.Session.Clients[c.ID].Name,
				}})
				c.Session.UnregisterConnection(c)
				return
			}
		}
	}

}

func (c *Connection) ReadEvent() (*Event, error) {
	c.Ws.SetPongHandler(func(string) error { c.Ws.SetReadDeadline(time.Now().Add(constants.PongWait)); return nil })
	_, msg, err := c.Ws.ReadMessage()
	if err != nil {
		return nil, err
	}
	m := &Event{}
	if err = json.Unmarshal(msg, &m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *Connection) Send(msg *Event) error {
	c.Ws.SetWriteDeadline(time.Now().Add(constants.WriteWait))
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
