package http

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/models"
	albumMocks "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/album/mocks"
	artistMocks "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/artist/mocks"
	logMocks "github.com/go-park-mail-ru/2023_1_Technokaif/pkg/logger/mocks"
)

var wrapRequestWithUser = func(r *http.Request, user *models.User) *http.Request {
	if user == nil {
		return r
	}
	ctx := context.WithValue(r.Context(), models.ContextKeyUserType{}, user)
	return r.WithContext(ctx)
}

var correctUser = models.User{
	ID: 1,
}

func TestAlbumDeliveryCreate(t *testing.T) {
	// Init
	type mockBehavior func(au *albumMocks.MockUsecase)

	c := gomock.NewController(t)

	alu := albumMocks.NewMockUsecase(c)
	aru := artistMocks.NewMockUsecase(c)

	l := logMocks.NewMockLogger(c)
	l.EXPECT().Error(gomock.Any()).AnyTimes()
	l.EXPECT().Info(gomock.Any()).AnyTimes()
	l.EXPECT().Errorf(gomock.Any(), gomock.Any()).AnyTimes()
	l.EXPECT().Infof(gomock.Any(), gomock.Any()).AnyTimes()

	h := NewHandler(alu, aru, l)

	// Routing
	r := chi.NewRouter()
	r.Post("/api/albums/", h.Create)

	// Test filling
	correctRequestBody := `{
		"name": "Горгород",
		"artistsID": [1],
		"description": "Антиутопия",
		"cover": "/albums/covers/gorgorod.png"
	}`

	correctArtistsID := []uint32{1}

	description := "Антиутопия"
	expectedCallAlbum := models.Album{
		Name:        "Горгород",
		Description: &description,
		CoverSrc:    "/albums/covers/gorgorod.png",
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
					expectedCallAlbum,
					correctArtistsID,
					correctUser.ID,
				).Return(uint32(1), nil)
			},
			expectedStatus:   200,
			expectedResponse: `{"id": 1}`,
		},
		{
			name:             "No User",
			user:             nil,
			mockBehavior:     func(au *albumMocks.MockUsecase) {},
			expectedStatus:   401,
			expectedResponse: `{"message": "unathorized"}`,
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
			expectedStatus:   400,
			expectedResponse: `{"message": "incorrect input body"}`,
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
			expectedStatus:   400,
			expectedResponse: `{"message": "incorrect input body"}`,
		},
		{
			name:        "User Has No Rights",
			user:        &correctUser,
			requestBody: correctRequestBody,
			mockBehavior: func(au *albumMocks.MockUsecase) {
				au.EXPECT().Create(
					expectedCallAlbum,
					correctArtistsID,
					correctUser.ID,
				).Return(uint32(0), &models.ForbiddenUserError{})
			},
			expectedStatus:   403,
			expectedResponse: `{"message": "no rights to crearte album"}`,
		},
		{
			name:        "Server Error",
			user:        &correctUser,
			requestBody: correctRequestBody,
			mockBehavior: func(au *albumMocks.MockUsecase) {
				au.EXPECT().Create(
					expectedCallAlbum,
					correctArtistsID,
					correctUser.ID,
				).Return(uint32(0), errors.New(""))
			},
			expectedStatus:   500,
			expectedResponse: `{"message": "can't create album"}`,
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			// Call mock
			tc.mockBehavior(alu)

			// Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/api/albums/", bytes.NewBufferString(tc.requestBody))
			r.ServeHTTP(w, wrapRequestWithUser(req, tc.user))

			// Test
			assert.Equal(t, tc.expectedStatus, w.Code)
			assert.JSONEq(t, tc.expectedResponse, w.Body.String())
		})
	}
}

