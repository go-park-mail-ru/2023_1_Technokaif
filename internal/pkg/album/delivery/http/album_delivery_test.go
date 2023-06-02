package http

import (
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/golang/mock/gomock"

	commonHTTP "github.com/go-park-mail-ru/2023_1_Technokaif/internal/common/http"
	commonTests "github.com/go-park-mail-ru/2023_1_Technokaif/internal/common/tests"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/models"
	albumMocks "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/album/mocks"
	artistMocks "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/artist/mocks"
)

var correctUser = models.User{
	ID: 1,
}

func TestAlbumDeliveryHTTP_Create(t *testing.T) {
	// Init
	type mockBehavior func(au *albumMocks.MockUsecase)

	c := gomock.NewController(t)

	alu := albumMocks.NewMockUsecase(c)
	aru := artistMocks.NewMockUsecase(c)

	l := commonTests.MockLogger(c)

	h := NewHandler(alu, aru, l)

	// Routing
	r := chi.NewRouter()
	r.Post("/api/albums/", h.Create)

	// Test filling
	correctRequestBody := `{
		"name": "Горгород",
		"artists": [1],
		"description": "Антиутопия",
		"cover": "/albums/covers/gorgorod.png"
	}`

	correctArtistsID := []uint32{1}

	description := "Антиутопия"
	expectedCallAlbum := models.Album{
		Name:        "Горгород",
		Description: &description,

		CoverSrc: "/albums/covers/gorgorod.png",
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
			mockBehavior: func(au *albumMocks.MockUsecase) {
				au.EXPECT().Create(
					gomock.Any(), expectedCallAlbum, correctArtistsID, correctUser.ID,
				).Return(uint32(1), nil)
			},
			expectedStatus:   http.StatusOK,
			expectedResponse: `{"id": 1}`,
		},
		{
			name:             "No User",
			user:             nil,
			mockBehavior:     func(au *albumMocks.MockUsecase) {},
			expectedStatus:   http.StatusUnauthorized,
			expectedResponse: commonTests.ErrorResponse(commonHTTP.UnathorizedUser),
		},
		{
			name: "Incorrect JSON",
			user: &correctUser,
			requestBody: `{
				"name": ,
				"artistsID": [1],
				"description": "Антиутопия",
				"cover": "/albums/covers/gorgorod.png"
			}`,
			mockBehavior:     func(au *albumMocks.MockUsecase) {},
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: commonTests.ErrorResponse(commonHTTP.IncorrectRequestBody),
		},
		{
			name: "Incorrect Body (no name)",
			user: &correctUser,
			requestBody: `{
				"artistsID": [1],
				"description": "Антиутопия",
				"cover": "/albums/covers/gorgorod.png"
			}`,
			mockBehavior:     func(au *albumMocks.MockUsecase) {},
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: commonTests.ErrorResponse(commonHTTP.IncorrectRequestBody),
		},
		{
			name:        "User Has No Rights",
			user:        &correctUser,
			requestBody: correctRequestBody,
			mockBehavior: func(au *albumMocks.MockUsecase) {
				au.EXPECT().Create(
					gomock.Any(), expectedCallAlbum, correctArtistsID, correctUser.ID,
				).Return(uint32(0), &models.ForbiddenUserError{})
			},
			expectedStatus:   http.StatusForbidden,
			expectedResponse: commonTests.ErrorResponse(albumCreateNorights),
		},
		{
			name:        "Server Error",
			user:        &correctUser,
			requestBody: correctRequestBody,
			mockBehavior: func(au *albumMocks.MockUsecase) {
				au.EXPECT().Create(
					gomock.Any(), expectedCallAlbum, correctArtistsID, correctUser.ID,
				).Return(uint32(0), errors.New(""))
			},
			expectedStatus:   http.StatusInternalServerError,
			expectedResponse: commonTests.ErrorResponse(albumCreateServerError),
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			// Call mock
			tc.mockBehavior(alu)

			commonTests.DeliveryTestPost(t, r, "/api/albums/", tc.requestBody, tc.expectedStatus, tc.expectedResponse,
				commonTests.WrapRequestWithUserNotNilFunc(tc.user))
		})
	}
}

