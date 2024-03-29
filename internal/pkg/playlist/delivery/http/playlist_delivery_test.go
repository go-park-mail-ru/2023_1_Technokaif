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

	commonHTTP "github.com/go-park-mail-ru/2023_1_Technokaif/internal/common/http"
	commonTests "github.com/go-park-mail-ru/2023_1_Technokaif/internal/common/tests"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/models"
	playlistMocks "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/playlist/mocks"
	trackMocks "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/track/mocks"
	userMocks "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/user/mocks"
)

var correctUser = models.User{
	ID: 1,
}

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
		BirthDate: birthDate,
		AvatarSrc: "/users/avatars/yarik_tri.png",
	}
}

func TestPlaylistDeliveryHTTP_Create(t *testing.T) {
	// Init
	type mockBehavior func(pu *playlistMocks.MockUsecase)

	c := gomock.NewController(t)

	pu := playlistMocks.NewMockUsecase(c)
	tu := trackMocks.NewMockUsecase(c)
	uu := userMocks.NewMockUsecase(c)

	l := commonTests.MockLogger(c)

	h := NewHandler(pu, tu, uu, l)

	// Routing
	r := chi.NewRouter()
	r.Post("/api/playlists/", h.Create)

	// Test filling
	correctRequestBody := `{
		"name": "Музыка для эпичной защиты",
		"users": [1],
		"description": "Ожидайте 3 июня"
	}`

	correctUsersID := []uint32{1}

	description := "Ожидайте 3 июня"
	expectedCallPlaylist := models.Playlist{
		Name:        "Музыка для эпичной защиты",
		Description: &description,
	}

	testTable := []struct {
		name             string
		user             *models.User
		requestBody      string
		mockBehavior     mockBehavior
		expectedStatus   int
		expectedResponse string
	}{
		{
			name:        "Common",
			user:        &correctUser,
			requestBody: correctRequestBody,
			mockBehavior: func(pu *playlistMocks.MockUsecase) {
				pu.EXPECT().Create(
					gomock.Any(), expectedCallPlaylist, correctUsersID, correctUser.ID,
				).Return(uint32(1), nil)
			},
			expectedStatus:   http.StatusOK,
			expectedResponse: `{"id": 1}`,
		},
		{
			name:             "No User",
			user:             nil,
			mockBehavior:     func(pu *playlistMocks.MockUsecase) {},
			expectedStatus:   http.StatusUnauthorized,
			expectedResponse: commonTests.ErrorResponse(commonHTTP.UnathorizedUser),
		},
		{
			name: "Incorrect JSON",
			user: &correctUser,
			requestBody: `{
				"name": ,
				"users": [1],
				"description": "Ожидайте 3 июня"
			}`,
			mockBehavior:     func(pu *playlistMocks.MockUsecase) {},
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: commonTests.ErrorResponse(commonHTTP.IncorrectRequestBody),
		},
		{
			name: "Incorrect Body (no name)",
			user: &correctUser,
			requestBody: `{
				"users": [1],
				"description": "Ожидайте 3 июня"
			}`,
			mockBehavior:     func(pu *playlistMocks.MockUsecase) {},
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: commonTests.ErrorResponse(commonHTTP.IncorrectRequestBody),
		},
		{
			name:        "User Has No Rights",
			user:        &correctUser,
			requestBody: correctRequestBody,
			mockBehavior: func(pu *playlistMocks.MockUsecase) {
				pu.EXPECT().Create(gomock.Any(), expectedCallPlaylist, correctUsersID, correctUser.ID).
					Return(uint32(0), &models.ForbiddenUserError{})
			},
			expectedStatus:   http.StatusForbidden,
			expectedResponse: commonTests.ErrorResponse(playlistCreateNorights),
		},
		{
			name:        "Server Error",
			user:        &correctUser,
			requestBody: correctRequestBody,
			mockBehavior: func(pu *playlistMocks.MockUsecase) {
				pu.EXPECT().Create(
					gomock.Any(), expectedCallPlaylist, correctUsersID, correctUser.ID,
				).Return(uint32(0), errors.New(""))
			},
			expectedStatus:   http.StatusInternalServerError,
			expectedResponse: commonTests.ErrorResponse(playlistCreateServerError),
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			// Call mock
			tc.mockBehavior(pu)

			commonTests.DeliveryTestPost(t, r, "/api/playlists/",
				tc.requestBody, tc.expectedStatus, tc.expectedResponse,
				commonTests.WrapRequestWithUserNotNilFunc(tc.user))
		})
	}
}

