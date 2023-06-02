package usecase

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Usecase struct{}

func NewUsecase() *Usecase {
	return &Usecase{}
}

var tokenSecret = os.Getenv("SECRET")

const csrfTokenTTL = 30 * time.Minute
const accessTokenTTL = 24 * 30 * time.Hour
const jwtParsingMaxTime = 3 * time.Second

type jwtAccessClaims struct {
	UserId      uint32 `json:"id"`
	UserVersion uint32 `json:"user_version"`
	jwt.RegisteredClaims
}

type jwtCSRFClaims struct {
	UserId uint32 `json:"id"`
	jwt.RegisteredClaims
}

func (u *Usecase) GenerateAccessToken(userID, userVersion uint32) (string, error) {
	claims := &jwtAccessClaims{
		userID,
		userVersion,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(accessTokenTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		},
	}

	signedToken, err := signToken(claims, tokenSecret)
	if err != nil {
		return "", fmt.Errorf("(usecase) failed to sign acess token: %w", err)
	}

	return signedToken, nil
}

func (u *Usecase) CheckAccessToken(acessToken string) (uint32, uint32, error) {
	token, err := checkToken(acessToken, tokenSecret, &jwtAccessClaims{})
	if err != nil {
		return 0, 0, fmt.Errorf("(usecase) invalid access token: %w", err)
	}

	claims, ok := token.Claims.(*jwtAccessClaims)
	if !ok {
		return 0, 0, errors.New("(usecase) token claims are not of type *tokenClaims")
	}

	now := time.Now().UTC()
	if claims.ExpiresAt.Time.Before(now) {
		return 0, 0, errors.New("(usecase) access token is expired")
	}

	return claims.UserId, claims.UserVersion, nil
}

func (u *Usecase) GenerateCSRFToken(userID uint32) (string, error) {
	claims := &jwtCSRFClaims{
		userID,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(csrfTokenTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		},
	}

	signedToken, err := signToken(claims, tokenSecret)
	if err != nil {
		return "", fmt.Errorf("(usecase) failed to sign acess token: %w", err)
	}

	return signedToken, nil
}

func (u *Usecase) CheckCSRFToken(acessToken string) (uint32, error) {
	token, err := checkToken(acessToken, tokenSecret, &jwtCSRFClaims{})
	if err != nil {
		return 0, fmt.Errorf("(usecase) invalid CSRF token")
	}

	claims, ok := token.Claims.(*jwtCSRFClaims)
	if !ok {
		return 0, errors.New("(usecase) token claims are not of type *tokenClaims")
	}

	now := time.Now().UTC()
	if claims.ExpiresAt.Time.Before(now) {
		return 0, errors.New("(usecase) csrf token is expired")
	}

	return claims.UserId, nil
}

func signToken(claims jwt.Claims, secret string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func checkToken(tokenStr string, secret string, claims jwt.Claims) (*jwt.Token, error) {
	token, err := jwt.ParseWithClaims(tokenStr, claims,
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("invalid token signing method")
			}
			return []byte(secret), nil
		}, jwt.WithLeeway(jwtParsingMaxTime))
	if err != nil {
		return nil, err
	}

	return token, nil
}
