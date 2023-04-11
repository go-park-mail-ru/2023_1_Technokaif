package http

import (
	"bytes"
	"errors"
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/models"

	artistMocks "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/artist/mocks"
	commonTests "github.com/go-park-mail-ru/2023_1_Technokaif/internal/common/tests"
)

var wrapRequestWithUser = commonTests.WrapRequestWithUserNotNil

var correctUser = models.User{
	ID: 1,
}

func TestArtistDeliveryCreate(t *testing.T) {
	// Init
	type mockBehavior func(au *artistMocks.MockUsecase)

	c := gomock.NewController(t)

	au := artistMocks.NewMockUsecase(c)

	l := commonTests.MockLogger(c)

	h := NewHandler(au, l)

	// Routing
	r := chi.NewRouter()
	r.Post("/api/artists/", h.Create)

	// Test filling
	correctRequestBody := `{
		"name": "YARIK",
		"cover": "/artists/covers/yarik.png"
	}`

	expectedCallArtist := models.Artist{
		Name:      "YARIK",
		UserID:    &correctUser.ID,
		AvatarSrc: "/artists/covers/yarik.png",
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
			mockBehavior: func(tu *artistMocks.MockUsecase) {
				tu.EXPECT().Create(
					expectedCallArtist,
				).Return(uint32(1), nil)
			},
			expectedStatus:   200,
			expectedResponse: `{"id": 1}`,
		},
		{
			name:             "No User",
			user:             nil,
			mockBehavior:     func(au *artistMocks.MockUsecase) {},
			expectedStatus:   401,
			expectedResponse: `{"message": "unathorized"}`,
		},
		{
			name: "Incorrect JSON",
			user: &correctUser,
			requestBody: `{
				"name":
				"cover": "/artists/covers/yarik.png"
			}`,
			mockBehavior:     func(tu *artistMocks.MockUsecase) {},
			expectedStatus:   400,
			expectedResponse: `{"message": "incorrect input body"}`,
		},
		{
			name:             "Incorrect body (no cover)",
			user:             &correctUser,
			requestBody:      `{"name": "YARIK"}`,
			mockBehavior:     func(tu *artistMocks.MockUsecase) {},
			expectedStatus:   400,
			expectedResponse: `{"message": "incorrect input body"}`,
		},
		{
			name:        "Server Error",
			user:        &correctUser,
			requestBody: correctRequestBody,
			mockBehavior: func(tu *artistMocks.MockUsecase) {
				tu.EXPECT().Create(
					expectedCallArtist,
				).Return(uint32(0), errors.New(""))
			},
			expectedStatus:   500,
			expectedResponse: `{"message": "can't create artist"}`,
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			// Call mock
			tc.mockBehavior(au)

			// Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/api/artists/", bytes.NewBufferString(tc.requestBody))
			r.ServeHTTP(w, wrapRequestWithUser(req, tc.user))

			// Test
			assert.Equal(t, tc.expectedStatus, w.Code)
			assert.JSONEq(t, tc.expectedResponse, w.Body.String())
		})
	}
}