func TestPlaylistDeliveryHTTP_Get(t *testing.T) {
	// Init
	type mockBehavior func(pu *playlistMocks.MockUsecase, uu *userMocks.MockUsecase)

	c := gomock.NewController(t)

	pu := playlistMocks.NewMockUsecase(c)
	tu := trackMocks.NewMockUsecase(c)
	uu := userMocks.NewMockUsecase(c)

	l := commonTests.MockLogger(c)

	h := NewHandler(pu, tu, uu, l)

	// Routing
	r := chi.NewRouter()
	r.Get("/api/playlists/{playlistID}/", h.Get)

	// Test filling
	const correctPlaylistID uint32 = 1
	correctPlaylistIDPath := fmt.Sprint(correctPlaylistID)

	description := "Ожидайте 3 июня"
	expectedReturnPlaylist := models.Playlist{
		ID:          correctPlaylistID,
		Name:        "Музыка для эпичной защиты",
		Description: &description,
		CoverSrc:    "/playlists/covers/epic.png",
	}

	expectedReturnUsers := []models.User{*getCorrectUser(t)}

	correctResponse := `{
		"id": 1,
		"name": "Музыка для эпичной защиты",
		"users": [
			{
				"id": 1,
				"email": "yarik1448kuzmin@gmail.com",
				"username": "yarik_tri",
				"firstName": "Yaroslav",
				"lastName": "Kuzmin",
				"birthDate": "2003-08-23T00:00:00Z",
				"avatarSrc": "/users/avatars/yarik_tri.png"
			}
		],
		"description": "Ожидайте 3 июня",
		"isLiked": true,
		"cover": "/playlists/covers/epic.png"
	}`

	testTable := []struct {
		name             string
		playlistIDPath   string
		user             *models.User
		mockBehavior     mockBehavior
		expectedStatus   int
		expectedResponse string
	}{
		{
			name:           "Common",
			playlistIDPath: correctPlaylistIDPath,
			user:           &correctUser,
			mockBehavior: func(pu *playlistMocks.MockUsecase, uu *userMocks.MockUsecase) {
				pu.EXPECT().GetByID(gomock.Any(), correctPlaylistID).Return(&expectedReturnPlaylist, nil)
				pu.EXPECT().IsLiked(gomock.Any(), correctPlaylistID, correctUser.ID).Return(true, nil)
				uu.EXPECT().GetByPlaylist(gomock.Any(), correctPlaylistID).Return(expectedReturnUsers, nil)
			},
			expectedStatus:   http.StatusOK,
			expectedResponse: correctResponse,
		},
		{
			name:             "Incorrect ID In Path",
			playlistIDPath:   "incorrect",
			mockBehavior:     func(pu *playlistMocks.MockUsecase, uu *userMocks.MockUsecase) {},
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: commonTests.ErrorResponse(commonHTTP.InvalidURLParameter),
		},
		{
			name:           "No Playlist To Get",
			playlistIDPath: correctPlaylistIDPath,
			user:           &correctUser,
			mockBehavior: func(pu *playlistMocks.MockUsecase, uu *userMocks.MockUsecase) {
				pu.EXPECT().GetByID(gomock.Any(), correctPlaylistID).Return(nil, &models.NoSuchPlaylistError{})
			},
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: commonTests.ErrorResponse(playlistNotFound),
		},
		{
			name:           "Playlists Issues",
			playlistIDPath: correctPlaylistIDPath,
			user:           &correctUser,
			mockBehavior: func(pu *playlistMocks.MockUsecase, uu *userMocks.MockUsecase) {
				pu.EXPECT().GetByID(gomock.Any(), correctPlaylistID).Return(nil, errors.New(""))
			},
			expectedStatus:   http.StatusInternalServerError,
			expectedResponse: commonTests.ErrorResponse(playlistGetServerError),
		},
		{
			name:           "Users Issues",
			playlistIDPath: correctPlaylistIDPath,
			user:           &correctUser,
			mockBehavior: func(pu *playlistMocks.MockUsecase, uu *userMocks.MockUsecase) {
				pu.EXPECT().GetByID(gomock.Any(), correctPlaylistID).Return(&expectedReturnPlaylist, nil)
				uu.EXPECT().GetByPlaylist(gomock.Any(), correctPlaylistID).Return(nil, errors.New(""))
			},
			expectedStatus:   http.StatusInternalServerError,
			expectedResponse: commonTests.ErrorResponse(playlistGetServerError),
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			// Call mock
			tc.mockBehavior(pu, uu)

			commonTests.DeliveryTestGet(t, r, "/api/playlists/"+tc.playlistIDPath+"/",
				tc.expectedStatus, tc.expectedResponse,
				commonTests.WrapRequestWithUserNotNilFunc(tc.user))
		})
	}
}

func TestPlaylistDeliveryHTTP_Update(t *testing.T) {
	// Init
	type mockBehavior func(pu *playlistMocks.MockUsecase)

	c := gomock.NewController(t)

	pu := playlistMocks.NewMockUsecase(c)
	tu := trackMocks.NewMockUsecase(c)
	uu := userMocks.NewMockUsecase(c)

	l := commonTests.MockLogger(c)

	h := NewHandler(pu, tu, uu, l)

	// Routing
	r := chi.NewRouter()
	r.Post("/api/playlists/{playlistID}/update", h.Update)

	// Test filling
	correctRequestBody := `{
		"id": 1,
		"name": "Музыка для эпичной защиты",
		"users": [1],
		"description": "Ожидайте 3 июня"
	}`

	correctUsersID := []uint32{1}

	description := "Ожидайте 3 июня"
	expectedCallPlaylist := models.Playlist{
		ID:          1,
		Name:        "Музыка для эпичной защиты",
		Description: &description,
	}

	const correctPlaylistID uint32 = 1
	correctPlaylistIDPath := fmt.Sprint(correctPlaylistID)

	testTable := []struct {
		name             string
		playlistIDPath   string
		user             *models.User
		requestBody      string
		mockBehavior     mockBehavior
		expectedStatus   int
		expectedResponse string
	}{
		{
			name:           "Common",
			playlistIDPath: correctPlaylistIDPath,
			user:           &correctUser,
			requestBody:    correctRequestBody,
			mockBehavior: func(pu *playlistMocks.MockUsecase) {
				pu.EXPECT().UpdateInfoAndMembers(
					gomock.Any(), expectedCallPlaylist, correctUsersID, correctUser.ID,
				).Return(nil)
			},
			expectedStatus:   http.StatusOK,
			expectedResponse: commonTests.OKResponse(playlistUpdatedSuccessfully),
		},
		{
			name:             "Incorrect ID In Path",
			playlistIDPath:   "incorrect",
			mockBehavior:     func(pu *playlistMocks.MockUsecase) {},
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: commonTests.ErrorResponse(commonHTTP.InvalidURLParameter),
		},
		{
			name:             "No User",
			playlistIDPath:   correctPlaylistIDPath,
			user:             nil,
			mockBehavior:     func(pu *playlistMocks.MockUsecase) {},
			expectedStatus:   http.StatusUnauthorized,
			expectedResponse: commonTests.ErrorResponse(commonHTTP.UnathorizedUser),
		},
		{
			name:           "Incorrect JSON",
			playlistIDPath: correctPlaylistIDPath,
			user:           &correctUser,
			requestBody: `{
				"name": ,
				"users": [1],
				"description": "Ожидайте ? июня"
			}`,
			mockBehavior:     func(pu *playlistMocks.MockUsecase) {},
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: commonTests.ErrorResponse(commonHTTP.IncorrectRequestBody),
		},
		{
			name:           "Incorrect Body (no name)",
			playlistIDPath: correctPlaylistIDPath,
			user:           &correctUser,
			requestBody: `{
				"users": [1],
				"description": "Ожидайте ? июня"
			}`,
			mockBehavior:     func(pu *playlistMocks.MockUsecase) {},
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: commonTests.ErrorResponse(commonHTTP.IncorrectRequestBody),
		},
		{
			name:           "User Has No Rights",
			playlistIDPath: correctPlaylistIDPath,
			user:           &correctUser,
			requestBody:    correctRequestBody,
			mockBehavior: func(pu *playlistMocks.MockUsecase) {
				pu.EXPECT().UpdateInfoAndMembers(
					gomock.Any(), expectedCallPlaylist, correctUsersID, correctUser.ID,
				).Return(&models.ForbiddenUserError{})
			},
			expectedStatus:   http.StatusForbidden,
			expectedResponse: commonTests.ErrorResponse(playlistUpdateNoRights),
		},
		{
			name:           "Server Error",
			playlistIDPath: correctPlaylistIDPath,
			user:           &correctUser,
			requestBody:    correctRequestBody,
			mockBehavior: func(pu *playlistMocks.MockUsecase) {
				pu.EXPECT().UpdateInfoAndMembers(
					gomock.Any(), expectedCallPlaylist, correctUsersID, correctUser.ID,
				).Return(errors.New(""))
			},
			expectedStatus:   http.StatusInternalServerError,
			expectedResponse: commonTests.ErrorResponse(playlistUpdateServerError),
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			// Call mock
			tc.mockBehavior(pu)

			commonTests.DeliveryTestPost(t, r, "/api/playlists/"+tc.playlistIDPath+"/update",
				tc.requestBody, tc.expectedStatus, tc.expectedResponse,
				commonTests.WrapRequestWithUserNotNilFunc(tc.user))
		})
	}
}

