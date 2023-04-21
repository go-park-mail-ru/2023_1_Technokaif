package auth

import (
	"context"

	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/models"
)

//go:generate mockgen -source=auth.go -destination=mocks/mock.go

// Usecase includes bussiness logics methods to work with authentication
type Usecase interface {
	// SignUpUser creates new user and returns it's id
	SignUpUser(user models.User) (uint32, error)

	// GetUserByCreds returns User if such User exists in repository
	GetUserByCreds(username, plainPassword string) (*models.User, error)

	// GetUserByAuthData returns User if such User exists in repository
	GetUserByAuthData(userID, userVersion uint32) (*models.User, error)

	// IncreaseUserVersion increases user's access token version
	IncreaseUserVersion(userID uint32) error

	ChangePassword(userID uint32, password string) error
}

// Agent ...
type Agent interface {
	SignUpUser(u models.User, context context.Context) (uint32, error)
}

// Repository includes DBMS-relatable methods to work with authentication
type Repository interface {
	GetUserByAuthData(userID, userVersion uint32) (*models.User, error)
	IncreaseUserVersion(userID uint32) error
	UpdatePassword(userID uint32, passwordHash, salt string) error
}

// Tables includes methods which return needed tables
// to work with auth on repository-layer
type Tables interface {
	Users() string
}