func TestAlbumDeliveryHTTP_Get(t *testing.T) {
	// Init
	type mockBehavior func(alu *albumMocks.MockUsecase, aru *artistMocks.MockUsecase)

	c := gomock.NewController(t)

	alu := albumMocks.NewMockUsecase(c)
	aru := artistMocks.NewMockUsecase(c)

	l := commonTests.MockLogger(c)

	h := NewHandler(alu, aru, l)

	// Routing
	r := chi.NewRouter()
	r.Get("/api/albums/{albumID}/", h.Get)

	// Test filling
	const correctAlbumID uint32 = 1
	correctAlbumIDPath := fmt.Sprint(correctAlbumID)

	description := "Антиутопия"
	expectedReturnAlbum := models.Album{
		ID:          correctAlbumID,
		Name:        "Горгород",
		Description: &description,
		CoverSrc:    "/albums/covers/gorgorod.png",
	}

	expectedReturnArtists := []models.Artist{
		{
			ID:        1,
			Name:      "Oxxxymiron",
			AvatarSrc: "/artists/avatars/oxxxymiron.png",
			Listens:   0,
		},
	}

	correctResponse := `{
		"id": 1,
		"name": "Горгород",
		"artists": [
			{
				"id": 1,
				"name": "Oxxxymiron",
				"isLiked": false,
				"cover": "/artists/avatars/oxxxymiron.png",
				"listens": 0
			}
		],
		"description": "Антиутопия",
		"isLiked": false,
		"cover": "/albums/covers/gorgorod.png"
	}`

	testTable := []struct {
		name             string
		albumIDPath      string
		user             *models.User
		mockBehavior     mockBehavior
		expectedStatus   int
		expectedResponse string
	}{
		{
			name:        "Common",
			albumIDPath: correctAlbumIDPath,
			user:        &correctUser,
			mockBehavior: func(alu *albumMocks.MockUsecase, aru *artistMocks.MockUsecase) {
				alu.EXPECT().GetByID(gomock.Any(), correctAlbumID).Return(&expectedReturnAlbum, nil)
				alu.EXPECT().IsLiked(gomock.Any(), correctAlbumID, correctUser.ID).Return(false, nil)
				aru.EXPECT().GetByAlbum(gomock.Any(), correctAlbumID).Return(expectedReturnArtists, nil)
				for _, a := range expectedReturnArtists {
					aru.EXPECT().IsLiked(gomock.Any(), a.ID, correctUser.ID).Return(false, nil)
				}
			},
			expectedStatus:   http.StatusOK,
			expectedResponse: correctResponse,
		},
		{
			name:             "Incorrect ID In Path",
			albumIDPath:      "incorrect",
			mockBehavior:     func(alu *albumMocks.MockUsecase, aru *artistMocks.MockUsecase) {},
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: commonTests.ErrorResponse(commonHTTP.InvalidURLParameter),
		},
		{
			name:        "No Album To Get",
			albumIDPath: correctAlbumIDPath,
			user:        &correctUser,
			mockBehavior: func(alu *albumMocks.MockUsecase, aru *artistMocks.MockUsecase) {
				alu.EXPECT().GetByID(gomock.Any(), correctAlbumID).Return(nil, &models.NoSuchAlbumError{})
			},
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: commonTests.ErrorResponse(albumNotFound),
		},
		{
			name:        "Albums Issues",
			albumIDPath: correctAlbumIDPath,
			user:        &correctUser,
			mockBehavior: func(alu *albumMocks.MockUsecase, aru *artistMocks.MockUsecase) {
				alu.EXPECT().GetByID(gomock.Any(), correctAlbumID).Return(nil, errors.New(""))
			},
			expectedStatus:   http.StatusInternalServerError,
			expectedResponse: commonTests.ErrorResponse(albumGetServerError),
		},
		{
			name:        "Artists Issues",
			albumIDPath: correctAlbumIDPath,
			user:        &correctUser,
			mockBehavior: func(alu *albumMocks.MockUsecase, aru *artistMocks.MockUsecase) {
				alu.EXPECT().GetByID(gomock.Any(), correctAlbumID).Return(&expectedReturnAlbum, nil)
				aru.EXPECT().GetByAlbum(gomock.Any(), correctAlbumID).Return(nil, errors.New(""))
			},
			expectedStatus:   http.StatusInternalServerError,
			expectedResponse: commonTests.ErrorResponse(albumGetServerError),
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			// Call mock
			tc.mockBehavior(alu, aru)

			commonTests.DeliveryTestGet(t, r, "/api/albums/"+tc.albumIDPath+"/",
				tc.expectedStatus, tc.expectedResponse,
				commonTests.WrapRequestWithUserNotNilFunc(tc.user))
		})
	}
}

