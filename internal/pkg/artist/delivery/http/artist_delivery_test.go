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

func TestArtistDeliveryCreate(t *testing.T) {
	type mockBehaviour func(au *artistMocks.MockUsecase)

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
		mockBehaviour    mockBehaviour
		expectedStatus   int
		expectedResponse string
	}{
		{
			name:        "Common",
			user:        &correctUser,
			requestBody: correctRequestBody,
			mockBehaviour: func(tu *artistMocks.MockUsecase) {
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
			mockBehaviour:    func(au *artistMocks.MockUsecase) {},
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
			mockBehaviour:    func(tu *artistMocks.MockUsecase) {},
			expectedStatus:   400,
			expectedResponse: `{"message": "incorrect input body"}`,
		},
		{
			name:             "Incorrect body (no cover)",
			user:             &correctUser,
			requestBody:      `{"name": "YARIK"}`,
			mockBehaviour:    func(tu *artistMocks.MockUsecase) {},
			expectedStatus:   400,
			expectedResponse: `{"message": "incorrect input body"}`,
		},
		{
			name:        "Server Error",
			user:        &correctUser,
			requestBody: correctRequestBody,
			mockBehaviour: func(tu *artistMocks.MockUsecase) {
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
			// Init
			c := gomock.NewController(t)
			defer c.Finish()

			au := artistMocks.NewMockUsecase(c)
			tc.mockBehaviour(au)

			l := logMocks.NewMockLogger(c)
			l.EXPECT().Error(gomock.Any()).AnyTimes()
			l.EXPECT().Info(gomock.Any()).AnyTimes()
			l.EXPECT().Errorf(gomock.Any(), gomock.Any()).AnyTimes()
			l.EXPECT().Infof(gomock.Any(), gomock.Any()).AnyTimes()

			h := NewHandler(au, l)

			// Routing
			r := chi.NewRouter()
			r.Post("/api/artists/", h.Create)

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
	type mockBehaviour func(au *artistMocks.MockUsecase)

	correctArtistID := uint32(1)
	correctArtistIDPath := fmt.Sprint(correctArtistID)

	expectedReturnArtist := models.Artist{
		ID:        1,
		Name:      "Oxxxymiron",
		AvatarSrc: "/avatars/artists/oxxxymiron.png",
	}

	correctResponse := `{
		"id": 1,
		"name": "Oxxxymiron",
		"cover": "/avatars/artists/oxxxymiron.png"
	}`

	testTable := []struct {
		name             string
		artistIDPath     string
		user             *models.User
		mockBehaviour    mockBehaviour
		expectedStatus   int
		expectedResponse string
	}{
		{
			name:         "Common",
			artistIDPath: correctArtistIDPath,
			user:         &correctUser,
			mockBehaviour: func(au *artistMocks.MockUsecase) {
				au.EXPECT().GetByID(correctArtistID).Return(&expectedReturnArtist, nil)
			},
			expectedStatus:   200,
			expectedResponse: correctResponse,
		},
		{
			name:             "Incorrect ID In Path",
			artistIDPath:     "incorrect",
			mockBehaviour:    func(au *artistMocks.MockUsecase) {},
			expectedStatus:   400,
			expectedResponse: `{"message": "invalid url parameter"}`,
		},
		{
			name:             "No User",
			artistIDPath:     correctArtistIDPath,
			user:             nil,
			mockBehaviour:    func(au *artistMocks.MockUsecase) {},
			expectedStatus:   401,
			expectedResponse: `{"message": "unathorized"}`,
		},
		{
			name:         "No Artist To Get",
			artistIDPath: correctArtistIDPath,
			user:         &correctUser,
			mockBehaviour: func(au *artistMocks.MockUsecase) {
				au.EXPECT().GetByID(correctArtistID).Return(nil, &models.NoSuchArtistError{})
			},
			expectedStatus:   400,
			expectedResponse: `{"message": "no such artist"}`,
		},
		{
			name:         "Server Error",
			artistIDPath: correctArtistIDPath,
			user:         &correctUser,
			mockBehaviour: func(au *artistMocks.MockUsecase) {
				au.EXPECT().GetByID(correctArtistID).Return(nil, errors.New(""))
			},
			expectedStatus:   500,
			expectedResponse: `{"message": "can't get artist"}`,
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			// Init
			c := gomock.NewController(t)
			defer c.Finish()

			au := artistMocks.NewMockUsecase(c)
			tc.mockBehaviour(au)

			l := logMocks.NewMockLogger(c)
			l.EXPECT().Error(gomock.Any()).AnyTimes()
			l.EXPECT().Info(gomock.Any()).AnyTimes()
			l.EXPECT().Errorf(gomock.Any(), gomock.Any()).AnyTimes()
			l.EXPECT().Infof(gomock.Any(), gomock.Any()).AnyTimes()

			h := NewHandler(au, l)

			// Routing
			r := chi.NewRouter()
			r.Get("/api/artists/{artistID}/", h.Get)

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
	type mockBehaviour func(au *artistMocks.MockUsecase)

	correctArtistID := uint32(1)
	correctArtistIDPath := fmt.Sprint(correctArtistID)

	testTable := []struct {
		name             string
		artistIDPath     string
		user             *models.User
		mockBehaviour    mockBehaviour
		expectedStatus   int
		expectedResponse string
	}{
		{
			name:         "Common",
			artistIDPath: correctArtistIDPath,
			user:         &correctUser,
			mockBehaviour: func(au *artistMocks.MockUsecase) {
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
			mockBehaviour:    func(au *artistMocks.MockUsecase) {},
			expectedStatus:   400,
			expectedResponse: `{"message": "invalid url parameter"}`,
		},
		{
			name:             "No User",
			artistIDPath:     correctArtistIDPath,
			user:             nil,
			mockBehaviour:    func(au *artistMocks.MockUsecase) {},
			expectedStatus:   401,
			expectedResponse: `{"message": "unathorized"}`,
		},
		{
			name:         "User Has No Rights",
			artistIDPath: correctArtistIDPath,
			user:         &correctUser,
			mockBehaviour: func(au *artistMocks.MockUsecase) {
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
			mockBehaviour: func(au *artistMocks.MockUsecase) {
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
			mockBehaviour: func(au *artistMocks.MockUsecase) {
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
			// Init
			c := gomock.NewController(t)
			defer c.Finish()

			au := artistMocks.NewMockUsecase(c)
			tc.mockBehaviour(au)

			l := logMocks.NewMockLogger(c)
			l.EXPECT().Error(gomock.Any()).AnyTimes()
			l.EXPECT().Info(gomock.Any()).AnyTimes()
			l.EXPECT().Errorf(gomock.Any(), gomock.Any()).AnyTimes()
			l.EXPECT().Infof(gomock.Any(), gomock.Any()).AnyTimes()

			h := NewHandler(au, l)

			// Routing
			r := chi.NewRouter()
			r.Delete("/api/artists/{artistID}/", h.Delete)

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
