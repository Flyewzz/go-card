package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/Flyewzz/go-card/game"
	"github.com/Flyewzz/go-card/interfaces"
)

type Room struct {
	Users      map[int]*User
	timer      *time.Timer
	register   chan *User
	unregister chan *User
	mtx        *sync.Mutex
	start      chan struct{}
	capacity   int
	game       interfaces.Game
}

func (room *Room) UserMessageHandler(user *User) {
	for {
		select {
		case err := <-user.ErrorCh:
			if err != nil {
				room.RemoveUser(user)
				room.game.Finish()
				return
			}
		}
	}
}

func NewRoom(capacity int, game interfaces.Game) *Room {
	return &Room{
		Users:      make(map[int]*User),
		timer:      nil,
		register:   make(chan *User),
		unregister: make(chan *User),
		mtx:        &sync.Mutex{},
		start:      make(chan struct{}, 1),
		capacity:   capacity,
		game:       game,
	}
}

func (r *Room) AddUser(u *User) {
	r.register <- u
}

func (room *Room) RemoveUser(u *User) error {
	var err error
	if _, ok := room.Users[u.ID]; !ok {
		err = errors.New("No user found")
		log.Println(err)
		return err
	}
	go u.Disconnect()
	room.mtx.Lock()
	delete(room.Users, u.ID)
	room.mtx.Unlock()
	log.Printf("User %d disconnected\n", u.ID)
	return nil
}

func (room *Room) Listen() {
	for {
		select {
		case u := <-room.register:
			u.ID = len(room.Users) + 1
			room.mtx.Lock()
			room.Users[u.ID] = u
			err := room.game.RegisterPlayer(u.ID)
			if err != nil {
				log.Println(err)
			}
			room.mtx.Unlock()
			go room.UserMessageHandler(u)
			log.Printf("User %d connected\n", u.ID)
			if len(room.Users) == room.capacity {
				room.start <- struct{}{}
			}
		case <-room.start:
			fmt.Println("Game started")
			go room.game.Listen()
			// Wait for game's ready to accept messages
			<-room.game.GetStatus()
			for _, u := range room.Users {
				player, _ := room.game.GetPlayerById(u.ID)
				go func(user *User, player *game.Player) {
					for msg := range player.MessageCh {
						log.Println("Ошибка игрока: ", msg.Type)
						err := u.Send(msg)
						log.Println("Ошибка отправки сокета: ", err.Error())
					}
					log.Println("Корректное завершение горутины чтения")
				}(u, player)
				go func(user *User, player *game.Player) {
					for msg := range user.MessageCh {
						var step game.Step
						err := json.Unmarshal(msg.Payload, &step)
						if err != nil {
							log.Println(err)
							continue
						}
						step.PlayerId = player.ID
						room.game.PushCard(step.PlayerId, step.CardId)
					}
				}(u, player)
			}

		}
	}
}
