package main

import (
	"context"
	"log"
	"time"
)

func gameCleaner(repository IRepository, context context.Context, ticker <-chan time.Time) {
	for {
		select {
		case <-context.Done():
			return
		case <-ticker:
			for _, game := range repository.GameSessions() {
				if game.isOver {
					repository.RemoveGame(game)
				}
			}
		}
	}
}

func gameSessionsCreator(repository IRepository, context context.Context, ticker <-chan time.Time) {
	var playerFirst *User
	var playerSecond *User
	for {
		select {
		case <-context.Done():
			return
		case <-ticker:
			playerFirst = nil
			playerSecond = nil
			for _, user := range repository.UsersInSearchInsertionOrder() {
				if playerFirst == nil {
					playerFirst = user
				} else {
					playerSecond = user
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
					repository.AddGame(game)

					repository.RemoveUserInSearch(playerFirst)
					repository.RemoveUserInSearch(playerSecond)

					playerFirst = nil
					playerSecond = nil
				}
			}
		}
	}
}
