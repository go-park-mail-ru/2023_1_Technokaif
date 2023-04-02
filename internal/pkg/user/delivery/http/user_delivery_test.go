package http

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/models"

	albumMocks "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/album/mocks"
	artistMocks "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/artist/mocks"
	trackMocks "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/track/mocks"
	userMocks "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/user/mocks"
	logMocks "github.com/go-park-mail-ru/2023_1_Technokaif/pkg/logger/mocks"
)

var wrapRequestWithUser = func(r *http.Request, user *models.User) *http.Request {
	if user == nil {
		return r
	}
	ctx := context.WithValue(r.Context(), models.ContextKeyUserType{}, user)
	return r.WithContext(ctx)
}

func getCorrectUser() *models.User {
	birthTime, err := time.Parse(time.RFC3339, "2003-08-23T00:00:00Z")
	if err != nil {
		log.Fatalf("can't Parse birth date: %v", err)
	}
	birthDate := models.Date{Time: birthTime}

	return &models.User{
		ID:        1,
		Username:  "yarik_tri",
		Email:     "yarik1448kuzmin@gmail.com",
		FirstName: "Yaroslav",
		LastName:  "Kuzmin",
		Sex:       models.Male,
		BirthDate: birthDate,
		AvatarSrc: "/users/avatars/yarik_tri.png",
	}
}

func getCorrectUserInfo() *models.User {
	birthTime, err := time.Parse(time.RFC3339, "2003-08-23T00:00:00Z")
	if err != nil {
		log.Fatalf("can't Parse birth date: %v", err)
	}
	birthDate := models.Date{Time: birthTime}

	return &models.User{
		ID:        1,
		Email:     "yarik1448kuzmin@gmail.com",
		FirstName: "Yaroslav",
		LastName:  "Kuzmin",
		Sex:       models.Male,
		BirthDate: birthDate,
	}
}

func TestUserDeliveryGet(t *testing.T) {
	// Init
	c := gomock.NewController(t)

	uu := userMocks.NewMockUsecase(c)
	tu := trackMocks.NewMockUsecase(c)
	alu := albumMocks.NewMockUsecase(c)
	aru := artistMocks.NewMockUsecase(c)

	l := logMocks.NewMockLogger(c)
	l.EXPECT().Error(gomock.Any()).AnyTimes()
	l.EXPECT().Info(gomock.Any()).AnyTimes()
	l.EXPECT().Errorf(gomock.Any(), gomock.Any()).AnyTimes()
	l.EXPECT().Infof(gomock.Any(), gomock.Any()).AnyTimes()

	h := NewHandler(uu, tu, alu, aru, l)

	// Routing
	r := chi.NewRouter()
	r.Get("/api/users/{userID}/", h.Get)

	// Test filling
	correctUserID := uint32(1)
	correctUserIDPath := fmt.Sprint(correctUserID)

	correctResponse := `{
		"id": 1,
		"username": "yarik_tri",
		"email": "yarik1448kuzmin@gmail.com",
		"firstName": "Yaroslav",
		"lastName": "Kuzmin",
		"sex": "M",
		"birthDate": "2003-08-23T00:00:00Z",
		"avatar": "/users/avatars/yarik_tri.png"
	}`

	testTable := []struct {
		name             string
		userIDPath       string
		user             *models.User
		expectedStatus   int
		expectedResponse string
	}{
		{
			name:             "Common",
			userIDPath:       correctUserIDPath,
			user:             getCorrectUser(),
			expectedStatus:   200,
			expectedResponse: correctResponse,
		},
		{
			name:             "Incorrect ID In Path",
			userIDPath:       "0",
			expectedStatus:   400,
			expectedResponse: `{"message": "invalid url parameter"}`,
		},
		{
			name:             "No User",
			userIDPath:       correctUserIDPath,
			user:             nil,
			expectedStatus:   401,
			expectedResponse: `{"message": "unathorized"}`,
		},
		{
			name:             "Forbidden",
			userIDPath:       correctUserIDPath,
			user:             &models.User{ID: 2},
			expectedStatus:   403,
			expectedResponse: `{"message": "user has no rights"}`,
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			// Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/api/users/"+tc.userIDPath+"/", nil)
			r.ServeHTTP(w, wrapRequestWithUser(req, tc.user))

			// Test
			assert.Equal(t, tc.expectedStatus, w.Code)
			assert.JSONEq(t, tc.expectedResponse, w.Body.String())
		})
	}
}

