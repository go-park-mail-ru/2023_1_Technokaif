package delivery

import (
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi"
	logMocks "github.com/go-park-mail-ru/2023_1_Technokaif/internal/logger/mocks"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/models"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/usecase"
	mocks "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/usecase/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestDeliveryFeed(t *testing.T) {
	type mockBehavior func(ar *mocks.MockArtist, tr *mocks.MockTrack, al *mocks.MockAlbum)

	testTable := []struct {
		name             string
		mockBehavior     mockBehavior
		expectedStatus   int
		expectedResponse string
	}{
		{
			name: "Common",
			mockBehavior: func(ar *mocks.MockArtist, t *mocks.MockTrack, al *mocks.MockAlbum) {
				ar.EXPECT().GetFeed().Return([]models.ArtistFeed{
					{ID: 1, Name: "SALUKI"},
					{ID: 2, Name: "ATL"},
				}, nil)
				t.EXPECT().GetFeed().Return([]models.TrackFeed{
					{
						ID:   1,
						Name: "LAGG OUT",
						Artists: []models.ArtistFeed{
							{ID: 1, Name: "SALUKI"},
							{ID: 2, Name: "ATL"},
						},
					},
				}, nil)
				al.EXPECT().GetFeed().Return([]models.AlbumFeed{}, nil)
			},
			expectedStatus: 200,
			expectedResponse: `{
				"artists": [
					{"id": 1, "name": "SALUKI"},
					{"id": 2, "name": "ATL"}
				],
				"tracks": [
					{
						"id": 1,
						"name": "LAGG OUT",
						"artists": [
							{"id": 1, "name": "SALUKI"},
							{"id": 2, "name": "ATL"}
						]
					}
				],
				"albums": []
			}`,
		},
		{
			name: "Artists Feed Error",
			mockBehavior: func(ar *mocks.MockArtist, t *mocks.MockTrack, al *mocks.MockAlbum) {
				ar.EXPECT().GetFeed().Return(nil, fmt.Errorf(""))
			},
			expectedStatus:   500,
			expectedResponse: `{"message: "error while getting artists"}`,
		},
		{
			name: "Tracks Feed Error",
			mockBehavior: func(ar *mocks.MockArtist, t *mocks.MockTrack, al *mocks.MockAlbum) {
				ar.EXPECT().GetFeed().Return(nil, nil)
				t.EXPECT().GetFeed().Return(nil, fmt.Errorf(""))
			},
			expectedStatus:   500,
			expectedResponse: `{"message: "error while getting tracks"}`,
		},
		{
			name: "Albums Feed Error",
			mockBehavior: func(ar *mocks.MockArtist, t *mocks.MockTrack, al *mocks.MockAlbum) {
				ar.EXPECT().GetFeed().Return(nil, nil)
				t.EXPECT().GetFeed().Return(nil, nil)
				al.EXPECT().GetFeed().Return(nil, fmt.Errorf(""))
			},
			expectedStatus:   500,
			expectedResponse: `{"message: "error while getting albums"}`,
		},
		// Test Encode Fail ?
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			// Init
			c := gomock.NewController(t)
			defer c.Finish()

			tr := mocks.NewMockTrack(c)
			ar := mocks.NewMockArtist(c)
			al := mocks.NewMockAlbum(c)

			tc.mockBehavior(ar, tr, al)

			u := &usecase.Usecase{
				Artist: ar,
				Track:  tr,
				Album:  al,
			}
			l := logMocks.NewMockLogger(c)
			l.EXPECT().Error(gomock.Any()).AnyTimes()
			l.EXPECT().Info(gomock.Any()).AnyTimes()
			l.EXPECT().Errorf(gomock.Any(), gomock.Any()).AnyTimes()
			l.EXPECT().Infof(gomock.Any(), gomock.Any()).AnyTimes()
			h := NewHandler(u, l)

			// Routing
			r := chi.NewRouter()
			r.Get("/feed", h.feed)

			// Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/feed", nil)
			r.ServeHTTP(w, req)

			var expResp FeedResponse
			var actualResp FeedResponse
			json.Unmarshal([]byte(tc.expectedResponse), &expResp)
			json.Unmarshal(w.Body.Bytes(), &actualResp)

			// Test
			assert.Equal(t, tc.expectedStatus, w.Code)
			assert.Equal(t, expResp, actualResp)
		})
	}
}
