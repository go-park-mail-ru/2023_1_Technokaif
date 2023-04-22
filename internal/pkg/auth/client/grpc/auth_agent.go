package grpc

import (
	"context"

	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/microservices/auth/proto"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/models"
)

type authAgent struct {
	client proto.AuthorizationClient
}

func NewAuthAgent(c proto.AuthorizationClient) *authAgent {
	return &authAgent{
		client: c,
	}
}

func (a *authAgent) SignUpUser(ctx context.Context, u models.User) (uint32, error) {
	msg := &proto.SignUpMsg{
		Username:  u.Username,
		Email:     u.Email,
		Password:  u.Password,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Sex:       string(u.Sex),
		BirthDate: u.BirthDate.Format("2006-01-02"), // TODO Format
	}
	resp, err := a.client.SignUpUser(ctx, msg)

	return resp.Id, err
}

func (a *authAgent) CheckCredentials() {

}
