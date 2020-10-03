package game

type Step struct {
	PlayerId int `json:"player_id,omitempty"`
	// CardId is a card's position on a player's hand
	CardId int `json:"card_id,omitempty"`
}
