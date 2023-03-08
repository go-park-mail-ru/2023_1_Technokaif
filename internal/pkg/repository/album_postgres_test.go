package repository

import (
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/logger"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestAlbumPostgresGetFeed(t *testing.T) {
	type mockBehavior func(tf []models.AlbumFeed)

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("%v", err)
	}
	defer db.Close()

	l, err := logger.NewFLogger()
	if err != nil {
		t.Errorf("%v", err)
	}

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

				mock.ExpectQuery("SELECT (.+) FROM " + AlbumsTable).
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
