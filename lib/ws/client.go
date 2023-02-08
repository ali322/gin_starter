package ws

import (
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

type Event string

const (
	OnlineEvent       Event = "online"
	OfflineEvent      Event = "offline"
	StatusEvent       Event = "status"
	StatusFailEvent   Event = "statusFail"
	StatusResultEvent Event = "statusResult"
)

type Message struct {
	Event Event       `json:"event"`
	Data  interface{} `json:"data"`
}

type Client struct {
	Key  string
	Conn *threadSafeConn
}

func (c *Client) ReadPump(cb func(message []byte) error) {
	defer func() {
		WebsocketManager.UnRegisterConn(c.Key)
		c.Conn.Close()
	}()
	for {
		_, raw, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("err: %v", err)
			}
			break
		}
		if len(raw) == 0 {
			continue
		}
		if err := cb(raw); err != nil {
			continue
		}
	}
}

type threadSafeConn struct {
	*websocket.Conn
	sync.Mutex
}

func (t *threadSafeConn) WriteJSON(v interface{}) error {
	t.Lock()
	defer t.Unlock()
	return t.Conn.WriteJSON(v)
}