func TestUserDeliveryUpdateInfo(t *testing.T) {
	// Init
	type mockBehavior func(uu *userMocks.MockUsecase, user *models.User)

	c := gomock.NewController(t)

	uu := userMocks.NewMockUsecase(c)
	tu := trackMocks.NewMockUsecase(c)
	alu := albumMocks.NewMockUsecase(c)
	aru := artistMocks.NewMockUsecase(c)

	l := logMocks.NewMockLogger(c)
	l.EXPECT().Error(gomock.Any()).AnyTimes()
	l.EXPECT().Info(gomock.Any()).AnyTimes()
	l.EXPECT().Errorf(gomock.Any(), gomock.Any()).AnyTimes()
	l.EXPECT().Infof(gomock.Any(), gomock.Any()).AnyTimes()

	h := NewHandler(uu, tu, alu, aru, l)

	// Routing
	r := chi.NewRouter()
	r.Post("/api/users/{userID}/update", h.UpdateInfo)

	// Test filling
	correctUserID := uint32(1)
	correctUserIDPath := fmt.Sprint(correctUserID)

	correctBody := `{
		"id": 1,
		"email": "yarik1448kuzmin@gmail.com",
		"firstName": "Yaroslav",
		"lastName": "Kuzmin",
		"sex": "M",
		"birthDate": "2003-08-23"
	}`

	testTable := []struct {
		name             string
		requestBody      string
		userIDPath       string
		user             *models.User
		mockBehavior     mockBehavior
		expectedStatus   int
		expectedResponse string
	}{
		{
			name:        "Common",
			userIDPath:  correctUserIDPath,
			user:        getCorrectUserInfo(),
			requestBody: correctBody,
			mockBehavior: func(uu *userMocks.MockUsecase, user *models.User) {
				uu.EXPECT().UpdateInfo(user).Return(nil)
			},
			expectedStatus:   200,
			expectedResponse: `{"status": "ok"}`,
		},
		{
			name:             "Incorrect ID In Path",
			userIDPath:       "0",
			mockBehavior:     func(uu *userMocks.MockUsecase, user *models.User) {},
			expectedStatus:   400,
			expectedResponse: `{"message": "invalid url parameter"}`,
		},
		{
			name:             "No User",
			userIDPath:       correctUserIDPath,
			user:             nil,
			mockBehavior:     func(uu *userMocks.MockUsecase, user *models.User) {},
			expectedStatus:   401,
			expectedResponse: `{"message": "unathorized"}`,
		},
		{
			name:             "Forbidden",
			userIDPath:       correctUserIDPath,
			user:             &models.User{ID: 2},
			mockBehavior:     func(uu *userMocks.MockUsecase, user *models.User) {},
			expectedStatus:   403,
			expectedResponse: `{"message": "user has no rights"}`,
		},
		{
			name:             "Incorrect Body",
			userIDPath:       correctUserIDPath,
			user:             getCorrectUserInfo(),
			requestBody:      `{"id": 1`,
			mockBehavior:     func(uu *userMocks.MockUsecase, user *models.User) {},
			expectedStatus:   400,
			expectedResponse: `{"message": "incorrect input body"}`,
		},
		{
			name:        "No Such User",
			userIDPath:  correctUserIDPath,
			user:        getCorrectUserInfo(),
			requestBody: correctBody,
			mockBehavior: func(uu *userMocks.MockUsecase, user *models.User) {
				uu.EXPECT().UpdateInfo(user).Return(&models.NoSuchUserError{})
			},
			expectedStatus:   400,
			expectedResponse: `{"message": "no user to update"}`,
		},
		{
			name:        "Server Error",
			userIDPath:  correctUserIDPath,
			user:        getCorrectUserInfo(),
			requestBody: correctBody,
			mockBehavior: func(uu *userMocks.MockUsecase, user *models.User) {
				uu.EXPECT().UpdateInfo(user).Return(errors.New(""))
			},
			expectedStatus:   500,
			expectedResponse: `{"message": "can't change user info"}`,
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			// Call mock
			tc.mockBehavior(uu, tc.user)

			// Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/api/users/"+tc.userIDPath+"/update", bytes.NewBufferString(tc.requestBody))
			r.ServeHTTP(w, wrapRequestWithUser(req, tc.user))

			// Test
			assert.Equal(t, tc.expectedStatus, w.Code)
			assert.JSONEq(t, tc.expectedResponse, w.Body.String())
		})
	}
}
