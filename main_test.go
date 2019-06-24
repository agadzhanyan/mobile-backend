package main

import (
	"context"
	"reflect"
	"testing"
	"time"
)

func MockUserWithRepository(repository IRepository) *User {
	return &User{
		generateUUID(),
		generateUUID(),
		"",
		nil,
		make(chan Message, 2),
		repository,
	}
}

func TestGame_Caht(t *testing.T) {
	repository := InmemoryRepository()

	mockUserFirst := MockUserWithRepository(repository)
	mockUserSecond := MockUserWithRepository(repository)

	mockUserFirst.resolveMessage(Message{Login, map[string]string{"username": "1"}})
	expectMsg := Message{
		LoginSuccess,
		map[string]string{
			"uuid":     mockUserFirst.uuid,
			"username": "1",
		},
	}
	gotMsg := <-mockUserFirst.writeChan
	if !reflect.DeepEqual(expectMsg, gotMsg) {
		t.Errorf("invalid write message\n expect: %v\n got %v", expectMsg, gotMsg)
	}

	mockUserSecond.resolveMessage(Message{Login, map[string]string{"username": "2"}})
	expectMsg = Message{
		LoginSuccess,
		map[string]string{
			"uuid":     mockUserSecond.uuid,
			"username": "2",
		},
	}
	gotMsg = <-mockUserSecond.writeChan
	if !reflect.DeepEqual(expectMsg, gotMsg) {
		t.Errorf("invalid write message\n expect: %v\n got %v", expectMsg, gotMsg)
	}

	mockUserFirst.resolveMessage(Message{MessageSend, map[string]string{"text": "Hi! It is gopher!"}})
	expectMsg = Message{
		MessageNew,
		map[string]string{
			"text": "Hi! It is gopher!",
		},
	}
	gotMsg = <-mockUserSecond.writeChan
	if !reflect.DeepEqual(expectMsg, gotMsg) {
		t.Errorf("invalid write message\n expect: %v\n got %v", expectMsg, gotMsg)
	}

	mockUserSecond.resolveMessage(Message{MessageSend, map[string]string{"text": "Hey! It is Elephant!"}})
	expectMsg = Message{
		MessageNew,
		map[string]string{
			"text": "Hey! It is Elephant!",
		},
	}
	gotMsg = <-mockUserFirst.writeChan
	if !reflect.DeepEqual(expectMsg, gotMsg) {
		t.Errorf("invalid write message\n expect: %v\n got %v", expectMsg, gotMsg)
	}
}

// when someone leave the game
func TestGame_GameOver(t *testing.T) {
	repository := InmemoryRepository()

	mockUserFirst := MockUserWithRepository(repository)
	mockUserSecond := MockUserWithRepository(repository)

	mockUserFirst.resolveMessage(Message{Login, map[string]string{"username": "1"}})
	expectMsg := Message{
		LoginSuccess,
		map[string]string{
			"uuid":     mockUserFirst.uuid,
			"username": "1",
		},
	}
	gotMsg := <-mockUserFirst.writeChan
	if !reflect.DeepEqual(expectMsg, gotMsg) {
		t.Errorf("invalid write message\n expect: %v\n got %v", expectMsg, gotMsg)
	}

	mockUserSecond.resolveMessage(Message{Login, map[string]string{"username": "2"}})
	expectMsg = Message{
		LoginSuccess,
		map[string]string{
			"uuid":     mockUserSecond.uuid,
			"username": "2",
		},
	}
	gotMsg = <-mockUserSecond.writeChan
	if !reflect.DeepEqual(expectMsg, gotMsg) {
		t.Errorf("invalid write message\n expect: %v\n got %v", expectMsg, gotMsg)
	}

	mockUserFirst.resolveMessage(Message{GameSearchOn, map[string]string{}})
	mockUserSecond.resolveMessage(Message{GameSearchOn, map[string]string{}})

	ctx, cancel := context.WithCancel(context.Background())
	go gameSessionsCreator(repository, ctx, time.Tick(1*time.Nanosecond))

	time.Sleep(200 * time.Millisecond)
	cancel()

	expectMsg = Message{
		GameSearchOff,
		map[string]string{},
	}
	gotMsg = <-mockUserFirst.writeChan
	if !reflect.DeepEqual(expectMsg, gotMsg) {
		t.Errorf("invalid write message\n expect: %v\n got %v", expectMsg, gotMsg)
	}
	gotMsg = <-mockUserSecond.writeChan
	if !reflect.DeepEqual(expectMsg, gotMsg) {
		t.Errorf("invalid write message\n expect: %v\n got %v", expectMsg, gotMsg)
	}

	expectMsg = Message{
		GameSearchStart,
		map[string]string{
			"gameUUID":      mockUserFirst.currentGameUUID,
			"crossUserUUID": mockUserFirst.uuid,
			"zeroUserUUID":  mockUserSecond.uuid,
		},
	}
	gotMsg = <-mockUserFirst.writeChan
	if !reflect.DeepEqual(expectMsg, gotMsg) {
		t.Errorf("invalid write message\n expect: %v\n got %v", expectMsg, gotMsg)
	}
	gotMsg = <-mockUserSecond.writeChan
	if !reflect.DeepEqual(expectMsg, gotMsg) {
		t.Errorf("invalid write message\n expect: %v\n got %v", expectMsg, gotMsg)
	}

	mockUserFirst.resolveMessage(Message{GameOver, map[string]string{}})
	expectMsg = Message{
		GameOver,
		map[string]string{},
	}
	gotMsg = <-mockUserSecond.writeChan
	if !reflect.DeepEqual(expectMsg, gotMsg) {
		t.Errorf("invalid write message\n expect: %v\n got %v", expectMsg, gotMsg)
	}
}

