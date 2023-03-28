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

	artistMocks "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/artist/mocks"
	trackMocks "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/track/mocks"
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

func TestTrackDeliveryCreate(t *testing.T) {
	type mockBehaviour func(tu *trackMocks.MockUsecase)

	correctRequestBody := `{
		"name": "Хит",
		"artistsID": [1],
		"cover": "/tracks/covers/hit.png",
		"record": "/tracks/records/hit.wav"
	}`

	correctArtistsID := []uint32{1}

	expectedCallTrack := models.Track{
		Name:      "Хит",
		CoverSrc:  "/tracks/covers/hit.png",
		RecordSrc: "/tracks/records/hit.wav",
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
			mockBehaviour: func(tu *trackMocks.MockUsecase) {
				tu.EXPECT().Create(
					expectedCallTrack,
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
			mockBehaviour:    func(tu *trackMocks.MockUsecase) {},
			expectedStatus:   401,
			expectedResponse: `{"message": "unathorized"}`,
		},
		{
			name: "Incorrect JSON",
			user: &correctUser,
			requestBody: `{
				"name":
				"artistsID": [1],
				"cover": "/tracks/covers/hit.png"
			}`,
			mockBehaviour:    func(tu *trackMocks.MockUsecase) {},
			expectedStatus:   400,
			expectedResponse: `{"message": "incorrect input body"}`,
		},
		{
			name: "Incorrect Body (no source)",
			user: &correctUser,
			requestBody: `{
				"name": "Хит",
				"artistsID": [1],
				"cover": "/tracks/covers/hit.png"
			}`,
			mockBehaviour:    func(tu *trackMocks.MockUsecase) {},
			expectedStatus:   400,
			expectedResponse: `{"message": "incorrect input body"}`,
		},
		{
			name: "Incorrect Body (albumID w/o albumPosition)",
			user: &correctUser,
			requestBody: `{
				"name": "Хит",
				"artistsID": [1],
				"albumID": 1,
				"cover": "/tracks/covers/gorgorod.png"
			}`,
			mockBehaviour:    func(au *trackMocks.MockUsecase) {},
			expectedStatus:   400,
			expectedResponse: `{"message": "incorrect input body"}`,
		},
		{
			name:        "User Has No Rights",
			user:        &correctUser,
			requestBody: correctRequestBody,
			mockBehaviour: func(tu *trackMocks.MockUsecase) {
				tu.EXPECT().Create(
					expectedCallTrack,
					correctArtistsID,
					correctUser.ID,
				).Return(uint32(0), &models.ForbiddenUserError{})
			},
			expectedStatus:   403,
			expectedResponse: `{"message": "no rights to create track"}`,
		},
		{
			name:        "Server Error",
			user:        &correctUser,
			requestBody: correctRequestBody,
			mockBehaviour: func(tu *trackMocks.MockUsecase) {
				tu.EXPECT().Create(
					expectedCallTrack,
					correctArtistsID,
					correctUser.ID,
				).Return(uint32(0), errors.New(""))
			},
			expectedStatus:   500,
			expectedResponse: `{"message": "can't create track"}`,
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			// Init
			c := gomock.NewController(t)
			defer c.Finish()

			tu := trackMocks.NewMockUsecase(c)
			au := artistMocks.NewMockUsecase(c)
			tc.mockBehaviour(tu)

			l := logMocks.NewMockLogger(c)
			l.EXPECT().Error(gomock.Any()).AnyTimes()
			l.EXPECT().Info(gomock.Any()).AnyTimes()
			l.EXPECT().Errorf(gomock.Any(), gomock.Any()).AnyTimes()
			l.EXPECT().Infof(gomock.Any(), gomock.Any()).AnyTimes()

			h := NewHandler(tu, au, l)

			// Routing
			r := chi.NewRouter()
			r.Post("/api/tracks/", h.Create)

			// Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/api/tracks/", bytes.NewBufferString(tc.requestBody))
			r.ServeHTTP(w, wrapRequestWithUser(req, tc.user))

			// Test
			assert.Equal(t, tc.expectedStatus, w.Code)
			assert.JSONEq(t, tc.expectedResponse, w.Body.String())
		})
	}
}

