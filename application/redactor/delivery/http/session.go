package http

import (
	"liokoredu/application/models"
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

func NewSession(code models.Solution) *Session {
	s := map[string]*session.Session{}
	for key, value := range code.SourceCode {
		s[key] = session.New(value.(string))
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

func (s *Session) GetDocument(filename string) string {
	s.lock.Lock()
	text := s.FileSessions[filename].Document
	s.lock.Unlock()

	return text
}

func (s *Session) UnregisterConnection(c *Connection) {
	s.lock.Lock()
	delete(s.Connections, c)
	var filename string
	for filename = range s.FileSessions {
		s.FileSessions[filename].RemoveClient(c.ID)
	}

	s.lock.Unlock()

	for filename = range s.FileSessions {
		c.Broadcast(&Event{"quit", map[string]interface{}{
			"client_id": c.ID,
			"username":  s.FileSessions[filename].Clients[c.ID].Name,
		}})
	}
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

			err := c.Send(&Event{"registered", c.ID})
			if err != nil {
				log.Println(username)
				break
			}
			c.Broadcast(&Event{"join", map[string]interface{}{
				"client_id": c.ID,
				"username":  username,
			}})
		case "op":
			// filename: "text", data: [revision, ops, selection?]
			source, ok := e.Data.(map[string]interface{})
			if !ok {
				log.Println("error getting source from op")
				break
			}
			filename, ok := source["filename"].(string)
			if !ok {
				log.Println("error converting filename to string:", filename)
				break
			}
			data, ok := source["data"].([]interface{})
			if !ok {
				log.Println("error getting data from op")
				break
			}
			if len(data) < 2 {
				break
			}
			// revision
			revf, ok := data[0].(float64)
			rev := int(revf)
			if !ok {
				break
			}
			// ops
			ops, ok := data[1].([]interface{})
			if !ok {
				break
			}
			top, err := operation.Unmarshal(ops)
			if err != nil {
				break
			}
			// selection (optional)
			if len(data) >= 3 {
				selm, ok := data[2].(map[string]interface{})
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

			err = c.Send(&Event{"ok", nil})
			if err != nil {
				break
			}

			if sel, ok := top2.Meta.(*selection.Selection); ok {
				s.FileSessions[filename].SetSelection(c.ID, sel)
				c.Broadcast(&Event{"op", map[string]interface{}{"data": []interface{}{c.ID, top2.Marshal(), sel.Marshal()}, "filename": filename}})
			} else {
				c.Broadcast(&Event{"op", map[string]interface{}{"data": []interface{}{c.ID, top2.Marshal()}, "filename": filename}})
			}
		case "sel":
			source, ok := e.Data.(map[string]interface{})
			if !ok {
				log.Println("error getting source from op")
				break
			}
			filename, ok := source["filename"].(string)
			if !ok {
				log.Println("error converting filename to string:", filename)
				break
			}
			data, ok := source["data"].(map[string]interface{})
			if !ok {
				log.Println("error getting data from op")
				break
			}

			sel, err := selection.Unmarshal(data)
			if err != nil {
				break
			}
			s.FileSessions[filename].SetSelection(c.ID, sel)
			c.Broadcast(&Event{"sel", map[string]interface{}{"data": []interface{}{c.ID, sel.Marshal()}, "filename": filename}})
			//[]interface{}{c.ID, sel.Marshal()}})
		}
	}
}
