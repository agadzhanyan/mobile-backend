package main

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
)

type User struct {
	uuid            string
	username        string
	currentGameUUID string
	ws              *websocket.Conn
	writeChan       chan Message
}

func (u *User) readLoop() {
	for {
		var message Message
		err := u.ws.ReadJSON(&message)
		log.Printf("message read: %v", message)
		if err != nil {
			log.Printf("error: %v", err)
			break
		}
		resolveMessage(u, message)
	}
	err := u.ws.Close()
	if err != nil {
		log.Fatal(err)
	}
}

func (u *User) writeLoop() {
	for {
		message := <-u.writeChan
		log.Printf("message sended: %v", message)
		err := u.ws.WriteJSON(message)
		if err != nil {
			log.Printf("error: %v", err)
			break
		}
	}
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
	}
	fmt.Println("Game started", g.uuid)
}

func (g *Game) GetField() map[string]string {
	field := map[string]string{}
	for k, v := range g.field {
		field[strconv.Itoa(k)] = string(v)
	}
	return field
}

func (g *Game) CheckWinner() (GameUnit, bool) {
	md5 := md5.New()
	jsonData, _ := json.Marshal(g.field)
	hash := fmt.Sprintf("%x", md5.Sum(jsonData))
	winner, ok := winHashMap[hash]
	if !ok {
		return EMPTY, false
	}
	return winner, true
}

var users = make(map[string]*User)
var usersInSearch = make(map[string]*User)
var gameSessions = make(map[string]*Game)

var wsUpgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Message struct {
	Type    string            `json:"type"`
	Payload map[string]string `json:"payload"`
}

func main() {
	go gameSessionCreator()
	go gameCleaner()
	http.HandleFunc("/", handleWebsocketConnections)
	log.Println("http server started on :17666")
	err := http.ListenAndServe(":17666", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func handleWebsocketConnections(w http.ResponseWriter, r *http.Request) {
	log.Println("New connection")
	ws, err := wsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	user := &User{
		generateUUID(),
		"<empty>",
		"",
		ws,
		make(chan Message),
	}
	go user.readLoop()
	go user.writeLoop()
}

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

func resolveMessage(user *User, message Message) {
	switch message.Type {
	case Login:
		user.username = message.Payload["username"]
		_, ok := users[user.uuid]
		if !ok {
			users[user.uuid] = user
		}
		user.writeChan <- Message{
			LoginSuccess,
			map[string]string{
				"uuid":     user.uuid,
				"username": user.username,
			},
		}
		break

	case GameSearchOn:
		usersInSearch[user.uuid] = user
		break

	case GameOver:
		_, ok := usersInSearch[user.uuid]
		if ok {
			delete(usersInSearch, user.uuid)
		}
		if user.currentGameUUID != "" {
			game, ok := gameSessions[user.currentGameUUID]
			if ok {
				message := Message{
					GameOver,
					map[string]string{},
				}
				for _, u := range game.users {
					u.writeChan <- message
				}
				game.isOver = true
			}
		}
		break

	case GameMove:
		position, err := strconv.Atoi(message.Payload["position"])
		if err != nil {
			return
		}
		game, ok := gameSessions[user.currentGameUUID]
		if !ok {
			return
		}
		if game.isOver {
			return
		}
		val, ok := game.field[position]
		if !ok || val != EMPTY {
			return
		}
		if (game.currentMoveUnit == CROSS) && (game.crossUser == user) {
			game.field[position] = CROSS
		} else if (game.currentMoveUnit == ZERO) && (game.zeroUser == user) {
			game.field[position] = ZERO
		}
		message = Message{
			GameMoved,
			game.GetField(),
		}
		winner, ok := game.CheckWinner()
		if ok {
			game.isOver = true
		}
		message = Message{
			GameWinner,
			map[string]string{
				"winner": string(winner),
			},
		}
		for _, user := range game.users {
			user.currentGameUUID = ""
			user.writeChan <- message
		}
		break
	}
}

func gameCleaner() {
	ticker := time.Tick(30 * time.Second)
	for {
		<-ticker
		for _, game := range gameSessions {
			if game.isOver {
				delete(gameSessions, game.uuid)
			}
		}
	}
}

func gameSessionCreator() {
	ticker := time.Tick(200 * time.Millisecond)
	var playerFirst *User
	var playerSecond *User
	for {
		<-ticker
		for _, user := range usersInSearch {
			if playerFirst == nil {
				playerFirst = user

				// added for testing
				// to remove
				playerSecond = user
			} else if playerSecond == nil {
				playerSecond = user
			} else {
				if rand.Intn(1) == 0 {
					playerFirst, playerSecond = playerSecond, playerFirst
				}
				game := &Game{
					generateUUID(),
					[]*User{playerFirst, playerSecond},
					playerFirst,
					playerSecond,
					CROSS,
					false,
					GameField{
						1: EMPTY,
						2: EMPTY,
						3: EMPTY,
						4: EMPTY,
						5: EMPTY,
						6: EMPTY,
						7: EMPTY,
						8: EMPTY,
						9: EMPTY,
					},
				}
				go game.Start()
				playerFirst.currentGameUUID = game.uuid
				playerSecond.currentGameUUID = game.uuid
				gameSessions[game.uuid] = game

				delete(usersInSearch, playerFirst.uuid)
				delete(usersInSearch, playerSecond.uuid)

				playerFirst = nil
				playerSecond = nil
			}

		}
	}
}
