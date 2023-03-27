package http

import (
	"bytes"
	"context"
	"errors"
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

func TestAlbumDeliveryCreate(t *testing.T) {
	type mockBehaviour func(tu *trackMocks.MockUsecase)

	correctUser := models.User{
		ID: 1,
	}

	correctRequestBody := `
	{
		"name": "Хит",
		"artistsID": [1],
		"cover": "/tracks/covers/hit.png",
		"record": "/tracks/records/hit.wav"
	}
	`

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
					[]uint32{1},
					uint32(1),
				).Return(uint32(1), nil)
			},
			expectedStatus: 200,
			expectedResponse: `
				{
					"id": 1
				}
			`,
		},
		{
			name: "Incorrect JSON",
			user: &correctUser,
			requestBody: `
			{
				"name":
				"artistsID": [1],
				"cover": "/tracks/covers/hit.png"
			}
			`,
			mockBehaviour:  func(tu *trackMocks.MockUsecase) {},
			expectedStatus: 400,
			expectedResponse: `
				{
					"message": "incorrect input body"
				}
			`,
		},
		{
			name: "Incorrect body (no source)",
			user: &correctUser,
			requestBody: `
			{
				"name": "Хит",
				"artistsID": [1],
				"cover": "/tracks/covers/hit.png"
			}
			`,
			mockBehaviour:  func(tu *trackMocks.MockUsecase) {},
			expectedStatus: 400,
			expectedResponse: `
				{
					"message": "incorrect input body"
				}
			`,
		},
		{
			name:        "User Has No Rights",
			user:        &correctUser,
			requestBody: correctRequestBody,
			mockBehaviour: func(tu *trackMocks.MockUsecase) {
				tu.EXPECT().Create(
					expectedCallTrack,
					[]uint32{1},
					uint32(1),
				).Return(uint32(0), &models.ForbiddenUserError{})
			},
			expectedStatus: 403,
			expectedResponse: `
				{
					"message": "no rights to create track"
				}
			`,
		},
		{
			name:        "Server Error",
			user:        &correctUser,
			requestBody: correctRequestBody,
			mockBehaviour: func(tu *trackMocks.MockUsecase) {
				tu.EXPECT().Create(
					expectedCallTrack,
					[]uint32{1},
					uint32(1),
				).Return(uint32(0), errors.New(""))
			},
			expectedStatus: 500,
			expectedResponse: `
				{
					"message": "can't create track"
				}
			`,
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