func TestAlbumDeliveryHTTP_Delete(t *testing.T) {
	// Init
	type mockBehavior func(au *albumMocks.MockUsecase)

	c := gomock.NewController(t)

	alu := albumMocks.NewMockUsecase(c)
	aru := artistMocks.NewMockUsecase(c)

	l := commonTests.MockLogger(c)

	h := NewHandler(alu, aru, l)

	// Routing
	r := chi.NewRouter()
	r.Delete("/api/albums/{albumID}/", h.Delete)

	const correctAlbumID uint32 = 1
	correctAlbumIDPath := fmt.Sprint(correctAlbumID)

	testTable := []struct {
		name             string
		albumIDPath      string
		user             *models.User
		mockBehavior     mockBehavior
		expectedStatus   int
		expectedResponse string
	}{
		{
			name:        "Common",
			albumIDPath: correctAlbumIDPath,
			user:        &correctUser,
			mockBehavior: func(au *albumMocks.MockUsecase) {
				au.EXPECT().Delete(
					gomock.Any(), correctAlbumID, correctUser.ID,
				).Return(nil)
			},
			expectedStatus:   http.StatusOK,
			expectedResponse: commonTests.OKResponse(albumDeletedSuccessfully),
		},
		{
			name:             "Incorrect ID In Path",
			albumIDPath:      "incorrect",
			mockBehavior:     func(au *albumMocks.MockUsecase) {},
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: commonTests.ErrorResponse(commonHTTP.InvalidURLParameter),
		},
		{
			name:             "No User",
			albumIDPath:      correctAlbumIDPath,
			user:             nil,
			mockBehavior:     func(au *albumMocks.MockUsecase) {},
			expectedStatus:   http.StatusUnauthorized,
			expectedResponse: commonTests.ErrorResponse(commonHTTP.UnathorizedUser),
		},
		{
			name:        "User Has No Rights",
			albumIDPath: correctAlbumIDPath,
			user:        &correctUser,
			mockBehavior: func(au *albumMocks.MockUsecase) {
				au.EXPECT().Delete(
					gomock.Any(), correctAlbumID, correctUser.ID,
				).Return(&models.ForbiddenUserError{})
			},
			expectedStatus:   http.StatusForbidden,
			expectedResponse: commonTests.ErrorResponse(albumDeleteNoRights),
		},
		{
			name:        "No Album To Delete",
			albumIDPath: correctAlbumIDPath,
			user:        &correctUser,
			mockBehavior: func(au *albumMocks.MockUsecase) {
				au.EXPECT().Delete(
					gomock.Any(), correctAlbumID, correctUser.ID,
				).Return(&models.NoSuchAlbumError{})
			},
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: commonTests.ErrorResponse(albumNotFound),
		},
		{
			name:        "Server Error",
			albumIDPath: correctAlbumIDPath,
			user:        &correctUser,
			mockBehavior: func(au *albumMocks.MockUsecase) {
				au.EXPECT().Delete(
					gomock.Any(), correctAlbumID, correctUser.ID,
				).Return(errors.New(""))
			},
			expectedStatus:   http.StatusInternalServerError,
			expectedResponse: commonTests.ErrorResponse(albumDeleteServerError),
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			// Call mock
			tc.mockBehavior(alu)

			commonTests.DeliveryTestDelete(t, r, "/api/albums/"+tc.albumIDPath+"/",
				tc.expectedStatus, tc.expectedResponse,
				commonTests.WrapRequestWithUserNotNilFunc(tc.user))
		})
	}
}