func TestPlaylistDeliveryHTTP_Delete(t *testing.T) {
	// Init
	type mockBehavior func(pu *playlistMocks.MockUsecase)

	c := gomock.NewController(t)

	pu := playlistMocks.NewMockUsecase(c)
	tu := trackMocks.NewMockUsecase(c)
	uu := userMocks.NewMockUsecase(c)

	l := commonTests.MockLogger(c)

	h := NewHandler(pu, tu, uu, l)

	// Routing
	r := chi.NewRouter()
	r.Delete("/api/playlists/{playlistID}/", h.Delete)

	const correctPlaylistID uint32 = 1
	correctPlaylistIDPath := fmt.Sprint(correctPlaylistID)

	testTable := []struct {
		name             string
		playlistIDPath   string
		user             *models.User
		mockBehavior     mockBehavior
		expectedStatus   int
		expectedResponse string
	}{
		{
			name:           "Common",
			playlistIDPath: correctPlaylistIDPath,
			user:           &correctUser,
			mockBehavior: func(pu *playlistMocks.MockUsecase) {
				pu.EXPECT().Delete(
					gomock.Any(), correctPlaylistID, correctUser.ID,
				).Return(nil)
			},
			expectedStatus:   http.StatusOK,
			expectedResponse: commonTests.OKResponse(playlistDeletedSuccessfully),
		},
		{
			name:             "Incorrect ID In Path",
			playlistIDPath:   "incorrect",
			mockBehavior:     func(pu *playlistMocks.MockUsecase) {},
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: commonTests.ErrorResponse(commonHTTP.InvalidURLParameter),
		},
		{
			name:             "No User",
			playlistIDPath:   correctPlaylistIDPath,
			user:             nil,
			mockBehavior:     func(pu *playlistMocks.MockUsecase) {},
			expectedStatus:   http.StatusUnauthorized,
			expectedResponse: commonTests.ErrorResponse(commonHTTP.UnathorizedUser),
		},
		{
			name:           "User Has No Rights",
			playlistIDPath: correctPlaylistIDPath,
			user:           &correctUser,
			mockBehavior: func(pu *playlistMocks.MockUsecase) {
				pu.EXPECT().Delete(
					gomock.Any(), correctPlaylistID, correctUser.ID,
				).Return(&models.ForbiddenUserError{})
			},
			expectedStatus:   http.StatusForbidden,
			expectedResponse: commonTests.ErrorResponse(playlistDeleteNoRights),
		},
		{
			name:           "No Playlist To Delete",
			playlistIDPath: correctPlaylistIDPath,
			user:           &correctUser,
			mockBehavior: func(pu *playlistMocks.MockUsecase) {
				pu.EXPECT().Delete(
					gomock.Any(), correctPlaylistID, correctUser.ID,
				).Return(&models.NoSuchPlaylistError{})
			},
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: commonTests.ErrorResponse(playlistNotFound),
		},
		{
			name:           "Server Error",
			playlistIDPath: correctPlaylistIDPath,
			user:           &correctUser,
			mockBehavior: func(pu *playlistMocks.MockUsecase) {
				pu.EXPECT().Delete(
					gomock.Any(), correctPlaylistID, correctUser.ID,
				).Return(errors.New(""))
			},
			expectedStatus:   http.StatusInternalServerError,
			expectedResponse: commonTests.ErrorResponse(playlistDeleteServerError),
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			// Call mock
			tc.mockBehavior(pu)

			commonTests.DeliveryTestDelete(t, r, "/api/playlists/"+tc.playlistIDPath+"/",
				tc.expectedStatus, tc.expectedResponse,
				commonTests.WrapRequestWithUserNotNilFunc(tc.user))
		})
	}
}