func TestAlbumDeliveryGet(t *testing.T) {
	// Init
	type mockBehavior func(alu *albumMocks.MockUsecase, aru *artistMocks.MockUsecase)

	c := gomock.NewController(t)

	alu := albumMocks.NewMockUsecase(c)
	aru := artistMocks.NewMockUsecase(c)

	l := logMocks.NewMockLogger(c)
	l.EXPECT().Error(gomock.Any()).AnyTimes()
	l.EXPECT().Info(gomock.Any()).AnyTimes()
	l.EXPECT().Errorf(gomock.Any(), gomock.Any()).AnyTimes()
	l.EXPECT().Infof(gomock.Any(), gomock.Any()).AnyTimes()

	h := NewHandler(alu, aru, l)

	// Routing
	r := chi.NewRouter()
	r.Get("/api/albums/{albumID}/", h.Get)

	// Test filling
	correctAlbumID := uint32(1)
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
		},
	}

	correctResponse := `{
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
				alu.EXPECT().GetByID(correctAlbumID).Return(&expectedReturnAlbum, nil)
				aru.EXPECT().GetByAlbum(correctAlbumID).Return(expectedReturnArtists, nil)
			},
			expectedStatus:   200,
			expectedResponse: correctResponse,
		},
		{
			name:             "Incorrect ID In Path",
			albumIDPath:      "incorrect",
			mockBehavior:     func(alu *albumMocks.MockUsecase, aru *artistMocks.MockUsecase) {},
			expectedStatus:   400,
			expectedResponse: `{"message": "invalid url parameter"}`,
		},
		{
			name:             "No User",
			albumIDPath:      correctAlbumIDPath,
			user:             nil,
			mockBehavior:     func(alu *albumMocks.MockUsecase, aru *artistMocks.MockUsecase) {},
			expectedStatus:   401,
			expectedResponse: `{"message": "unathorized"}`,
		},
		{
			name:        "No Album To Get",
			albumIDPath: correctAlbumIDPath,
			user:        &correctUser,
			mockBehavior: func(alu *albumMocks.MockUsecase, aru *artistMocks.MockUsecase) {
				alu.EXPECT().GetByID(correctAlbumID).Return(nil, &models.NoSuchAlbumError{})
			},
			expectedStatus:   400,
			expectedResponse: `{"message": "no such album"}`,
		},
		{
			name:        "Albums Issues",
			albumIDPath: correctAlbumIDPath,
			user:        &correctUser,
			mockBehavior: func(alu *albumMocks.MockUsecase, aru *artistMocks.MockUsecase) {
				alu.EXPECT().GetByID(correctAlbumID).Return(nil, errors.New(""))
			},
			expectedStatus:   500,
			expectedResponse: `{"message": "can't get album"}`,
		},
		{
			name:        "Artists Issues",
			albumIDPath: correctAlbumIDPath,
			user:        &correctUser,
			mockBehavior: func(alu *albumMocks.MockUsecase, aru *artistMocks.MockUsecase) {
				alu.EXPECT().GetByID(correctAlbumID).Return(&expectedReturnAlbum, nil)
				aru.EXPECT().GetByAlbum(correctAlbumID).Return(nil, errors.New(""))
			},
			expectedStatus:   500,
			expectedResponse: `{"message": "can't get album"}`,
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			// Call mock
			tc.mockBehavior(alu, aru)

			// Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/api/albums/"+tc.albumIDPath+"/", nil)
			r.ServeHTTP(w, wrapRequestWithUser(req, tc.user))

			// Test
			assert.Equal(t, tc.expectedStatus, w.Code)
			assert.JSONEq(t, tc.expectedResponse, w.Body.String())
		})
	}
}

func TestAlbumDeliveryDelete(t *testing.T) {
	// Init
	type mockBehavior func(au *albumMocks.MockUsecase)

	c := gomock.NewController(t)

	alu := albumMocks.NewMockUsecase(c)
	aru := artistMocks.NewMockUsecase(c)

	l := logMocks.NewMockLogger(c)
	l.EXPECT().Error(gomock.Any()).AnyTimes()
	l.EXPECT().Info(gomock.Any()).AnyTimes()
	l.EXPECT().Errorf(gomock.Any(), gomock.Any()).AnyTimes()
	l.EXPECT().Infof(gomock.Any(), gomock.Any()).AnyTimes()

	h := NewHandler(alu, aru, l)

	// Routing
	r := chi.NewRouter()
	r.Delete("/api/albums/{albumID}/", h.Delete)

	correctAlbumID := uint32(1)
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
					correctAlbumID,
					correctUser.ID,
				).Return(nil)
			},
			expectedStatus:   200,
			expectedResponse: `{"status": "ok"}`,
		},
		{
			name:             "Incorrect ID In Path",
			albumIDPath:      "incorrect",
			mockBehavior:     func(au *albumMocks.MockUsecase) {},
			expectedStatus:   400,
			expectedResponse: `{"message": "invalid url parameter"}`,
		},
		{
			name:             "No User",
			albumIDPath:      correctAlbumIDPath,
			user:             nil,
			mockBehavior:     func(au *albumMocks.MockUsecase) {},
			expectedStatus:   401,
			expectedResponse: `{"message": "unathorized"}`,
		},
		{
			name:        "User Has No Rights",
			albumIDPath: correctAlbumIDPath,
			user:        &correctUser,
			mockBehavior: func(au *albumMocks.MockUsecase) {
				au.EXPECT().Delete(
					correctAlbumID,
					correctUser.ID,
				).Return(&models.ForbiddenUserError{})
			},
			expectedStatus:   403,
			expectedResponse: `{"message": "no rights to delete album"}`,
		},
		{
			name:        "No Album To Delete",
			albumIDPath: correctAlbumIDPath,
			user:        &correctUser,
			mockBehavior: func(au *albumMocks.MockUsecase) {
				au.EXPECT().Delete(
					correctAlbumID,
					correctUser.ID,
				).Return(&models.NoSuchAlbumError{})
			},
			expectedStatus:   400,
			expectedResponse: `{"message": "no such album"}`,
		},
		{
			name:        "Server Error",
			albumIDPath: correctAlbumIDPath,
			user:        &correctUser,
			mockBehavior: func(au *albumMocks.MockUsecase) {
				au.EXPECT().Delete(
					correctAlbumID,
					correctUser.ID,
				).Return(errors.New(""))
			},
			expectedStatus:   500,
			expectedResponse: `{"message": "can't delete album"}`,
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			// Call mock
			tc.mockBehavior(alu)

			// Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("DELETE", "/api/albums/"+tc.albumIDPath+"/", nil)
			r.ServeHTTP(w, wrapRequestWithUser(req, tc.user))

			// Test
			assert.Equal(t, tc.expectedStatus, w.Code)
			assert.JSONEq(t, tc.expectedResponse, w.Body.String())
		})
	}
}

func TestAlbumDeliveryFeed(t *testing.T) {
	// Init
	type mockBehavior func(alu *albumMocks.MockUsecase, aru *artistMocks.MockUsecase)

	c := gomock.NewController(t)

	alu := albumMocks.NewMockUsecase(c)
	aru := artistMocks.NewMockUsecase(c)

	l := logMocks.NewMockLogger(c)
	l.EXPECT().Error(gomock.Any()).AnyTimes()
	l.EXPECT().Info(gomock.Any()).AnyTimes()
	l.EXPECT().Errorf(gomock.Any(), gomock.Any()).AnyTimes()
	l.EXPECT().Infof(gomock.Any(), gomock.Any()).AnyTimes()

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
					"cover": "/artists/avatars/oxxxymiron.png"
				}
			],
			"description": "Антиутопия",
			"cover": "/albums/covers/gorgorod.png"
		},
		{
			"id": 2,
			"name": "Стыд или Слава",
			"artists": [
				{
					"id": 2,
					"name": "SALUKI",
					"cover": "/artists/avatars/saluki.png"
				},
				{
					"id": 3,
					"name": "104",
					"cover": "/artists/avatars/104.png"
				}
			],
			"description": "Крутой альбом от крутого дуета",
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
				alu.EXPECT().GetFeed().Return(expectedReturnAlbums, nil)
				aru.EXPECT().GetByAlbum(expectedReturnAlbums[0].ID).Return(expectedReturnArtists[0:1], nil)
				aru.EXPECT().GetByAlbum(expectedReturnAlbums[1].ID).Return(expectedReturnArtists[1:3], nil)
			},
			expectedStatus:   200,
			expectedResponse: correctResponse,
		},
		{
			name: "No Albums",
			mockBehavior: func(alu *albumMocks.MockUsecase, aru *artistMocks.MockUsecase) {
				alu.EXPECT().GetFeed().Return([]models.Album{}, nil)
			},
			expectedStatus:   200,
			expectedResponse: `[]`,
		},
		{
			name: "Albums Issues",
			mockBehavior: func(alu *albumMocks.MockUsecase, aru *artistMocks.MockUsecase) {
				alu.EXPECT().GetFeed().Return(nil, errors.New(""))
			},
			expectedStatus:   500,
			expectedResponse: `{"message": "can't get albums"}`,
		},
		{
			name: "Artists Issues",
			mockBehavior: func(alu *albumMocks.MockUsecase, aru *artistMocks.MockUsecase) {
				alu.EXPECT().GetFeed().Return(expectedReturnAlbums, nil)
				aru.EXPECT().GetByAlbum(expectedReturnAlbums[0].ID).Return(nil, errors.New(""))
			},
			expectedStatus:   500,
			expectedResponse: `{"message": "can't get albums"}`,
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			// Call mock
			tc.mockBehavior(alu, aru)

			// Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/api/albums/feed", nil)
			r.ServeHTTP(w, req)

			// Test
			assert.Equal(t, tc.expectedStatus, w.Code)
			assert.JSONEq(t, tc.expectedResponse, w.Body.String())
		})
	}
}

func TestAlbumDeliveryLike(t *testing.T) {
	// Init
	type mockBehavior func(au *albumMocks.MockUsecase)

	c := gomock.NewController(t)

	alu := albumMocks.NewMockUsecase(c)
	aru := artistMocks.NewMockUsecase(c)

	l := logMocks.NewMockLogger(c)
	l.EXPECT().Error(gomock.Any()).AnyTimes()
	l.EXPECT().Info(gomock.Any()).AnyTimes()
	l.EXPECT().Errorf(gomock.Any(), gomock.Any()).AnyTimes()
	l.EXPECT().Infof(gomock.Any(), gomock.Any()).AnyTimes()

	h := NewHandler(alu, aru, l)

	// Routing
	r := chi.NewRouter()
	r.Get("/api/albums/{albumID}/like", h.Like)

	correctAlbumID := uint32(1)
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
				au.EXPECT().SetLike(correctAlbumID, correctUser.ID).Return(true, nil)
			},
			expectedStatus:   200,
			expectedResponse: `{"status": "ok"}`,
		},
		{
			name:        "Already liked (Anyway Success)",
			albumIDPath: correctAlbumIDPath,
			user:        &correctUser,
			mockBehavior: func(au *albumMocks.MockUsecase) {
				au.EXPECT().SetLike(correctAlbumID, correctUser.ID).Return(false, nil)
			},
			expectedStatus:   200,
			expectedResponse: `{"status": "already liked"}`,
		},
		{
			name:             "Incorrect ID In Path",
			albumIDPath:      "0",
			user:             &correctUser,
			mockBehavior:     func(au *albumMocks.MockUsecase) {},
			expectedStatus:   400,
			expectedResponse: `{"message": "invalid url parameter"}`,
		},
		{
			name:             "No User",
			albumIDPath:      correctAlbumIDPath,
			user:             nil,
			mockBehavior:     func(au *albumMocks.MockUsecase) {},
			expectedStatus:   401,
			expectedResponse: `{"message": "unathorized"}`,
		},
		{
			name:        "No Album To Like",
			albumIDPath: correctAlbumIDPath,
			user:        &correctUser,
			mockBehavior: func(au *albumMocks.MockUsecase) {
				au.EXPECT().SetLike(correctAlbumID, correctUser.ID).Return(false, &models.NoSuchAlbumError{})
			},
			expectedStatus:   400,
			expectedResponse: `{"message": "no such album"}`,
		},
		{
			name:        "Server Error",
			albumIDPath: correctAlbumIDPath,
			user:        &correctUser,
			mockBehavior: func(au *albumMocks.MockUsecase) {
				au.EXPECT().SetLike(correctAlbumID, correctUser.ID).Return(false, errors.New(""))
			},
			expectedStatus:   500,
			expectedResponse: `{"message": "can't set like"}`,
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			// Call mock
			tc.mockBehavior(alu)

			// Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/api/albums/"+tc.albumIDPath+"/like", nil)
			r.ServeHTTP(w, wrapRequestWithUser(req, tc.user))

			// Test
			assert.Equal(t, tc.expectedStatus, w.Code)
			assert.JSONEq(t, tc.expectedResponse, w.Body.String())
		})
	}
}

func TestAlbumDeliveryUnLike(t *testing.T) {
	// Init
	type mockBehavior func(au *albumMocks.MockUsecase)

	c := gomock.NewController(t)

	alu := albumMocks.NewMockUsecase(c)
	aru := artistMocks.NewMockUsecase(c)

	l := logMocks.NewMockLogger(c)
	l.EXPECT().Error(gomock.Any()).AnyTimes()
	l.EXPECT().Info(gomock.Any()).AnyTimes()
	l.EXPECT().Errorf(gomock.Any(), gomock.Any()).AnyTimes()
	l.EXPECT().Infof(gomock.Any(), gomock.Any()).AnyTimes()

	h := NewHandler(alu, aru, l)

	// Routing
	r := chi.NewRouter()
	r.Get("/api/albums/{albumID}/like", h.UnLike)

	correctAlbumID := uint32(1)
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
				au.EXPECT().UnLike(correctAlbumID, correctUser.ID).Return(true, nil)
			},
			expectedStatus:   200,
			expectedResponse: `{"status": "ok"}`,
		},
		{
			name:        "Wasn't Liked (Anyway Success)",
			albumIDPath: correctAlbumIDPath,
			user:        &correctUser,
			mockBehavior: func(au *albumMocks.MockUsecase) {
				au.EXPECT().UnLike(correctAlbumID, correctUser.ID).Return(false, nil)
			},
			expectedStatus:   200,
			expectedResponse: `{"status": "wasn't liked"}`,
		},
		{
			name:             "Incorrect ID In Path",
			albumIDPath:      "0",
			user:             &correctUser,
			mockBehavior:     func(au *albumMocks.MockUsecase) {},
			expectedStatus:   400,
			expectedResponse: `{"message": "invalid url parameter"}`,
		},
		{
			name:             "No User",
			albumIDPath:      correctAlbumIDPath,
			user:             nil,
			mockBehavior:     func(au *albumMocks.MockUsecase) {},
			expectedStatus:   401,
			expectedResponse: `{"message": "unathorized"}`,
		},
		{
			name:        "No Album To Unlike",
			albumIDPath: correctAlbumIDPath,
			user:        &correctUser,
			mockBehavior: func(au *albumMocks.MockUsecase) {
				au.EXPECT().UnLike(correctAlbumID, correctUser.ID).Return(false, &models.NoSuchAlbumError{})
			},
			expectedStatus:   400,
			expectedResponse: `{"message": "no such album"}`,
		},
		{
			name:        "Server Error",
			albumIDPath: correctAlbumIDPath,
			user:        &correctUser,
			mockBehavior: func(au *albumMocks.MockUsecase) {
				au.EXPECT().UnLike(correctAlbumID, correctUser.ID).Return(false, errors.New(""))
			},
			expectedStatus:   500,
			expectedResponse: `{"message": "can't remove like"}`,
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			// Call mock
			tc.mockBehavior(alu)

			// Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/api/albums/"+tc.albumIDPath+"/like", nil)
			r.ServeHTTP(w, wrapRequestWithUser(req, tc.user))

			// Test
			assert.Equal(t, tc.expectedStatus, w.Code)
			assert.JSONEq(t, tc.expectedResponse, w.Body.String())
		})
	}
}