func TestArtistDeliveryGet(t *testing.T) {
	// Init
	type mockBehavior func(au *artistMocks.MockUsecase)

	c := gomock.NewController(t)

	au := artistMocks.NewMockUsecase(c)

	l := commonTests.MockLogger(c)

	h := NewHandler(au, l)

	// Routing
	r := chi.NewRouter()
	r.Get("/api/artists/{artistID}/", h.Get)

	// Test filling
	correctArtistID := uint32(1)
	correctArtistIDPath := fmt.Sprint(correctArtistID)

	expectedReturnArtist := models.Artist{
		ID:        1,
		Name:      "Oxxxymiron",
		AvatarSrc: "/artists/avatars/oxxxymiron.png",
	}

	correctResponse := `{
		"id": 1,
		"name": "Oxxxymiron",
		"cover": "/artists/avatars/oxxxymiron.png"
	}`

	testTable := []struct {
		name             string
		artistIDPath     string
		user             *models.User
		mockBehavior     mockBehavior
		expectedStatus   int
		expectedResponse string
	}{
		{
			name:         "Common",
			artistIDPath: correctArtistIDPath,
			user:         &correctUser,
			mockBehavior: func(au *artistMocks.MockUsecase) {
				au.EXPECT().GetByID(correctArtistID).Return(&expectedReturnArtist, nil)
			},
			expectedStatus:   200,
			expectedResponse: correctResponse,
		},
		{
			name:             "Incorrect ID In Path",
			artistIDPath:     "0",
			mockBehavior:     func(au *artistMocks.MockUsecase) {},
			expectedStatus:   400,
			expectedResponse: `{"message": "invalid url parameter"}`,
		},
		{
			name:         "No Artist To Get",
			artistIDPath: correctArtistIDPath,
			user:         &correctUser,
			mockBehavior: func(au *artistMocks.MockUsecase) {
				au.EXPECT().GetByID(correctArtistID).Return(nil, &models.NoSuchArtistError{})
			},
			expectedStatus:   400,
			expectedResponse: `{"message": "no such artist"}`,
		},
		{
			name:         "Server Error",
			artistIDPath: correctArtistIDPath,
			user:         &correctUser,
			mockBehavior: func(au *artistMocks.MockUsecase) {
				au.EXPECT().GetByID(correctArtistID).Return(nil, errors.New(""))
			},
			expectedStatus:   500,
			expectedResponse: `{"message": "can't get artist"}`,
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			// Call mock
			tc.mockBehavior(au)

			// Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/api/artists/"+tc.artistIDPath+"/", nil)
			r.ServeHTTP(w, wrapRequestWithUser(req, tc.user))

			// Test
			assert.Equal(t, tc.expectedStatus, w.Code)
			assert.JSONEq(t, tc.expectedResponse, w.Body.String())
		})
	}
}

func TestArtistDeliveryDelete(t *testing.T) {
	// Init
	type mockBehavior func(au *artistMocks.MockUsecase)

	c := gomock.NewController(t)

	au := artistMocks.NewMockUsecase(c)

	l := commonTests.MockLogger(c)

	h := NewHandler(au, l)

	// Routing
	r := chi.NewRouter()
	r.Delete("/api/artists/{artistID}/", h.Delete)

	// Test filing
	correctArtistID := uint32(1)
	correctArtistIDPath := fmt.Sprint(correctArtistID)

	testTable := []struct {
		name             string
		artistIDPath     string
		user             *models.User
		mockBehavior     mockBehavior
		expectedStatus   int
		expectedResponse string
	}{
		{
			name:         "Common",
			artistIDPath: correctArtistIDPath,
			user:         &correctUser,
			mockBehavior: func(au *artistMocks.MockUsecase) {
				au.EXPECT().Delete(
					correctArtistID,
					correctUser.ID,
				).Return(nil)
			},
			expectedStatus:   200,
			expectedResponse: `{"status": "ok"}`,
		},
		{
			name:             "Incorrect ID In Path",
			artistIDPath:     "incorrect",
			mockBehavior:     func(au *artistMocks.MockUsecase) {},
			expectedStatus:   400,
			expectedResponse: `{"message": "invalid url parameter"}`,
		},
		{
			name:             "No User",
			artistIDPath:     correctArtistIDPath,
			user:             nil,
			mockBehavior:     func(au *artistMocks.MockUsecase) {},
			expectedStatus:   401,
			expectedResponse: `{"message": "unathorized"}`,
		},
		{
			name:         "User Has No Rights",
			artistIDPath: correctArtistIDPath,
			user:         &correctUser,
			mockBehavior: func(au *artistMocks.MockUsecase) {
				au.EXPECT().Delete(
					correctArtistID,
					correctUser.ID,
				).Return(&models.ForbiddenUserError{})
			},
			expectedStatus:   403,
			expectedResponse: `{"message": "no rights to delete artist"}`,
		},
		{
			name:         "No Artist To Delete",
			artistIDPath: correctArtistIDPath,
			user:         &correctUser,
			mockBehavior: func(au *artistMocks.MockUsecase) {
				au.EXPECT().Delete(
					correctArtistID,
					correctUser.ID,
				).Return(&models.NoSuchArtistError{})
			},
			expectedStatus:   400,
			expectedResponse: `{"message": "no such artist"}`,
		},
		{
			name:         "Server Error",
			artistIDPath: correctArtistIDPath,
			user:         &correctUser,
			mockBehavior: func(au *artistMocks.MockUsecase) {
				au.EXPECT().Delete(
					correctArtistID,
					correctUser.ID,
				).Return(errors.New(""))
			},
			expectedStatus:   500,
			expectedResponse: `{"message": "can't delete artist"}`,
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			// Call mock
			tc.mockBehavior(au)

			// Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("DELETE", "/api/artists/"+tc.artistIDPath+"/", nil)
			r.ServeHTTP(w, wrapRequestWithUser(req, tc.user))

			// Test
			assert.Equal(t, tc.expectedStatus, w.Code)
			assert.JSONEq(t, tc.expectedResponse, w.Body.String())
		})
	}
}

