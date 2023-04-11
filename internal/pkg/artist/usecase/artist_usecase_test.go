package usecase

import (
	"errors"
	"testing"

	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/models"
	artistMocks "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/artist/mocks"
	commonTests "github.com/go-park-mail-ru/2023_1_Technokaif/internal/common/tests"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestArtistUsecaseCreate(t *testing.T) {
	type mockBehavior func(ar *artistMocks.MockRepository, artist models.Artist)

	c := gomock.NewController(t)

	au := artistMocks.NewMockRepository(c)

	l := commonTests.MockLogger(c)

	u := NewUsecase(au, l)

	userID := uint32(1)
	correctArtist := models.Artist{
		ID:        1,
		UserID:    &userID,
		Name:      "Oxxxymiron",
		AvatarSrc: "/artists/avatars/1.png",
	}

	testTable := []struct {
		name         string
		artist       models.Artist
		mockBehavior mockBehavior
		expectError  bool
	}{
		{
			name:   "Common",
			artist: correctArtist,
			mockBehavior: func(ar *artistMocks.MockRepository, artist models.Artist) {
				ar.EXPECT().Insert(artist).Return(correctArtist.ID, nil)
			},
		},
		{
			name:   "Insert Error",
			artist: correctArtist,
			mockBehavior: func(ar *artistMocks.MockRepository, artist models.Artist) {
				ar.EXPECT().Insert(artist).Return(uint32(0), errors.New(""))
			},
			expectError: true,
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockBehavior(au, tc.artist)

			artistID, err := u.Create(tc.artist)

			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.Equal(t, correctArtist.ID, artistID)
			}
		})
	}
}
