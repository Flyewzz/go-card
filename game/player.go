package game

type Player struct {
	ID        int
	Cards     *Deck
	MessageCh chan *Message
}
