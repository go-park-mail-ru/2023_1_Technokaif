package token


//go:generate mockgen -source=token.go -destination=mocks/mock.go
type Usecase interface {
	GenerateAccessToken(userID, userVersion uint32) (string, error)
	CheckAccessToken(acessToken string) (uint32, uint32, error)
	GenerateCSRFToken(userID uint32) (string, error)
	CheckCSRFToken(csrfToken string) (uint32, error)
}