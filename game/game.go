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

func (game *Game) PushCard(playerId int, cardPosition int) error {
	step := Step{
		PlayerId: playerId,
		Position: cardPosition,
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

func (game *Game) Print() {
	fmt.Printf("Очередь игрока %d\n", game.turn.ID)
	fmt.Println("Содержимое колоды:")
	fmt.Println(game.deck.Print())
	fmt.Println("Содержимое поля:")
	fmt.Println(game.field.Print())
	for i := 0; i < len(game.Players); i++ {
		fmt.Printf("Рука игрока %d\n", i+1)
		fmt.Println(game.Players[i].Cards.Print())
	}
}
func (game *Game) Listen() {
	game.deck.Shuffle()
	game.turn = game.ChooseRandomPlayer()
	game.CardsDeal()
	game.Print()
	game.UpdateCards()
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
				player := game.turn
				game.mtx.Lock()
				err := player.Cards.MoveCard(step.Position, game.field)
				game.mtx.Unlock()
				if err != nil {
					msg.Type = "moving card error"
					msg.Payload = nil
					player.MessageCh <- msg
					continue
				}
				log.Printf("Card %d moved\n", step.Position)

				game.NextStep()
				game.UpdateCards()
				game.Print()
			}
		case <-game.FinishCh:
			finishMsg := Message{
				Type: "finish",
			}
			game.SendBroadcast(func(player *Player) *Message {
				return &finishMsg
			})
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
	// TODO: Количество карт должно быть вариативным (не всегда строго 6)
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

func (game *Game) SendToPlayer(msg *Message, player *Player) {
	player.MessageCh <- msg
}

func (game *Game) SendBroadcast(transform func(player *Player) *Message) {
	for _, player := range game.Players {
		game.SendToPlayer(transform(player), player)
	}
}

func (game *Game) UpdateCards() {
	game.SendBroadcast(func(player *Player) *Message {
		type Payload struct {
			Cards []*Card `json:"cards"`
		}
		payload := Payload{
			Cards: player.Cards.Cards,
		}
		payloadData, _ := json.Marshal(payload)
		msg := &Message{
			Type:    "update_cards",
			Payload: payloadData,
		}
		return msg
	})
}
