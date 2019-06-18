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

func (u *User) resolveMessage(message Message) {
	switch message.Type {
	case Login:
		u.username = message.Payload["username"]
		Repository.AddUser(u)
		u.writeChan <- Message{
			LoginSuccess,
			map[string]string{
				"uuid":     u.uuid,
				"username": u.username,
			},
		}
		break

	case GameSearchOn:
		Repository.AddUserInSearch(u)
		break

	case GameOver:
		Repository.RemoveUserInSearch(u)
		if u.currentGameUUID != "" {
			game := Repository.GameByUUID(u.currentGameUUID)
			if game != nil {
				game.isOver = true
				message := Message{
					GameOver,
					map[string]string{},
				}
				for _, u := range game.users {
					if u.uuid != u.uuid {
						u.writeChan <- message
					}
				}
			}
		}
		break

	case GameMove:
		position, err := strconv.Atoi(message.Payload["position"])
		if err != nil {
			return
		}
		game := Repository.GameByUUID(u.currentGameUUID)
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
			log.Println("CROSS MOVER")
			game.field[position] = CROSS
			game.currentMoveUnit = ZERO
		} else if (game.currentMoveUnit == ZERO) && (game.zeroUser == u) {
			log.Println("ZERO MOVER")
			game.field[position] = ZERO
			game.currentMoveUnit = CROSS
		}
		message = Message{
			GameMoved,
			game.GetField(),
		}
		for _, user := range game.users {
			user.writeChan <- message
			break // dev
		}
		winner, ok := game.CheckWinner()
		if !ok {
			return
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
			break // dev
		}
		break
	}
}
