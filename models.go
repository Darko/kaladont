package main

const objectPrefix = "kaladont:"

// Room struct
type Room struct {
	ID                string   `json:"id"`
	Players           []Player `json:"players,omitempty"`
	CurrentPlayer     Player   `json:"currentPlayer,omitempty"`
	NextWordStartWith string   `json:"nextWordStartWith,omitempty"`
	Words             []Word   `json:"words,omitempty"`
}

// Word struct
type Word struct {
	Word   string `json:"word"`
	Player string `json:"player"`
}

// Player struct
type Player struct {
	Name  string `json:"name"`
	Score int32  `json:"score"`
}
