package interfaces

import (
	"github.com/Flyewzz/go-card/game"
)

type Game interface {
	PushCard(playerId int, cardPosition int) error
	Listen()
	GetPlayerById(id int) (*game.Player, error)
	RegisterPlayer(id int) error
	Finish()
	GetStatus() chan string
	GetState() string
}
