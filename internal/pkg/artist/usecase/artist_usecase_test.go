package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/models"
	artistMocks "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/artist/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

var ctx = context.Background()

func TestArtistUsecase_Create(t *testing.T) {
	type mockBehavior func(ar *artistMocks.MockRepository, artist models.Artist)

	c := gomock.NewController(t)

	au := artistMocks.NewMockRepository(c)

	u := NewUsecase(au)

	var userID uint32 = 1
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
				ar.EXPECT().Insert(ctx, artist).Return(correctArtist.ID, nil)
			},
		},
		{
			name:   "Insert Issue",
			artist: correctArtist,
			mockBehavior: func(ar *artistMocks.MockRepository, artist models.Artist) {
				ar.EXPECT().Insert(ctx, artist).Return(uint32(0), errors.New(""))
			},
			expectError: true,
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockBehavior(au, tc.artist)

			artistID, err := u.Create(ctx, tc.artist)

			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, correctArtist.ID, artistID)
			}
		})
	}
}

func TestArtistUsecase_Delete(t *testing.T) {
	type mockBehavior func(ar *artistMocks.MockRepository, artistID, userID uint32)

	c := gomock.NewController(t)

	au := artistMocks.NewMockRepository(c)

	u := NewUsecase(au)

	var correctArtistID uint32 = 1
	var correctUserID uint32 = 1
	correctArtist := &models.Artist{
		ID:        correctArtistID,
		UserID:    &correctUserID,
		Name:      "Oxxxymiron",
		AvatarSrc: "/artists/avatars/1.png",
	}

	testTable := []struct {
		name             string
		artistID         uint32
		userID           uint32
		mockBehavior     mockBehavior
		expectError      bool
		expectedErrorMsg string
	}{
		{
			name:     "Common",
			artistID: correctArtistID,
			userID:   correctUserID,
			mockBehavior: func(ar *artistMocks.MockRepository, artistID, userID uint32) {
				ar.EXPECT().GetByID(ctx, artistID).Return(correctArtist, nil)
				ar.EXPECT().DeleteByID(ctx, artistID).Return(nil)
			},
		},
		{
			name:     "No Such Artist",
			artistID: correctArtistID,
			userID:   correctUserID,
			mockBehavior: func(ar *artistMocks.MockRepository, artistID, userID uint32) {
				ar.EXPECT().GetByID(ctx, artistID).Return(nil, errors.New(""))
			},
			expectError:      true,
			expectedErrorMsg: "can't find artist",
		},
		{
			name:     "No User",
			artistID: correctArtistID,
			userID:   correctUserID,
			mockBehavior: func(ar *artistMocks.MockRepository, artistID, userID uint32) {
				ar.EXPECT().GetByID(ctx, artistID).Return(&models.Artist{UserID: nil}, nil)
			},
			expectError:      true,
			expectedErrorMsg: "artist can't be deleted",
		},
		{
			name:     "User Has No Rights",
			artistID: correctArtistID,
			userID:   uint32(2),
			mockBehavior: func(ar *artistMocks.MockRepository, artistID, userID uint32) {
				ar.EXPECT().GetByID(ctx, artistID).Return(correctArtist, nil)
			},
			expectError:      true,
			expectedErrorMsg: "artist can't be deleted",
		},
		{
			name:     "Delete Issue",
			artistID: correctArtistID,
			userID:   correctUserID,
			mockBehavior: func(ar *artistMocks.MockRepository, artistID, userID uint32) {
				ar.EXPECT().GetByID(ctx, artistID).Return(correctArtist, nil)
				ar.EXPECT().DeleteByID(ctx, artistID).Return(errors.New(""))
			},
			expectError:      true,
			expectedErrorMsg: "can't delete artist",
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockBehavior(au, tc.artistID, tc.userID)

			err := u.Delete(ctx, tc.artistID, tc.userID)

			if tc.expectError {
				assert.ErrorContains(t, err, tc.expectedErrorMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
