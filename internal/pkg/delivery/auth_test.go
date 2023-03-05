package delivery

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/logger"
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
			"email": "yarik1448kuzmin@gmail.com", "firstName: "Yaroslav", "lastName": "Kuzmin",
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
				a.EXPECT().CreateUser(u).Return(0, fmt.Errorf(""))
			},
			expectedStatus:   400,
			expectedResponse: `{"message": "incorrect input body"}`,
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
			l, _ := logger.NewFLogger()
			h := NewHandler(u, l)

			// Routing
			r := chi.NewRouter()
			r.Post("/signup", h.signUp)

			// Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/signup", bytes.NewBufferString(tc.requestBody))
			r.ServeHTTP(w, req)

			var expResp signUpResponse
			var actualResp signUpResponse
			json.Unmarshal([]byte(tc.expectedResponse), &expResp)
			json.Unmarshal(w.Body.Bytes(), &actualResp)

			// Test
			assert.Equal(t, tc.expectedStatus, w.Code)
			assert.Equal(t, expResp, actualResp)
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
			l, _ := logger.NewFLogger()
			h := NewHandler(u, l)

			// Routing
			r := chi.NewRouter()
			r.Post("/login", h.login)

			// Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/login", bytes.NewBufferString(tc.requestBody))
			r.ServeHTTP(w, req)

			var expResp loginResponse
			var actualResp loginResponse
			json.Unmarshal([]byte(tc.expectedResponse), &expResp)
			json.Unmarshal(w.Body.Bytes(), &actualResp)

			// Test
			assert.Equal(t, tc.expectedStatus, w.Code)
			assert.Equal(t, expResp, actualResp)
		})
	}
}

func TestDelivery_logout(t *testing.T) {
	type mockBehavior func(a *mocks.MockAuth)

	testTable := []struct {
		name             string
		mockBehavior     mockBehavior
		expectedStatus   int
		expectedResponse string
	}{
		// {
		// 	name:             "Common",
		// 	mockBehavior:     func(a *mocks.MockAuth) {
		// 		a.EXPECT().ChangeUserVersion()
		// 	},
		// 	expectedStatus:   200,
		// 	expectedResponse: `{"status": "ok"}`,
		// },
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			// Init
			c := gomock.NewController(t)
			defer c.Finish()

			a := mocks.NewMockAuth(c)

			tc.mockBehavior(a)

			u := &usecase.Usecase{
				Auth: a,
			}
			l, _ := logger.NewFLogger()
			h := NewHandler(u, l)

			// Routing
			r := chi.NewRouter()
			r.With(h.Authorization).Get("/logout", h.logout)

			// Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/logout", nil)
			r.ServeHTTP(w, req)

			var expResp logoutResponse
			var actualResp logoutResponse
			json.Unmarshal([]byte(tc.expectedResponse), &expResp)
			json.Unmarshal(w.Body.Bytes(), &actualResp)

			// Test
			assert.Equal(t, tc.expectedStatus, w.Code)
			assert.Equal(t, expResp, actualResp)
		})
	}
}
