package webui

import (
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
	"utils"
)

type WebSocketWrapper struct {
	Push []chan PushData
}

type PushData struct {
	UserCount int
	Call1Id   string
	Call2Id   string
}

func (c WebSocketWrapper) WebSocketHandler(w http.ResponseWriter, r *http.Request) {
	push := make(chan PushData, 5)
	c.Push = append(c.Push, push)
	fmt.Println("Incoming web socket request")

	conn, err := websocket.Upgrade(w, r, nil, 1024, 1024)

	if _, ok := err.(websocket.HandshakeError); ok {
		http.Error(w, "Not a websocket handshake", 400)
		return
	} else if err != nil {
		panic(err)
	}

	data := PushData{UserCount: utils.GetUserCount()}

	err = conn.WriteJSON(data)

	if err != nil {

	}

	go func() {
		var err error = nil
		for err == nil {
			p := <-push

			err = conn.WriteJSON(p)

		}
		for i, p := range c.Push {
			if p == push {
				c.Push[i], c.Push = c.Push[len(c.Push)-1], c.Push[:len(c.Push)-1]
				close(p)
				break
			}
		}
	}()
}