func TestPlaylistDeliveryHTTP_AddTrack(t *testing.T) {
	// Init
	type mockBehavior func(pu *playlistMocks.MockUsecase)

	c := gomock.NewController(t)

	pu := playlistMocks.NewMockUsecase(c)
	tu := trackMocks.NewMockUsecase(c)
	uu := userMocks.NewMockUsecase(c)

	l := commonTests.MockLogger(c)

	h := NewHandler(pu, tu, uu, l)

	// Routing
	r := chi.NewRouter()
	r.Post("/api/playlists/{playlistID}/tracks/{trackID}", h.AddTrack)

	const correctPlaylistID uint32 = 1
	correctPlaylistIDPath := fmt.Sprint(correctPlaylistID)
	const correctTrackID uint32 = 1
	correctTrackIDPath := fmt.Sprint(correctTrackID)

	testTable := []struct {
		name             string
		playlistIDPath   string
		trackIDPath      string
		user             *models.User
		mockBehavior     mockBehavior
		expectedStatus   int
		expectedResponse string
	}{
		{
			name:           "Common",
			playlistIDPath: correctPlaylistIDPath,
			trackIDPath:    correctTrackIDPath,
			user:           &correctUser,
			mockBehavior: func(pu *playlistMocks.MockUsecase) {
				pu.EXPECT().AddTrack(
					gomock.Any(), correctPlaylistID, correctTrackID, correctUser.ID,
				).Return(nil)
			},
			expectedStatus:   http.StatusOK,
			expectedResponse: commonTests.OKResponse(playlistTrackAddedSuccessfully),
		},
		{
			name:             "Incorrect Playlist ID In Path",
			playlistIDPath:   "incorrect",
			trackIDPath:      correctTrackIDPath,
			mockBehavior:     func(pu *playlistMocks.MockUsecase) {},
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: commonTests.ErrorResponse(commonHTTP.InvalidURLParameter),
		},
		{
			name:             "Incorrect Track ID In Path",
			playlistIDPath:   correctPlaylistIDPath,
			trackIDPath:      "0",
			mockBehavior:     func(pu *playlistMocks.MockUsecase) {},
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: commonTests.ErrorResponse(commonHTTP.InvalidURLParameter),
		},
		{
			name:             "No User",
			playlistIDPath:   correctPlaylistIDPath,
			trackIDPath:      correctTrackIDPath,
			user:             nil,
			mockBehavior:     func(pu *playlistMocks.MockUsecase) {},
			expectedStatus:   http.StatusUnauthorized,
			expectedResponse: commonTests.ErrorResponse(commonHTTP.UnathorizedUser),
		},
		{
			name:           "User Has No Rights",
			playlistIDPath: correctPlaylistIDPath,
			trackIDPath:    correctTrackIDPath,
			user:           &correctUser,
			mockBehavior: func(pu *playlistMocks.MockUsecase) {
				pu.EXPECT().AddTrack(
					gomock.Any(), correctTrackID, correctPlaylistID, correctUser.ID,
				).Return(&models.ForbiddenUserError{})
			},
			expectedStatus:   http.StatusForbidden,
			expectedResponse: commonTests.ErrorResponse(playlistAddTrackNoRights),
		},
		{
			name:           "No Playlist",
			playlistIDPath: correctPlaylistIDPath,
			trackIDPath:    correctTrackIDPath,
			user:           &correctUser,
			mockBehavior: func(pu *playlistMocks.MockUsecase) {
				pu.EXPECT().AddTrack(
					gomock.Any(), correctTrackID, correctPlaylistID, correctUser.ID,
				).Return(&models.NoSuchPlaylistError{})
			},
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: commonTests.ErrorResponse(playlistNotFound),
		},
		{
			name:           "No Track",
			playlistIDPath: correctPlaylistIDPath,
			trackIDPath:    correctTrackIDPath,
			user:           &correctUser,
			mockBehavior: func(pu *playlistMocks.MockUsecase) {
				pu.EXPECT().AddTrack(
					gomock.Any(), correctTrackID, correctPlaylistID, correctUser.ID,
				).Return(&models.NoSuchTrackError{})
			},
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: commonTests.ErrorResponse(trackNotFound),
		},
		{
			name:           "Server Error",
			playlistIDPath: correctPlaylistIDPath,
			trackIDPath:    correctTrackIDPath,
			user:           &correctUser,
			mockBehavior: func(pu *playlistMocks.MockUsecase) {
				pu.EXPECT().AddTrack(
					gomock.Any(), correctTrackID, correctPlaylistID, correctUser.ID,
				).Return(errors.New(""))
			},
			expectedStatus:   http.StatusInternalServerError,
			expectedResponse: commonTests.ErrorResponse(playlistAddTrackServerError),
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			// Call mock
			tc.mockBehavior(pu)

			commonTests.DeliveryTestPost(t, r,
				"/api/playlists/"+tc.playlistIDPath+"/tracks/"+tc.trackIDPath,
				"", tc.expectedStatus, tc.expectedResponse,
				commonTests.WrapRequestWithUserNotNilFunc(tc.user))
		})
	}
}