func TestAlbumDeliveryHTTP_Feed(t *testing.T) {
	// Init
	type mockBehavior func(alu *albumMocks.MockUsecase, aru *artistMocks.MockUsecase)

	c := gomock.NewController(t)

	alu := albumMocks.NewMockUsecase(c)
	aru := artistMocks.NewMockUsecase(c)

	l := commonTests.MockLogger(c)

	h := NewHandler(alu, aru, l)

	// Routing
	r := chi.NewRouter()
	r.Get("/api/albums/feed", h.Feed)

	descriptionID1 := "Антиутопия"
	descriptionID2 := "Крутой альбом от крутого дуета"
	expectedReturnAlbums := []models.Album{
		{
			ID:          1,
			Name:        "Горгород",
			Description: &descriptionID1,
			CoverSrc:    "/albums/covers/gorgorod.png",
		},
		{
			ID:          2,
			Name:        "Стыд или Слава",
			Description: &descriptionID2,
			CoverSrc:    "/albums/covers/shameorglory.png",
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
		{
			ID:        3,
			Name:      "104",
			AvatarSrc: "/artists/avatars/104.png",
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
					"isLiked": false,
					"cover": "/artists/avatars/oxxxymiron.png",
					"listens": 0
				}
			],
			"description": "Антиутопия",
			"isLiked": false,
			"cover": "/albums/covers/gorgorod.png"
		},
		{
			"id": 2,
			"name": "Стыд или Слава",
			"artists": [
				{
					"id": 2,
					"name": "SALUKI",
					"isLiked": false,
					"cover": "/artists/avatars/saluki.png",
					"listens": 0
				},
				{
					"id": 3,
					"name": "104",
					"isLiked": false,
					"cover": "/artists/avatars/104.png",
					"listens": 0
				}
			],
			"description": "Крутой альбом от крутого дуета",
			"isLiked": false,
			"cover": "/albums/covers/shameorglory.png"
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
			mockBehavior: func(alu *albumMocks.MockUsecase, aru *artistMocks.MockUsecase) {
				alu.EXPECT().GetFeed(gomock.Any()).Return(expectedReturnAlbums, nil)
				aru.EXPECT().GetByAlbum(gomock.Any(), expectedReturnAlbums[0].ID).Return(expectedReturnArtists[0:1], nil)
				aru.EXPECT().GetByAlbum(gomock.Any(), expectedReturnAlbums[1].ID).Return(expectedReturnArtists[1:3], nil)
			},
			expectedStatus:   http.StatusOK,
			expectedResponse: correctResponse,
		},
		{
			name: "No Albums",
			mockBehavior: func(alu *albumMocks.MockUsecase, aru *artistMocks.MockUsecase) {
				alu.EXPECT().GetFeed(gomock.Any()).Return([]models.Album{}, nil)
			},
			expectedStatus:   http.StatusOK,
			expectedResponse: `[]`,
		},
		{
			name: "Albums Issues",
			mockBehavior: func(alu *albumMocks.MockUsecase, aru *artistMocks.MockUsecase) {
				alu.EXPECT().GetFeed(gomock.Any()).Return(nil, errors.New(""))
			},
			expectedStatus:   http.StatusInternalServerError,
			expectedResponse: commonTests.ErrorResponse(albumsGetServerError),
		},
		{
			name: "Artists Issues",
			mockBehavior: func(alu *albumMocks.MockUsecase, aru *artistMocks.MockUsecase) {
				alu.EXPECT().GetFeed(gomock.Any()).Return(expectedReturnAlbums, nil)
				aru.EXPECT().GetByAlbum(gomock.Any(), expectedReturnAlbums[0].ID).Return(nil, errors.New(""))
			},
			expectedStatus:   http.StatusInternalServerError,
			expectedResponse: commonTests.ErrorResponse(albumsGetServerError),
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			// Call mock
			tc.mockBehavior(alu, aru)

			commonTests.DeliveryTestGet(t, r, "/api/albums/feed",
				tc.expectedStatus, tc.expectedResponse,
				func(req *http.Request) *http.Request { return req })
		})
	}
}