func TestGame_WithWinner(t *testing.T) {
	repository := InmemoryRepository()

	mockUserFirst := MockUserWithRepository(repository)
	mockUserSecond := MockUserWithRepository(repository)

	mockUserFirst.resolveMessage(Message{Login, map[string]string{"username": "1"}})
	expectMsg := Message{
		LoginSuccess,
		map[string]string{
			"uuid":     mockUserFirst.uuid,
			"username": "1",
		},
	}
	gotMsg := <-mockUserFirst.writeChan
	if !reflect.DeepEqual(expectMsg, gotMsg) {
		t.Errorf("invalid write message\n expect: %v\n got %v", expectMsg, gotMsg)
	}

	mockUserSecond.resolveMessage(Message{Login, map[string]string{"username": "2"}})
	expectMsg = Message{
		LoginSuccess,
		map[string]string{
			"uuid":     mockUserSecond.uuid,
			"username": "2",
		},
	}
	gotMsg = <-mockUserSecond.writeChan
	if !reflect.DeepEqual(expectMsg, gotMsg) {
		t.Errorf("invalid write message\n expect: %v\n got %v", expectMsg, gotMsg)
	}

	mockUserFirst.resolveMessage(Message{GameSearchOn, map[string]string{}})
	mockUserSecond.resolveMessage(Message{GameSearchOn, map[string]string{}})

	ctx, cancel := context.WithCancel(context.Background())
	go gameSessionsCreator(repository, ctx, time.Tick(1*time.Nanosecond))

	time.Sleep(200 * time.Millisecond)
	cancel()

	expectMsg = Message{
		GameSearchOff,
		map[string]string{},
	}
	gotMsg = <-mockUserFirst.writeChan
	if !reflect.DeepEqual(expectMsg, gotMsg) {
		t.Errorf("invalid write message\n expect: %v\n got %v", expectMsg, gotMsg)
	}
	gotMsg = <-mockUserSecond.writeChan
	if !reflect.DeepEqual(expectMsg, gotMsg) {
		t.Errorf("invalid write message\n expect: %v\n got %v", expectMsg, gotMsg)
	}

	expectMsg = Message{
		GameSearchStart,
		map[string]string{
			"gameUUID":      mockUserFirst.currentGameUUID,
			"crossUserUUID": mockUserFirst.uuid,
			"zeroUserUUID":  mockUserSecond.uuid,
		},
	}
	gotMsg = <-mockUserFirst.writeChan
	if !reflect.DeepEqual(expectMsg, gotMsg) {
		t.Errorf("invalid write message\n expect: %v\n got %v", expectMsg, gotMsg)
	}
	gotMsg = <-mockUserSecond.writeChan
	if !reflect.DeepEqual(expectMsg, gotMsg) {
		t.Errorf("invalid write message\n expect: %v\n got %v", expectMsg, gotMsg)
	}

	mockUserFirst.resolveMessage(Message{GameMove, map[string]string{"position": "1"}})
	// check duplicate
	mockUserFirst.resolveMessage(Message{GameMove, map[string]string{"position": "1"}})
	// check empty position
	mockUserFirst.resolveMessage(Message{GameMove, map[string]string{"failedparam": "failedvalue"}})
	expectMsg = Message{
		GameMoved,
		map[string]string{
			"1": "CROSS",
			"2": "EMPTY",
			"3": "EMPTY",
			"4": "EMPTY",
			"5": "EMPTY",
			"6": "EMPTY",
			"7": "EMPTY",
			"8": "EMPTY",
			"9": "EMPTY",
		},
	}
	gotMsg = <-mockUserFirst.writeChan
	if !reflect.DeepEqual(expectMsg, gotMsg) {
		t.Errorf("invalid write message\n expect: %v\n got %v", expectMsg, gotMsg)
	}
	gotMsg = <-mockUserSecond.writeChan
	if !reflect.DeepEqual(expectMsg, gotMsg) {
		t.Errorf("invalid write message\n expect: %v\n got %v", expectMsg, gotMsg)
	}

	mockUserSecond.resolveMessage(Message{GameMove, map[string]string{"position": "7"}})
	expectMsg = Message{
		GameMoved,
		map[string]string{
			"1": "CROSS",
			"2": "EMPTY",
			"3": "EMPTY",
			"4": "EMPTY",
			"5": "EMPTY",
			"6": "EMPTY",
			"7": "ZERO",
			"8": "EMPTY",
			"9": "EMPTY",
		},
	}
	gotMsg = <-mockUserFirst.writeChan
	if !reflect.DeepEqual(expectMsg, gotMsg) {
		t.Errorf("invalid write message\n expect: %v\n got %v", expectMsg, gotMsg)
	}
	gotMsg = <-mockUserSecond.writeChan
	if !reflect.DeepEqual(expectMsg, gotMsg) {
		t.Errorf("invalid write message\n expect: %v\n got %v", expectMsg, gotMsg)
	}

	mockUserFirst.resolveMessage(Message{GameMove, map[string]string{"position": "2"}})
	expectMsg = Message{
		GameMoved,
		map[string]string{
			"1": "CROSS",
			"2": "CROSS",
			"3": "EMPTY",
			"4": "EMPTY",
			"5": "EMPTY",
			"6": "EMPTY",
			"7": "ZERO",
			"8": "EMPTY",
			"9": "EMPTY",
		},
	}
	gotMsg = <-mockUserFirst.writeChan
	if !reflect.DeepEqual(expectMsg, gotMsg) {
		t.Errorf("invalid write message\n expect: %v\n got %v", expectMsg, gotMsg)
	}
	gotMsg = <-mockUserSecond.writeChan
	if !reflect.DeepEqual(expectMsg, gotMsg) {
		t.Errorf("invalid write message\n expect: %v\n got %v", expectMsg, gotMsg)
	}

	mockUserSecond.resolveMessage(Message{GameMove, map[string]string{"position": "8"}})
	expectMsg = Message{
		GameMoved,
		map[string]string{
			"1": "CROSS",
			"2": "CROSS",
			"3": "EMPTY",
			"4": "EMPTY",
			"5": "EMPTY",
			"6": "EMPTY",
			"7": "ZERO",
			"8": "ZERO",
			"9": "EMPTY",
		},
	}
	gotMsg = <-mockUserFirst.writeChan
	if !reflect.DeepEqual(expectMsg, gotMsg) {
		t.Errorf("invalid write message\n expect: %v\n got %v", expectMsg, gotMsg)
	}
	gotMsg = <-mockUserSecond.writeChan
	if !reflect.DeepEqual(expectMsg, gotMsg) {
		t.Errorf("invalid write message\n expect: %v\n got %v", expectMsg, gotMsg)
	}

	mockUserFirst.resolveMessage(Message{GameMove, map[string]string{"position": "3"}})
	expectMsg = Message{
		GameMoved,
		map[string]string{
			"1": "CROSS",
			"2": "CROSS",
			"3": "CROSS",
			"4": "EMPTY",
			"5": "EMPTY",
			"6": "EMPTY",
			"7": "ZERO",
			"8": "ZERO",
			"9": "EMPTY",
		},
	}
	gotMsg = <-mockUserFirst.writeChan
	if !reflect.DeepEqual(expectMsg, gotMsg) {
		t.Errorf("invalid write message\n expect: %v\n got %v", expectMsg, gotMsg)
	}
	gotMsg = <-mockUserSecond.writeChan
	if !reflect.DeepEqual(expectMsg, gotMsg) {
		t.Errorf("invalid write message\n expect: %v\n got %v", expectMsg, gotMsg)
	}

	expectMsg = Message{
		GameWinner,
		map[string]string{
			"winner": "CROSS",
		},
	}
	gotMsg = <-mockUserFirst.writeChan
	if !reflect.DeepEqual(expectMsg, gotMsg) {
		t.Errorf("invalid write message\n expect: %v\n got %v", expectMsg, gotMsg)
	}
	gotMsg = <-mockUserSecond.writeChan
	if !reflect.DeepEqual(expectMsg, gotMsg) {
		t.Errorf("invalid write message\n expect: %v\n got %v", expectMsg, gotMsg)
	}
}

