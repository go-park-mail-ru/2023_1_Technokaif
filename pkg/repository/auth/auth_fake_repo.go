package repository

import "github.com/go-park-mail-ru/2023_1_Technokaif/models"

type AuthFake struct {
	// fake repo
}

func (a *AuthFake) CreateUser(user models.User) (int, error) {
	return 0, nil
}

func (a *AuthFake) GetUser(username, password string) (models.User, error) {
	return models.User{
		Username:  "manFromFakeRepo1337",
		FirstName: "Vasya",
		LastName:  "Pupkin",
	}, nil
}
