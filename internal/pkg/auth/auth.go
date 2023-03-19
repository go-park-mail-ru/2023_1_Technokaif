package auth

import (
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/models"
)

// AuthUsecase describes which methods have to be implemented by auth-service
type AuthUsecase interface {
	// LoginUser finds user in repository and returns its access token
	LoginUser(username, password string) (string, error)

	// CreateUser creates new user and returns it's id
	CreateUser(user models.User) (uint32, error)

	// GetUserID returns User if such User exists in repository
	GetUserByCreds(username, password string) (*models.User, error)

	// GetUserByAuthData returns User if such User exists in repository
	GetUserByAuthData(userID, userVersion uint32) (*models.User, error)

	// GenerateAccessToken returns token created with user's username and password
	GenerateAccessToken(userID, userVersion uint32) (string, error)

	// CheckAccessToken validates accessToken
	CheckAccessToken(accessToken string) (uint32, uint32, error)

	// IncreaseUserVersion increases user's access token version
	IncreaseUserVersion(userID uint32) error
}

// AuthRepository includes DBMS-relatable methods for authentication
type AuthRepository interface {
	// CreateUser inserts new user into DB and return it's id
	// or error if it already exists
	CreateUser(user models.User) (uint32, error)

	// GetUserByUsername returns models.User if it's entry in DB exists or error otherwise
	GetUserByUsername(username string) (*models.User, error)

	GetUserByAuthData(userID, userVersion uint32) (*models.User, error)

	IncreaseUserVersion(userID uint32) error
}
