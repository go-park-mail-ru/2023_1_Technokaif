package usecase

import (
	"errors"
	"testing"

	commonTests "github.com/go-park-mail-ru/2023_1_Technokaif/internal/common/tests"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/models"
	albumMocks "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/album/mocks"
	artistMocks "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/artist/mocks"
	trackMocks "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/track/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestTrackUsecaseCreate(t *testing.T) {
	type mockBehavior func(tr *trackMocks.MockRepository, ar *artistMocks.MockRepository,
		track models.Track, artistsID []uint32, userID uint32)

	c := gomock.NewController(t)

	tr := trackMocks.NewMockRepository(c)
	arr := artistMocks.NewMockRepository(c)
	alr := albumMocks.NewMockRepository(c)

	l := commonTests.MockLogger(c)

	u := NewUsecase(tr, arr, alr, l)

	var correctUserID uint32 = 1
	correctArtists := []models.Artist{
		{
			ID:        1,
			UserID:    &correctUserID,
			Name:      "Oxxxymiron",
			AvatarSrc: "/artists/avatars/1.png",
		},
	}

	correctTrack := models.Track{
		ID:        1,
		Name:      "Горгород",
		CoverSrc:  "/tracks/covers/1.png",
		RecordSrc: "tracks/records/1.wav",
	}

	testTable := []struct {
		name             string
		album            models.Track
		userID           uint32
		artistsID        []uint32
		mockBehavior     mockBehavior
		expectError      bool
		expectedErrorMsg string
	}{
		{
			name:      "Common",
			album:     correctTrack,
			userID:    correctUserID,
			artistsID: []uint32{1},
			mockBehavior: func(tr *trackMocks.MockRepository, arr *artistMocks.MockRepository,
				track models.Track, artistsID []uint32, userID uint32) {

				for ind, id := range artistsID {
					arr.EXPECT().GetByID(id).Return(&correctArtists[ind], nil)
				}
				tr.EXPECT().Insert(track, artistsID).Return(correctTrack.ID, nil)
			},
		},
		{
			name:      "Forbidden User",
			album:     correctTrack,
			userID:    uint32(2),
			artistsID: []uint32{1},
			mockBehavior: func(tr *trackMocks.MockRepository, arr *artistMocks.MockRepository,
				track models.Track, artistsID []uint32, userID uint32) {

				for ind, id := range artistsID {
					arr.EXPECT().GetByID(id).Return(&correctArtists[ind], nil)
				}
			},
			expectError:      true,
			expectedErrorMsg: "user has no rights",
		},
		{
			name:      "Artist Issue",
			album:     correctTrack,
			userID:    correctUserID,
			artistsID: []uint32{1},
			mockBehavior: func(tr *trackMocks.MockRepository, arr *artistMocks.MockRepository,
				track models.Track, artistsID []uint32, userID uint32) {

				for _, id := range artistsID {
					arr.EXPECT().GetByID(id).Return(nil, errors.New(""))
				}
			},
			expectError:      true,
			expectedErrorMsg: "can't get artist",
		},
		{
			name:      "Insert Issue",
			album:     correctTrack,
			userID:    correctUserID,
			artistsID: []uint32{1},
			mockBehavior: func(tr *trackMocks.MockRepository, arr *artistMocks.MockRepository,
				track models.Track, artistsID []uint32, userID uint32) {

				for ind, id := range artistsID {
					arr.EXPECT().GetByID(id).Return(&correctArtists[ind], nil)
				}
				tr.EXPECT().Insert(track, artistsID).Return(uint32(0), errors.New(""))
			},
			expectError:      true,
			expectedErrorMsg: "can't insert track",
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockBehavior(tr, arr, tc.album, tc.artistsID, tc.userID)

			albumID, err := u.Create(tc.album, tc.artistsID, tc.userID)

			if tc.expectError {
				assert.ErrorContains(t, err, tc.expectedErrorMsg)
			} else {
				assert.Equal(t, albumID, tc.album.ID)
			}
		})
	}
}
