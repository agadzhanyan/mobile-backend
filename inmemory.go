package main

import "sync"

func InmemoryRepository() IRepository {
	return &GameRepository{
		make(map[string]*User),
		make(map[string]*User),
		[]string{},
		make(map[string]*Game),
		&sync.RWMutex{},
		&sync.RWMutex{},
		&sync.RWMutex{},
	}
}

type GameRepository struct {
	users                                             map[string]*User
	usersInSearch                                     map[string]*User
	usersInSearchKeys                                 []string
	gameSessions                                      map[string]*Game
	usersMutex, usersInSearchMutex, gameSessionsMutex *sync.RWMutex
}

func (gr *GameRepository) UserByUUID(uuid string) *User {
	gr.usersMutex.Lock()
	user, ok := gr.users[uuid]
	gr.usersMutex.Unlock()
	if ok {
		return user
	}
	return nil
}

func (gr *GameRepository) GameByUUID(uuid string) *Game {
	gr.gameSessionsMutex.Lock()
	game, ok := gr.gameSessions[uuid]
	gr.gameSessionsMutex.Unlock()
	if ok {
		return game
	}
	return nil
}

func (gr *GameRepository) GameSessions() map[string]*Game {
	return gr.gameSessions
}

func (gr *GameRepository) UsersInSearch() map[string]*User {
	return gr.usersInSearch
}

func (gr *GameRepository) UsersInSearchInsertionOrder() []*User {
	slice := []*User{}
	for _, key := range gr.usersInSearchKeys {
		slice = append(slice, gr.usersInSearch[key])
	}
	return slice
}

func (gr *GameRepository) Users() map[string]*User {
	return gr.users
}

func (gr *GameRepository) AddUser(user *User) {
	gr.usersMutex.RLock()
	_, ok := gr.users[user.uuid]
	gr.usersMutex.RUnlock()
	if !ok {
		gr.usersMutex.Lock()
		gr.users[user.uuid] = user
		gr.usersMutex.Unlock()
	}
}

func (gr *GameRepository) RemoveUser(user *User) {
	gr.usersMutex.RLock()
	_, ok := gr.users[user.uuid]
	var newSlice []string
	for _, v := range gr.usersInSearchKeys {
		if v == user.uuid {
			continue
		}
		newSlice = append(newSlice, v)
	}

	gr.usersMutex.RUnlock()
	if ok {
		gr.usersMutex.Lock()
		delete(gr.users, user.uuid)
		gr.usersMutex.Unlock()
	}
}

func (gr *GameRepository) AddUserInSearch(user *User) {
	gr.usersInSearchMutex.RLock()
	_, ok := gr.usersInSearch[user.uuid]
	gr.usersInSearchKeys = append(gr.usersInSearchKeys, user.uuid)
	gr.usersInSearchMutex.RUnlock()
	if !ok {
		gr.usersInSearchMutex.Lock()
		gr.usersInSearch[user.uuid] = user
		gr.usersInSearchMutex.Unlock()
	}
}

func (gr *GameRepository) RemoveUserInSearch(user *User) {
	gr.usersInSearchMutex.RLock()
	_, ok := gr.usersInSearch[user.uuid]
	gr.usersInSearchMutex.RUnlock()
	if ok {
		gr.usersInSearchMutex.Lock()
		delete(gr.usersInSearch, user.uuid)
		gr.usersInSearchMutex.Unlock()
	}
}

func (gr *GameRepository) AddGame(game *Game) {
	gr.gameSessionsMutex.RLock()
	_, ok := gr.gameSessions[game.uuid]
	gr.gameSessionsMutex.RUnlock()
	if !ok {
		gr.gameSessionsMutex.Lock()
		gr.gameSessions[game.uuid] = game
		gr.gameSessionsMutex.Unlock()
	}
}

func (gr *GameRepository) RemoveGame(game *Game) {
	gr.gameSessionsMutex.RLock()
	_, ok := gr.gameSessions[game.uuid]
	gr.gameSessionsMutex.RUnlock()
	if ok {
		gr.gameSessionsMutex.Lock()
		delete(gr.gameSessions, game.uuid)
		gr.gameSessionsMutex.Unlock()
	}
}
