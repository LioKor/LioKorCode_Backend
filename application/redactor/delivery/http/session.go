package http

import (
	"log"
	"strconv"
	"sync"

	"github.com/nitrous-io/ot.go/ot/operation"
	"github.com/nitrous-io/ot.go/ot/selection"
	"github.com/nitrous-io/ot.go/ot/session"
)

type Session struct {
	nextConnID  int
	Connections map[*Connection]struct{}

	EventChan chan ConnEvent

	lock sync.Mutex

	FileSessions map[string]*session.Session
}

func NewSession(document map[string]interface{}) *Session {
	s := map[string]*session.Session{}
	for filename, text := range document {
		s[filename] = session.New(text.(string))
	}
	return &Session{
		Connections:  map[*Connection]struct{}{},
		EventChan:    make(chan ConnEvent),
		FileSessions: s,
	}
}

func (s *Session) RegisterConnection(c *Connection) {
	s.lock.Lock()
	id := strconv.Itoa(s.nextConnID)
	c.ID = id
	s.nextConnID++
	s.Connections[c] = struct{}{}
	for filename := range s.FileSessions {
		s.FileSessions[filename].AddClient(c.ID)
	}
	s.lock.Unlock()
}

func (s *Session) UnregisterConnection(c *Connection) string {
	s.lock.Lock()
	delete(s.Connections, c)
	var filename string
	for filename = range s.FileSessions {
		s.FileSessions[filename].RemoveClient(c.ID)
	}

	s.lock.Unlock()
	return filename
}

func (s *Session) HandleEvents() {
	// this method should run in a single go routine
	for {
		e, ok := <-s.EventChan
		if !ok {
			return
		}
		log.Println("aaa")

		c := e.Conn
		log.Println(e.Name)
		switch e.Name {
		case "join":
			data, ok := e.Data.(map[string]interface{})
			log.Println(data, ok)
			if !ok {
				break
			}
			username, ok := data["username"].(string)
			log.Println(username, ok)
			if !ok || username == "" {
				log.Println(username)
				break
			}

			for filename := range s.FileSessions {
				s.FileSessions[filename].SetName(c.ID, username)
			}

			log.Println(c.ID)

			err := c.Send(&Event{"all", "registered", c.ID})
			if err != nil {
				log.Println(username)
				break
			}
			c.Broadcast(&Event{"all", "join", map[string]interface{}{
				"client_id": c.ID,
				"username":  username,
			}})
		case "op":
			// data: [filename revision, ops, selection? ]
			data, ok := e.Data.([]interface{})
			if !ok {
				break
			}
			if len(data) < 3 {
				break
			}
			// filename
			filename, ok := data[0].(string)
			if !ok {
				break
			}
			// revision
			revf, ok := data[1].(float64)
			rev := int(revf)
			if !ok {
				break
			}
			// ops
			ops, ok := data[2].([]interface{})
			if !ok {
				break
			}
			top, err := operation.Unmarshal(ops)
			if err != nil {
				break
			}
			// selection (optional)
			if len(data) >= 4 {
				selm, ok := data[3].(map[string]interface{})
				if !ok {
					break
				}
				sel, err := selection.Unmarshal(selm)
				if err != nil {
					break
				}
				top.Meta = sel
			}

			top2, err := s.FileSessions[filename].AddOperation(rev, top)
			if err != nil {
				break
			}

			err = c.Send(&Event{filename, "ok", nil})
			if err != nil {
				break
			}

			if sel, ok := top2.Meta.(*selection.Selection); ok {
				s.FileSessions[filename].SetSelection(c.ID, sel)
				c.Broadcast(&Event{filename, "op", []interface{}{c.ID, top2.Marshal(), sel.Marshal()}})
			} else {
				c.Broadcast(&Event{filename, "op", []interface{}{c.ID, top2.Marshal()}})
			}
		case "sel":
			data, ok := e.Data.([]interface{})
			//data, ok := e.Data.(map[string]interface{})
			if !ok {
				break
			}
			filename, ok := data[0].(string)
			if !ok {
				break
			}
			selm, ok := data[1].(map[string]interface{})
			if !ok {
				break
			}
			sel, err := selection.Unmarshal(selm)
			if err != nil {
				break
			}
			s.FileSessions[filename].SetSelection(c.ID, sel)
			c.Broadcast(&Event{filename, "sel", []interface{}{c.ID, sel.Marshal()}})
		}
	}
}
