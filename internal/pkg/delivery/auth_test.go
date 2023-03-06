package delivery

import (
	"bytes"
	"context"
	"fmt"
	"net/http/httptest"
	"net/http"
	"testing"
	"time"

	"github.com/go-chi/chi"
	logMocks "github.com/go-park-mail-ru/2023_1_Technokaif/internal/logger/mocks"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/models"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/usecase"
	mocks "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/usecase/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestDelivery_signUp(t *testing.T) {
	type mockBehavior func(a *mocks.MockAuth, u models.User)

	correctTestRequestBody := `{"username": "yarik_tri", "password": "Love1234",
	"email": "yarik1448kuzmin@gmail.com", "firstName": "Yaroslav", "lastName": "Kuzmin",
	"birthDate": "2003-08-23", "sex": "M"}`

	birthTime, err := time.Parse(time.RFC3339, "2003-08-23T00:00:00Z")
	if err != nil {
		t.Errorf("can't Parse birth date: %v", err)
	}
	birthDate := models.Date{birthTime}

	correctTestUser := models.User{
		Username:  "yarik_tri",
		Password:  "Love1234",
		Email:     "yarik1448kuzmin@gmail.com",
		FirstName: "Yaroslav",
		LastName:  "Kuzmin",
		BirhDate:  birthDate,
		Sex:       models.Male,
	}

	testTable := []struct {
		name             string
		requestBody      string
		userFromBody     models.User
		mockBehavior     mockBehavior
		expectedStatus   int
		expectedResponse string
	}{
		{
			name:         "Common",
			requestBody:  correctTestRequestBody,
			userFromBody: correctTestUser,
			mockBehavior: func(a *mocks.MockAuth, u models.User) {
				a.EXPECT().CreateUser(u).Return(1, nil)
			},
			expectedStatus:   200,
			expectedResponse: `{"id": 1}`,
		},
		{
			// Missing one quote (after firstName)
			name: "Incorrect request body",
			requestBody: `{"username": "yarik_tri", "password": "Love1234",
			"email": "yarik1448kuzmin@gmail.com", "firstName: "Yaroslav, "lastName": "Kuzmin",
			"birthDate": "2003-08-23", "sex": "M"}`,
			userFromBody:     correctTestUser,
			mockBehavior:     func(a *mocks.MockAuth, u models.User) {},
			expectedStatus:   400,
			expectedResponse: `{"message": "incorrect input body"}`,
		},
		{
			// These tests aren't tests of validation but delivery-layer
			// So, check only one example of validation error (short username)
			// to convince the error is caught
			name: "Validation Error",
			requestBody: `{"username": "yar", "password": "Love1234",
			"email": "yarik1448kuzmin@gmail.com", "firstName": "Yaroslav", "lastName": "Kuzmin",
			"birthDate": "2003-08-23", "sex": "M"}`,
			userFromBody:     models.User{},
			mockBehavior:     func(a *mocks.MockAuth, u models.User) {},
			expectedStatus:   400,
			expectedResponse: `{"message": "incorrect input body"}`,
		},
		{
			name:         "Creating Error",
			requestBody:  correctTestRequestBody,
			userFromBody: correctTestUser,
			mockBehavior: func(a *mocks.MockAuth, u models.User) {
				a.EXPECT().CreateUser(u).Return(0, fmt.Errorf("user already exists"))
			},
			expectedStatus:   400,
			expectedResponse: `{"message": "user already exists"}`,
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			// Init
			c := gomock.NewController(t)
			defer c.Finish()

			a := mocks.NewMockAuth(c)

			tc.mockBehavior(a, tc.userFromBody)

			u := &usecase.Usecase{
				Auth: a,
			}

			l := logMocks.NewMockLogger(c)
			l.EXPECT().Error(gomock.Any()).AnyTimes()
			l.EXPECT().Info(gomock.Any()).AnyTimes()

			h := NewHandler(u, l)

			// Routing
			r := chi.NewRouter()
			r.Post("/signup", h.signUp)

			// Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/signup", bytes.NewBufferString(tc.requestBody))
			r.ServeHTTP(w, req)

			// Test
			assert.Equal(t, tc.expectedStatus, w.Code)
			assert.JSONEq(t, tc.expectedResponse, w.Body.String())
		})
	}
}