func TestAlbumDeliveryHTTP_GetFavorite(t *testing.T) {
	type mockBehavior func(alu *albumMocks.MockUsecase, aru *artistMocks.MockUsecase, userID uint32)

	c := gomock.NewController(t)

	alu := albumMocks.NewMockUsecase(c)
	aru := artistMocks.NewMockUsecase(c)

	l := commonTests.MockLogger(c)

	h := NewHandler(alu, aru, l)

	// Routing
	r := chi.NewRouter()
	r.Get("/api/users/{userID}/favorite/albums", h.GetFavorite)

	// Test filling
	const correctUserID uint32 = 1
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
					"isLiked": false,
					"cover": "/artists/avatars/oxxxymiron.png",
					"listens": 0
				}
			],
			"description": "Антиутопия",
			"isLiked": true,
			"cover": "/albums/covers/gorgorod.png"
		},
		{
			"id": 2,
			"name": "Властелин Калек",
			"artists": [
				{
					"id": 2,
					"name": "SALUKI",
					"isLiked": false,
					"cover": "/artists/avatars/saluki.png",
					"listens": 0
				}
			],
			"description": "Стиль",
			"isLiked": true,
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
			user: &correctUser,
			mockBehavior: func(alu *albumMocks.MockUsecase, au *artistMocks.MockUsecase, userID uint32) {
				alu.EXPECT().GetLikedByUser(gomock.Any(), userID).Return(expectedReturnAlbums, nil)
				for ind, album := range expectedReturnAlbums {
					alu.EXPECT().IsLiked(gomock.Any(), album.ID, correctUserID).Return(true, nil)
					au.EXPECT().GetByAlbum(gomock.Any(), album.ID).Return(expectedReturnArtists[ind:ind+1], nil)
					for _, a := range expectedReturnArtists[ind : ind+1] {
						au.EXPECT().IsLiked(gomock.Any(), a.ID, correctUserID).Return(false, nil)
					}
				}
			},
			expectedStatus:   http.StatusOK,
			expectedResponse: correctResponse,
		},
		{
			name: "Albums Issue",
			user: &correctUser,
			mockBehavior: func(alu *albumMocks.MockUsecase, au *artistMocks.MockUsecase, userID uint32) {
				alu.EXPECT().GetLikedByUser(gomock.Any(), userID).Return(nil, errors.New(""))
			},
			expectedStatus:   http.StatusInternalServerError,
			expectedResponse: commonTests.ErrorResponse(albumsGetServerError),
		},
		{
			name: "Artists Issue",
			user: &correctUser,
			mockBehavior: func(alu *albumMocks.MockUsecase, au *artistMocks.MockUsecase, userID uint32) {
				alu.EXPECT().GetLikedByUser(gomock.Any(), userID).Return(expectedReturnAlbums, nil)
				au.EXPECT().GetByAlbum(gomock.Any(), expectedReturnAlbums[0].ID).Return(nil, errors.New(""))
			},
			expectedStatus:   http.StatusInternalServerError,
			expectedResponse: commonTests.ErrorResponse(albumsGetServerError),
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			// Call mock
			tc.mockBehavior(alu, aru, tc.user.ID)

			commonTests.DeliveryTestGet(t, r, "/api/users/"+correctUserIDPath+"/favorite/albums",
				tc.expectedStatus, tc.expectedResponse,
				commonTests.WrapRequestWithUserNotNilFunc(tc.user))
		})
	}
}

