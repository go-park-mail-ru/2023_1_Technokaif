package http

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/models"

	commonTests "github.com/go-park-mail-ru/2023_1_Technokaif/internal/common/tests"
	albumMocks "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/album/mocks"
	artistMocks "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/artist/mocks"
	trackMocks "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/track/mocks"
	userMocks "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/user/mocks"
)

func getCorrectUser(t *testing.T) *models.User {
	birthTime, err := time.Parse(time.RFC3339, "2003-08-23T00:00:00Z")
	if err != nil {
		require.NoError(t, err, "can't Parse birth date")
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

func getCorrectUserInfo(t *testing.T) *models.User {
	birthTime, err := time.Parse(time.RFC3339, "2003-08-23T00:00:00Z")
	if err != nil {
		require.NoError(t, err, "can't Parse birth date")
	}
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

func TestUserDeliveryGet(t *testing.T) {
	// Init
	c := gomock.NewController(t)

	uu := userMocks.NewMockUsecase(c)
	tu := trackMocks.NewMockUsecase(c)
	alu := albumMocks.NewMockUsecase(c)
	aru := artistMocks.NewMockUsecase(c)

	l := commonTests.MockLogger(c)

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
			expectedStatus:   200,
			expectedResponse: correctResponse,
		},
		{
			name:             "Server error",
			userIDPath:       correctUserIDPath,
			user:             nil,
			expectedStatus:   500,
			expectedResponse: `{"message": "can't get user"}`,
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			commonTests.DeliveryTestGet(t, r, "/api/users/"+tc.userIDPath+"/", tc.expectedStatus, tc.expectedResponse,
				commonTests.WrapRequestWithUserNotNilFunc(tc.user))
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

	l := commonTests.MockLogger(c)

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
				uu.EXPECT().UpdateInfo(user).Return(nil)
			},
			expectedStatus:   200,
			expectedResponse: `{"status": "ok"}`,
		},
		{
			name:             "Incorrect Body",
			userIDPath:       correctUserIDPath,
			user:             getCorrectUserInfo(t),
			requestBody:      `{"id": 1`,
			mockBehavior:     func(uu *userMocks.MockUsecase, user *models.User) {},
			expectedStatus:   400,
			expectedResponse: `{"message": "incorrect input body"}`,
		},
		{
			name:        "No Such User",
			userIDPath:  correctUserIDPath,
			user:        getCorrectUserInfo(t),
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
			user:        getCorrectUserInfo(t),
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

			commonTests.DeliveryTestPost(t, r, "/api/users/"+tc.userIDPath+"/update", tc.requestBody, tc.expectedStatus, tc.expectedResponse,
				commonTests.WrapRequestWithUserNotNilFunc(tc.user))
		})
	}
}