func TestPlaylistDeliveryHTTP_DeleteTrack(t *testing.T) {
	// Init
	type mockBehavior func(pu *playlistMocks.MockUsecase)

	c := gomock.NewController(t)

	pu := playlistMocks.NewMockUsecase(c)
	tu := trackMocks.NewMockUsecase(c)
	uu := userMocks.NewMockUsecase(c)

	l := commonTests.MockLogger(c)

	h := NewHandler(pu, tu, uu, l)

	// Routing
	r := chi.NewRouter()
	r.Delete("/api/playlists/{playlistID}/tracks/{trackID}", h.DeleteTrack)

	const correctPlaylistID uint32 = 1
	correctPlaylistIDPath := fmt.Sprint(correctPlaylistID)
	const correctTrackID uint32 = 1
	correctTrackIDPath := fmt.Sprint(correctTrackID)

	testTable := []struct {
		name             string
		playlistIDPath   string
		trackIDPath      string
		user             *models.User
		mockBehavior     mockBehavior
		expectedStatus   int
		expectedResponse string
	}{
		{
			name:           "Common",
			playlistIDPath: correctPlaylistIDPath,
			trackIDPath:    correctTrackIDPath,
			user:           &correctUser,
			mockBehavior: func(pu *playlistMocks.MockUsecase) {
				pu.EXPECT().DeleteTrack(
					gomock.Any(), correctTrackID, correctPlaylistID, correctUser.ID,
				).Return(nil)
			},
			expectedStatus:   http.StatusOK,
			expectedResponse: commonTests.OKResponse(playlistTrackAddedSuccessfully),
		},
		{
			name:             "Incorrect Playlist ID In Path",
			playlistIDPath:   "incorrect",
			trackIDPath:      correctTrackIDPath,
			mockBehavior:     func(pu *playlistMocks.MockUsecase) {},
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: commonTests.ErrorResponse(commonHTTP.InvalidURLParameter),
		},
		{
			name:             "Incorrect Track ID In Path",
			playlistIDPath:   correctPlaylistIDPath,
			trackIDPath:      "0",
			mockBehavior:     func(pu *playlistMocks.MockUsecase) {},
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: commonTests.ErrorResponse(commonHTTP.InvalidURLParameter),
		},
		{
			name:             "No User",
			playlistIDPath:   correctPlaylistIDPath,
			trackIDPath:      correctTrackIDPath,
			user:             nil,
			mockBehavior:     func(pu *playlistMocks.MockUsecase) {},
			expectedStatus:   http.StatusUnauthorized,
			expectedResponse: commonTests.ErrorResponse(commonHTTP.UnathorizedUser),
		},
		{
			name:           "User Has No Rights",
			playlistIDPath: correctPlaylistIDPath,
			trackIDPath:    correctTrackIDPath,
			user:           &correctUser,
			mockBehavior: func(pu *playlistMocks.MockUsecase) {
				pu.EXPECT().DeleteTrack(
					gomock.Any(), correctTrackID, correctPlaylistID, correctUser.ID,
				).Return(&models.ForbiddenUserError{})
			},
			expectedStatus:   http.StatusForbidden,
			expectedResponse: commonTests.ErrorResponse(playlistDeleteTrackNoRights),
		},
		{
			name:           "No Playlist",
			playlistIDPath: correctPlaylistIDPath,
			trackIDPath:    correctTrackIDPath,
			user:           &correctUser,
			mockBehavior: func(pu *playlistMocks.MockUsecase) {
				pu.EXPECT().DeleteTrack(
					gomock.Any(), correctTrackID, correctPlaylistID, correctUser.ID,
				).Return(&models.NoSuchPlaylistError{})
			},
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: commonTests.ErrorResponse(playlistNotFound),
		},
		{
			name:           "No Track",
			playlistIDPath: correctPlaylistIDPath,
			trackIDPath:    correctTrackIDPath,
			user:           &correctUser,
			mockBehavior: func(pu *playlistMocks.MockUsecase) {
				pu.EXPECT().DeleteTrack(
					gomock.Any(), correctTrackID, correctPlaylistID, correctUser.ID,
				).Return(&models.NoSuchTrackError{})
			},
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: commonTests.ErrorResponse(trackNotFound),
		},
		{
			name:           "Server Error",
			playlistIDPath: correctPlaylistIDPath,
			trackIDPath:    correctTrackIDPath,
			user:           &correctUser,
			mockBehavior: func(pu *playlistMocks.MockUsecase) {
				pu.EXPECT().DeleteTrack(
					gomock.Any(), correctTrackID, correctPlaylistID, correctUser.ID,
				).Return(errors.New(""))
			},
			expectedStatus:   http.StatusInternalServerError,
			expectedResponse: commonTests.ErrorResponse(playlistDeleteTrackServerError),
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			// Call mock
			tc.mockBehavior(pu)

			commonTests.DeliveryTestDelete(t, r,
				"/api/playlists/"+tc.playlistIDPath+"/tracks/"+tc.trackIDPath,
				tc.expectedStatus, tc.expectedResponse,
				commonTests.WrapRequestWithUserNotNilFunc(tc.user))
		})
	}
}

func TestPlaylistDeliveryHTTP_Feed(t *testing.T) {
	// Init
	type mockBehavior func(pu *playlistMocks.MockUsecase, uu *userMocks.MockUsecase)

	c := gomock.NewController(t)

	pu := playlistMocks.NewMockUsecase(c)
	tu := trackMocks.NewMockUsecase(c)
	uu := userMocks.NewMockUsecase(c)

	l := commonTests.MockLogger(c)

	h := NewHandler(pu, tu, uu, l)

	// Routing
	r := chi.NewRouter()
	r.Get("/api/playlists/feed", h.Feed)

	descriptionID1 := "Ожидайте 3 июня"
	descriptionID2 := "Если вдруг решил отдохнуть"
	expectedReturnPlaylists := []models.Playlist{
		{
			ID:          1,
			Name:        "Музыка для эпичной защиты",
			Description: &descriptionID1,
			CoverSrc:    "/playlists/covers/epic.png",
		},
		{
			ID:          2,
			Name:        "Для чилла",
			Description: &descriptionID2,
			CoverSrc:    "/playlists/covers/chill.png",
		},
	}

	expectedReturnUsers := []models.User{*getCorrectUser(t)}

	correctResponse := `[
		{
			"id": 1,
			"name": "Музыка для эпичной защиты",
			"users": [
				{
					"id": 1,
					"email": "yarik1448kuzmin@gmail.com",
					"username": "yarik_tri",
					"firstName": "Yaroslav",
					"lastName": "Kuzmin",
					"birthDate": "2003-08-23T00:00:00Z",
					"avatarSrc": "/users/avatars/yarik_tri.png"
				}
			],
			"description": "Ожидайте 3 июня",
			"isLiked": false,
			"cover": "/playlists/covers/epic.png"
		},
		{
			"id": 2,
			"name": "Для чилла",
			"users": [
				{
					"id": 1,
					"email": "yarik1448kuzmin@gmail.com",
					"username": "yarik_tri",
					"firstName": "Yaroslav",
					"lastName": "Kuzmin",
					"birthDate": "2003-08-23T00:00:00Z",
					"avatarSrc": "/users/avatars/yarik_tri.png"
				}
			],
			"description": "Если вдруг решил отдохнуть",
			"isLiked": false,
			"cover": "/playlists/covers/chill.png"
		}
	]`

	testTable := []struct {
		name             string
		mockBehavior     mockBehavior
		expectedStatus   int
		expectedResponse string
	}{
		{
			name: "Common",
			mockBehavior: func(pu *playlistMocks.MockUsecase, uu *userMocks.MockUsecase) {
				pu.EXPECT().GetFeed(gomock.Any()).Return(expectedReturnPlaylists, nil)
				for _, p := range expectedReturnPlaylists {
					// Makes up only for 1:1 users:playlists
					uu.EXPECT().GetByPlaylist(gomock.Any(), p.ID).Return(expectedReturnUsers[0:], nil)
				}
			},
			expectedStatus:   http.StatusOK,
			expectedResponse: correctResponse,
		},
		{
			name: "No Playlists",
			mockBehavior: func(pu *playlistMocks.MockUsecase, uu *userMocks.MockUsecase) {
				pu.EXPECT().GetFeed(gomock.Any()).Return([]models.Playlist{}, nil)
			},
			expectedStatus:   http.StatusOK,
			expectedResponse: `[]`,
		},
		{
			name: "Playlists Issues",
			mockBehavior: func(pu *playlistMocks.MockUsecase, uu *userMocks.MockUsecase) {
				pu.EXPECT().GetFeed(gomock.Any()).Return(nil, errors.New(""))
			},
			expectedStatus:   http.StatusInternalServerError,
			expectedResponse: commonTests.ErrorResponse(playlistsGetServerError),
		},
		{
			name: "Users Issues",
			mockBehavior: func(pu *playlistMocks.MockUsecase, uu *userMocks.MockUsecase) {
				pu.EXPECT().GetFeed(gomock.Any()).Return(expectedReturnPlaylists, nil)
				uu.EXPECT().GetByPlaylist(gomock.Any(), expectedReturnPlaylists[0].ID).Return(nil, errors.New(""))
			},
			expectedStatus:   http.StatusInternalServerError,
			expectedResponse: commonTests.ErrorResponse(playlistsGetServerError),
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			// Call mock
			tc.mockBehavior(pu, uu)

			commonTests.DeliveryTestGet(t, r, "/api/playlists/feed",
				tc.expectedStatus, tc.expectedResponse,
				func(req *http.Request) *http.Request { return req })
		})
	}
}

