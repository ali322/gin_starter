package ws

import (
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var Upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

var WebsocketManager = websocketManager{
	Clients:    make([]*Client, 0),
	Register:   make(chan *Client, 128),
	UnRegister: make(chan *Client, 128),
}

type websocketManager struct {
	Clients              []*Client
	Register, UnRegister chan *Client
	Locker               sync.Mutex
}

func (s *websocketManager) Start() {
	for {
		select {
		case client := <-s.Register:
			s.Locker.Lock()
			s.Clients = append(s.Clients, client)
			s.notifyClients(client, OnlineEvent)
			s.Locker.Unlock()
		case client := <-s.UnRegister:
			s.Locker.Lock()
			for i := 0; i < len(s.Clients); i++ {
				if s.Clients[i].Key == client.Key {
					s.Clients = append(s.Clients[0:i], s.Clients[i+1:]...)
				}
			}
			s.notifyClients(client, OfflineEvent)
			s.Locker.Unlock()
		}
	}
}

func (s *websocketManager) notifyClients(c *Client, event Event) {
	for i := 0; i < len(s.Clients); i++ {
		if s.Clients[i].Key != c.Key {
			// s.Clients[i].Conn.WriteJSON(keys)
			s.SendTo(&Message{Event: event, Data: c.Key}, s.Clients[i].Key)
		}
	}
}

func (s *websocketManager) RegisterConn(key string, unsafeConn *websocket.Conn) *Client {
	conn := &threadSafeConn{unsafeConn, sync.Mutex{}}
	c := &Client{
		Key: key, Conn: conn,
	}
	s.Register <- c
	// go c.ReadPump()
	return c
}

func (s *websocketManager) FindClient(key string) *Client {
	for _, c := range s.Clients {
		if c.Key == key {
			return c
		}
	}
	return nil
}

func (s *websocketManager) FindClientKeys() []string {
	keys := make([]string, 0)
	for i := 0; i < len(s.Clients); i++ {
		keys = append(keys, s.Clients[i].Key)
	}
	return keys
}

func (s *websocketManager) Contains(key string) bool {
	for _, c := range s.Clients {
		if c.Key == key {
			return true
		}
	}
	return false
}

func (s *websocketManager) Send(msg *Message) {
	for _, c := range s.Clients {
		c.Conn.WriteJSON(msg)
	}
}

func (s *websocketManager) SendTo(msg *Message, key string) {
	for _, c := range s.Clients {
		if c.Key == key {
			c.Conn.WriteJSON(msg)
		}
	}
}

func (s *websocketManager) UnRegisterConn(key string) {
	c := s.FindClient(key)
	s.UnRegister <- c
}
