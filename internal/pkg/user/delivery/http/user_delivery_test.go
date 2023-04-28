package http

import (
	"errors"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/models"

	commonHttp "github.com/go-park-mail-ru/2023_1_Technokaif/internal/common/http"
	commonTests "github.com/go-park-mail-ru/2023_1_Technokaif/internal/common/tests"
	userMocks "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/user/mocks"
)

func getCorrectUser(t *testing.T) *models.User {
	birthTime, err := time.Parse(time.RFC3339, "2003-08-23T00:00:00Z")
	require.NoError(t, err, "can't Parse birth date")

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

func getCorrectUserInfo(t *testing.T) *models.User {
	birthTime, err := time.Parse(time.RFC3339, "2003-08-23T00:00:00Z")
	require.NoError(t, err, "can't Parse birth date")

	birthDate := models.Date{Time: birthTime}

	return &models.User{
		ID:        1,
		Email:     "yarik1448kuzmin@gmail.com",
		FirstName: "Yaroslav",
		LastName:  "Kuzmin",
		Sex:       models.Male,
		BirthDate: birthDate,
		AvatarSrc: "/users/avatars/yarik_tri.png",
	}
}

func TestUserDeliveryHTTP_Get(t *testing.T) {
	// Init
	c := gomock.NewController(t)

	uu := userMocks.NewMockUsecase(c)

	l := commonTests.MockLogger(c)

	h := NewHandler(uu, l)

	// Routing
	r := chi.NewRouter()
	r.Get("/api/users/{userID}/", h.Get)

	// Test filling
	const correctUserID uint32 = 1
	correctUserIDPath := fmt.Sprint(correctUserID)

	correctResponse := `{
		"id": 1,
		"username": "yarik_tri",
		"email": "yarik1448kuzmin@gmail.com",
		"firstName": "Yaroslav",
		"lastName": "Kuzmin",
		"sex": "M",
		"birthDate": "2003-08-23T00:00:00Z",
		"avatarSrc": "/users/avatars/yarik_tri.png"
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
			user:             getCorrectUser(t),
			expectedStatus:   http.StatusOK,
			expectedResponse: correctResponse,
		},
		{
			name:             "Server error",
			userIDPath:       correctUserIDPath,
			user:             nil,
			expectedStatus:   http.StatusInternalServerError,
			expectedResponse: commonTests.ErrorResponse(userGetServerError),
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			commonTests.DeliveryTestGet(t, r, "/api/users/"+tc.userIDPath+"/", tc.expectedStatus, tc.expectedResponse,
				commonTests.WrapRequestWithUserNotNilFunc(tc.user))
		})
	}
}

func TestUserDeliveryHTTP_UpdateInfo(t *testing.T) {
	// Init
	type mockBehavior func(uu *userMocks.MockUsecase, user *models.User)

	c := gomock.NewController(t)

	uu := userMocks.NewMockUsecase(c)

	l := commonTests.MockLogger(c)

	h := NewHandler(uu, l)

	// Routing
	r := chi.NewRouter()
	r.Post("/api/users/{userID}/update", h.UpdateInfo)

	// Test filling
	const correctUserID uint32 = 1
	correctUserIDPath := fmt.Sprint(correctUserID)

	correctBody := `{
		"id": 1,
		"email": "yarik1448kuzmin@gmail.com",
		"firstName": "Yaroslav",
		"lastName": "Kuzmin",
		"sex": "M",
		"birthDate": "2003-08-23",
		"avatarSrc": "/users/avatars/yarik_tri.png"
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
			user:        getCorrectUserInfo(t),
			requestBody: correctBody,
			mockBehavior: func(uu *userMocks.MockUsecase, user *models.User) {
				uu.EXPECT().UpdateInfo(gomock.Any(), user).Return(nil)
			},
			expectedStatus:   http.StatusOK,
			expectedResponse: commonTests.OKResponse(userUpdatedInfoSuccessfully),
		},
		{
			name:             "Incorrect Body",
			userIDPath:       correctUserIDPath,
			user:             getCorrectUserInfo(t),
			requestBody:      `{"id": 1`,
			mockBehavior:     func(uu *userMocks.MockUsecase, user *models.User) {},
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: commonTests.ErrorResponse(commonHttp.IncorrectRequestBody),
		},
		{
			name:        "No Such User",
			userIDPath:  correctUserIDPath,
			user:        getCorrectUserInfo(t),
			requestBody: correctBody,
			mockBehavior: func(uu *userMocks.MockUsecase, user *models.User) {
				uu.EXPECT().UpdateInfo(gomock.Any(), user).Return(&models.NoSuchUserError{})
			},
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: commonTests.ErrorResponse(userNotFound),
		},
		{
			name:        "Server Error",
			userIDPath:  correctUserIDPath,
			user:        getCorrectUserInfo(t),
			requestBody: correctBody,
			mockBehavior: func(uu *userMocks.MockUsecase, user *models.User) {
				uu.EXPECT().UpdateInfo(gomock.Any(), user).Return(errors.New(""))
			},
			expectedStatus:   http.StatusInternalServerError,
			expectedResponse: commonTests.ErrorResponse(userUpdateInfoServerError),
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			// Call mock
			tc.mockBehavior(uu, tc.user)

			commonTests.DeliveryTestPost(t, r, "/api/users/"+tc.userIDPath+"/update", tc.requestBody, tc.expectedStatus, tc.expectedResponse,
				commonTests.WrapRequestWithUserNotNilFunc(tc.user))
		})
	}
}