func TestDelivery_login(t *testing.T) {
	type mockBehavior func(a *mocks.MockAuth, l loginInput)

	correctTestRequestBody := `{"username": "yarik_tri", "password": "Love1234"}`
	correctTestLogin := loginInput{
		Username: "yarik_tri",
		Password: "Love1234",
	}

	testTable := []struct {
		name             string
		requestBody      string
		loginFromBody    loginInput
		mockBehavior     mockBehavior
		expectedStatus   int
		expectedResponse string
	}{
		{
			name:          "Common",
			requestBody:   correctTestRequestBody,
			loginFromBody: correctTestLogin,
			mockBehavior: func(a *mocks.MockAuth, l loginInput) {
				a.EXPECT().LoginUser(l.Username, l.Password).Return("jwt.access.token", nil)
			},
			expectedStatus:   200,
			expectedResponse: `{"jwt": "jwt.access.token"}`,
		},
		{
			// Missing one quote (after username)
			name:             "Incorrect Request Body",
			requestBody:      `{"username": "yarik_tri, "password": "Love1234"}`,
			loginFromBody:    correctTestLogin,
			mockBehavior:     func(a *mocks.MockAuth, l loginInput) {},
			expectedStatus:   400,
			expectedResponse: `{"message": "incorrect input body"}`,
		},
		{
			// These tests aren't tests of validation but delivery-layer
			// So, check only one example of validation error (no password)
			// to convince the error is caught
			name:             "Validation Error",
			requestBody:      `{"username": "yarik_tri"}`,
			loginFromBody:    loginInput{},
			mockBehavior:     func(a *mocks.MockAuth, l loginInput) {},
			expectedStatus:   400,
			expectedResponse: `{"message": "incorrect input body"}`,
		},
		{
			name:          "Login Error",
			requestBody:   correctTestRequestBody,
			loginFromBody: correctTestLogin,
			mockBehavior: func(a *mocks.MockAuth, l loginInput) {
				a.EXPECT().LoginUser(correctTestLogin.Username,
					correctTestLogin.Password).Return("", fmt.Errorf(""))
			},
			expectedStatus:   400,
			expectedResponse: `{"message": "can't login user"}`,
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			// Init
			c := gomock.NewController(t)
			defer c.Finish()

			a := mocks.NewMockAuth(c)

			tc.mockBehavior(a, tc.loginFromBody)

			u := &usecase.Usecase{
				Auth: a,
			}

			l := logMocks.NewMockLogger(c)
			l.EXPECT().Error(gomock.Any()).AnyTimes()
			l.EXPECT().Info(gomock.Any()).AnyTimes()

			h := NewHandler(u, l)

			// Routing
			r := chi.NewRouter()
			r.Post("/login", h.login)

			// Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/login", bytes.NewBufferString(tc.requestBody))
			r.ServeHTTP(w, req)

			assert.Equal(t, tc.expectedStatus, w.Code)
			assert.JSONEq(t, tc.expectedResponse, w.Body.String())
		})
	}
}

func TestDelivery_logout(t *testing.T) {  // TODO maybe mock GetUserFromAuthorization
	correctTestUser := &models.User{
		ID: 1,
	}
	wrapRequestWithUser := func(r *http.Request, user *models.User) *http.Request {
		if user == nil {
			return r
		}
		ctx := context.WithValue(r.Context(), contextValueUser, user)
		return r.WithContext(ctx)
	}

	type mockBehavior func(a *mocks.MockAuth, user *models.User)

	testTable := []struct {
		name             	string
		user	 			*models.User
		mockBehavior     	mockBehavior
		expectedStatus   	int
		expectedResponse 	string
	}{
		{
		 	name:             "Common",
			user:			  correctTestUser,	
		 	mockBehavior:     func(a *mocks.MockAuth, user *models.User) {
		 		a.EXPECT().IncreaseUserVersion(user.ID).Return(nil)
		 	},
		 	expectedStatus:   200,
		 	expectedResponse: `{"status": "ok"}`,
		},
		{
			name:             "No user in request",
		   	user:			  nil,	
			mockBehavior:     func(a *mocks.MockAuth, user *models.User) {},
			expectedStatus:   400,
			expectedResponse: `{"message": "invalid token"}`,
	   	},
		{
			name:             "Failed to increase user version",
		   	user:			  correctTestUser,	
			mockBehavior:     func(a *mocks.MockAuth, user *models.User) {
				a.EXPECT().IncreaseUserVersion(user.ID).Return(fmt.Errorf(""))
			},
			expectedStatus:   400,
			expectedResponse: `{"message": "failed to log out"}`,
	   	},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			// Init
			c := gomock.NewController(t)
			defer c.Finish()

			a := mocks.NewMockAuth(c)

			tc.mockBehavior(a, tc.user)

			u := &usecase.Usecase{
				Auth: a,
			}

			l := logMocks.NewMockLogger(c)
			l.EXPECT().Error(gomock.Any()).AnyTimes()
			l.EXPECT().Info(gomock.Any()).AnyTimes()

			h := NewHandler(u, l)

			// Routing
			r := chi.NewRouter()
			r.Get("/logout", h.logout)

			// Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/logout", nil)

			r.ServeHTTP(w, wrapRequestWithUser(req, tc.user))

			// Test
			assert.Equal(t, tc.expectedStatus, w.Code)
			assert.JSONEq(t, tc.expectedResponse, w.Body.String())
		})
	}
}
