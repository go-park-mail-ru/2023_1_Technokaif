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
	type mockBehaviour func(au *albumMocks.MockUsecase)

	correctRequestBody := `{
		"name": "Горгород",
		"artistsID": [1],
		"Description": "Антиутопия",
		"cover": "/covers/albums/gorgorod.png"
	}`

	correctArtistsID := []uint32{1}

	description := "Антиутопия"
	expectedCallAlbum := models.Album{
		Name:        "Горгород",
		Description: &description,
		CoverSrc:    "/covers/albums/gorgorod.png",
	}

	testTable := []struct {
		name             string
		user             *models.User
		requestBody      string
		mockBehaviour    mockBehaviour
		expectedStatus   int
		expectedResponse string
	}{
		{
			name:        "Common",
			user:        &correctUser,
			requestBody: correctRequestBody,
			mockBehaviour: func(au *albumMocks.MockUsecase) {
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
			mockBehaviour:    func(au *albumMocks.MockUsecase) {},
			expectedStatus:   401,
			expectedResponse: `{"message": "unathorized"}`,
		},
		{
			name: "Incorrect JSON",
			user: &correctUser,
			requestBody: `{
				"name": ,
				"artistsID": [1],
				"Description": "Антиутопия",
				"cover": "/covers/albums/gorgorod.png"
			}`,
			mockBehaviour:    func(au *albumMocks.MockUsecase) {},
			expectedStatus:   400,
			expectedResponse: `{"message": "incorrect input body"}`,
		},
		{
			name: "Incorrect body (no name)",
			user: &correctUser,
			requestBody: `{
				"artistsID": [1],
				"Description": "Антиутопия",
				"cover": "/covers/albums/gorgorod.png"
			}`,
			mockBehaviour:    func(au *albumMocks.MockUsecase) {},
			expectedStatus:   400,
			expectedResponse: `{"message": "incorrect input body"}`,
		},
		{
			name:        "User Has No Rights",
			user:        &correctUser,
			requestBody: correctRequestBody,
			mockBehaviour: func(au *albumMocks.MockUsecase) {
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
			mockBehaviour: func(au *albumMocks.MockUsecase) {
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
			// Init
			c := gomock.NewController(t)
			defer c.Finish()

			alu := albumMocks.NewMockUsecase(c)
			aru := artistMocks.NewMockUsecase(c)
			tc.mockBehaviour(alu)

			l := logMocks.NewMockLogger(c)
			l.EXPECT().Error(gomock.Any()).AnyTimes()
			l.EXPECT().Info(gomock.Any()).AnyTimes()
			l.EXPECT().Errorf(gomock.Any(), gomock.Any()).AnyTimes()
			l.EXPECT().Infof(gomock.Any(), gomock.Any()).AnyTimes()

			h := NewHandler(alu, aru, l)

			// Routing
			r := chi.NewRouter()
			r.Post("/api/albums/", h.Create)

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
	type mockBehaviour func(alu *albumMocks.MockUsecase, aru *artistMocks.MockUsecase)

	correctAlbumID := uint32(1)
	correctAlbumIDPath := fmt.Sprint(correctAlbumID)

	description := "Антиутопия"
	expectedReturnAlbum := models.Album{
		ID:          correctAlbumID,
		Name:        "Горгород",
		Description: &description,
		CoverSrc:    "/covers/albums/gorgorod.png",
	}

	expectedReturnArtists := []models.Artist{
		{
			ID:        1,
			Name:      "Oxxxymiron",
			AvatarSrc: "/avatars/artists/oxxxymiron.png",
		},
	}

	correctResponse := `{
		"id": 1,
		"name": "Горгород",
		"artists": [
			{
				"id": 1,
				"name": "Oxxxymiron",
				"cover": "/avatars/artists/oxxxymiron.png"
			}
		],
		"description": "Антиутопия",
		"cover": "/covers/albums/gorgorod.png"
	}`

	testTable := []struct {
		name             string
		albumIDPath      string
		user             *models.User
		mockBehaviour    mockBehaviour
		expectedStatus   int
		expectedResponse string
	}{
		{
			name:        "Common",
			albumIDPath: correctAlbumIDPath,
			user:        &correctUser,
			mockBehaviour: func(alu *albumMocks.MockUsecase, aru *artistMocks.MockUsecase) {
				alu.EXPECT().GetByID(correctAlbumID).Return(&expectedReturnAlbum, nil)
				aru.EXPECT().GetByAlbum(correctAlbumID).Return(expectedReturnArtists, nil)
			},
			expectedStatus:   200,
			expectedResponse: correctResponse,
		},
		{
			name:             "Incorrect ID In Path",
			albumIDPath:      "incorrect",
			mockBehaviour:    func(alu *albumMocks.MockUsecase, aru *artistMocks.MockUsecase) {},
			expectedStatus:   400,
			expectedResponse: `{"message": "invalid url parameter"}`,
		},
		{
			name:             "No User",
			albumIDPath:      correctAlbumIDPath,
			user:             nil,
			mockBehaviour:    func(alu *albumMocks.MockUsecase, aru *artistMocks.MockUsecase) {},
			expectedStatus:   401,
			expectedResponse: `{"message": "unathorized"}`,
		},
		{
			name:        "No Album To Get",
			albumIDPath: correctAlbumIDPath,
			user:        &correctUser,
			mockBehaviour: func(alu *albumMocks.MockUsecase, aru *artistMocks.MockUsecase) {
				alu.EXPECT().GetByID(correctAlbumID).Return(nil, &models.NoSuchAlbumError{})
			},
			expectedStatus:   400,
			expectedResponse: `{"message": "no such album"}`,
		},
		{
			name:        "Albums Issues",
			albumIDPath: correctAlbumIDPath,
			user:        &correctUser,
			mockBehaviour: func(alu *albumMocks.MockUsecase, aru *artistMocks.MockUsecase) {
				alu.EXPECT().GetByID(correctAlbumID).Return(nil, errors.New(""))
			},
			expectedStatus:   500,
			expectedResponse: `{"message": "can't get album"}`,
		},
		{
			name:        "Artists Issues",
			albumIDPath: correctAlbumIDPath,
			user:        &correctUser,
			mockBehaviour: func(alu *albumMocks.MockUsecase, aru *artistMocks.MockUsecase) {
				alu.EXPECT().GetByID(correctAlbumID).Return(&expectedReturnAlbum, nil)
				aru.EXPECT().GetByAlbum(correctAlbumID).Return(nil, errors.New(""))
			},
			expectedStatus:   500,
			expectedResponse: `{"message": "can't get album"}`,
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			// Init
			c := gomock.NewController(t)
			defer c.Finish()

			alu := albumMocks.NewMockUsecase(c)
			aru := artistMocks.NewMockUsecase(c)
			tc.mockBehaviour(alu, aru)

			l := logMocks.NewMockLogger(c)
			l.EXPECT().Error(gomock.Any()).AnyTimes()
			l.EXPECT().Info(gomock.Any()).AnyTimes()
			l.EXPECT().Errorf(gomock.Any(), gomock.Any()).AnyTimes()
			l.EXPECT().Infof(gomock.Any(), gomock.Any()).AnyTimes()

			h := NewHandler(alu, aru, l)

			// Routing
			r := chi.NewRouter()
			r.Get("/api/albums/{albumID}/", h.Get)

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
	type mockBehaviour func(au *albumMocks.MockUsecase)

	correctAlbumID := uint32(1)
	correctAlbumIDPath := fmt.Sprint(correctAlbumID)

	testTable := []struct {
		name             string
		albumIDPath      string
		user             *models.User
		mockBehaviour    mockBehaviour
		expectedStatus   int
		expectedResponse string
	}{
		{
			name:        "Common",
			albumIDPath: correctAlbumIDPath,
			user:        &correctUser,
			mockBehaviour: func(au *albumMocks.MockUsecase) {
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
			mockBehaviour:    func(au *albumMocks.MockUsecase) {},
			expectedStatus:   400,
			expectedResponse: `{"message": "invalid url parameter"}`,
		},
		{
			name:             "No User",
			albumIDPath:      correctAlbumIDPath,
			user:             nil,
			mockBehaviour:    func(au *albumMocks.MockUsecase) {},
			expectedStatus:   401,
			expectedResponse: `{"message": "unathorized"}`,
		},
		{
			name:        "User Has No Rights",
			albumIDPath: correctAlbumIDPath,
			user:        &correctUser,
			mockBehaviour: func(au *albumMocks.MockUsecase) {
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
			mockBehaviour: func(au *albumMocks.MockUsecase) {
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
			mockBehaviour: func(au *albumMocks.MockUsecase) {
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
			// Init
			c := gomock.NewController(t)
			defer c.Finish()

			alu := albumMocks.NewMockUsecase(c)
			aru := artistMocks.NewMockUsecase(c)
			tc.mockBehaviour(alu)

			l := logMocks.NewMockLogger(c)
			l.EXPECT().Error(gomock.Any()).AnyTimes()
			l.EXPECT().Info(gomock.Any()).AnyTimes()
			l.EXPECT().Errorf(gomock.Any(), gomock.Any()).AnyTimes()
			l.EXPECT().Infof(gomock.Any(), gomock.Any()).AnyTimes()

			h := NewHandler(alu, aru, l)

			// Routing
			r := chi.NewRouter()
			r.Delete("/api/albums/{albumID}/", h.Delete)

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
