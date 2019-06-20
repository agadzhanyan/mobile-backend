package main

import (
	"context"
	"reflect"
	"testing"
	"time"
)

func Test_gameCleaner(t *testing.T) {
	mockUserFirst := MockUser()
	mockUserSecond := MockUser()

	mockGame := MockGame(mockUserFirst, mockUserSecond)
	mockGame.isOver = true

	repository := InmemoryRepository()
	repository.AddUser(mockUserFirst)
	repository.AddUser(mockUserSecond)
	repository.AddGame(mockGame)

	ctx, cancel := context.WithCancel(context.Background())
	go gameCleaner(repository, ctx, time.Tick(1*time.Nanosecond))

	time.Sleep(200 * time.Millisecond)
	cancel()

	if !reflect.DeepEqual(map[string]*Game{}, repository.GameSessions()) {
		t.Errorf("Error while running cleaner: game wasn't removed")
	}
}

func Test_gameSessionsCreator(t *testing.T) {
	mockUserFirst := MockUser()
	mockUserSecond := MockUser()

	repository := InmemoryRepository()
	repository.AddUser(mockUserFirst)
	repository.AddUser(mockUserSecond)
	repository.AddUserInSearch(mockUserFirst)
	repository.AddUserInSearch(mockUserSecond)

	ctx, cancel := context.WithCancel(context.Background())
	go gameSessionsCreator(repository, ctx, time.Tick(1*time.Nanosecond))

	time.Sleep(200 * time.Millisecond)
	cancel()

	if len(repository.GameSessions()) != 1 {
		t.Errorf("Game count should be 1, got %v", len(repository.GameSessions()))
	}
}
