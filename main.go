package main

import (
	"context"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"time"
)

type IRepository interface {
	UserByUUID(uuid string) *User
	GameByUUID(uuid string) *Game
	GameSessions() map[string]*Game
	UsersInSearch() map[string]*User
	UsersInSearchInsertionOrder() []*User
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
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	tickerGSC := time.Tick(15 * time.Second)
	tickerCleaner := time.Tick(15 * time.Second)

	go gameSessionsCreator(Repository, ctx, tickerGSC)
	go gameCleaner(Repository, ctx, tickerCleaner)

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
		make(chan Message, 2),
		Repository,
	}
	go user.readLoop()
	go user.writeLoop()
}
