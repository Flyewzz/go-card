package game

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"
)

type Game struct {
	deck      *Deck
	field     *Deck
	State     string
	Players   []*Player
	turn      *Player
	MessageCh chan *Message
	FinishCh  chan struct{}
	random    *rand.Rand
	mtx       *sync.Mutex
	Ready     chan string
}

func (game *Game) RegisterPlayer(id int) error {
	if _, err := game.GetPlayerById(id); err != nil {
		player := &Player{
			ID:        id,
			Cards:     &Deck{},
			MessageCh: make(chan *Message, 1),
		}
		game.mtx.Lock()
		game.Players = append(game.Players, player)
		game.mtx.Unlock()
		return nil
	}
	return errors.New("Player with such id already exists")
}

func (game *Game) GetPlayerById(id int) (*Player, error) {
	for _, player := range game.Players {
		if player.ID == id {
			return player, nil
		}
	}
	return nil, errors.New("No player found")
}

func NewGame(deck *Deck) *Game {
	return &Game{
		deck:      deck,
		field:     NewDeck([]*Card{}),
		State:     "init",
		Players:   []*Player{},
		MessageCh: make(chan *Message),
		FinishCh:  make(chan struct{}, 1),
		random:    rand.New(rand.NewSource(time.Now().Unix())),
		mtx:       &sync.Mutex{},
		Ready:     make(chan string, 1),
	}
}

func (game *Game) PushCard(playerId int, cardID int) error {
	step := Step{
		PlayerId: playerId,
		CardId:   cardID,
	}
	stepData, err := json.Marshal(step)
	if err != nil {
		return err
	}
	msg := &Message{
		Type:    "push_card",
		Payload: stepData,
	}
	game.MessageCh <- msg
	return nil
}

func (game *Game) Listen() {
	game.deck.Shuffle()
	game.turn = game.ChooseRandomPlayer()
	fmt.Printf("Очередь игрока %d\n", game.turn.ID)
	game.CardsDeal()
	fmt.Println("Содержимое колоды:")
	fmt.Println(game.deck.Print())
	for i := 0; i < len(game.Players); i++ {
		fmt.Printf("Рука игрока %d\n", i+1)
		fmt.Println(game.Players[i].Cards.Print())
	}
	game.Init()
	for {
		select {
		case msg := <-game.MessageCh:
			switch msg.Type {
			case "push_card":
				var step Step
				json.Unmarshal(msg.Payload, &step)
				if step.PlayerId != game.turn.ID {
					msg.Type = "step turn error"
					msg.Payload = nil
					player, _ := game.GetPlayerById(step.PlayerId)
					player.MessageCh <- msg
					continue
				}
				player, _ := game.GetPlayerById(game.turn.ID)
				game.mtx.Lock()
				err := player.Cards.MoveCard(step.CardId, game.field)
				game.mtx.Unlock()
				if err != nil {
					msg.Type = "moving card error"
					msg.Payload = nil
					player.MessageCh <- msg
					continue
				}
				game.NextStep()
			}
		case <-game.FinishCh:
			finishMsg := Message{
				Type: "finish",
			}
			for _, player := range game.Players {
				player.MessageCh <- &finishMsg
			}
			game.Players = []*Player{}
			game.State = "finished"
			log.Println("Game finished")
			return
		}
	}
}

func (game *Game) Finish() {
	game.FinishCh <- struct{}{}
}

func (game *Game) ChooseRandomPlayer() *Player {
	playersCount := len(game.Players)
	if playersCount == 0 {
		return nil
	}
	return game.Players[game.random.Intn(playersCount)]
}
func (game *Game) GetTurn() *Player {
	return game.turn
}

func (game *Game) Init() {
	game.State = "in progress"
	game.Ready <- "step"
}

func (game *Game) giveCardsToPlayer(player *Player) {
	for i := 0; i < 6; i++ {
		game.mtx.Lock()
		game.deck.MoveCard(0, player.Cards)
		game.mtx.Unlock()
	}
}

func (game *Game) NextStep() {
	game.turn = game.getNextPlayer(game.turn.ID)
}

func (game *Game) getNextPlayer(id int) *Player {
	var playerSliceIndex int
	for index, player := range game.Players {
		if player.ID == id {
			playerSliceIndex = index
			break
		}
	}
	playerSliceIndex++
	if playerSliceIndex == len(game.Players) {
		playerSliceIndex = 0
	}
	return game.Players[playerSliceIndex]
}

func (game *Game) CardsDeal() {
	iterCount := len(game.Players)
	currentPlayer := game.GetTurn()
	for i := 0; i < iterCount; i++ {
		game.giveCardsToPlayer(currentPlayer)
		currentPlayer = game.getNextPlayer(currentPlayer.ID)
	}
}

func (game *Game) GetStatus() chan string {
	return game.Ready
}

func (game *Game) GetState() string {
	return game.State
}
