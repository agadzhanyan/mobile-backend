package main

import (
	"reflect"
	"sync"
	"testing"
)

func MockUser() *User {
	return &User{
		generateUUID(),
		generateUUID(),
		"",
		nil,
		make(chan Message),
		nil,
	}
}

func MockGame(crossUser, zeroUser *User) *Game {
	return &Game{
		generateUUID(),
		[]*User{crossUser, zeroUser},
		crossUser,
		zeroUser,
		"CROSS",
		false,
		GameField{},
	}
}

func TestGameRepository_UserByUUID(t *testing.T) {
	type fields struct {
		users              map[string]*User
		usersInSearch      map[string]*User
		gameSessions       map[string]*Game
		usersMutex         *sync.RWMutex
		usersInSearchMutex *sync.RWMutex
		gameSessionsMutex  *sync.RWMutex
	}
	type args struct {
		uuid string
	}
	mockUser := MockUser()
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *User
	}{
		{
			"user exists",
			fields{
				map[string]*User{
					mockUser.uuid: mockUser,
				},
				map[string]*User{},
				map[string]*Game{},
				&sync.RWMutex{},
				&sync.RWMutex{},
				&sync.RWMutex{},
			},
			args{
				mockUser.uuid,
			},
			mockUser,
		},
		{
			"user doesn't exists",
			fields{
				map[string]*User{
					mockUser.uuid: mockUser,
				},
				map[string]*User{},
				map[string]*Game{},
				&sync.RWMutex{},
				&sync.RWMutex{},
				&sync.RWMutex{},
			},
			args{
				generateUUID(),
			},
			nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gr := &GameRepository{
				users:              tt.fields.users,
				usersInSearch:      tt.fields.usersInSearch,
				gameSessions:       tt.fields.gameSessions,
				usersMutex:         tt.fields.usersMutex,
				usersInSearchMutex: tt.fields.usersInSearchMutex,
				gameSessionsMutex:  tt.fields.gameSessionsMutex,
			}
			if got := gr.UserByUUID(tt.args.uuid); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GameRepository.UserByUUID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGameRepository_GameByUUID(t *testing.T) {
	type fields struct {
		users              map[string]*User
		usersInSearch      map[string]*User
		gameSessions       map[string]*Game
		usersMutex         *sync.RWMutex
		usersInSearchMutex *sync.RWMutex
		gameSessionsMutex  *sync.RWMutex
	}
	type args struct {
		uuid string
	}
	mockCrossUser, mockZeroUser := MockUser(), MockUser()
	mockGame := MockGame(mockCrossUser, mockZeroUser)
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *Game
	}{
		{
			"game found",
			fields{
				map[string]*User{},
				map[string]*User{},
				map[string]*Game{
					mockGame.uuid: mockGame,
				},
				&sync.RWMutex{},
				&sync.RWMutex{},
				&sync.RWMutex{},
			},
			args{
				mockGame.uuid,
			},
			mockGame,
		},
		{
			"no games found",
			fields{
				map[string]*User{},
				map[string]*User{},
				map[string]*Game{
					mockGame.uuid: mockGame,
				},
				&sync.RWMutex{},
				&sync.RWMutex{},
				&sync.RWMutex{},
			},
			args{
				generateUUID(),
			},
			nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gr := &GameRepository{
				users:              tt.fields.users,
				usersInSearch:      tt.fields.usersInSearch,
				gameSessions:       tt.fields.gameSessions,
				usersMutex:         tt.fields.usersMutex,
				usersInSearchMutex: tt.fields.usersInSearchMutex,
				gameSessionsMutex:  tt.fields.gameSessionsMutex,
			}
			if got := gr.GameByUUID(tt.args.uuid); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GameRepository.GameByUUID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGameRepository_GameSessions(t *testing.T) {
	type fields struct {
		users              map[string]*User
		usersInSearch      map[string]*User
		gameSessions       map[string]*Game
		usersMutex         *sync.RWMutex
		usersInSearchMutex *sync.RWMutex
		gameSessionsMutex  *sync.RWMutex
	}
	mockCrossUser, mockZeroUser := MockUser(), MockUser()
	mockGame := MockGame(mockCrossUser, mockZeroUser)
	tests := []struct {
		name   string
		fields fields
		want   map[string]*Game
	}{
		{
			"no games found",
			fields{
				map[string]*User{},
				map[string]*User{},
				map[string]*Game{
					mockGame.uuid: mockGame,
				},
				&sync.RWMutex{},
				&sync.RWMutex{},
				&sync.RWMutex{},
			},
			map[string]*Game{
				mockGame.uuid: mockGame,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gr := &GameRepository{
				users:              tt.fields.users,
				usersInSearch:      tt.fields.usersInSearch,
				gameSessions:       tt.fields.gameSessions,
				usersMutex:         tt.fields.usersMutex,
				usersInSearchMutex: tt.fields.usersInSearchMutex,
				gameSessionsMutex:  tt.fields.gameSessionsMutex,
			}
			if got := gr.GameSessions(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GameRepository.GameSessions() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGameRepository_UsersInSearch(t *testing.T) {
	type fields struct {
		users              map[string]*User
		usersInSearch      map[string]*User
		gameSessions       map[string]*Game
		usersMutex         *sync.RWMutex
		usersInSearchMutex *sync.RWMutex
		gameSessionsMutex  *sync.RWMutex
	}
	mockUser := MockUser()
	tests := []struct {
		name   string
		fields fields
		want   map[string]*User
	}{
		{
			"user exists",
			fields{
				map[string]*User{},
				map[string]*User{
					mockUser.uuid: mockUser,
				},
				map[string]*Game{},
				&sync.RWMutex{},
				&sync.RWMutex{},
				&sync.RWMutex{},
			},
			map[string]*User{
				mockUser.uuid: mockUser,
			},
		},
		{
			"user doesn't exists",
			fields{
				map[string]*User{},
				map[string]*User{},
				map[string]*Game{},
				&sync.RWMutex{},
				&sync.RWMutex{},
				&sync.RWMutex{},
			},
			map[string]*User{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gr := &GameRepository{
				users:              tt.fields.users,
				usersInSearch:      tt.fields.usersInSearch,
				gameSessions:       tt.fields.gameSessions,
				usersMutex:         tt.fields.usersMutex,
				usersInSearchMutex: tt.fields.usersInSearchMutex,
				gameSessionsMutex:  tt.fields.gameSessionsMutex,
			}
			if got := gr.UsersInSearch(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GameRepository.UsersInSearch() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGameRepository_UsersInSearchInsertionOrder(t *testing.T) {
	type fields struct {
		users              map[string]*User
		usersInSearch      map[string]*User
		usersInSerachKeys  []string
		gameSessions       map[string]*Game
		usersMutex         *sync.RWMutex
		usersInSearchMutex *sync.RWMutex
		gameSessionsMutex  *sync.RWMutex
	}
	mockUser := MockUser()
	tests := []struct {
		name   string
		fields fields
		want   []*User
	}{
		{
			"user exists",
			fields{
				map[string]*User{},
				map[string]*User{
					mockUser.uuid: mockUser,
				},
				[]string{
					mockUser.uuid,
				},
				map[string]*Game{},
				&sync.RWMutex{},
				&sync.RWMutex{},
				&sync.RWMutex{},
			},
			[]*User{
				mockUser,
			},
		},
		{
			"user doesn't exists",
			fields{
				map[string]*User{},
				map[string]*User{},
				[]string{},
				map[string]*Game{},
				&sync.RWMutex{},
				&sync.RWMutex{},
				&sync.RWMutex{},
			},
			[]*User{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gr := &GameRepository{
				users:              tt.fields.users,
				usersInSearch:      tt.fields.usersInSearch,
				usersInSearchKeys:  tt.fields.usersInSerachKeys,
				gameSessions:       tt.fields.gameSessions,
				usersMutex:         tt.fields.usersMutex,
				usersInSearchMutex: tt.fields.usersInSearchMutex,
				gameSessionsMutex:  tt.fields.gameSessionsMutex,
			}
			if got := gr.UsersInSearchInsertionOrder(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GameRepository.UsersInSearchInsertionOrder() = %T, want %T", got, tt.want)
			}
		})
	}
}

func TestGameRepository_AddUser(t *testing.T) {
	type fields struct {
		users              map[string]*User
		usersInSearch      map[string]*User
		gameSessions       map[string]*Game
		usersMutex         *sync.RWMutex
		usersInSearchMutex *sync.RWMutex
		gameSessionsMutex  *sync.RWMutex
	}
	type args struct {
		user *User
	}
	mockUser := MockUser()
	tests := []struct {
		name   string
		fields fields
		args   args
		want   map[string]*User
	}{
		{
			"user added successfully",
			fields{
				map[string]*User{},
				map[string]*User{},
				map[string]*Game{},
				&sync.RWMutex{},
				&sync.RWMutex{},
				&sync.RWMutex{},
			},
			args{
				mockUser,
			},
			map[string]*User{
				mockUser.uuid: mockUser,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gr := &GameRepository{
				users:              tt.fields.users,
				usersInSearch:      tt.fields.usersInSearch,
				gameSessions:       tt.fields.gameSessions,
				usersMutex:         tt.fields.usersMutex,
				usersInSearchMutex: tt.fields.usersInSearchMutex,
				gameSessionsMutex:  tt.fields.gameSessionsMutex,
			}
			gr.AddUser(tt.args.user)
			if got := gr.Users(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GameRepository.AddUser() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGameRepository_RemoveUser(t *testing.T) {
	type fields struct {
		users              map[string]*User
		usersInSearch      map[string]*User
		gameSessions       map[string]*Game
		usersMutex         *sync.RWMutex
		usersInSearchMutex *sync.RWMutex
		gameSessionsMutex  *sync.RWMutex
	}
	type args struct {
		user *User
	}
	mockUser := MockUser()
	mockUserAdditional := MockUser()
	tests := []struct {
		name   string
		fields fields
		args   args
		want   map[string]*User
	}{
		{
			"user removed successfully",
			fields{
				map[string]*User{
					mockUser.uuid: mockUser,
				},
				map[string]*User{},
				map[string]*Game{},
				&sync.RWMutex{},
				&sync.RWMutex{},
				&sync.RWMutex{},
			},
			args{
				mockUser,
			},
			map[string]*User{},
		},
		{
			"no users removed",
			fields{
				map[string]*User{
					mockUser.uuid: mockUser,
				},
				map[string]*User{},
				map[string]*Game{},
				&sync.RWMutex{},
				&sync.RWMutex{},
				&sync.RWMutex{},
			},
			args{
				mockUserAdditional,
			},
			map[string]*User{
				mockUser.uuid: mockUser,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gr := &GameRepository{
				users:              tt.fields.users,
				usersInSearch:      tt.fields.usersInSearch,
				gameSessions:       tt.fields.gameSessions,
				usersMutex:         tt.fields.usersMutex,
				usersInSearchMutex: tt.fields.usersInSearchMutex,
				gameSessionsMutex:  tt.fields.gameSessionsMutex,
			}
			gr.RemoveUser(tt.args.user)
			if got := gr.Users(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GameRepository.RemoveUser() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGameRepository_AddUserInSearch(t *testing.T) {
	type fields struct {
		users              map[string]*User
		usersInSearch      map[string]*User
		gameSessions       map[string]*Game
		usersMutex         *sync.RWMutex
		usersInSearchMutex *sync.RWMutex
		gameSessionsMutex  *sync.RWMutex
	}
	type args struct {
		user *User
	}
	mockUser := MockUser()
	tests := []struct {
		name   string
		fields fields
		args   args
		want   map[string]*User
	}{
		{
			"no users removed",
			fields{
				map[string]*User{},
				map[string]*User{},
				map[string]*Game{},
				&sync.RWMutex{},
				&sync.RWMutex{},
				&sync.RWMutex{},
			},
			args{
				mockUser,
			},
			map[string]*User{
				mockUser.uuid: mockUser,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gr := &GameRepository{
				users:              tt.fields.users,
				usersInSearch:      tt.fields.usersInSearch,
				gameSessions:       tt.fields.gameSessions,
				usersMutex:         tt.fields.usersMutex,
				usersInSearchMutex: tt.fields.usersInSearchMutex,
				gameSessionsMutex:  tt.fields.gameSessionsMutex,
			}
			gr.AddUserInSearch(tt.args.user)
			gr.AddUser(tt.args.user)
			if got := gr.UsersInSearch(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GameRepository.AddUserInSearch() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGameRepository_RemoveUserInSearch(t *testing.T) {
	type fields struct {
		users              map[string]*User
		usersInSearch      map[string]*User
		gameSessions       map[string]*Game
		usersMutex         *sync.RWMutex
		usersInSearchMutex *sync.RWMutex
		gameSessionsMutex  *sync.RWMutex
	}
	type args struct {
		user *User
	}
	mockUser := MockUser()
	mockUserAdditional := MockUser()
	tests := []struct {
		name   string
		fields fields
		args   args
		want   map[string]*User
	}{
		{
			"user removed",
			fields{
				map[string]*User{},
				map[string]*User{
					mockUser.uuid: mockUser,
				},
				map[string]*Game{},
				&sync.RWMutex{},
				&sync.RWMutex{},
				&sync.RWMutex{},
			},
			args{
				mockUser,
			},
			map[string]*User{},
		},
		{
			"no users removed",
			fields{
				map[string]*User{},
				map[string]*User{
					mockUser.uuid: mockUser,
				},
				map[string]*Game{},
				&sync.RWMutex{},
				&sync.RWMutex{},
				&sync.RWMutex{},
			},
			args{
				mockUserAdditional,
			},
			map[string]*User{
				mockUser.uuid: mockUser,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gr := &GameRepository{
				users:              tt.fields.users,
				usersInSearch:      tt.fields.usersInSearch,
				gameSessions:       tt.fields.gameSessions,
				usersMutex:         tt.fields.usersMutex,
				usersInSearchMutex: tt.fields.usersInSearchMutex,
				gameSessionsMutex:  tt.fields.gameSessionsMutex,
			}
			gr.RemoveUserInSearch(tt.args.user)
			gr.AddUser(tt.args.user)
			if got := gr.UsersInSearch(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GameRepository.RemoveUserInSearch() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGameRepository_AddGame(t *testing.T) {
	type fields struct {
		users              map[string]*User
		usersInSearch      map[string]*User
		gameSessions       map[string]*Game
		usersMutex         *sync.RWMutex
		usersInSearchMutex *sync.RWMutex
		gameSessionsMutex  *sync.RWMutex
	}
	type args struct {
		game *Game
	}
	mockCrossUser, mockZeroUser := MockUser(), MockUser()
	mockGame := MockGame(mockCrossUser, mockZeroUser)
	tests := []struct {
		name   string
		fields fields
		args   args
		want   map[string]*Game
	}{
		{
			"game added",
			fields{
				map[string]*User{},
				map[string]*User{},
				map[string]*Game{},
				&sync.RWMutex{},
				&sync.RWMutex{},
				&sync.RWMutex{},
			},
			args{
				mockGame,
			},
			map[string]*Game{
				mockGame.uuid: mockGame,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gr := &GameRepository{
				users:              tt.fields.users,
				usersInSearch:      tt.fields.usersInSearch,
				gameSessions:       tt.fields.gameSessions,
				usersMutex:         tt.fields.usersMutex,
				usersInSearchMutex: tt.fields.usersInSearchMutex,
				gameSessionsMutex:  tt.fields.gameSessionsMutex,
			}
			gr.AddGame(tt.args.game)
			if got := gr.GameSessions(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GameRepository.AddGame() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGameRepository_RemoveGame(t *testing.T) {
	type fields struct {
		users              map[string]*User
		usersInSearch      map[string]*User
		gameSessions       map[string]*Game
		usersMutex         *sync.RWMutex
		usersInSearchMutex *sync.RWMutex
		gameSessionsMutex  *sync.RWMutex
	}
	type args struct {
		game *Game
	}
	mockCrossUser, mockZeroUser := MockUser(), MockUser()
	mockGame := MockGame(mockCrossUser, mockZeroUser)
	mockGameAdditional := MockGame(mockCrossUser, mockZeroUser)
	tests := []struct {
		name   string
		fields fields
		args   args
		want   map[string]*Game
	}{
		{
			"game removed",
			fields{
				map[string]*User{},
				map[string]*User{},
				map[string]*Game{
					mockGame.uuid: mockGame,
				},
				&sync.RWMutex{},
				&sync.RWMutex{},
				&sync.RWMutex{},
			},
			args{
				mockGame,
			},
			map[string]*Game{},
		},
		{
			"no games removed",
			fields{
				map[string]*User{},
				map[string]*User{},
				map[string]*Game{
					mockGame.uuid: mockGame,
				},
				&sync.RWMutex{},
				&sync.RWMutex{},
				&sync.RWMutex{},
			},
			args{
				mockGameAdditional,
			},
			map[string]*Game{
				mockGame.uuid: mockGame,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gr := &GameRepository{
				users:              tt.fields.users,
				usersInSearch:      tt.fields.usersInSearch,
				gameSessions:       tt.fields.gameSessions,
				usersMutex:         tt.fields.usersMutex,
				usersInSearchMutex: tt.fields.usersInSearchMutex,
				gameSessionsMutex:  tt.fields.gameSessionsMutex,
			}
			gr.RemoveGame(tt.args.game)
			if got := gr.GameSessions(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GameRepository.RemoveGame() = %v, want %v", got, tt.want)
			}
		})
	}
}
