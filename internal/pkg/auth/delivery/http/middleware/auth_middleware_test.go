package middleware

import (
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	commonHttp "github.com/go-park-mail-ru/2023_1_Technokaif/internal/common/http"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/models"
	authMocks "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/auth/mocks"
	tokenMocks "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/token/mocks"
	commonTests "github.com/go-park-mail-ru/2023_1_Technokaif/internal/common/tests"
)

func TestAuthDeliveryAuthorization(t *testing.T) {
	type mockBehavior func(a *authMocks.MockUsecase, t *tokenMocks.MockUsecase, token string, user models.User)

	testTable := []struct {
		name                   string
		cookieName             string
		cookieValue            string
		token                  string
		mockBehavior           mockBehavior
		expectingUserInContext bool
		expectedUser           models.User
		expectingResponse      bool
		expectedStatus         int
		expectedResponse       string
	}{
		{
			name:        "Ok",
			cookieName:  commonHttp.AcessTokenCookieName,
			cookieValue: "token",
			mockBehavior: func(a *authMocks.MockUsecase, t *tokenMocks.MockUsecase, token string, user models.User) {
				t.EXPECT().CheckAccessToken(token).Return(user.ID, user.Version, nil)
				a.EXPECT().GetUserByAuthData(user.ID, user.Version).Return(&user, nil)
			},
			expectingUserInContext: true,
			expectedUser:           models.User{ID: uint32(rand.Intn(100)), Version: uint32(rand.Intn(100))},
			expectingResponse:      false,
		},
		{
			name:                   "Wrong cookie name",
			cookieName:             "Wrong-Access-Token",
			cookieValue:            "token",
			mockBehavior:           func(a *authMocks.MockUsecase, t *tokenMocks.MockUsecase, token string, user models.User) {},
			expectingUserInContext: false,
			expectingResponse:      false,
		},
		{
			name:                   "Empty cookies",
			cookieName:             commonHttp.AcessTokenCookieName,
			cookieValue:            "",
			mockBehavior:           func(a *authMocks.MockUsecase, t *tokenMocks.MockUsecase, token string, user models.User) {},
			expectingUserInContext: false,
			expectingResponse:      false,
		},
		{
			name:        "Incorrect token sign",
			cookieName:  commonHttp.AcessTokenCookieName,
			cookieValue: "token",
			mockBehavior: func(a *authMocks.MockUsecase, t *tokenMocks.MockUsecase, token string, user models.User) {
				t.EXPECT().CheckAccessToken(token).Return(uint32(0), uint32(0), fmt.Errorf(""))
			},
			expectingUserInContext: false,
			expectingResponse:      true,
			expectedStatus:         400,
			expectedResponse:       `{"message": "token check failed"}`,
		},
		{
			name:        "Auth failed",
			cookieName:  commonHttp.AcessTokenCookieName,
			cookieValue: "token",
			mockBehavior: func(a *authMocks.MockUsecase, t *tokenMocks.MockUsecase, token string, user models.User) {
				randVal := uint32(rand.Intn(100))

				t.EXPECT().CheckAccessToken(token).Return(randVal, randVal, nil)
				a.EXPECT().GetUserByAuthData(randVal, randVal).Return(&user, &models.NoSuchUserError{})
			},
			expectingUserInContext: false,
			expectingResponse:      true,
			expectedStatus:         400,
			expectedResponse:       `{"message": "auth data check failed"}`,
		},
		{
			name:        "Server error",
			cookieName:  commonHttp.AcessTokenCookieName,
			cookieValue: "token",
			mockBehavior: func(a *authMocks.MockUsecase, t *tokenMocks.MockUsecase, token string, user models.User) {
				randVal := uint32(rand.Intn(100))

				t.EXPECT().CheckAccessToken(token).Return(randVal, randVal, nil)
				a.EXPECT().GetUserByAuthData(randVal, randVal).Return(&user, errors.New("server error"))
			},
			expectingUserInContext: false,
			expectingResponse:      true,
			expectedStatus:         500,
			expectedResponse:       `{"message": "server failed to check authorization"}`,
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			authMockUsecase := authMocks.NewMockUsecase(c)
			tokenMockUsecase := tokenMocks.NewMockUsecase(c)

			tc.mockBehavior(authMockUsecase, tokenMockUsecase, tc.cookieValue, tc.expectedUser)

			l := commonTests.MockLogger(c)

			h := NewMiddleware(authMockUsecase, tokenMockUsecase, l)

			r := chi.NewRouter()
			r.With(h.Authorization).Get("/auth", func(w http.ResponseWriter, r *http.Request) {
				u, err := commonHttp.GetUserFromRequest(r)

				// Asserts
				if tc.expectingUserInContext {
					got := u
					expected := &tc.expectedUser

					assert.NoError(t, err)
					assert.Equal(t, got, expected)
				} else {
					assert.Error(t, err)
				}
			})

			// Init Test Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/auth", nil)
			req.AddCookie(&http.Cookie{
				Name:  tc.cookieName,
				Value: tc.cookieValue,
			})
			r.ServeHTTP(w, req)

			if tc.expectingResponse {
				assert.Equal(t, tc.expectedStatus, w.Code)
				assert.JSONEq(t, tc.expectedResponse, w.Body.String())
			} else {
				assert.Equal(t, 200, w.Code)
			}
		})
	}
}