func TestTrackDeliveryGet(t *testing.T) {
	type mockBehaviour func(tu *trackMocks.MockUsecase, au *artistMocks.MockUsecase)

	correctTrackID := uint32(1)
	correctTrackIDPath := fmt.Sprint(correctTrackID)

	expectedReturnTrack := models.Track{
		ID:        correctTrackID,
		Name:      "Хит",
		CoverSrc:  "/tracks/covers/hit.png",
		RecordSrc: "/tracks/records/hit.wav",
		Listens:   99999999,
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
		"name": "Хит",
		"artists": [
			{
				"id": 1,
				"name": "Oxxxymiron",
				"cover": "/artists/avatars/oxxxymiron.png"
			}
		],
		"cover": "/tracks/covers/hit.png",
		"record": "/tracks/records/hit.wav",
		"listens": 99999999
	}`

	testTable := []struct {
		name             string
		trackIDPath      string
		user             *models.User
		mockBehaviour    mockBehaviour
		expectedStatus   int
		expectedResponse string
	}{
		{
			name:        "Common",
			trackIDPath: correctTrackIDPath,
			user:        &correctUser,
			mockBehaviour: func(tu *trackMocks.MockUsecase, au *artistMocks.MockUsecase) {
				tu.EXPECT().GetByID(correctTrackID).Return(&expectedReturnTrack, nil)
				au.EXPECT().GetByTrack(correctTrackID).Return(expectedReturnArtists, nil)
			},
			expectedStatus:   200,
			expectedResponse: correctResponse,
		},
		{
			name:             "Incorrect ID In Path",
			trackIDPath:      "-5",
			mockBehaviour:    func(tu *trackMocks.MockUsecase, au *artistMocks.MockUsecase) {},
			expectedStatus:   400,
			expectedResponse: `{"message": "invalid url parameter"}`,
		},
		{
			name:             "No User",
			trackIDPath:      correctTrackIDPath,
			user:             nil,
			mockBehaviour:    func(tu *trackMocks.MockUsecase, au *artistMocks.MockUsecase) {},
			expectedStatus:   401,
			expectedResponse: `{"message": "unathorized"}`,
		},
		{
			name:        "No Track To Get",
			trackIDPath: correctTrackIDPath,
			user:        &correctUser,
			mockBehaviour: func(tu *trackMocks.MockUsecase, au *artistMocks.MockUsecase) {
				tu.EXPECT().GetByID(correctTrackID).Return(nil, &models.NoSuchTrackError{})
			},
			expectedStatus:   400,
			expectedResponse: `{"message": "no such track"}`,
		},
		{
			name:        "Tracks Issues",
			trackIDPath: correctTrackIDPath,
			user:        &correctUser,
			mockBehaviour: func(tu *trackMocks.MockUsecase, au *artistMocks.MockUsecase) {
				tu.EXPECT().GetByID(correctTrackID).Return(nil, errors.New(""))
			},
			expectedStatus:   500,
			expectedResponse: `{"message": "can't get track"}`,
		},
		{
			name:        "Artists Issues",
			trackIDPath: correctTrackIDPath,
			user:        &correctUser,
			mockBehaviour: func(tu *trackMocks.MockUsecase, au *artistMocks.MockUsecase) {
				tu.EXPECT().GetByID(correctTrackID).Return(&expectedReturnTrack, nil)
				au.EXPECT().GetByTrack(correctTrackID).Return(nil, errors.New(""))
			},
			expectedStatus:   500,
			expectedResponse: `{"message": "can't get track"}`,
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			// Init
			c := gomock.NewController(t)
			defer c.Finish()

			tu := trackMocks.NewMockUsecase(c)
			au := artistMocks.NewMockUsecase(c)
			tc.mockBehaviour(tu, au)

			l := logMocks.NewMockLogger(c)
			l.EXPECT().Error(gomock.Any()).AnyTimes()
			l.EXPECT().Info(gomock.Any()).AnyTimes()
			l.EXPECT().Errorf(gomock.Any(), gomock.Any()).AnyTimes()
			l.EXPECT().Infof(gomock.Any(), gomock.Any()).AnyTimes()

			h := NewHandler(tu, au, l)

			// Routing
			r := chi.NewRouter()
			r.Get("/api/tracks/{trackID}/", h.Get)

			// Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/api/tracks/"+tc.trackIDPath+"/", nil)
			r.ServeHTTP(w, wrapRequestWithUser(req, tc.user))

			// Test
			assert.Equal(t, tc.expectedStatus, w.Code)
			assert.JSONEq(t, tc.expectedResponse, w.Body.String())
		})
	}
}

func TestTrackDeliveryDelete(t *testing.T) {
	type mockBehaviour func(au *trackMocks.MockUsecase)

	correctTrackID := uint32(1)
	correctTrackIDPath := fmt.Sprint(correctTrackID)

	testTable := []struct {
		name             string
		trackIDPath      string
		user             *models.User
		mockBehaviour    mockBehaviour
		expectedStatus   int
		expectedResponse string
	}{
		{
			name:        "Common",
			trackIDPath: correctTrackIDPath,
			user:        &correctUser,
			mockBehaviour: func(au *trackMocks.MockUsecase) {
				au.EXPECT().Delete(
					correctTrackID,
					correctUser.ID,
				).Return(nil)
			},
			expectedStatus:   200,
			expectedResponse: `{"status": "ok"}`,
		},
		{
			name:             "Incorrect ID In Path",
			trackIDPath:      "incorrect",
			mockBehaviour:    func(au *trackMocks.MockUsecase) {},
			expectedStatus:   400,
			expectedResponse: `{"message": "invalid url parameter"}`,
		},
		{
			name:             "No User",
			trackIDPath:      correctTrackIDPath,
			user:             nil,
			mockBehaviour:    func(au *trackMocks.MockUsecase) {},
			expectedStatus:   401,
			expectedResponse: `{"message": "unathorized"}`,
		},
		{
			name:        "User Has No Rights",
			trackIDPath: correctTrackIDPath,
			user:        &correctUser,
			mockBehaviour: func(au *trackMocks.MockUsecase) {
				au.EXPECT().Delete(
					correctTrackID,
					correctUser.ID,
				).Return(&models.ForbiddenUserError{})
			},
			expectedStatus:   403,
			expectedResponse: `{"message": "no rights to delete track"}`,
		},
		{
			name:        "No Track To Delete",
			trackIDPath: correctTrackIDPath,
			user:        &correctUser,
			mockBehaviour: func(au *trackMocks.MockUsecase) {
				au.EXPECT().Delete(
					correctTrackID,
					correctUser.ID,
				).Return(&models.NoSuchTrackError{})
			},
			expectedStatus:   400,
			expectedResponse: `{"message": "no such track"}`,
		},
		{
			name:        "Server Error",
			trackIDPath: correctTrackIDPath,
			user:        &correctUser,
			mockBehaviour: func(au *trackMocks.MockUsecase) {
				au.EXPECT().Delete(
					correctTrackID,
					correctUser.ID,
				).Return(errors.New(""))
			},
			expectedStatus:   500,
			expectedResponse: `{"message": "can't delete track"}`,
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			// Init
			c := gomock.NewController(t)
			defer c.Finish()

			tu := trackMocks.NewMockUsecase(c)
			au := artistMocks.NewMockUsecase(c)
			tc.mockBehaviour(tu)

			l := logMocks.NewMockLogger(c)
			l.EXPECT().Error(gomock.Any()).AnyTimes()
			l.EXPECT().Info(gomock.Any()).AnyTimes()
			l.EXPECT().Errorf(gomock.Any(), gomock.Any()).AnyTimes()
			l.EXPECT().Infof(gomock.Any(), gomock.Any()).AnyTimes()

			h := NewHandler(tu, au, l)

			// Routing
			r := chi.NewRouter()
			r.Delete("/api/tracks/{trackID}/", h.Delete)

			// Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("DELETE", "/api/tracks/"+tc.trackIDPath+"/", nil)
			r.ServeHTTP(w, wrapRequestWithUser(req, tc.user))

			// Test
			assert.Equal(t, tc.expectedStatus, w.Code)
			assert.JSONEq(t, tc.expectedResponse, w.Body.String())
		})
	}
}

func TestTrackDeliveryFeed(t *testing.T) {
	type mockBehaviour func(tu *trackMocks.MockUsecase, au *artistMocks.MockUsecase)

	expectedReturnTracks := []models.Track{
		{
			ID:        1,
			Name:      "Накануне",
			CoverSrc:  "/tracks/covers/nakanune.png",
			RecordSrc: "/tracks/records/nakanune.wav",
			Listens:   2700000,
		},
		{
			ID:        2,
			Name:      "LAGG OUT",
			CoverSrc:  "/tracks/covers/laggout.png",
			RecordSrc: "/tracks/records/laggout.wav",
			Listens:   4500000,
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
			Name:      "ATL",
			AvatarSrc: "/artists/avatars/atl.png",
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
					"cover": "/artists/avatars/oxxxymiron.png"
				}
			],
			"cover": "/tracks/covers/nakanune.png",
			"record": "/tracks/records/nakanune.wav",
			"listens": 2700000
		},
		{
			"id": 2,
			"name": "LAGG OUT",
			"artists": [
				{
					"id": 2,
					"name": "SALUKI",
					"cover": "/artists/avatars/saluki.png"
				},
				{
					"id": 3,
					"name": "ATL",
					"cover": "/artists/avatars/atl.png"
				}
			],
			"cover": "/tracks/covers/laggout.png",
			"record": "/tracks/records/laggout.wav",
			"listens": 4500000
		}
	]`

	testTable := []struct {
		name             string
		mockBehaviour    mockBehaviour
		expectedStatus   int
		expectedResponse string
	}{
		{
			name: "Common",
			mockBehaviour: func(tu *trackMocks.MockUsecase, au *artistMocks.MockUsecase) {
				tu.EXPECT().GetFeed().Return(expectedReturnTracks, nil)
				au.EXPECT().GetByTrack(expectedReturnTracks[0].ID).Return(expectedReturnArtists[0:1], nil)
				au.EXPECT().GetByTrack(expectedReturnTracks[1].ID).Return(expectedReturnArtists[1:3], nil)
			},
			expectedStatus:   200,
			expectedResponse: correctResponse,
		},
		{
			name: "No Tracks",
			mockBehaviour: func(tu *trackMocks.MockUsecase, au *artistMocks.MockUsecase) {
				tu.EXPECT().GetFeed().Return([]models.Track{}, nil)
			},
			expectedStatus:   200,
			expectedResponse: `[]`,
		},
		{
			name: "Tracks Issues",
			mockBehaviour: func(tu *trackMocks.MockUsecase, au *artistMocks.MockUsecase) {
				tu.EXPECT().GetFeed().Return(nil, errors.New(""))
			},
			expectedStatus:   500,
			expectedResponse: `{"message": "can't get tracks"}`,
		},
		{
			name: "Artists Issues",
			mockBehaviour: func(tu *trackMocks.MockUsecase, au *artistMocks.MockUsecase) {
				tu.EXPECT().GetFeed().Return(expectedReturnTracks, nil)
				au.EXPECT().GetByTrack(expectedReturnTracks[0].ID).Return(nil, errors.New(""))
			},
			expectedStatus:   500,
			expectedResponse: `{"message": "can't get tracks"}`,
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			// Init
			c := gomock.NewController(t)
			defer c.Finish()

			tu := trackMocks.NewMockUsecase(c)
			au := artistMocks.NewMockUsecase(c)
			tc.mockBehaviour(tu, au)

			l := logMocks.NewMockLogger(c)
			l.EXPECT().Error(gomock.Any()).AnyTimes()
			l.EXPECT().Info(gomock.Any()).AnyTimes()
			l.EXPECT().Errorf(gomock.Any(), gomock.Any()).AnyTimes()
			l.EXPECT().Infof(gomock.Any(), gomock.Any()).AnyTimes()

			h := NewHandler(tu, au, l)

			// Routing
			r := chi.NewRouter()
			r.Get("/api/tracks/feed", h.Feed)

			// Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/api/tracks/feed", nil)
			r.ServeHTTP(w, req)

			// Test
			assert.Equal(t, tc.expectedStatus, w.Code)
			assert.JSONEq(t, tc.expectedResponse, w.Body.String())
		})
	}
}
