package auth

import (
	"context"

	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/models"
)

//go:generate mockgen -source=auth.go -destination=mocks/mock.go

// Usecase includes bussiness logics methods to work with authentication
type Usecase interface {
	// SignUpUser creates new user and returns it's id
	SignUpUser(ctx context.Context, user models.User) (uint32, error)

	// GetUserByCreds returns User if such User exists in repository
	GetUserByCreds(ctx context.Context, username, plainPassword string) (*models.User, error)

	// GetUserByAuthData returns User if such User exists in repository
	GetUserByAuthData(ctx context.Context, userID, userVersion uint32) (*models.User, error)

	// IncreaseUserVersion increases user's access token version
	IncreaseUserVersion(ctx context.Context, userID uint32) error

	ChangePassword(ctx context.Context, userID uint32, password string) error
}

// Repository includes DBMS-relatable methods to work with authentication
type Repository interface {
	GetUserByAuthData(ctx context.Context, userID, userVersion uint32) (*models.User, error)
	IncreaseUserVersion(ctx context.Context, userID uint32) error
	UpdatePassword(ctx context.Context, userID uint32, passwordHash, salt string) error
}

// Tables includes methods which return needed tables
// to work with auth on repository-layer
type Tables interface {
	Users() string
}
