package usecase

import (
	"context"
	"fmt"

	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/models"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/auth"
)

// Usecase implements auth.Usecase
type Usecase struct {
	authAgent auth.Agent
}

func NewUsecase(aa auth.Agent) *Usecase {
	return &Usecase{
		authAgent: aa,
	}
}

func (u *Usecase) SignUpUser(ctx context.Context, user models.User) (uint32, error) {
	userId, err := u.authAgent.SignUpUser(context.Background(), user) // TODO request context
	if err != nil {
		return 0, fmt.Errorf("(usecase) can't sign up user: %w", err)
	}
	return userId, nil
}

func (u *Usecase) GetUserByCreds(ctx context.Context, username, password string) (*models.User, error) {
	user, err := u.authAgent.GetUserByCreds(context.Background(), username, password)
	if err != nil {
		return nil, fmt.Errorf("(usecase) can't get user: %w", err)
	}
	return user, err
}

func (u *Usecase) GetUserByAuthData(ctx context.Context, userID, userVersion uint32) (*models.User, error) {
	user, err := u.authAgent.GetUserByAuthData(context.Background(), userID, userVersion)
	if err != nil {
		return nil, fmt.Errorf("(usecase) can't get user: %w", err)
	}
	return user, err
}

func (u *Usecase) IncreaseUserVersion(ctx context.Context, userID uint32) error {
	if err := u.authAgent.IncreaseUserVersion(context.Background(), userID); err != nil {
		return fmt.Errorf("(usecase) failed to update user version: %w", err)
	}

	return nil
}

func (u *Usecase) ChangePassword(ctx context.Context, userID uint32, password string) error {
	if err := u.authAgent.ChangePassword(context.Background(), userID, password); err != nil {
		return fmt.Errorf("(usecase) failed to cahnge password: %w", err)
	}

	return nil
}