func TestUserDeliveryGetFavoriteTracks(t *testing.T) {
	type mockBehavior func(tu *trackMocks.MockUsecase, au *artistMocks.MockUsecase, userID uint32)

	c := gomock.NewController(t)

	uu := userMocks.NewMockUsecase(c)
	tu := trackMocks.NewMockUsecase(c)
	alu := albumMocks.NewMockUsecase(c)
	aru := artistMocks.NewMockUsecase(c)

	l := commonTests.MockLogger(c)

	h := NewHandler(uu, tu, alu, aru, l)

	// Routing
	r := chi.NewRouter()
	r.Get("/api/users/{userID}/tracks", h.GetFavouriteTracks)

	// Test filling
	correctUserID := uint32(1)
	correctUserIDPath := fmt.Sprint(correctUserID)

	expectedReturnTracks := []models.Track{
		{
			ID:        1,
			Name:      "Накануне",
			CoverSrc:  "/tracks/covers/1.png",
			Listens:   2700000,
			RecordSrc: "/tracks/records/1.wav",
		},
		{
			ID:        2,
			Name:      "LAGG OUT",
			CoverSrc:  "/tracks/covers/2.png",
			Listens:   4500000,
			RecordSrc: "/tracks/records/2.wav",
		},
	}

	expectedReturnArtists := []models.Artist{
		{
			ID:        1,
			Name:      "Oxxxymiron",
			AvatarSrc: "/artists/avatars/1.png",
		},
		{
			ID:        2,
			Name:      "SALUKI",
			AvatarSrc: "/artists/avatars/2.png",
		},
	}

	correctResponse := `[
		{
			"id": 1,
			"name": "Накануне",
			"artists": [
				{
					"id": 1,
					"name": "Oxxxymiron",
					"cover": "/artists/avatars/1.png"
				}
			],
			"cover": "/tracks/covers/1.png",
			"listens": 2700000,
			"isLiked": true,
			"recordSrc": "/tracks/records/1.wav"
		},
		{
			"id": 2,
			"name": "LAGG OUT",
			"artists": [
				{
					"id": 2,
					"name": "SALUKI",
					"cover": "/artists/avatars/2.png"
				}
			],
			"cover": "/tracks/covers/2.png",
			"listens": 4500000,
			"isLiked": true,
			"recordSrc": "/tracks/records/2.wav"
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
			user: getCorrectUser(t),
			mockBehavior: func(tu *trackMocks.MockUsecase, au *artistMocks.MockUsecase, userID uint32) {
				tu.EXPECT().GetLikedByUser(userID).Return(expectedReturnTracks, nil)
				for ind, track := range expectedReturnTracks {
					au.EXPECT().GetByTrack(track.ID).Return(expectedReturnArtists[ind:ind+1], nil)
					tu.EXPECT().IsLiked(track.ID, userID).Return(true, nil)
				}
			},
			expectedStatus:   200,
			expectedResponse: correctResponse,
		},
		{
			name: "Tracks Issue",
			user: getCorrectUser(t),
			mockBehavior: func(tu *trackMocks.MockUsecase, au *artistMocks.MockUsecase, userID uint32) {
				tu.EXPECT().GetLikedByUser(userID).Return(nil, errors.New(""))
			},
			expectedStatus:   500,
			expectedResponse: `{"message": "can't get favorite tracks"}`,
		},
		{
			name: "Artists Issue",
			user: getCorrectUser(t),
			mockBehavior: func(tu *trackMocks.MockUsecase, au *artistMocks.MockUsecase, userID uint32) {
				tu.EXPECT().GetLikedByUser(userID).Return(expectedReturnTracks, nil)
				au.EXPECT().GetByTrack(expectedReturnTracks[0].ID).Return(nil, errors.New(""))
			},
			expectedStatus:   500,
			expectedResponse: `{"message": "can't get favorite tracks"}`,
		},
		{
			name: "Likes Issue",
			user: getCorrectUser(t),
			mockBehavior: func(tu *trackMocks.MockUsecase, au *artistMocks.MockUsecase, userID uint32) {
				tu.EXPECT().GetLikedByUser(userID).Return(expectedReturnTracks, nil)
				au.EXPECT().GetByTrack(expectedReturnTracks[0].ID).Return(expectedReturnArtists[0:1], nil)
				tu.EXPECT().IsLiked(expectedReturnTracks[0].ID, userID).Return(false, errors.New(""))
			},
			expectedStatus:   500,
			expectedResponse: `{"message": "can't get favorite tracks"}`,
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			// Call mock
			tc.mockBehavior(tu, aru, tc.user.ID)

			commonTests.DeliveryTestGet(t, r, "/api/users/"+correctUserIDPath+"/tracks", tc.expectedStatus, tc.expectedResponse,
				commonTests.WrapRequestWithUserNotNilFunc(tc.user))
		})
	}
}

func TestUserDeliveryGetFavoriteAlbums(t *testing.T) {
	type mockBehavior func(alu *albumMocks.MockUsecase, aru *artistMocks.MockUsecase, userID uint32)

	c := gomock.NewController(t)

	uu := userMocks.NewMockUsecase(c)
	tu := trackMocks.NewMockUsecase(c)
	alu := albumMocks.NewMockUsecase(c)
	aru := artistMocks.NewMockUsecase(c)

	l := commonTests.MockLogger(c)

	h := NewHandler(uu, tu, alu, aru, l)

	// Routing
	r := chi.NewRouter()
	r.Get("/api/users/{userID}/albums", h.GetFavouriteAlbums)

	// Test filling
	correctUserID := uint32(1)
	correctUserIDPath := fmt.Sprint(correctUserID)

	descriptionID1 := "Антиутопия"
	descriptionID2 := "Стиль"
	expectedReturnAlbums := []models.Album{
		{
			ID:          1,
			Name:        "Горгород",
			Description: &descriptionID1,
			CoverSrc:    "/albums/covers/gorgorod.png",
		},
		{
			ID:          2,
			Name:        "Властелин Калек",
			Description: &descriptionID2,
			CoverSrc:    "/albums/covers/vlkal.png",
		},
	}

	expectedReturnArtists := []models.Artist{
		{
			ID:        1,
			Name:      "Oxxxymiron",
			AvatarSrc: "/artists/avatars/oxxxymiron.png",
		},
		{
			ID:        2,
			Name:      "SALUKI",
			AvatarSrc: "/artists/avatars/saluki.png",
		},
	}

	correctResponse := `[
		{
			"id": 1,
			"name": "Горгород",
			"artists": [
				{
					"id": 1,
					"name": "Oxxxymiron",
					"cover": "/artists/avatars/oxxxymiron.png"
				}
			],
			"description": "Антиутопия",
			"cover": "/albums/covers/gorgorod.png"
		},
		{
			"id": 2,
			"name": "Властелин Калек",
			"artists": [
				{
					"id": 2,
					"name": "SALUKI",
					"cover": "/artists/avatars/saluki.png"
				}
			],
			"description": "Стиль",
			"cover": "/albums/covers/vlkal.png"
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
			user: getCorrectUser(t),
			mockBehavior: func(alu *albumMocks.MockUsecase, au *artistMocks.MockUsecase, userID uint32) {
				alu.EXPECT().GetLikedByUser(userID).Return(expectedReturnAlbums, nil)
				for ind, track := range expectedReturnAlbums {
					au.EXPECT().GetByAlbum(track.ID).Return(expectedReturnArtists[ind:ind+1], nil)
				}
			},
			expectedStatus:   200,
			expectedResponse: correctResponse,
		},
		{
			name: "Albums Issue",
			user: getCorrectUser(t),
			mockBehavior: func(alu *albumMocks.MockUsecase, au *artistMocks.MockUsecase, userID uint32) {
				alu.EXPECT().GetLikedByUser(userID).Return(nil, errors.New(""))
			},
			expectedStatus:   500,
			expectedResponse: `{"message": "can't get favorite albums"}`,
		},
		{
			name: "Artists Issue",
			user: getCorrectUser(t),
			mockBehavior: func(alu *albumMocks.MockUsecase, au *artistMocks.MockUsecase, userID uint32) {
				alu.EXPECT().GetLikedByUser(userID).Return(expectedReturnAlbums, nil)
				au.EXPECT().GetByAlbum(expectedReturnAlbums[0].ID).Return(nil, errors.New(""))
			},
			expectedStatus:   500,
			expectedResponse: `{"message": "can't get favorite albums"}`,
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			// Call mock
			tc.mockBehavior(alu, aru, tc.user.ID)

			commonTests.DeliveryTestGet(t, r, "/api/users/"+correctUserIDPath+"/albums", tc.expectedStatus, tc.expectedResponse,
				commonTests.WrapRequestWithUserNotNilFunc(tc.user))
		})
	}
}

