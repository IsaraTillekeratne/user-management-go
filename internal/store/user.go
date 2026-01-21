package store

import (
	"example.com/user-management/internal/model"
	"github.com/google/uuid"
)

type UserStore struct {
	users map[uuid.UUID]model.User
}

func NewUserStore() *UserStore {
	return &UserStore{
		users: make(map[uuid.UUID]model.User),
	}
}

func (store *UserStore) CreateUser(user model.User) model.User {
	user.UserId = uuid.New()
	store.users[user.UserId] = user
	return user
}

func (store *UserStore) GetAllUsers() []model.User {
	users := make([]model.User, 0, len(store.users))
	for _, v := range store.users {
		users = append(users, v)
	}
	return users
}

func (store *UserStore) GetUserById(userId uuid.UUID) (model.User, bool) {
	u, ok := store.users[userId]
	return u, ok
}

func (store *UserStore) UpdateUser(user model.User, userId uuid.UUID) (model.User, bool) {
	_, ok := store.users[userId]
	if !ok {
		return model.User{}, false
	}
	user.UserId = userId
	store.users[userId] = user
	return store.users[userId], true
}

func (store *UserStore) DeleteUser(userId uuid.UUID) bool {
	_, ok := store.users[userId]
	if !ok {
		return false
	}
	delete(store.users, userId)
	return true
}
