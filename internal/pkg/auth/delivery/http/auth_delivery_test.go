package http

import (
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	commonHTTP "github.com/go-park-mail-ru/2023_1_Technokaif/internal/common/http"
	commonTests "github.com/go-park-mail-ru/2023_1_Technokaif/internal/common/tests"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/models"
	authMocks "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/auth/mocks"
	tokenMocks "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/token/mocks"
)

func TestDeliverySignUp(t *testing.T) {
	// Init
	type mockBehavior func(a *authMocks.MockUsecase, u models.User)

	c := gomock.NewController(t)

	authMockUsecase := authMocks.NewMockUsecase(c)
	tokenMockUsecase := tokenMocks.NewMockUsecase(c)

	l := commonTests.MockLogger(c)

	h := NewHandler(authMockUsecase, tokenMockUsecase, l)

	// Routing
	r := chi.NewRouter()
	r.Post("/signup", h.SignUp)

	// Test filling
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
				a.EXPECT().SignUpUser(gomock.Any(), u).Return(uint32(1), nil)
			},
			expectedStatus:   http.StatusOK,
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
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: commonTests.ErrorResponse(commonHTTP.IncorrectRequestBody),
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
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: commonTests.ErrorResponse(commonHTTP.IncorrectRequestBody),
		},
		{
			name:         "Creating existing user Error",
			requestBody:  correctTestRequestBody,
			userFromBody: correctTestUser,
			mockBehavior: func(a *authMocks.MockUsecase, u models.User) {
				a.EXPECT().SignUpUser(gomock.Any(), u).Return(uint32(0), &models.UserAlreadyExistsError{})
			},
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: commonTests.ErrorResponse(userAlreadyExists),
		},
		{
			name:         "Creating database Error",
			requestBody:  correctTestRequestBody,
			userFromBody: correctTestUser,
			mockBehavior: func(a *authMocks.MockUsecase, u models.User) {
				a.EXPECT().SignUpUser(gomock.Any(), u).Return(uint32(0), fmt.Errorf("database query error"))
			},
			expectedStatus:   http.StatusInternalServerError,
			expectedResponse: commonTests.ErrorResponse(userSignUpServerError),
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			// Call mock
			tc.mockBehavior(authMockUsecase, tc.userFromBody)

			commonTests.DeliveryTestPost(t, r, "/signup", tc.requestBody, tc.expectedStatus, tc.expectedResponse,
				commonTests.NoWrapUserFunc())
		})
	}
}

func TestDeliveryLogin(t *testing.T) {
	// Init
	type mockBehavior func(a *authMocks.MockUsecase, t *tokenMocks.MockUsecase, l loginInput)

	c := gomock.NewController(t)

	authMockUsecase := authMocks.NewMockUsecase(c)
	tokenMockUsecase := tokenMocks.NewMockUsecase(c)

	l := commonTests.MockLogger(c)

	h := NewHandler(authMockUsecase, tokenMockUsecase, l)

	// Routing
	r := chi.NewRouter()
	r.Post("/login", h.Login)

	// Test filling
	correctTestRequestBody := `{"username": "yarik_tri", "password": "Love1234"}`
	correctTestLogin := loginInput{
		Username: "yarik_tri",
		Password: "Love1234",
	}

	correctCookieName := commonHTTP.AccessTokenCookieName
	randomUserID := uint32(rand.Intn(100))

	testTable := []struct {
		name                string
		requestBody         string
		loginFromBody       loginInput
		mockBehavior        mockBehavior
		expectedStatus      int
		expectedResponse    string
		expectingCookie     bool
		expectedCookieValue string
	}{
		{
			name:          "Common",
			requestBody:   correctTestRequestBody,
			loginFromBody: correctTestLogin,
			mockBehavior: func(a *authMocks.MockUsecase, t *tokenMocks.MockUsecase, l loginInput) {
				user := &models.User{ID: randomUserID, Version: uint32(rand.Intn(100))}

				a.EXPECT().GetUserByCreds(gomock.Any(), l.Username, l.Password).Return(user, nil)
				t.EXPECT().GenerateAccessToken(user.ID, user.Version).Return("token", nil)
			},
			expectedStatus:      http.StatusOK,
			expectedResponse:    fmt.Sprintf(`{"id": %d}`, randomUserID),
			expectingCookie:     true,
			expectedCookieValue: "token",
		},
		{
			// Missing one quote (after username)
			name:             "Incorrect Request Body",
			requestBody:      `{"username": "yarik_tri, "password": "Love1234"}`,
			loginFromBody:    correctTestLogin,
			mockBehavior:     func(a *authMocks.MockUsecase, t *tokenMocks.MockUsecase, l loginInput) {},
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: commonTests.ErrorResponse(commonHTTP.IncorrectRequestBody),
			expectingCookie:  false,
		},
		{
			// These tests aren't tests of validation but delivery-layer
			// So, check only one example of validation error (no password)
			// to convince the error is caught
			name:             "Validation Error",
			requestBody:      `{"username": "yarik_tri"}`,
			loginFromBody:    loginInput{},
			mockBehavior:     func(a *authMocks.MockUsecase, t *tokenMocks.MockUsecase, l loginInput) {},
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: commonTests.ErrorResponse(commonHTTP.IncorrectRequestBody),
			expectingCookie:  false,
		},
		{
			name:          "No Such User",
			requestBody:   correctTestRequestBody,
			loginFromBody: correctTestLogin,
			mockBehavior: func(a *authMocks.MockUsecase, t *tokenMocks.MockUsecase, l loginInput) {
				a.EXPECT().GetUserByCreds(gomock.Any(), l.Username, l.Password).
					Return(&models.User{}, &models.NoSuchUserError{})
			},
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: commonTests.ErrorResponse(userNotFound),
			expectingCookie:  false,
		},
		{
			name:          "Incorrect Password",
			requestBody:   correctTestRequestBody,
			loginFromBody: correctTestLogin,
			mockBehavior: func(a *authMocks.MockUsecase, t *tokenMocks.MockUsecase, l loginInput) {
				a.EXPECT().GetUserByCreds(gomock.Any(), l.Username, l.Password).
					Return(&models.User{}, &models.IncorrectPasswordError{})
			},
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: commonTests.ErrorResponse(passwordMismatch),
			expectingCookie:  false,
		},
		{
			name:          "Getting User Server Error",
			requestBody:   correctTestRequestBody,
			loginFromBody: correctTestLogin,
			mockBehavior: func(a *authMocks.MockUsecase, t *tokenMocks.MockUsecase, l loginInput) {
				a.EXPECT().GetUserByCreds(gomock.Any(), l.Username, l.Password).
					Return(&models.User{}, errors.New("database error"))
			},
			expectedStatus:   http.StatusInternalServerError,
			expectedResponse: commonTests.ErrorResponse(userLoginServerError),
			expectingCookie:  false,
		},
		{
			name:          "Generating token Server Error",
			requestBody:   correctTestRequestBody,
			loginFromBody: correctTestLogin,
			mockBehavior: func(a *authMocks.MockUsecase, t *tokenMocks.MockUsecase, l loginInput) {
				user := &models.User{ID: uint32(rand.Intn(100)), Version: uint32(rand.Intn(100))}

				a.EXPECT().GetUserByCreds(gomock.Any(), l.Username, l.Password).Return(user, nil)
				t.EXPECT().GenerateAccessToken(user.ID, user.Version).
					Return("", errors.New("generating token error"))
			},
			expectedStatus:   http.StatusInternalServerError,
			expectedResponse: commonTests.ErrorResponse(userLoginServerError),
			expectingCookie:  false,
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			// Call mock
			tc.mockBehavior(authMockUsecase, tokenMockUsecase, tc.loginFromBody)

			w := commonTests.DeliveryTestPost(t, r, "/login", tc.requestBody, tc.expectedStatus, tc.expectedResponse,
				commonTests.NoWrapUserFunc())

			if tc.expectingCookie {
				assert.Equal(t, correctCookieName, w.Result().Cookies()[0].Name)
				assert.Equal(t, tc.expectedCookieValue, w.Result().Cookies()[0].Value)
			}
		})
	}
}

