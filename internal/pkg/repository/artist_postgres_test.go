package repository

import (
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	logMocks "github.com/go-park-mail-ru/2023_1_Technokaif/internal/logger/mocks"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/models"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestArtistPostgresGetFeed(t *testing.T) {
	type mockBehavior func(tf []models.ArtistFeed)

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

	r := NewArtistPostgres(db, l)

	commonAF := []models.ArtistFeed{
		{
			ID:        1,
			Name:      "SALUKI",
			AvatarSrc: "/artists/saluki.png",
		},
		{
			ID:        2,
			Name:      "Oxxxymiron",
			AvatarSrc: "/artists/oxxxymiron.jpg",
		},
	}

	testTable := []struct {
		name          string
		actualFeed    []models.ArtistFeed
		mockBehavior  mockBehavior
		expectedFeed  []models.ArtistFeed
		expectError   bool
		expectedError string
	}{
		{
			name:       "Common",
			actualFeed: commonAF,
			mockBehavior: func(af []models.ArtistFeed) {
				rows := sqlmock.NewRows([]string{"id", "name", "cover_src"})

				for _, artist := range af {
					rows.AddRow(artist.ID, artist.Name, artist.AvatarSrc)
				}

				mock.ExpectQuery("SELECT (.+) FROM " + ArtistsTable).
					WillReturnRows(rows)
			},
			expectedFeed: commonAF,
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockBehavior(tc.actualFeed)

			af, err := r.GetFeed()
			if tc.expectError {
				assert.EqualError(t, err, tc.expectedError)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedFeed, af)
			}
		})
	}
}
