package middleware

/* import (
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
	logMocks "github.com/go-park-mail-ru/2023_1_Technokaif/pkg/logger/mocks"
)

func TestDelivery_authorization(t *testing.T) { // TODO maybe without h.getUserFromAuthorization
	type mockBehavior func(a *authMocks.MockUsecase, t *tokenMocks.MockUsecase, token string, user models.User)

	testTable := []struct {
		name                     string
		cookieName               string
		cookieValue              string
		// token                    string
		mockBehavior             mockBehavior
		expectingError           bool
		expectedUser             models.User
		expectingHttpResponse500 bool
	}{
		{
			name:        "Ok",
			cookieName:  commonHttp.AcessTokenCookieName,
			cookieValue: "token",
			mockBehavior: func(a *authMocks.MockUsecase, t *tokenMocks.MockUsecase, token string, user models.User) {
				t.EXPECT().CheckAccessToken(token).Return(user.ID, user.Version, nil)
				a.EXPECT().GetUserByAuthData(user.ID, user.Version).Return(&user, nil)
			},
			expectingError:           false,
			expectedUser:             models.User{ID: 1, Version: 2},
			expectingHttpResponse500: false,
		},
		{
			name:        "Wrong Cookie Name",
			cookieName:  "Wrong-Access-Token",
			cookieValue: "token",
			mockBehavior: func(a *authMocks.MockUsecase, t *tokenMocks.MockUsecase, token string, user models.User) {},
			expectingError:           true,
			expectingHttpResponse500: false,
		},
		{
			name:        "Incorrect token sign",
			cookieName:  commonHttp.AcessTokenCookieName,
			cookieValue: "token",
			mockBehavior: func(a *authMocks.MockUsecase, t *tokenMocks.MockUsecase, token string, user models.User) {
				t.EXPECT().CheckAccessToken(token).Return(uint32(0), uint32(0), fmt.Errorf(""))
			},
			expectingError:           true,
			expectingHttpResponse500: false,
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
			expectingError:           true,
			expectingHttpResponse500: false,
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
			expectingError:           false,
			expectingHttpResponse500: true,
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			authMockUsecase := authMocks.NewMockUsecase(c)
			tokenMockUsecase := tokenMocks.NewMockUsecase(c)

			tc.mockBehavior(authMockUsecase, tokenMockUsecase, tc.cookieValue, tc.expectedUser)

			l := logMocks.NewMockLogger(c)
			l.EXPECT().Error(gomock.Any()).AnyTimes()
			l.EXPECT().Info(gomock.Any()).AnyTimes()
			l.EXPECT().Errorf(gomock.Any(), gomock.Any()).AnyTimes()
			l.EXPECT().Infof(gomock.Any(), gomock.Any()).AnyTimes()

			h := NewMiddleware(authMockUsecase, tokenMockUsecase, l)

			r := chi.NewRouter()
			r.With(h.Authorization).Get("/auth", func(w http.ResponseWriter, r *http.Request) {
				u, err := commonHttp.GetUserFromRequest(r)

				// Asserts
				if tc.expectingError {
					assert.Error(t, err)
				} else {
					got := u
					expected := &tc.expectedUser

					assert.NoError(t, err)
					assert.Equal(t, got, expected)
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

			if tc.expectingHttpResponse500 {
				assert.Equal(t, 500, w.Code)
			} else {
				assert.Equal(t, 200, w.Code)
			}
		})
	}
}
*/