func TestAlbumDeliveryHTTP_Like(t *testing.T) {
	// Init
	type mockBehavior func(au *albumMocks.MockUsecase)

	c := gomock.NewController(t)

	alu := albumMocks.NewMockUsecase(c)
	aru := artistMocks.NewMockUsecase(c)

	l := commonTests.MockLogger(c)

	h := NewHandler(alu, aru, l)

	// Routing
	r := chi.NewRouter()
	r.Get("/api/albums/{albumID}/like", h.Like)

	const correctAlbumID uint32 = 1
	correctAlbumIDPath := fmt.Sprint(correctAlbumID)

	testTable := []struct {
		name             string
		albumIDPath      string
		user             *models.User
		mockBehavior     mockBehavior
		expectedStatus   int
		expectedResponse string
	}{
		{
			name:        "Common",
			albumIDPath: correctAlbumIDPath,
			user:        &correctUser,
			mockBehavior: func(au *albumMocks.MockUsecase) {
				au.EXPECT().SetLike(gomock.Any(), correctAlbumID, correctUser.ID).Return(true, nil)
			},
			expectedStatus:   http.StatusOK,
			expectedResponse: commonTests.OKResponse(commonHTTP.LikeSuccess),
		},
		{
			name:        "Already liked (Anyway Success)",
			albumIDPath: correctAlbumIDPath,
			user:        &correctUser,
			mockBehavior: func(au *albumMocks.MockUsecase) {
				au.EXPECT().SetLike(gomock.Any(), correctAlbumID, correctUser.ID).Return(false, nil)
			},
			expectedStatus:   http.StatusOK,
			expectedResponse: commonTests.OKResponse(commonHTTP.LikeAlreadyExists),
		},
		{
			name:             "Incorrect ID In Path",
			albumIDPath:      "0",
			user:             &correctUser,
			mockBehavior:     func(au *albumMocks.MockUsecase) {},
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: commonTests.ErrorResponse(commonHTTP.InvalidURLParameter),
		},
		{
			name:             "No User",
			albumIDPath:      correctAlbumIDPath,
			user:             nil,
			mockBehavior:     func(au *albumMocks.MockUsecase) {},
			expectedStatus:   http.StatusUnauthorized,
			expectedResponse: commonTests.ErrorResponse(commonHTTP.UnathorizedUser),
		},
		{
			name:        "No Album To Like",
			albumIDPath: correctAlbumIDPath,
			user:        &correctUser,
			mockBehavior: func(au *albumMocks.MockUsecase) {
				au.EXPECT().SetLike(
					gomock.Any(), correctAlbumID, correctUser.ID,
				).Return(false, &models.NoSuchAlbumError{})
			},
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: commonTests.ErrorResponse(albumNotFound),
		},
		{
			name:        "Server Error",
			albumIDPath: correctAlbumIDPath,
			user:        &correctUser,
			mockBehavior: func(au *albumMocks.MockUsecase) {
				au.EXPECT().SetLike(
					gomock.Any(), correctAlbumID, correctUser.ID,
				).Return(false, errors.New(""))
			},
			expectedStatus:   http.StatusInternalServerError,
			expectedResponse: commonTests.ErrorResponse(commonHTTP.SetLikeServerError),
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			// Call mock
			tc.mockBehavior(alu)

			commonTests.DeliveryTestGet(t, r, "/api/albums/"+tc.albumIDPath+"/like",
				tc.expectedStatus, tc.expectedResponse,
				commonTests.WrapRequestWithUserNotNilFunc(tc.user))
		})
	}
}