func TestPlaylistDeliveryHTTP_GetByUser(t *testing.T) {
	type mockBehavior func(pu *playlistMocks.MockUsecase, uu *userMocks.MockUsecase, userID uint32)

	c := gomock.NewController(t)

	pu := playlistMocks.NewMockUsecase(c)
	tu := trackMocks.NewMockUsecase(c)
	uu := userMocks.NewMockUsecase(c)

	l := commonTests.MockLogger(c)

	h := NewHandler(pu, tu, uu, l)

	// Routing
	r := chi.NewRouter()
	r.Get("/api/users/{userID}/playlists", h.GetByUser)

	// Test filling
	const correctUserID uint32 = 1
	correctUserIDPath := fmt.Sprint(correctUserID)

	descriptionID1 := "Ожидайте 3 июня"
	descriptionID2 := "Если вдруг решил отдохнуть"
	expectedReturnPlaylists := []models.Playlist{
		{
			ID:          1,
			Name:        "Музыка для эпичной защиты",
			Description: &descriptionID1,
			CoverSrc:    "/playlists/covers/epic.png",
		},
		{
			ID:          2,
			Name:        "Для чилла",
			Description: &descriptionID2,
			CoverSrc:    "/playlists/covers/chill.png",
		},
	}

	expectedReturnUsers := []models.User{*getCorrectUser(t)}

	correctResponse := `[
		{
			"id": 1,
			"name": "Музыка для эпичной защиты",
			"users": [
				{
					"id": 1,
					"email": "yarik1448kuzmin@gmail.com",
					"username": "yarik_tri",
					"firstName": "Yaroslav",
					"lastName": "Kuzmin",
					"birthDate": "2003-08-23T00:00:00Z",
					"avatarSrc": "/users/avatars/yarik_tri.png"
				}
			],
			"description": "Ожидайте 3 июня",
			"isLiked": true,
			"cover": "/playlists/covers/epic.png"
		},
		{
			"id": 2,
			"name": "Для чилла",
			"users": [
				{
					"id": 1,
					"email": "yarik1448kuzmin@gmail.com",
					"username": "yarik_tri",
					"firstName": "Yaroslav",
					"lastName": "Kuzmin",
					"birthDate": "2003-08-23T00:00:00Z",
					"avatarSrc": "/users/avatars/yarik_tri.png"
				}
			],
			"description": "Если вдруг решил отдохнуть",
			"isLiked": true,
			"cover": "/playlists/covers/chill.png"
		}
	]`

	testTable := []struct {
		name             string
		user             *models.User
		mockBehavior     mockBehavior
		expectedStatus   int
		expectedResponse string
	}{
		{
			name: "Common",
			user: &correctUser,
			mockBehavior: func(pu *playlistMocks.MockUsecase, uu *userMocks.MockUsecase, userID uint32) {
				pu.EXPECT().GetByUser(gomock.Any(), userID).Return(expectedReturnPlaylists, nil)
				for _, playlist := range expectedReturnPlaylists {
					pu.EXPECT().IsLiked(gomock.Any(), playlist.ID, correctUserID).Return(true, nil)
					uu.EXPECT().GetByPlaylist(gomock.Any(), playlist.ID).Return(expectedReturnUsers, nil)
				}
			},
			expectedStatus:   http.StatusOK,
			expectedResponse: correctResponse,
		},
		{
			name: "Playlists Issue",
			user: &correctUser,
			mockBehavior: func(pu *playlistMocks.MockUsecase, uu *userMocks.MockUsecase, userID uint32) {
				pu.EXPECT().GetByUser(gomock.Any(), userID).Return(nil, errors.New(""))
			},
			expectedStatus:   http.StatusInternalServerError,
			expectedResponse: commonTests.ErrorResponse(playlistsGetServerError),
		},
		{
			name: "Users Issue",
			user: &correctUser,
			mockBehavior: func(pu *playlistMocks.MockUsecase, uu *userMocks.MockUsecase, userID uint32) {
				pu.EXPECT().GetByUser(gomock.Any(), userID).Return(expectedReturnPlaylists, nil)
				uu.EXPECT().GetByPlaylist(gomock.Any(), expectedReturnPlaylists[0].ID).Return(nil, errors.New(""))
			},
			expectedStatus:   http.StatusInternalServerError,
			expectedResponse: commonTests.ErrorResponse(playlistsGetServerError),
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			// Call mock
			tc.mockBehavior(pu, uu, tc.user.ID)

			commonTests.DeliveryTestGet(t, r,
				"/api/users/"+correctUserIDPath+"/playlists",
				tc.expectedStatus, tc.expectedResponse,
				commonTests.WrapRequestWithUserNotNilFunc(tc.user))
		})
	}
}

