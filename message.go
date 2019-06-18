package main

const (
	Login        = "Login"
	LoginSuccess = "LoginSuccess"

	GameSearchOn    = "GameSearchOn"
	GameSearchOff   = "GameSearchOff"
	GameSearchStart = "GameSearchStart"

	GameOver   = "GameOver"
	GameMove   = "GameMove"
	GameMoved  = "GameMoved"
	GameWinner = "GameWinner"
)

type Message struct {
	Type    string            `json:"type"`
	Payload map[string]string `json:"payload"`
}
