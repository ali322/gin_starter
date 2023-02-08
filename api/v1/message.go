package v1

import (
	"app/lib/ws"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func messager(c *gin.Context) {
	unsafeConn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		_ = c.Error(err)
		return
	}
	token := c.Query("token")
	client := ws.WebsocketManager.RegisterConn(token, unsafeConn)
	go client.ReadPump(func(raw []byte) error {
		message := ws.Message{}
		if err := json.Unmarshal(raw, &message); err != nil {
			return err
		}
		switch message.Event {
		case ws.StatusEvent:
			data, ok := message.Data.([]interface{})
			if !ok {
				client.Conn.WriteJSON(&ws.Message{Event: ws.StatusFailEvent, Data: "data 类型不正确"})
				return errors.New("unsupported data type")
			}
			result := make(map[string]bool)
			for _, v := range data {
				k := fmt.Sprintf("%v", v)
				if ws.WebsocketManager.Contains(k) {
					result[k] = true
				} else {
					result[k] = false
				}
			}
			client.Conn.WriteJSON(&ws.Message{Event: ws.StatusResultEvent, Data: result})
		}
		return nil
	})
	c.Status(http.StatusOK)
}
