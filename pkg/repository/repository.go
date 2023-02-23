package repository

import "github.com/go-park-mail-ru/2023_1_Technokaif/models"

type Repository struct {
	Auth
	// Other
}

type Auth interface {
	CreateUser(user models.User) (int, error)
	GetUser(username, password string) (models.User, error)
}
