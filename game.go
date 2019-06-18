package main

import (
	"log"
	"strconv"
)

var winSlice = [][3]int{
	{1, 2, 3},
	{4, 5, 6},
	{7, 8, 9},
	{1, 4, 7},
	{2, 5, 8},
	{3, 6, 9},
	{1, 5, 9},
	{3, 5, 7},
}

type GameUnit string

const (
	EMPTY GameUnit = "EMPTY"
	CROSS GameUnit = "CROSS"
	ZERO  GameUnit = "ZERO"
)

type GameField map[int]GameUnit

type Game struct {
	uuid            string
	users           []*User
	crossUser       *User
	zeroUser        *User
	currentMoveUnit GameUnit
	isOver          bool
	field           GameField
}

func (g *Game) Start() {
	messageGameSearchOff := Message{
		GameSearchOff,
		map[string]string{},
	}
	messageGameStart := Message{
		GameSearchStart,
		map[string]string{
			"gameUUID":      g.uuid,
			"crossUserUUID": g.crossUser.uuid,
			"zeroUserUUID":  g.zeroUser.uuid,
		},
	}
	for _, user := range g.users {
		user.writeChan <- messageGameSearchOff
		user.writeChan <- messageGameStart
		break // dev
	}
	log.Println("Game started", g.uuid)
}

func (g *Game) GetField() map[string]string {
	field := map[string]string{}
	for k, v := range g.field {
		field[strconv.Itoa(k)] = string(v)
	}
	return field
}

func (g *Game) CheckDraw() bool {
	for _, value := range g.field {
		if value == EMPTY {
			return false
		}
	}
	return true
}

func (g *Game) CheckWinner() (GameUnit, bool) {
	for _, combination := range winSlice {
		if g.field[combination[0]] == EMPTY {
			continue
		}
		if g.field[combination[0]] == g.field[combination[1]] &&
			g.field[combination[1]] == g.field[combination[2]] {
			return g.field[combination[0]], true
		}
	}
	return EMPTY, false
}