func TestDeliveryLogout(t *testing.T) {
	// Init
	type mockBehavior func(a *authMocks.MockUsecase, user *models.User)

	c := gomock.NewController(t)

	authMockUsecase := authMocks.NewMockUsecase(c)
	tokenMockUsecase := tokenMocks.NewMockUsecase(c)

	l := commonTests.MockLogger(c)

	h := NewHandler(authMockUsecase, tokenMockUsecase, l)

	// Routing
	r := chi.NewRouter()
	r.Get("/logout", h.Logout)

	// Test filling
	correctTestUser := &models.User{
		ID:      1,
		Version: 1,
	}

	testTable := []struct {
		name                 string
		user                 *models.User
		doWrap               bool
		mockBehavior         mockBehavior
		expectedStatus       int
		expectedResponse     string
		expectingCookieReset bool
	}{
		{
			name: "Common",
			user: correctTestUser,
			mockBehavior: func(a *authMocks.MockUsecase, user *models.User) {
				a.EXPECT().IncreaseUserVersion(gomock.Any(), user.ID).Return(nil)
			},
			expectedStatus:       http.StatusOK,
			expectedResponse:     commonTests.OKResponse(userLogedOutSuccessfully),
			doWrap:               true,
			expectingCookieReset: true,
		},
		{
			name:             "No user in request",
			user:             nil,
			mockBehavior:     func(a *authMocks.MockUsecase, user *models.User) {},
			expectedStatus:   http.StatusUnauthorized,
			expectedResponse: commonTests.ErrorResponse(invalidToken),
			doWrap:           false,
		},
		{
			name:             "Nil user in request",
			user:             nil,
			mockBehavior:     func(a *authMocks.MockUsecase, user *models.User) {},
			expectedStatus:   http.StatusUnauthorized,
			expectedResponse: commonTests.ErrorResponse(invalidToken),
			doWrap:           true,
		},
		{
			name: "Failed to increase user version",
			user: correctTestUser,
			mockBehavior: func(a *authMocks.MockUsecase, user *models.User) {
				a.EXPECT().IncreaseUserVersion(gomock.Any(), user.ID).Return(fmt.Errorf("database error"))
			},
			expectedStatus:   http.StatusInternalServerError,
			expectedResponse: commonTests.ErrorResponse(userLogoutServerError),
			doWrap:           true,
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			// Call mock
			tc.mockBehavior(authMockUsecase, tc.user)

			w := commonTests.DeliveryTestGet(t, r, "/logout", tc.expectedStatus, tc.expectedResponse,
				commonTests.WrapRequestWithUserFunc(tc.user, tc.doWrap))

			if tc.expectingCookieReset {
				assert.Equal(t, commonHTTP.AccessTokenCookieName, w.Result().Cookies()[0].Name)
				assert.Equal(t, "", w.Result().Cookies()[0].Value)
			}
		})
	}
}
