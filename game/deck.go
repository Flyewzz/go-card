package game

import (
	"errors"
	"math/rand"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Deck struct {
	mtx   *sync.Mutex
	Cards []*Card
}

func NewDeck(cards []*Card) *Deck {
	return &Deck{
		Cards: cards,
		mtx:   &sync.Mutex{},
	}
}

func (this *Deck) Push(card *Card) {
	this.Cards = append([]*Card{card}, this.Cards...)
}

func (this *Deck) GetCard(position int) (*Card, error) {
	if len(this.Cards) == 0 {
		return nil, errors.New("Deck is empty")
	}
	if position < 0 || position >= len(this.Cards) {
		return nil, errors.New("Wrong position")
	}
	card := this.Cards[position]
	return card, nil
}

func (this *Deck) Pop(position int) (*Card, error) {
	card, err := this.GetCard(position)
	if err != nil {
		return nil, err
	}
	if len(this.Cards) == 1 {
		this.Cards = []*Card{}
	} else {
		this.Cards = append(this.Cards[:position], this.Cards[position+1:]...)
	}
	return card, nil
}

func (this *Deck) Shuffle() {
	r := rand.New(rand.NewSource(time.Now().Unix()))

	for n := len(this.Cards); n > 0; n-- {
		randIndex := r.Intn(n)
		this.Cards[n-1], this.Cards[randIndex] = this.Cards[randIndex], this.Cards[n-1]
	}
}

func (this *Deck) MoveCard(position int, deck *Deck) error {
	card, err := this.Pop(position)
	if err != nil {
		return err
	}
	deck.Push(card)
	return nil
}

func (this *Deck) Print() string {
	var ids []string
	for _, v := range this.Cards {
		ids = append(ids, strconv.Itoa(v.ID))
	}
	return strings.Join(ids, " ")
}
