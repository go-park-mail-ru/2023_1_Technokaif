package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/models"
	albumMocks "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/album/mocks"
	artistMocks "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/artist/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

var ctx = context.Background()

func TestAlbumUsecase_Create(t *testing.T) {
	type mockBehavior func(alr *albumMocks.MockRepository, arr *artistMocks.MockRepository,
		album models.Album, artistsID []uint32, userID uint32)

	c := gomock.NewController(t)

	alr := albumMocks.NewMockRepository(c)
	arr := artistMocks.NewMockRepository(c)

	u := NewUsecase(alr, arr)

	var correctUserID uint32 = 1
	correctArtists := []models.Artist{
		{
			ID:        1,
			UserID:    &correctUserID,
			Name:      "Oxxxymiron",
			AvatarSrc: "/artists/avatars/1.png",
		},
	}

	correctAlbum := models.Album{
		ID:       1,
		Name:     "Горгород",
		CoverSrc: "/albums/covers/1.png",
	}

	testTable := []struct {
		name             string
		album            models.Album
		userID           uint32
		artistsID        []uint32
		mockBehavior     mockBehavior
		expectError      bool
		expectedErrorMsg string
	}{
		{
			name:      "Common",
			album:     correctAlbum,
			userID:    correctUserID,
			artistsID: []uint32{1},
			mockBehavior: func(alr *albumMocks.MockRepository, arr *artistMocks.MockRepository,
				album models.Album, artistsID []uint32, userID uint32) {

				for ind, id := range artistsID {
					arr.EXPECT().GetByID(ctx, id).Return(&correctArtists[ind], nil)
				}
				alr.EXPECT().Insert(ctx, album, artistsID).Return(correctAlbum.ID, nil)
			},
		},
		{
			name:      "Forbidden User",
			album:     correctAlbum,
			userID:    uint32(2),
			artistsID: []uint32{1},
			mockBehavior: func(alr *albumMocks.MockRepository, arr *artistMocks.MockRepository,
				album models.Album, artistsID []uint32, userID uint32) {

				for ind, id := range artistsID {
					arr.EXPECT().GetByID(ctx, id).Return(&correctArtists[ind], nil)
				}
			},
			expectError:      true,
			expectedErrorMsg: "user has no rights",
		},
		{
			name:      "Artist Issue",
			album:     correctAlbum,
			userID:    correctUserID,
			artistsID: []uint32{1},
			mockBehavior: func(alr *albumMocks.MockRepository, arr *artistMocks.MockRepository,
				album models.Album, artistsID []uint32, userID uint32) {

				for _, id := range artistsID {
					arr.EXPECT().GetByID(ctx, id).Return(nil, errors.New(""))
				}
			},
			expectError:      true,
			expectedErrorMsg: "can't get artist",
		},
		{
			name:      "Insert Issue",
			album:     correctAlbum,
			userID:    correctUserID,
			artistsID: []uint32{1},
			mockBehavior: func(alr *albumMocks.MockRepository, arr *artistMocks.MockRepository,
				album models.Album, artistsID []uint32, userID uint32) {

				for ind, id := range artistsID {
					arr.EXPECT().GetByID(ctx, id).Return(&correctArtists[ind], nil)
				}
				alr.EXPECT().Insert(ctx, album, artistsID).Return(uint32(0), errors.New(""))
			},
			expectError:      true,
			expectedErrorMsg: "can't insert album",
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockBehavior(alr, arr, tc.album, tc.artistsID, tc.userID)

			albumID, err := u.Create(ctx, tc.album, tc.artistsID, tc.userID)

			if tc.expectError {
				assert.ErrorContains(t, err, tc.expectedErrorMsg)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, albumID, tc.album.ID)
			}
		})
	}
}

func TestAlbumUsecase_Delete(t *testing.T) {
	type mockBehavior func(alr *albumMocks.MockRepository,
		arr *artistMocks.MockRepository, albumID, userID uint32)

	c := gomock.NewController(t)

	alr := albumMocks.NewMockRepository(c)
	arr := artistMocks.NewMockRepository(c)

	u := NewUsecase(alr, arr)

	var correctUserID uint32 = 1
	var correctAlbumID uint32 = 1

	correctArtists := []models.Artist{
		{
			ID:        1,
			UserID:    &correctUserID,
			Name:      "Oxxxymiron",
			AvatarSrc: "/artists/avatars/1.png",
		},
	}

	testTable := []struct {
		name             string
		albumID          uint32
		userID           uint32
		mockBehavior     mockBehavior
		expectError      bool
		expectedErrorMsg string
	}{
		{
			name:    "Common",
			albumID: correctAlbumID,
			userID:  correctUserID,
			mockBehavior: func(alr *albumMocks.MockRepository,
				arr *artistMocks.MockRepository, albumID, userID uint32) {

				alr.EXPECT().Check(ctx, albumID).Return(nil)
				arr.EXPECT().GetByAlbum(ctx, albumID).Return(correctArtists, nil)
				alr.EXPECT().DeleteByID(ctx, albumID).Return(nil)
			},
		},
		{
			name:    "No Such Album",
			albumID: correctAlbumID,
			userID:  correctUserID,
			mockBehavior: func(alr *albumMocks.MockRepository,
				arr *artistMocks.MockRepository, albumID, userID uint32) {

				alr.EXPECT().Check(ctx, albumID).Return(errors.New(""))
			},
			expectError:      true,
			expectedErrorMsg: "can't find album",
		},
		{
			name:    "Artists Issue",
			albumID: correctAlbumID,
			userID:  correctUserID,
			mockBehavior: func(alr *albumMocks.MockRepository,
				arr *artistMocks.MockRepository, albumID, userID uint32) {

				alr.EXPECT().Check(ctx, albumID).Return(nil)
				arr.EXPECT().GetByAlbum(ctx, albumID).Return(nil, errors.New(""))
			},
			expectError:      true,
			expectedErrorMsg: "can't get artists",
		},
		{
			name:    "User Has No Rights",
			albumID: correctAlbumID,
			userID:  uint32(2),
			mockBehavior: func(alr *albumMocks.MockRepository,
				arr *artistMocks.MockRepository, albumID, userID uint32) {

				alr.EXPECT().Check(ctx, albumID).Return(nil)
				arr.EXPECT().GetByAlbum(ctx, albumID).Return(correctArtists, nil)
			},
			expectError:      true,
			expectedErrorMsg: "album can't be deleted",
		},
		{
			name:    "Delete Issue",
			albumID: correctAlbumID,
			userID:  correctUserID,
			mockBehavior: func(alr *albumMocks.MockRepository,
				arr *artistMocks.MockRepository, albumID, userID uint32) {

				alr.EXPECT().Check(ctx, albumID).Return(nil)
				arr.EXPECT().GetByAlbum(ctx, albumID).Return(correctArtists, nil)
				alr.EXPECT().DeleteByID(ctx, albumID).Return(errors.New(""))
			},
			expectError:      true,
			expectedErrorMsg: "can't delete album",
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockBehavior(alr, arr, tc.albumID, tc.userID)

			err := u.Delete(ctx, tc.albumID, tc.userID)

			if tc.expectError {
				assert.ErrorContains(t, err, tc.expectedErrorMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