func TestArtistDeliveryFeed(t *testing.T) {
	// Init
	type mockBehavior func(au *artistMocks.MockUsecase)

	c := gomock.NewController(t)

	au := artistMocks.NewMockUsecase(c)

	l := commonTests.MockLogger(c)

	h := NewHandler(au, l)

	// Routing
	r := chi.NewRouter()
	r.Get("/api/artists/feed", h.Feed)

	// Test filling
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
		{
			ID:        4,
			Name:      "104",
			AvatarSrc: "/artists/avatars/104.png",
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
		},
		{
			"id": 3,
			"name": "ATL",
			"cover": "/artists/avatars/atl.png"
		},
		{
			"id": 4,
			"name": "104",
			"cover": "/artists/avatars/104.png"
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
			mockBehavior: func(au *artistMocks.MockUsecase) {
				au.EXPECT().GetFeed().Return(expectedReturnArtists, nil)
			},
			expectedStatus:   200,
			expectedResponse: correctResponse,
		},
		{
			name: "No Artists",
			mockBehavior: func(au *artistMocks.MockUsecase) {
				au.EXPECT().GetFeed().Return([]models.Artist{}, nil)
			},
			expectedStatus:   200,
			expectedResponse: `[]`,
		},
		{
			name: "Server Error",
			mockBehavior: func(au *artistMocks.MockUsecase) {
				au.EXPECT().GetFeed().Return(nil, errors.New(""))
			},
			expectedStatus:   500,
			expectedResponse: `{"message": "can't get artists"}`,
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			// Call mock
			tc.mockBehavior(au)

			// Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/api/artists/feed", nil)
			r.ServeHTTP(w, req)

			// Test
			assert.Equal(t, tc.expectedStatus, w.Code)
			assert.JSONEq(t, tc.expectedResponse, w.Body.String())
		})
	}
}

func TestArtistDeliveryLike(t *testing.T) {
	// Init
	type mockBehavior func(au *artistMocks.MockUsecase)

	c := gomock.NewController(t)

	au := artistMocks.NewMockUsecase(c)

	l := commonTests.MockLogger(c)

	h := NewHandler(au, l)

	// Routing
	r := chi.NewRouter()
	r.Get("/api/artists/{artistID}/like", h.Like)

	// Test filling
	correctArtistID := uint32(1)
	correctArtistIDPath := fmt.Sprint(correctArtistID)

	testTable := []struct {
		name             string
		artistIDPath     string
		user             *models.User
		mockBehavior     mockBehavior
		expectedStatus   int
		expectedResponse string
	}{
		{
			name:         "Common",
			artistIDPath: correctArtistIDPath,
			user:         &correctUser,
			mockBehavior: func(au *artistMocks.MockUsecase) {
				au.EXPECT().SetLike(correctArtistID, correctUser.ID).Return(true, nil)
			},
			expectedStatus:   200,
			expectedResponse: `{"status": "ok"}`,
		},
		{
			name:         "Already liked (Anyway Success)",
			artistIDPath: correctArtistIDPath,
			user:         &correctUser,
			mockBehavior: func(au *artistMocks.MockUsecase) {
				au.EXPECT().SetLike(correctArtistID, correctUser.ID).Return(false, nil)
			},
			expectedStatus:   200,
			expectedResponse: `{"status": "already liked"}`,
		},
		{
			name:             "Incorrect ID In Path",
			artistIDPath:     "0",
			user:             &correctUser,
			mockBehavior:     func(au *artistMocks.MockUsecase) {},
			expectedStatus:   400,
			expectedResponse: `{"message": "invalid url parameter"}`,
		},
		{
			name:             "No User",
			artistIDPath:     correctArtistIDPath,
			user:             nil,
			mockBehavior:     func(au *artistMocks.MockUsecase) {},
			expectedStatus:   401,
			expectedResponse: `{"message": "unathorized"}`,
		},
		{
			name:         "No Artist To Like",
			artistIDPath: correctArtistIDPath,
			user:         &correctUser,
			mockBehavior: func(au *artistMocks.MockUsecase) {
				au.EXPECT().SetLike(correctArtistID, correctUser.ID).Return(false, &models.NoSuchArtistError{})
			},
			expectedStatus:   400,
			expectedResponse: `{"message": "no such artist"}`,
		},
		{
			name:         "Server Error",
			artistIDPath: correctArtistIDPath,
			user:         &correctUser,
			mockBehavior: func(au *artistMocks.MockUsecase) {
				au.EXPECT().SetLike(correctArtistID, correctUser.ID).Return(false, errors.New(""))
			},
			expectedStatus:   500,
			expectedResponse: `{"message": "can't set like"}`,
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			// Call mock
			tc.mockBehavior(au)

			// Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/api/artists/"+tc.artistIDPath+"/like", nil)
			r.ServeHTTP(w, wrapRequestWithUser(req, tc.user))

			// Test
			assert.Equal(t, tc.expectedStatus, w.Code)
			assert.JSONEq(t, tc.expectedResponse, w.Body.String())
		})
	}
}

