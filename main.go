package main

import (
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

type IRepository interface {
	UserByUUID(uuid string) *User
	GameByUUID(uuid string) *Game
	GameSessions() map[string]*Game
	UsersInSearch() map[string]*User
	Users() map[string]*User
	AddUser(user *User)
	RemoveUser(user *User)
	AddUserInSearch(user *User)
	RemoveUserInSearch(user *User)
	AddGame(game *Game)
	RemoveGame(game *Game)
}

var wsUpgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var Repository IRepository

func init() {
	Repository = InmemoryRepository()
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

func gameCleaner() {
	ticker := time.Tick(30 * time.Second)
	for {
		<-ticker
		for _, game := range Repository.GameSessions() {
			if game.isOver {
				Repository.RemoveGame(game)
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
		playerFirst = nil
		playerSecond = nil
		for _, user := range Repository.UsersInSearch() {
			if playerFirst == nil {
				playerFirst = user

				// added for testing
				// to remove
				playerSecond = user
			} else if playerSecond == nil {
				playerSecond = user
			}
			{
				// for dev purpose
				// } else {
				if rand.Intn(1) == 0 {
					playerFirst, playerSecond = playerSecond, playerFirst
				}
				log.Println("Creating the game...")
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
				Repository.AddGame(game)

				Repository.RemoveUserInSearch(playerFirst)
				Repository.RemoveUserInSearch(playerSecond)

				playerFirst = nil
				playerSecond = nil
			}

		}
	}
}