func TestPlaylistDeliveryHTTP_GetFavorite(t *testing.T) {
	type mockBehavior func(pu *playlistMocks.MockUsecase, uu *userMocks.MockUsecase, userID uint32)

	c := gomock.NewController(t)

	pu := playlistMocks.NewMockUsecase(c)
	tu := trackMocks.NewMockUsecase(c)
	uu := userMocks.NewMockUsecase(c)

	l := commonTests.MockLogger(c)

	h := NewHandler(pu, tu, uu, l)

	// Routing
	r := chi.NewRouter()
	r.Get("/api/users/{userID}/favorite/playlists", h.GetFavorite)

	// Test filling
	const correctUserID uint32 = 1
	correctUserIDPath := fmt.Sprint(correctUserID)

	descriptionID1 := "Ожидайте 3 июня"
	descriptionID2 := "Если вдруг решил отдохнуть"
	expectedReturnPlaylists := []models.Playlist{
		{
			ID:          1,
			Name:        "Музыка для эпичной защиты",
			Description: &descriptionID1,
			CoverSrc:    "/playlists/covers/epic.png",
		},
		{
			ID:          2,
			Name:        "Для чилла",
			Description: &descriptionID2,
			CoverSrc:    "/playlists/covers/chill.png",
		},
	}

	expectedReturnUsers := []models.User{*getCorrectUser(t)}

	correctResponse := `[
		{
			"id": 1,
			"name": "Музыка для эпичной защиты",
			"users": [
				{
					"id": 1,
					"email": "yarik1448kuzmin@gmail.com",
					"username": "yarik_tri",
					"firstName": "Yaroslav",
					"lastName": "Kuzmin",
					"birthDate": "2003-08-23T00:00:00Z",
					"avatarSrc": "/users/avatars/yarik_tri.png"
				}
			],
			"description": "Ожидайте 3 июня",
			"isLiked": true,
			"cover": "/playlists/covers/epic.png"
		},
		{
			"id": 2,
			"name": "Для чилла",
			"users": [
				{
					"id": 1,
					"email": "yarik1448kuzmin@gmail.com",
					"username": "yarik_tri",
					"firstName": "Yaroslav",
					"lastName": "Kuzmin",
					"birthDate": "2003-08-23T00:00:00Z",
					"avatarSrc": "/users/avatars/yarik_tri.png"
				}
			],
			"description": "Если вдруг решил отдохнуть",
			"isLiked": true,
			"cover": "/playlists/covers/chill.png"
		}
	]`

	testTable := []struct {
		name             string
		user             *models.User
		mockBehavior     mockBehavior
		expectedStatus   int
		expectedResponse string
	}{
		{
			name: "Common",
			user: &correctUser,
			mockBehavior: func(pu *playlistMocks.MockUsecase, uu *userMocks.MockUsecase, userID uint32) {
				pu.EXPECT().GetLikedByUser(gomock.Any(), userID).Return(expectedReturnPlaylists, nil)
				for _, playlist := range expectedReturnPlaylists {
					pu.EXPECT().IsLiked(gomock.Any(), playlist.ID, correctUserID).Return(true, nil)
					uu.EXPECT().GetByPlaylist(gomock.Any(), playlist.ID).Return(expectedReturnUsers, nil)
				}
			},
			expectedStatus:   http.StatusOK,
			expectedResponse: correctResponse,
		},
		{
			name: "Playlists Issue",
			user: &correctUser,
			mockBehavior: func(pu *playlistMocks.MockUsecase, uu *userMocks.MockUsecase, userID uint32) {
				pu.EXPECT().GetLikedByUser(gomock.Any(), userID).Return(nil, errors.New(""))
			},
			expectedStatus:   http.StatusInternalServerError,
			expectedResponse: commonTests.ErrorResponse(playlistsGetServerError),
		},
		{
			name: "Users Issue",
			user: &correctUser,
			mockBehavior: func(pu *playlistMocks.MockUsecase, uu *userMocks.MockUsecase, userID uint32) {
				pu.EXPECT().GetLikedByUser(gomock.Any(), userID).Return(expectedReturnPlaylists, nil)
				uu.EXPECT().GetByPlaylist(gomock.Any(), expectedReturnPlaylists[0].ID).Return(nil, errors.New(""))
			},
			expectedStatus:   http.StatusInternalServerError,
			expectedResponse: commonTests.ErrorResponse(playlistsGetServerError),
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			// Call mock
			tc.mockBehavior(pu, uu, tc.user.ID)

			commonTests.DeliveryTestGet(t, r,
				"/api/users/"+correctUserIDPath+"/favorite/playlists",
				tc.expectedStatus, tc.expectedResponse,
				commonTests.WrapRequestWithUserNotNilFunc(tc.user))
		})
	}
}

