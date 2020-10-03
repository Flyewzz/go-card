package controllers

import (
	"net/http"

	. "github.com/Flyewzz/go-card/models"
	"github.com/gorilla/websocket"
)

type HandlerData struct {
	Rooms map[string]*Room
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func (hd *HandlerData) ConnectToRoom(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		HandleError(&w, r, 500, err.Error())
		return
	}

	user := NewUser(conn)
	go user.Listen()
	hd.Rooms["first"].AddUser(user)
}
