package game

type Card struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Ability struct {
	Card
}

type Character struct {
	Card
	Power     int `json:"power"`
	Intellect int `json:"intellect"`
	Stealth   int `json:"stealth"`
	Humor     int `json:"humor"`
	Cost      int `json:"cost"`
}