func TestGame_Draw(t *testing.T) {
	repository := InmemoryRepository()

	mockUserFirst := MockUserWithRepository(repository)
	mockUserSecond := MockUserWithRepository(repository)

	mockUserFirst.resolveMessage(Message{Login, map[string]string{"username": "1"}})
	expectMsg := Message{
		LoginSuccess,
		map[string]string{
			"uuid":     mockUserFirst.uuid,
			"username": "1",
		},
	}
	gotMsg := <-mockUserFirst.writeChan
	if !reflect.DeepEqual(expectMsg, gotMsg) {
		t.Errorf("invalid write message\n expect: %v\n got %v", expectMsg, gotMsg)
	}

	mockUserSecond.resolveMessage(Message{Login, map[string]string{"username": "2"}})
	expectMsg = Message{
		LoginSuccess,
		map[string]string{
			"uuid":     mockUserSecond.uuid,
			"username": "2",
		},
	}
	gotMsg = <-mockUserSecond.writeChan
	if !reflect.DeepEqual(expectMsg, gotMsg) {
		t.Errorf("invalid write message\n expect: %v\n got %v", expectMsg, gotMsg)
	}

	mockUserFirst.resolveMessage(Message{GameSearchOn, map[string]string{}})
	mockUserSecond.resolveMessage(Message{GameSearchOn, map[string]string{}})

	ctx, cancel := context.WithCancel(context.Background())
	go gameSessionsCreator(repository, ctx, time.Tick(1*time.Nanosecond))

	time.Sleep(200 * time.Millisecond)
	cancel()

	expectMsg = Message{
		GameSearchOff,
		map[string]string{},
	}
	gotMsg = <-mockUserFirst.writeChan
	if !reflect.DeepEqual(expectMsg, gotMsg) {
		t.Errorf("invalid write message\n expect: %v\n got %v", expectMsg, gotMsg)
	}
	gotMsg = <-mockUserSecond.writeChan
	if !reflect.DeepEqual(expectMsg, gotMsg) {
		t.Errorf("invalid write message\n expect: %v\n got %v", expectMsg, gotMsg)
	}

	expectMsg = Message{
		GameSearchStart,
		map[string]string{
			"gameUUID":      mockUserFirst.currentGameUUID,
			"crossUserUUID": mockUserFirst.uuid,
			"zeroUserUUID":  mockUserSecond.uuid,
		},
	}
	gotMsg = <-mockUserFirst.writeChan
	if !reflect.DeepEqual(expectMsg, gotMsg) {
		t.Errorf("invalid write message\n expect: %v\n got %v", expectMsg, gotMsg)
	}
	gotMsg = <-mockUserSecond.writeChan
	if !reflect.DeepEqual(expectMsg, gotMsg) {
		t.Errorf("invalid write message\n expect: %v\n got %v", expectMsg, gotMsg)
	}

	mockUserFirst.resolveMessage(Message{GameMove, map[string]string{"position": "5"}})
	expectMsg = Message{
		GameMoved,
		map[string]string{
			"1": "EMPTY",
			"2": "EMPTY",
			"3": "EMPTY",
			"4": "EMPTY",
			"5": "CROSS",
			"6": "EMPTY",
			"7": "EMPTY",
			"8": "EMPTY",
			"9": "EMPTY",
		},
	}
	gotMsg = <-mockUserFirst.writeChan
	if !reflect.DeepEqual(expectMsg, gotMsg) {
		t.Errorf("invalid write message\n expect: %v\n got %v", expectMsg, gotMsg)
	}
	gotMsg = <-mockUserSecond.writeChan
	if !reflect.DeepEqual(expectMsg, gotMsg) {
		t.Errorf("invalid write message\n expect: %v\n got %v", expectMsg, gotMsg)
	}

	mockUserSecond.resolveMessage(Message{GameMove, map[string]string{"position": "3"}})
	expectMsg = Message{
		GameMoved,
		map[string]string{
			"1": "EMPTY",
			"2": "EMPTY",
			"3": "ZERO",
			"4": "EMPTY",
			"5": "CROSS",
			"6": "EMPTY",
			"7": "EMPTY",
			"8": "EMPTY",
			"9": "EMPTY",
		},
	}
	gotMsg = <-mockUserFirst.writeChan
	if !reflect.DeepEqual(expectMsg, gotMsg) {
		t.Errorf("invalid write message\n expect: %v\n got %v", expectMsg, gotMsg)
	}
	gotMsg = <-mockUserSecond.writeChan
	if !reflect.DeepEqual(expectMsg, gotMsg) {
		t.Errorf("invalid write message\n expect: %v\n got %v", expectMsg, gotMsg)
	}

	mockUserFirst.resolveMessage(Message{GameMove, map[string]string{"position": "9"}})
	expectMsg = Message{
		GameMoved,
		map[string]string{
			"1": "EMPTY",
			"2": "EMPTY",
			"3": "ZERO",
			"4": "EMPTY",
			"5": "CROSS",
			"6": "EMPTY",
			"7": "EMPTY",
			"8": "EMPTY",
			"9": "CROSS",
		},
	}
	gotMsg = <-mockUserFirst.writeChan
	if !reflect.DeepEqual(expectMsg, gotMsg) {
		t.Errorf("invalid write message\n expect: %v\n got %v", expectMsg, gotMsg)
	}
	gotMsg = <-mockUserSecond.writeChan
	if !reflect.DeepEqual(expectMsg, gotMsg) {
		t.Errorf("invalid write message\n expect: %v\n got %v", expectMsg, gotMsg)
	}

	mockUserSecond.resolveMessage(Message{GameMove, map[string]string{"position": "1"}})
	expectMsg = Message{
		GameMoved,
		map[string]string{
			"1": "ZERO",
			"2": "EMPTY",
			"3": "ZERO",
			"4": "EMPTY",
			"5": "CROSS",
			"6": "EMPTY",
			"7": "EMPTY",
			"8": "EMPTY",
			"9": "CROSS",
		},
	}
	gotMsg = <-mockUserFirst.writeChan
	if !reflect.DeepEqual(expectMsg, gotMsg) {
		t.Errorf("invalid write message\n expect: %v\n got %v", expectMsg, gotMsg)
	}
	gotMsg = <-mockUserSecond.writeChan
	if !reflect.DeepEqual(expectMsg, gotMsg) {
		t.Errorf("invalid write message\n expect: %v\n got %v", expectMsg, gotMsg)
	}

	mockUserFirst.resolveMessage(Message{GameMove, map[string]string{"position": "2"}})
	expectMsg = Message{
		GameMoved,
		map[string]string{
			"1": "ZERO",
			"2": "CROSS",
			"3": "ZERO",
			"4": "EMPTY",
			"5": "CROSS",
			"6": "EMPTY",
			"7": "EMPTY",
			"8": "EMPTY",
			"9": "CROSS",
		},
	}
	gotMsg = <-mockUserFirst.writeChan
	if !reflect.DeepEqual(expectMsg, gotMsg) {
		t.Errorf("invalid write message\n expect: %v\n got %v", expectMsg, gotMsg)
	}
	gotMsg = <-mockUserSecond.writeChan
	if !reflect.DeepEqual(expectMsg, gotMsg) {
		t.Errorf("invalid write message\n expect: %v\n got %v", expectMsg, gotMsg)
	}

	mockUserSecond.resolveMessage(Message{GameMove, map[string]string{"position": "8"}})
	expectMsg = Message{
		GameMoved,
		map[string]string{
			"1": "ZERO",
			"2": "CROSS",
			"3": "ZERO",
			"4": "EMPTY",
			"5": "CROSS",
			"6": "EMPTY",
			"7": "EMPTY",
			"8": "ZERO",
			"9": "CROSS",
		},
	}
	gotMsg = <-mockUserFirst.writeChan
	if !reflect.DeepEqual(expectMsg, gotMsg) {
		t.Errorf("invalid write message\n expect: %v\n got %v", expectMsg, gotMsg)
	}
	gotMsg = <-mockUserSecond.writeChan
	if !reflect.DeepEqual(expectMsg, gotMsg) {
		t.Errorf("invalid write message\n expect: %v\n got %v", expectMsg, gotMsg)
	}

	mockUserFirst.resolveMessage(Message{GameMove, map[string]string{"position": "4"}})
	expectMsg = Message{
		GameMoved,
		map[string]string{
			"1": "ZERO",
			"2": "CROSS",
			"3": "ZERO",
			"4": "CROSS",
			"5": "CROSS",
			"6": "EMPTY",
			"7": "EMPTY",
			"8": "ZERO",
			"9": "CROSS",
		},
	}
	gotMsg = <-mockUserFirst.writeChan
	if !reflect.DeepEqual(expectMsg, gotMsg) {
		t.Errorf("invalid write message\n expect: %v\n got %v", expectMsg, gotMsg)
	}
	gotMsg = <-mockUserSecond.writeChan
	if !reflect.DeepEqual(expectMsg, gotMsg) {
		t.Errorf("invalid write message\n expect: %v\n got %v", expectMsg, gotMsg)
	}

	mockUserSecond.resolveMessage(Message{GameMove, map[string]string{"position": "6"}})
	expectMsg = Message{
		GameMoved,
		map[string]string{
			"1": "ZERO",
			"2": "CROSS",
			"3": "ZERO",
			"4": "CROSS",
			"5": "CROSS",
			"6": "ZERO",
			"7": "EMPTY",
			"8": "ZERO",
			"9": "CROSS",
		},
	}
	gotMsg = <-mockUserFirst.writeChan
	if !reflect.DeepEqual(expectMsg, gotMsg) {
		t.Errorf("invalid write message\n expect: %v\n got %v", expectMsg, gotMsg)
	}
	gotMsg = <-mockUserSecond.writeChan
	if !reflect.DeepEqual(expectMsg, gotMsg) {
		t.Errorf("invalid write message\n expect: %v\n got %v", expectMsg, gotMsg)
	}

	mockUserFirst.resolveMessage(Message{GameMove, map[string]string{"position": "7"}})
	expectMsg = Message{
		GameMoved,
		map[string]string{
			"1": "ZERO",
			"2": "CROSS",
			"3": "ZERO",
			"4": "CROSS",
			"5": "CROSS",
			"6": "ZERO",
			"7": "CROSS",
			"8": "ZERO",
			"9": "CROSS",
		},
	}
	gotMsg = <-mockUserFirst.writeChan
	if !reflect.DeepEqual(expectMsg, gotMsg) {
		t.Errorf("invalid write message\n expect: %v\n got %v", expectMsg, gotMsg)
	}
	gotMsg = <-mockUserSecond.writeChan
	if !reflect.DeepEqual(expectMsg, gotMsg) {
		t.Errorf("invalid write message\n expect: %v\n got %v", expectMsg, gotMsg)
	}

	expectMsg = Message{
		GameDraw,
		map[string]string{},
	}
	gotMsg = <-mockUserFirst.writeChan
	if !reflect.DeepEqual(expectMsg, gotMsg) {
		t.Errorf("invalid write message\n expect: %v\n got %v", expectMsg, gotMsg)
	}
	gotMsg = <-mockUserSecond.writeChan
	if !reflect.DeepEqual(expectMsg, gotMsg) {
		t.Errorf("invalid write message\n expect: %v\n got %v", expectMsg, gotMsg)
	}
}
