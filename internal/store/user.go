package store

import (
	"context"
	"database/sql"
	"errors"

	"example.com/user-management/internal/db"
	"example.com/user-management/internal/model"
	"github.com/google/uuid"
)

type UserStore struct {
	queries *db.Queries
}

func NewUserStore(dbConn *sql.DB) *UserStore {
	return &UserStore{
		queries: db.New(dbConn),
	}
}

func (store *UserStore) CreateUser(user model.User) (model.User, error) {

	createdUser, err := store.queries.CreateUser(context.Background(),
		db.CreateUserParams{
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Email:     user.Email,
			Phone:     user.Phone,
			Age: sql.NullInt32{
				Int32: int32(user.Age),
				Valid: user.Age > 0,
			},
			Status: string(user.Status),
		},
	)

	if err != nil {
		return model.User{}, err
	}

	return mapDbUserToModel(&createdUser), nil
}

func (store *UserStore) GetAllUsers() ([]model.User, error) {

	dbUsers, err := store.queries.GetAllUsers(context.Background())
	if err != nil {
		return nil, err
	}

	users := make([]model.User, len(dbUsers))
	for i, u := range dbUsers {
		users[i] = mapDbUserToModel(&u)
	}

	return users, nil
}

func (store *UserStore) GetUserById(userId uuid.UUID) (model.User, bool, error) {
	dbUser, err := store.queries.GetUserByID(context.Background(), userId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.User{}, false, nil
		}
		return model.User{}, false, err
	}

	return mapDbUserToModel(&dbUser), true, nil
}

func (store *UserStore) UpdateUser(user model.User, userId uuid.UUID) (model.User, bool, error) {
	dbUser, err := store.queries.UpdateUser(
		context.Background(),
		db.UpdateUserParams{
			UserID:    userId,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Email:     user.Email,
			Phone:     user.Phone,
			Age: sql.NullInt32{
				Int32: int32(user.Age),
				Valid: user.Age > 0,
			},
			Status: string(user.Status),
		},
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.User{}, false, nil
		}
		return model.User{}, false, err
	}

	return mapDbUserToModel(&dbUser), true, nil

}

func (store *UserStore) DeleteUser(userId uuid.UUID) (bool, error) {
	err := store.queries.DeleteUser(context.Background(), userId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func mapDbUserToModel(dbUser *db.User) model.User {
	return model.User{
		UserId:    dbUser.UserID,
		FirstName: dbUser.FirstName,
		LastName:  dbUser.LastName,
		Email:     dbUser.Email,
		Phone:     dbUser.Phone,
		Age:       int(dbUser.Age.Int32),
		Status:    model.Status(dbUser.Status),
	}
}