func TestAlbumDeliveryHTTP_UnLike(t *testing.T) {
	// Init
	type mockBehavior func(au *albumMocks.MockUsecase)

	c := gomock.NewController(t)

	alu := albumMocks.NewMockUsecase(c)
	aru := artistMocks.NewMockUsecase(c)

	l := commonTests.MockLogger(c)

	h := NewHandler(alu, aru, l)

	// Routing
	r := chi.NewRouter()
	r.Get("/api/albums/{albumID}/unlike", h.UnLike)

	const correctAlbumID uint32 = 1
	correctAlbumIDPath := fmt.Sprint(correctAlbumID)

	testTable := []struct {
		name             string
		albumIDPath      string
		user             *models.User
		mockBehavior     mockBehavior
		expectedStatus   int
		expectedResponse string
	}{
		{
			name:        "Common",
			albumIDPath: correctAlbumIDPath,
			user:        &correctUser,
			mockBehavior: func(au *albumMocks.MockUsecase) {
				au.EXPECT().UnLike(gomock.Any(), correctAlbumID, correctUser.ID).Return(true, nil)
			},
			expectedStatus:   http.StatusOK,
			expectedResponse: commonTests.OKResponse(commonHTTP.UnLikeSuccess),
		},
		{
			name:        "Wasn't Liked (Anyway Success)",
			albumIDPath: correctAlbumIDPath,
			user:        &correctUser,
			mockBehavior: func(au *albumMocks.MockUsecase) {
				au.EXPECT().UnLike(gomock.Any(), correctAlbumID, correctUser.ID).Return(false, nil)
			},
			expectedStatus:   http.StatusOK,
			expectedResponse: commonTests.OKResponse(commonHTTP.LikeDoesntExist),
		},
		{
			name:             "Incorrect ID In Path",
			albumIDPath:      "0",
			user:             &correctUser,
			mockBehavior:     func(au *albumMocks.MockUsecase) {},
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: commonTests.ErrorResponse(commonHTTP.InvalidURLParameter),
		},
		{
			name:             "No User",
			albumIDPath:      correctAlbumIDPath,
			user:             nil,
			mockBehavior:     func(au *albumMocks.MockUsecase) {},
			expectedStatus:   http.StatusUnauthorized,
			expectedResponse: commonTests.ErrorResponse(commonHTTP.UnathorizedUser),
		},
		{
			name:        "No Album To Unlike",
			albumIDPath: correctAlbumIDPath,
			user:        &correctUser,
			mockBehavior: func(au *albumMocks.MockUsecase) {
				au.EXPECT().UnLike(
					gomock.Any(), correctAlbumID, correctUser.ID,
				).Return(false, &models.NoSuchAlbumError{})
			},
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: commonTests.ErrorResponse(albumNotFound),
		},
		{
			name:        "Server Error",
			albumIDPath: correctAlbumIDPath,
			user:        &correctUser,
			mockBehavior: func(au *albumMocks.MockUsecase) {
				au.EXPECT().UnLike(
					gomock.Any(), correctAlbumID, correctUser.ID,
				).Return(false, errors.New(""))
			},
			expectedStatus:   http.StatusInternalServerError,
			expectedResponse: commonTests.ErrorResponse(commonHTTP.DeleteLikeServerError),
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			// Call mock
			tc.mockBehavior(alu)

			commonTests.DeliveryTestGet(t, r, "/api/albums/"+tc.albumIDPath+"/unlike",
				tc.expectedStatus, tc.expectedResponse,
				commonTests.WrapRequestWithUserNotNilFunc(tc.user))
		})
	}
}
