package auth

import (
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/models"
)

//go:generate mockgen -source=auth.go -destination=mocks/mock.go

// Usecase includes bussiness logics methods to work with authentication
type Usecase interface {
	// CreateUser creates new user and returns it's id
	SignUpUser(user models.User) (uint32, error)

	// GetUserID returns User if such User exists in repository
	GetUserByCreds(username, plainPassword string) (*models.User, error)

	// LoginUser finds user in repository and returns its access token
	LoginUser(username, plainPassword string) (string, error)

	// GetUserByAuthData returns User if such User exists in repository
	GetUserByAuthData(userID, userVersion uint32) (*models.User, error)

	// GenerateAccessToken returns token created with user's username and password
	GenerateAccessToken(userID, userVersion uint32) (string, error)

	// CheckAccessToken validates accessToken
	CheckAccessToken(accessToken string) (uint32, uint32, error)

	// IncreaseUserVersion increases user's access token version
	IncreaseUserVersion(userID uint32) error

	ChangePassword(userID uint32, password string) error
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