func TestPlaylistDeliveryHTTP_Like(t *testing.T) {
	// Init
	type mockBehavior func(pu *playlistMocks.MockUsecase)

	c := gomock.NewController(t)

	pu := playlistMocks.NewMockUsecase(c)
	tu := trackMocks.NewMockUsecase(c)
	uu := userMocks.NewMockUsecase(c)

	l := commonTests.MockLogger(c)

	h := NewHandler(pu, tu, uu, l)

	// Routing
	r := chi.NewRouter()
	r.Get("/api/playlists/{playlistID}/like", h.Like)

	const correctPlaylistID uint32 = 1
	correctPlaylistIDPath := fmt.Sprint(correctPlaylistID)

	testTable := []struct {
		name             string
		playlistIDPath   string
		user             *models.User
		mockBehavior     mockBehavior
		expectedStatus   int
		expectedResponse string
	}{
		{
			name:           "Common",
			playlistIDPath: correctPlaylistIDPath,
			user:           &correctUser,
			mockBehavior: func(pu *playlistMocks.MockUsecase) {
				pu.EXPECT().SetLike(gomock.Any(), correctPlaylistID, correctUser.ID).Return(true, nil)
			},
			expectedStatus:   http.StatusOK,
			expectedResponse: commonTests.OKResponse(commonHTTP.LikeSuccess),
		},
		{
			name:           "Already liked (Anyway Success)",
			playlistIDPath: correctPlaylistIDPath,
			user:           &correctUser,
			mockBehavior: func(pu *playlistMocks.MockUsecase) {
				pu.EXPECT().SetLike(gomock.Any(), correctPlaylistID, correctUser.ID).Return(false, nil)
			},
			expectedStatus:   http.StatusOK,
			expectedResponse: commonTests.OKResponse(commonHTTP.LikeAlreadyExists),
		},
		{
			name:             "Incorrect ID In Path",
			playlistIDPath:   "0",
			user:             &correctUser,
			mockBehavior:     func(pu *playlistMocks.MockUsecase) {},
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: commonTests.ErrorResponse(commonHTTP.InvalidURLParameter),
		},
		{
			name:             "No User",
			playlistIDPath:   correctPlaylistIDPath,
			user:             nil,
			mockBehavior:     func(pu *playlistMocks.MockUsecase) {},
			expectedStatus:   http.StatusUnauthorized,
			expectedResponse: commonTests.ErrorResponse(commonHTTP.UnathorizedUser),
		},
		{
			name:           "No Playlist To Like",
			playlistIDPath: correctPlaylistIDPath,
			user:           &correctUser,
			mockBehavior: func(pu *playlistMocks.MockUsecase) {
				pu.EXPECT().SetLike(gomock.Any(), correctPlaylistID, correctUser.ID).Return(false, &models.NoSuchPlaylistError{})
			},
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: commonTests.ErrorResponse(playlistNotFound),
		},
		{
			name:           "Server Error",
			playlistIDPath: correctPlaylistIDPath,
			user:           &correctUser,
			mockBehavior: func(pu *playlistMocks.MockUsecase) {
				pu.EXPECT().SetLike(gomock.Any(), correctPlaylistID, correctUser.ID).Return(false, errors.New(""))
			},
			expectedStatus:   http.StatusInternalServerError,
			expectedResponse: commonTests.ErrorResponse(commonHTTP.SetLikeServerError),
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			// Call mock
			tc.mockBehavior(pu)

			commonTests.DeliveryTestGet(t, r, "/api/playlists/"+tc.playlistIDPath+"/like",
				tc.expectedStatus, tc.expectedResponse,
				commonTests.WrapRequestWithUserNotNilFunc(tc.user))
		})
	}
}

func TestPlaylistDeliveryHTTP_UnLike(t *testing.T) {
	// Init
	type mockBehavior func(pu *playlistMocks.MockUsecase)

	c := gomock.NewController(t)

	pu := playlistMocks.NewMockUsecase(c)
	tu := trackMocks.NewMockUsecase(c)
	uu := userMocks.NewMockUsecase(c)

	l := commonTests.MockLogger(c)

	h := NewHandler(pu, tu, uu, l)

	// Routing
	r := chi.NewRouter()
	r.Get("/api/playlists/{playlistID}/unlike", h.UnLike)

	const correctPlaylistID uint32 = 1
	correctPlaylistIDPath := fmt.Sprint(correctPlaylistID)

	testTable := []struct {
		name             string
		playlistIDPath   string
		user             *models.User
		mockBehavior     mockBehavior
		expectedStatus   int
		expectedResponse string
	}{
		{
			name:           "Common",
			playlistIDPath: correctPlaylistIDPath,
			user:           &correctUser,
			mockBehavior: func(pu *playlistMocks.MockUsecase) {
				pu.EXPECT().UnLike(gomock.Any(), correctPlaylistID, correctUser.ID).Return(true, nil)
			},
			expectedStatus:   http.StatusOK,
			expectedResponse: commonTests.OKResponse(commonHTTP.UnLikeSuccess),
		},
		{
			name:           "Wasn't Liked (Anyway Success)",
			playlistIDPath: correctPlaylistIDPath,
			user:           &correctUser,
			mockBehavior: func(pu *playlistMocks.MockUsecase) {
				pu.EXPECT().UnLike(gomock.Any(), correctPlaylistID, correctUser.ID).Return(false, nil)
			},
			expectedStatus:   http.StatusOK,
			expectedResponse: commonTests.OKResponse(commonHTTP.LikeDoesntExist),
		},
		{
			name:             "Incorrect ID In Path",
			playlistIDPath:   "0",
			user:             &correctUser,
			mockBehavior:     func(pu *playlistMocks.MockUsecase) {},
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: commonTests.ErrorResponse(commonHTTP.InvalidURLParameter),
		},
		{
			name:             "No User",
			playlistIDPath:   correctPlaylistIDPath,
			user:             nil,
			mockBehavior:     func(pu *playlistMocks.MockUsecase) {},
			expectedStatus:   http.StatusUnauthorized,
			expectedResponse: commonTests.ErrorResponse(commonHTTP.UnathorizedUser),
		},
		{
			name:           "No Playlist To Unlike",
			playlistIDPath: correctPlaylistIDPath,
			user:           &correctUser,
			mockBehavior: func(pu *playlistMocks.MockUsecase) {
				pu.EXPECT().UnLike(gomock.Any(), correctPlaylistID, correctUser.ID).
					Return(false, &models.NoSuchPlaylistError{})
			},
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: commonTests.ErrorResponse(playlistNotFound),
		},
		{
			name:           "Server Error",
			playlistIDPath: correctPlaylistIDPath,
			user:           &correctUser,
			mockBehavior: func(pu *playlistMocks.MockUsecase) {
				pu.EXPECT().UnLike(gomock.Any(), correctPlaylistID, correctUser.ID).
					Return(false, errors.New(""))
			},
			expectedStatus:   http.StatusInternalServerError,
			expectedResponse: commonTests.ErrorResponse(commonHTTP.DeleteLikeServerError),
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			// Call mock
			tc.mockBehavior(pu)

			commonTests.DeliveryTestGet(t, r, "/api/playlists/"+tc.playlistIDPath+"/unlike",
				tc.expectedStatus, tc.expectedResponse,
				commonTests.WrapRequestWithUserNotNilFunc(tc.user))
		})
	}
}