func TestArtistDeliveryUnLike(t *testing.T) {
	// Init
	type mockBehavior func(au *artistMocks.MockUsecase)

	c := gomock.NewController(t)

	au := artistMocks.NewMockUsecase(c)

	l := commonTests.MockLogger(c)

	h := NewHandler(au, l)

	// Routing
	r := chi.NewRouter()
	r.Get("/api/artists/{artistID}/like", h.UnLike)

	// Test filling
	correctArtistID := uint32(1)
	correctArtistIDPath := fmt.Sprint(correctArtistID)

	testTable := []struct {
		name             string
		artistIDPath     string
		user             *models.User
		mockBehavior     mockBehavior
		expectedStatus   int
		expectedResponse string
	}{
		{
			name:         "Common",
			artistIDPath: correctArtistIDPath,
			user:         &correctUser,
			mockBehavior: func(au *artistMocks.MockUsecase) {
				au.EXPECT().UnLike(correctArtistID, correctUser.ID).Return(true, nil)
			},
			expectedStatus:   200,
			expectedResponse: `{"status": "ok"}`,
		},
		{
			name:         "Wasn't Liked (Anyway Success)",
			artistIDPath: correctArtistIDPath,
			user:         &correctUser,
			mockBehavior: func(au *artistMocks.MockUsecase) {
				au.EXPECT().UnLike(correctArtistID, correctUser.ID).Return(false, nil)
			},
			expectedStatus:   200,
			expectedResponse: `{"status": "wasn't liked"}`,
		},
		{
			name:             "Incorrect ID In Path",
			artistIDPath:     "0",
			user:             &correctUser,
			mockBehavior:     func(au *artistMocks.MockUsecase) {},
			expectedStatus:   400,
			expectedResponse: `{"message": "invalid url parameter"}`,
		},
		{
			name:             "No User",
			artistIDPath:     correctArtistIDPath,
			user:             nil,
			mockBehavior:     func(au *artistMocks.MockUsecase) {},
			expectedStatus:   401,
			expectedResponse: `{"message": "unathorized"}`,
		},
		{
			name:         "No Artist To Unlike",
			artistIDPath: correctArtistIDPath,
			user:         &correctUser,
			mockBehavior: func(au *artistMocks.MockUsecase) {
				au.EXPECT().UnLike(correctArtistID, correctUser.ID).Return(false, &models.NoSuchArtistError{})
			},
			expectedStatus:   400,
			expectedResponse: `{"message": "no such artist"}`,
		},
		{
			name:         "Server Error",
			artistIDPath: correctArtistIDPath,
			user:         &correctUser,
			mockBehavior: func(au *artistMocks.MockUsecase) {
				au.EXPECT().UnLike(correctArtistID, correctUser.ID).Return(false, errors.New(""))
			},
			expectedStatus:   500,
			expectedResponse: `{"message": "can't remove like"}`,
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			// Call mock
			tc.mockBehavior(au)

			// Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/api/artists/"+tc.artistIDPath+"/like", nil)
			r.ServeHTTP(w, wrapRequestWithUser(req, tc.user))

			// Test
			assert.Equal(t, tc.expectedStatus, w.Code)
			assert.JSONEq(t, tc.expectedResponse, w.Body.String())
		})
	}
}
