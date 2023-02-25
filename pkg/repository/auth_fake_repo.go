package repository

import (
	"errors"

	"github.com/go-park-mail-ru/2023_1_Technokaif/models"
)

type AuthFake struct {
	db []models.User
}

func (a *AuthFake) CreateUser(user models.User) (int, error) {
	a.db = append(a.db, user)
	return len(a.db) - 1, nil // No errors in fake DBMS!!!
}

func (a *AuthFake) GetUser(username, password string) (models.User, error) {

	for _, u := range a.db {
		if u.Username == username {
			return u, nil
		}
	}
	return models.User{}, errors.New("Error while getting user from fake repo lol")
}

func NewAuthFake() *AuthFake {
	return &AuthFake{
		db: []models.User{
			{
				ID:       1,
				Username: "yarik_tri",
			},
		},
	}
}
