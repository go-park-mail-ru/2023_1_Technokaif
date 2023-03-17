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

func TestAlbumPostgresGetFeed(t *testing.T) {
	type mockBehavior func(tf []models.AlbumFeed)

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

	r := NewAlbumPostgres(db, l)

	commonAF := []models.AlbumFeed{
		{
			ID:   1,
			Name: "Shame or Glory",
			Artists: []models.ArtistFeed{
				{ID: 1, Name: "SALUKI"},
				{ID: 2, Name: "104"},
			},
			Description: "Cool album of cool artists",
		},
		{
			ID:   2,
			Name: "Lord of the Cripples",
			Artists: []models.ArtistFeed{
				{ID: 1, Name: "SALUKI"},
			},
			CoverSrc: "/albums/lordofthecripples.png",
		},
	}

	testTable := []struct {
		name          string
		actualFeed    []models.AlbumFeed
		mockBehavior  mockBehavior
		expectedFeed  []models.AlbumFeed
		expectError   bool
		expectedError string
	}{
		{
			name:       "Common",
			actualFeed: commonAF,
			mockBehavior: func(af []models.AlbumFeed) {
				rows := sqlmock.NewRows([]string{
					"al.id", "al.name", "al.cover_src", "ar.id", "ar.name", "al.description",
				})

				for _, album := range af {
					for _, artist := range album.Artists {
						rows.AddRow(album.ID, album.Name, album.CoverSrc,
							artist.ID, artist.Name, album.Description)
					}
				}

				mock.ExpectQuery("SELECT (.+) FROM " + albumsTable).
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
				require.ElementsMatch(t, tc.expectedFeed, af)
			}
		})
	}
}
