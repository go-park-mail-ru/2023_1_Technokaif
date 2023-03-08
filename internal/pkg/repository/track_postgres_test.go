package repository

import (
	"testing"

	"github.com/golang/mock/gomock"
	logMocks "github.com/go-park-mail-ru/2023_1_Technokaif/internal/logger/mocks"
	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestTrackPostgresGetFeed(t *testing.T) {
	type mockBehavior func(tf []models.TrackFeed)

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("%v", err)
	}
	defer db.Close()

	c := gomock.NewController(t)
	defer c.Finish()
	l := logMocks.NewMockLogger(c)
	l.EXPECT().Error(gomock.Any()).AnyTimes()
	l.EXPECT().Info(gomock.Any()).AnyTimes()

	if err != nil {
		t.Errorf("%v", err)
	}

	r := NewTrackPostgres(db, l)

	commonTF := []models.TrackFeed{
		{
			ID:   1,
			Name: "LAGG OUT",
			Artists: []models.ArtistFeed{
				{ID: 1, Name: "SALUKI"},
				{ID: 2, Name: "ATL"},
			},
		},
		{
			ID:   2,
			Name: "Another track",
			Artists: []models.ArtistFeed{
				{ID: 1, Name: "SALUKI"},
				{ID: 2, Name: "ATL"},
			},
		},
	}

	testTable := []struct {
		name          string
		actualFeed    []models.TrackFeed
		mockBehavior  mockBehavior
		expectedFeed  []models.TrackFeed
		expectError   bool
		expectedError string
	}{
		{
			name:       "Common",
			actualFeed: commonTF,
			mockBehavior: func(tf []models.TrackFeed) {
				rows := sqlmock.NewRows([]string{"t.id", "t.name", "a.id", "a.name"})

				for _, track := range tf {
					for _, artist := range track.Artists {
						rows.AddRow(track.ID, track.Name, artist.ID, artist.Name)
					}
				}

				mock.ExpectQuery("SELECT (.+) FROM " + TracksTable).
					WillReturnRows(rows)
			},
			expectedFeed: commonTF,
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockBehavior(tc.actualFeed)

			tf, err := r.GetFeed()
			if tc.expectError {
				assert.EqualError(t, err, tc.expectedError)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedFeed, tf)
			}
		})
	}
}
