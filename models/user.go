package models

import (
	"encoding/json"
	"log"

	"github.com/Flyewzz/go-card/game"
	"github.com/gorilla/websocket"
)

type User struct {
	ID        int `json:"id"`
	Conn      *websocket.Conn
	ErrorCh   chan error
	MessageCh chan *game.Message
	stop      chan struct{}
}

func NewUser(conn *websocket.Conn) *User {
	user := &User{
		Conn:    conn,
		ErrorCh: make(chan error),
		stop:    make(chan struct{}, 1),
	}
	go user.Listen()
	return user
}

func (u *User) Disconnect() error {
	u.stop <- struct{}{}
	return u.Conn.Close()
}

func (u *User) Send(msg *game.Message) error {
	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	err = u.Conn.WriteJSON(data)
	return err
}

func (u *User) Listen() {
	for {
		select {
		case <-u.stop:
			close(u.MessageCh)
			return
		default:
			_, data, err := u.Conn.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
					u.ErrorCh <- err
					return
				}
				u.ErrorCh <- nil
			}
			var msg game.Message
			err = json.Unmarshal(data, &msg)
			if err != nil {
				log.Println(err)
				continue
			}
			u.MessageCh <- &msg
		}
	}

}
