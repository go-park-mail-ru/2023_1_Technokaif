package http

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/models"
	authMocks "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/auth/mocks"
	logMocks "github.com/go-park-mail-ru/2023_1_Technokaif/pkg/logger/mocks"
)

func TestDeliverySignUp(t *testing.T) {
	type mockBehavior func(a *authMocks.MockUsecase, u models.User)

	correctTestRequestBody := `{"username": "yarik_tri", "password": "Love1234",
	"email": "yarik1448kuzmin@gmail.com", "firstName": "Yaroslav", "lastName": "Kuzmin",
	"birthDate": "2003-08-23", "sex": "M"}`

	birthTime, err := time.Parse(time.RFC3339, "2003-08-23T00:00:00Z")
	if err != nil {
		t.Errorf("can't Parse birth date: %v", err)
	}
	birthDate := models.Date{Time: birthTime}

	correctTestUser := models.User{
		Username:  "yarik_tri",
		Password:  "Love1234",
		Email:     "yarik1448kuzmin@gmail.com",
		FirstName: "Yaroslav",
		LastName:  "Kuzmin",
		BirthDate: birthDate,
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
			mockBehavior: func(a *authMocks.MockUsecase, u models.User) {
				a.EXPECT().SignUpUser(u).Return(uint32(1), nil)
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
			mockBehavior:     func(a *authMocks.MockUsecase, u models.User) {},
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
			mockBehavior:     func(a *authMocks.MockUsecase, u models.User) {},
			expectedStatus:   400,
			expectedResponse: `{"message": "incorrect input body"}`,
		},
		{
			name:         "Creating existing user Error",
			requestBody:  correctTestRequestBody,
			userFromBody: correctTestUser,
			mockBehavior: func(a *authMocks.MockUsecase, u models.User) {
				a.EXPECT().SignUpUser(u).Return(uint32(0), &models.UserAlreadyExistsError{})
			},
			expectedStatus:   400,
			expectedResponse: `{"message": "user already exists"}`,
		},
		{
			name:         "Creating database Error",
			requestBody:  correctTestRequestBody,
			userFromBody: correctTestUser,
			mockBehavior: func(a *authMocks.MockUsecase, u models.User) {
				a.EXPECT().SignUpUser(u).Return(uint32(0), fmt.Errorf("database query error"))
			},
			expectedStatus:   500,
			expectedResponse: `{"message": "server failed to sign up user"}`,
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			// Init
			c := gomock.NewController(t)
			defer c.Finish()

			authMockUsecase := authMocks.NewMockUsecase(c)

			tc.mockBehavior(authMockUsecase, tc.userFromBody)

			l := logMocks.NewMockLogger(c)
			l.EXPECT().Error(gomock.Any()).AnyTimes()
			l.EXPECT().Info(gomock.Any()).AnyTimes()
			l.EXPECT().Errorf(gomock.Any(), gomock.Any()).AnyTimes()
			l.EXPECT().Infof(gomock.Any(), gomock.Any()).AnyTimes()

			h := NewHandler(authMockUsecase, l)

			// Routing
			r := chi.NewRouter()
			r.Post("/signup", h.SignUp)

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

func TestDeliveryLogin(t *testing.T) {
	type mockBehavior func(a *authMocks.MockUsecase, l loginInput)

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
			mockBehavior: func(a *authMocks.MockUsecase, l loginInput) {
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
			mockBehavior:     func(a *authMocks.MockUsecase, l loginInput) {},
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
			mockBehavior:     func(a *authMocks.MockUsecase, l loginInput) {},
			expectedStatus:   400,
			expectedResponse: `{"message": "incorrect input body"}`,
		},
		{
			name:          "Login Error",
			requestBody:   correctTestRequestBody,
			loginFromBody: correctTestLogin,
			mockBehavior: func(a *authMocks.MockUsecase, l loginInput) {
				a.EXPECT().LoginUser(correctTestLogin.Username,
					correctTestLogin.Password).Return("", &models.NoSuchUserError{})
			},
			expectedStatus:   400,
			expectedResponse: `{"message": "can't login user"}`,
		},
		{
			name:          "Server Error",
			requestBody:   correctTestRequestBody,
			loginFromBody: correctTestLogin,
			mockBehavior: func(a *authMocks.MockUsecase, l loginInput) {
				a.EXPECT().LoginUser(correctTestLogin.Username,
					correctTestLogin.Password).Return("", fmt.Errorf("database error"))
			},
			expectedStatus:   500,
			expectedResponse: `{"message": "server failed to login user"}`,
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			// Init
			c := gomock.NewController(t)
			defer c.Finish()

			authMockUsecase := authMocks.NewMockUsecase(c)

			tc.mockBehavior(authMockUsecase, tc.loginFromBody)

			l := logMocks.NewMockLogger(c)
			l.EXPECT().Error(gomock.Any()).AnyTimes()
			l.EXPECT().Info(gomock.Any()).AnyTimes()
			l.EXPECT().Errorf(gomock.Any(), gomock.Any()).AnyTimes()
			l.EXPECT().Infof(gomock.Any(), gomock.Any()).AnyTimes()

			h := NewHandler(authMockUsecase, l)

			// Routing
			r := chi.NewRouter()
			r.Post("/login", h.Login)

			// Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/login", bytes.NewBufferString(tc.requestBody))
			r.ServeHTTP(w, req)

			assert.Equal(t, tc.expectedStatus, w.Code)
			assert.JSONEq(t, tc.expectedResponse, w.Body.String())
		})
	}
}

func TestDeliveryLogout(t *testing.T) {
	correctTestUser := &models.User{
		ID:      1,
		Version: 1,
	}
	testWrapRequestWithUser := func(r *http.Request, user *models.User, doWrap bool) *http.Request {
		if !doWrap {
			return r
		}
		ctx := context.WithValue(r.Context(), models.ContextKeyUserType{}, user)
		return r.WithContext(ctx)
	}

	type mockBehavior func(a *authMocks.MockUsecase, user *models.User)

	testTable := []struct {
		name             string
		user             *models.User
		mockBehavior     mockBehavior
		expectedStatus   int
		expectedResponse string
		doWrap           bool
	}{
		{
			name: "Common",
			user: correctTestUser,
			mockBehavior: func(a *authMocks.MockUsecase, user *models.User) {
				a.EXPECT().IncreaseUserVersion(user.ID).Return(nil)
			},
			expectedStatus:   200,
			expectedResponse: `{"status": "ok"}`,
			doWrap:           true,
		},
		{
			name:             "No user in request",
			user:             nil,
			mockBehavior:     func(a *authMocks.MockUsecase, user *models.User) {},
			expectedStatus:   400,
			expectedResponse: `{"message": "invalid token"}`,
			doWrap:           false,
		},
		{
			name:             "Nil user in request",
			user:             nil,
			mockBehavior:     func(a *authMocks.MockUsecase, user *models.User) {},
			expectedStatus:   400,
			expectedResponse: `{"message": "invalid token"}`,
			doWrap:           true,
		},
		{
			name: "Failed to increase user version",
			user: correctTestUser,
			mockBehavior: func(a *authMocks.MockUsecase, user *models.User) {
				a.EXPECT().IncreaseUserVersion(user.ID).Return(fmt.Errorf("database error"))
			},
			expectedStatus:   500,
			expectedResponse: `{"message": "failed to log out"}`,
			doWrap:           true,
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			// Init
			c := gomock.NewController(t)
			defer c.Finish()

			authMockUsecase := authMocks.NewMockUsecase(c)

			tc.mockBehavior(authMockUsecase, tc.user)

			l := logMocks.NewMockLogger(c)
			l.EXPECT().Error(gomock.Any()).AnyTimes()
			l.EXPECT().Info(gomock.Any()).AnyTimes()
			l.EXPECT().Errorf(gomock.Any(), gomock.Any()).AnyTimes()
			l.EXPECT().Infof(gomock.Any(), gomock.Any()).AnyTimes()

			h := NewHandler(authMockUsecase, l)

			// Routing
			r := chi.NewRouter()
			r.Get("/logout", h.Logout)

			// Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/logout", nil)

			r.ServeHTTP(w, testWrapRequestWithUser(req, tc.user, tc.doWrap))

			// Test
			assert.Equal(t, tc.expectedStatus, w.Code)
			assert.JSONEq(t, tc.expectedResponse, w.Body.String())
		})
	}
}