func TestUserDeliveryGetFavoriteArtists(t *testing.T) {
	type mockBehavior func(aru *artistMocks.MockUsecase, userID uint32)

	c := gomock.NewController(t)

	uu := userMocks.NewMockUsecase(c)
	tu := trackMocks.NewMockUsecase(c)
	alu := albumMocks.NewMockUsecase(c)
	aru := artistMocks.NewMockUsecase(c)

	l := commonTests.MockLogger(c)

	h := NewHandler(uu, tu, alu, aru, l)

	// Routing
	r := chi.NewRouter()
	r.Get("/api/users/{userID}/artists", h.GetFavouriteArtists)

	// Test filling
	correctUserID := uint32(1)
	correctUserIDPath := fmt.Sprint(correctUserID)

	expectedReturnArtists := []models.Artist{
		{
			ID:        1,
			Name:      "Oxxxymiron",
			AvatarSrc: "/artists/avatars/oxxxymiron.png",
		},
		{
			ID:        2,
			Name:      "SALUKI",
			AvatarSrc: "/artists/avatars/saluki.png",
		},
	}

	correctResponse := `[
		{
			"id": 1,
			"name": "Oxxxymiron",
			"cover": "/artists/avatars/oxxxymiron.png"
		},
		{
			"id": 2,
			"name": "SALUKI",
			"cover": "/artists/avatars/saluki.png"
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
			user: getCorrectUser(t),
			mockBehavior: func(au *artistMocks.MockUsecase, userID uint32) {
				au.EXPECT().GetLikedByUser(userID).Return(expectedReturnArtists, nil)
			},
			expectedStatus:   200,
			expectedResponse: correctResponse,
		},
		{
			name: "Artists Issue",
			user: getCorrectUser(t),
			mockBehavior: func(au *artistMocks.MockUsecase, userID uint32) {
				au.EXPECT().GetLikedByUser(userID).Return(nil, errors.New(""))
			},
			expectedStatus:   500,
			expectedResponse: `{"message": "can't get favorite artists"}`,
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			// Call mock
			tc.mockBehavior(aru, tc.user.ID)

			commonTests.DeliveryTestGet(t, r, "/api/users/"+correctUserIDPath+"/artists", tc.expectedStatus, tc.expectedResponse,
				commonTests.WrapRequestWithUserNotNilFunc(tc.user))
		})
	}
}
