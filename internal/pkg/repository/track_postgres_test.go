package repository

import (
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	logMocks "github.com/go-park-mail-ru/2023_1_Technokaif/internal/logger/mocks"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/models"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTrackPostgresGetFeed(t *testing.T) {
	type mockBehavior func(tf []models.TrackFeed)

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("%v", err)
	}
	defer db.Close()

	c := gomock.NewController(t)
	defer c.Finish()
	l := logMocks.NewMockLogger(c)
	l.EXPECT().Error(gomock.Any()).AnyTimes()
	l.EXPECT().Info(gomock.Any()).AnyTimes()
	l.EXPECT().Errorf(gomock.Any(), gomock.Any()).AnyTimes()
	l.EXPECT().Infof(gomock.Any(), gomock.Any()).AnyTimes()

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
				rows := sqlmock.NewRows([]string{"t.id", "t.name", "a.id", "a.name", "t.cover_src"})

				for _, track := range tf {
					for _, artist := range track.Artists {
						rows.AddRow(track.ID, track.Name, artist.ID, artist.Name, track.CoverSrc)
					}
				}

				mock.ExpectQuery("SELECT (.+) FROM " + tracksTable).
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
				require.ElementsMatch(t, tc.expectedFeed, tf)
			}
		})
	}
}
