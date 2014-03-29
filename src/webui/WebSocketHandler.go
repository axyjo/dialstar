package webui

import (
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
	"time"
)

type Data struct {
	Field       string
	UserCount   int
	Conferences []Conference
}

type Conference struct {
	StartTime time.Time
	Caller1   string
	Caller2   string
}

type WebSocketWrapper struct {
	Push chan bool
}

func (c WebSocketWrapper) WebSocketHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Incoming web socket request")

	conn, err := websocket.Upgrade(w, r, nil, 1024, 1024)

	if _, ok := err.(websocket.HandshakeError); ok {
		http.Error(w, "Not a websocket handshake", 400)
		return
	} else if err != nil {
		panic(err)
	}

	data := getData()

	err = conn.WriteJSON(data)

	if err != nil {

	}
}

func getUserCount() int {
	return 15
}

func getConferences() []Conference {
	conference1 := Conference{StartTime: time.Now(),
		Caller1: "Foo",
		Caller2: "Bar"}
	conference2 := Conference{StartTime: time.Now(),
		Caller1: "Gitesh",
		Caller2: "Dhir"}
	return []Conference{conference1, conference2}
}

func getData() Data {
	return Data{
		UserCount:   getUserCount(),
		Conferences: getConferences(),
	}
}
