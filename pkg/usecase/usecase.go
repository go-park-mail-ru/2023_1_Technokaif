package usecase

import (
	"github.com/go-park-mail-ru/2023_1_Technokaif/models"
	"github.com/go-park-mail-ru/2023_1_Technokaif/pkg/repository"
)

// Usecase implements all current app's services
type Usecase struct {
	Auth
}

// Auth describes which methods have to be implemented by auth-service
type Auth interface {

	// CreateUser creates new entity of user and returns it's id
	CreateUser(user models.User) (int, error)

	// GetUserID gets User's ID if such User exists
	GetUserID(username, password string) (uint, error)

	// GenerateToken returns token created with user's username and password
	GenerateToken(userID uint) (string, error)
}

func NewUsecase(r *repository.Repository) *Usecase {
	return &Usecase{
		Auth: NewAuthUsecase(r.Auth),
	}
}
