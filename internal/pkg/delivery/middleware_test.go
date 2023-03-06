package delivery

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"math/rand"

	"github.com/go-chi/chi"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/models"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/usecase"

	logMocks "github.com/go-park-mail-ru/2023_1_Technokaif/internal/logger/mocks"
	mocks "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/usecase/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestDelivery_authorization(t *testing.T) {  // TODO maybe without h.getUserFromAuthorization
	type mockBehavior func(r *mocks.MockAuth, token string, user models.User)

	testTable := []struct {
		name                 string
		headerName           string
		headerValue          string
		token                string				 
		mockBehavior         mockBehavior
		expectingError   	 bool
		expectedUser		 models.User
	}{
		{
			name:         	"Ok",
			headerName:   	"Authorization",
			headerValue:  	"Bearer token",
			token:        	"token",
			mockBehavior: 	func(r *mocks.MockAuth, token string, user models.User) {
				r.EXPECT().CheckAccessToken(token).Return(user.ID, user.Version, nil)
				r.EXPECT().GetUserByAuthData(user.ID, user.Version).Return(&user, nil)
			},
			expectingError: false,
			expectedUser:	models.User{ID:1, Version: 2},
		},
		{
			name:         	"Missing Bearer",
			headerName:   	"Authorization",
			headerValue:  	"token",
			token:        	"token",
			mockBehavior: 	func(r *mocks.MockAuth, token string, user models.User) {},
			expectingError: true,
			expectedUser:	models.User{},
		},
		{
			name:         	"Missing token",
			headerName:   	"Authorization",
			headerValue:  	"Bearer",
			token:        	"",
			mockBehavior: 	func(r *mocks.MockAuth, token string, user models.User) {},
			expectingError: true,
			expectedUser:	models.User{},
		},
		{
			name:         	"Missing token with space",
			headerName:   	"Authorization",
			headerValue:  	"Bearer  ",
			token:        	"",
			mockBehavior: 	func(r *mocks.MockAuth, token string, user models.User) {},
			expectingError: true,
			expectedUser:	models.User{},
		},
		{
			name:         	"Incorrect token sign",
			headerName:   	"Authorization",
			headerValue:  	"Bearer token",
			token:        	"token",
			mockBehavior: 	func(r *mocks.MockAuth, token string, user models.User) {
				r.EXPECT().CheckAccessToken(token).Return(uint(0), uint(0), fmt.Errorf(""))
			},
			expectingError: true,
			expectedUser:	models.User{},
		},
		{
			name:         	"Auth failed",
			headerName:   	"Authorization",
			headerValue:  	"Bearer token",
			token:        	"token",
			mockBehavior: 	func(r *mocks.MockAuth, token string, user models.User) {
				randVal := uint(rand.Intn(100))

				r.EXPECT().CheckAccessToken(token).Return(randVal, randVal, nil)
				r.EXPECT().GetUserByAuthData(randVal, randVal).Return(&user, fmt.Errorf(""))
			},
			expectingError: true,
			expectedUser:	models.User{},
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			a := mocks.NewMockAuth(c)

			tc.mockBehavior(a, tc.token, tc.expectedUser)

			u := &usecase.Usecase{
				Auth: a,
			}

			l := logMocks.NewMockLogger(c)
			l.EXPECT().Error(gomock.Any()).AnyTimes()
			l.EXPECT().Info(gomock.Any()).AnyTimes()

			h := NewHandler(u, l)
			r := chi.NewRouter()
			r.With(h.authorization).Get("/auth", func(w http.ResponseWriter, r *http.Request) {
				u, err := h.getUserFromAuthorization(r)
				
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
			req.Header.Set(tc.headerName, tc.headerValue)
			r.ServeHTTP(w, req)
		})
	}
}