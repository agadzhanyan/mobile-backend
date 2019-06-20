package main

import (
	"github.com/gorilla/websocket"
	"log"
	"strconv"
)

type User struct {
	uuid            string
	username        string
	currentGameUUID string
	ws              *websocket.Conn
	writeChan       chan Message
	repository      IRepository
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
		u.resolveMessage(message)
	}
	err := u.ws.Close()
	if err != nil {
		log.Printf("error: %v", err)
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

func (u *User) close() {
	u.repository.RemoveUserInSearch(u)
	u.repository.RemoveUser(u)
}

func (u *User) resolveMessage(message Message) {
	switch message.Type {
	case Login:
		u.username = message.Payload["username"]
		u.repository.AddUser(u)
		u.writeChan <- Message{
			LoginSuccess,
			map[string]string{
				"uuid":     u.uuid,
				"username": u.username,
			},
		}
		break

	case GameSearchOn:
		u.repository.AddUserInSearch(u)
		break

	case GameSearchOff:
		u.repository.RemoveUserInSearch(u)
		break

	case GameOver:
		if u.currentGameUUID != "" {
			game := u.repository.GameByUUID(u.currentGameUUID)
			if game != nil {
				game.GameOver()
			}
		}
		break

	case GameMove:
		position, err := strconv.Atoi(message.Payload["position"])
		if err != nil {
			return
		}
		game := u.repository.GameByUUID(u.currentGameUUID)
		if game == nil {
			return
		}
		if game.isOver {
			return
		}
		val, ok := game.field[position]
		if !ok || val != EMPTY {
			return
		}
		if (game.currentMoveUnit == CROSS) && (game.crossUser == u) {
			log.Printf("CROSS MOVED: %v\n", position)
			game.field[position] = CROSS
			game.currentMoveUnit = ZERO
		} else if (game.currentMoveUnit == ZERO) && (game.zeroUser == u) {
			log.Printf("ZERO MOVED: %v\n", position)
			game.field[position] = ZERO
			game.currentMoveUnit = CROSS
		}
		message = Message{
			GameMoved,
			game.GetField(),
		}
		for _, user := range game.users {
			user.writeChan <- message
		}
		winner, ok := game.CheckWinner()
		if ok {
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
			return
		}
		if game.CheckDraw() {
			message = Message{
				GameDraw,
				map[string]string{},
			}
			for _, user := range game.users {
				user.currentGameUUID = ""
				user.writeChan <- message
			}
		}
		break
	}
}
