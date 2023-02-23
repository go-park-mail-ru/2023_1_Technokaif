package usecase

import "github.com/go-park-mail-ru/2023_1_Technokaif/models"

// Usecase implements all current app's services
type Usecase struct {
	Auth
	// Other services
}

// Auth describes which methods have to be implemented by auth-service
type Auth interface {

	// CreateUser creates new entity of user and returns it's id
	CreateUser(user models.User) (int, error)

	// GenerateToken returns token created with user's username and password
	GenerateToken(username, password string) (string, error)
}
