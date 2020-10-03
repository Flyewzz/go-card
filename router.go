package main

import (
	"github.com/Flyewzz/go-card/controllers"
	"github.com/Flyewzz/go-card/game"
	. "github.com/Flyewzz/go-card/models"
	"github.com/gorilla/mux"
)

func GetRouter() *mux.Router {

	router := mux.NewRouter()
	router.StrictSlash(true)

	hd := controllers.HandlerData{
		Rooms: make(map[string]*Room),
	}

	// TODO Needed to delete soon
	var deck = &game.Deck{}
	for i := 0; i < 15; i++ {
		deck.Push(&game.Card{
			ID: i,
		})
	}
	//*******************************

	game := game.NewGame(deck)
	room := NewRoom(2, game)

	go room.Listen()

	hd.Rooms["first"] = room

	router.HandleFunc("/", hd.ConnectToRoom)

	return router
}
