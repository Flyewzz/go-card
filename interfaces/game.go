package interfaces

import (
	. "github.com/Flyewzz/go-card/game"
)

type Game interface {
	PushCard(playerId int, cardID int) error
	Listen()
	GetPlayerById(id int) (*Player, error)
	RegisterPlayer(id int) error
	Finish()
	GetStatus() chan string
}
