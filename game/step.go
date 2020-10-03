package game

type Step struct {
	PlayerId int `json:"player_id,omitempty"`
	// Card position is a card's position on a player's hand
	Position int `json:"position,omitempty"`
}
