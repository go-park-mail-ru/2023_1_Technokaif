package usecase

import (
	"context"
	"errors"
	"io"
	"testing"

	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/models"
	playlistMocks "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/playlist/mocks"
	trackMocks "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/track/mocks"
	userMocks "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/user/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

var ctx = context.Background()

type coverSaverMock struct{}

func (cs coverSaverMock) Save(ctx context.Context,
	cover io.Reader, objectName string, size int64) error {
	return nil
}

func TestPlaylistUsecase_Create(t *testing.T) {
	type mockBehavior func(pr *playlistMocks.MockRepository, ur *userMocks.MockRepository,
		playlist models.Playlist, usersID []uint32, userID uint32)

	c := gomock.NewController(t)

	pr := playlistMocks.NewMockRepository(c)
	tr := trackMocks.NewMockRepository(c)
	ur := userMocks.NewMockRepository(c)
	cs := coverSaverMock{}

	u := NewUsecase(pr, tr, ur, cs)

	var correctUserID uint32 = 1
	correctUsers := []models.User{
		{
			ID: correctUserID,
		},
	}

	correctPlaylist := models.Playlist{
		ID:   1,
		Name: "Жара",
	}

	correctUsersID := []uint32{correctUserID}

	testTable := []struct {
		name             string
		playlist         models.Playlist
		userID           uint32
		mockBehavior     mockBehavior
		expectError      bool
		expectedErrorMsg string
	}{
		{
			name:   "Common",
			userID: correctUserID,
			mockBehavior: func(pr *playlistMocks.MockRepository, ur *userMocks.MockRepository,
				playlist models.Playlist, usersID []uint32, userID uint32) {

				for ind, id := range usersID {
					ur.EXPECT().GetByID(ctx, id).Return(&correctUsers[ind], nil)
				}
				pr.EXPECT().Insert(ctx, playlist, usersID).Return(correctPlaylist.ID, nil)
			},
		},
		{
			name:   "Forbidden User",
			userID: uint32(2),
			mockBehavior: func(pr *playlistMocks.MockRepository, ur *userMocks.MockRepository,
				playlist models.Playlist, usersID []uint32, userID uint32) {

				for ind, id := range usersID {
					ur.EXPECT().GetByID(ctx, id).Return(&correctUsers[ind], nil)
				}
			},
			expectError:      true,
			expectedErrorMsg: "playlist can't be created",
		},
		{
			name:   "User Issue",
			userID: correctUserID,
			mockBehavior: func(pr *playlistMocks.MockRepository, ur *userMocks.MockRepository,
				playlist models.Playlist, usersID []uint32, userID uint32) {

				ur.EXPECT().GetByID(ctx, usersID[0]).Return(nil, errors.New(""))
			},
			expectError:      true,
			expectedErrorMsg: "can't get user",
		},
		{
			name:   "Insert Issue",
			userID: correctUserID,
			mockBehavior: func(pr *playlistMocks.MockRepository, ur *userMocks.MockRepository,
				playlist models.Playlist, usersID []uint32, userID uint32) {

				for ind, id := range usersID {
					ur.EXPECT().GetByID(ctx, id).Return(&correctUsers[ind], nil)
				}
				pr.EXPECT().Insert(ctx, playlist, usersID).Return(uint32(0), errors.New(""))
			},
			expectError:      true,
			expectedErrorMsg: "can't insert playlist",
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockBehavior(pr, ur, correctPlaylist, correctUsersID, tc.userID)

			playlistID, err := u.Create(ctx, correctPlaylist, correctUsersID, tc.userID)

			if tc.expectError {
				assert.ErrorContains(t, err, tc.expectedErrorMsg)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, playlistID, correctPlaylist.ID)
			}
		})
	}
}

func TestPlaylistUsecase_Delete(t *testing.T) {
	type mockBehavior func(pr *playlistMocks.MockRepository,
		ur *userMocks.MockRepository, playlistID, userID uint32)

	c := gomock.NewController(t)

	pr := playlistMocks.NewMockRepository(c)
	tr := trackMocks.NewMockRepository(c)
	ur := userMocks.NewMockRepository(c)
	cs := coverSaverMock{}

	u := NewUsecase(pr, tr, ur, cs)

	var correctUserID uint32 = 1
	var correctPlaylistID uint32 = 1

	correctUsers := []models.User{
		{
			ID: 1,
		},
	}

	testTable := []struct {
		name             string
		userID           uint32
		mockBehavior     mockBehavior
		expectError      bool
		expectedErrorMsg string
	}{
		{
			name:   "Common",
			userID: correctUserID,
			mockBehavior: func(pr *playlistMocks.MockRepository,
				ur *userMocks.MockRepository, playlistID, userID uint32) {

				pr.EXPECT().Check(ctx, playlistID).Return(nil)
				ur.EXPECT().GetByPlaylist(ctx, playlistID).Return(correctUsers, nil)
				pr.EXPECT().DeleteByID(ctx, playlistID).Return(nil)
			},
		},
		{
			name:   "No Such Playlist",
			userID: correctUserID,
			mockBehavior: func(pr *playlistMocks.MockRepository,
				ur *userMocks.MockRepository, playlistID, userID uint32) {

				pr.EXPECT().Check(ctx, playlistID).Return(errors.New(""))
			},
			expectError:      true,
			expectedErrorMsg: "can't find playlist",
		},
		{
			name:   "Forbidden User",
			userID: uint32(2),
			mockBehavior: func(pr *playlistMocks.MockRepository,
				ur *userMocks.MockRepository, playlistID, userID uint32) {

				pr.EXPECT().Check(ctx, playlistID).Return(nil)
				ur.EXPECT().GetByPlaylist(ctx, playlistID).Return(correctUsers, nil)
			},
			expectError:      true,
			expectedErrorMsg: "playlist can't be deleted",
		},
		{
			name:   "Users Issue",
			userID: correctUserID,
			mockBehavior: func(pr *playlistMocks.MockRepository,
				ur *userMocks.MockRepository, playlistID, userID uint32) {

				pr.EXPECT().Check(ctx, playlistID).Return(nil)
				ur.EXPECT().GetByPlaylist(ctx, playlistID).Return(nil, errors.New(""))
			},
			expectError:      true,
			expectedErrorMsg: "can't get authors",
		},
		{
			name:   "Delete Issue",
			userID: correctUserID,
			mockBehavior: func(pr *playlistMocks.MockRepository,
				ur *userMocks.MockRepository, playlistID, userID uint32) {

				pr.EXPECT().Check(ctx, playlistID).Return(nil)
				ur.EXPECT().GetByPlaylist(ctx, playlistID).Return(correctUsers, nil)
				pr.EXPECT().DeleteByID(ctx, playlistID).Return(errors.New(""))
			},
			expectError:      true,
			expectedErrorMsg: "can't delete playlist",
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockBehavior(pr, ur, correctPlaylistID, tc.userID)

			err := u.Delete(ctx, correctPlaylistID, tc.userID)

			if tc.expectError {
				assert.ErrorContains(t, err, tc.expectedErrorMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestPlaylistUsecase_AddTrack(t *testing.T) {
	type mockBehavior func(pr *playlistMocks.MockRepository, ur *userMocks.MockRepository,
		tr *trackMocks.MockRepository, playlistID, trackID, userID uint32)

	c := gomock.NewController(t)

	pr := playlistMocks.NewMockRepository(c)
	tr := trackMocks.NewMockRepository(c)
	ur := userMocks.NewMockRepository(c)
	cs := coverSaverMock{}

	u := NewUsecase(pr, tr, ur, cs)

	var correctUserID uint32 = 1
	var correctPlaylistID uint32 = 1
	var correctTrackID uint32 = 1

	correctUsers := []models.User{
		{
			ID: 1,
		},
	}

	testTable := []struct {
		name             string
		userID           uint32
		mockBehavior     mockBehavior
		expectError      bool
		expectedErrorMsg string
	}{
		{
			name:   "Common",
			userID: correctUserID,
			mockBehavior: func(pr *playlistMocks.MockRepository, ur *userMocks.MockRepository,
				tr *trackMocks.MockRepository, playlistID, trackID, userID uint32) {

				pr.EXPECT().Check(ctx, playlistID).Return(nil)
				tr.EXPECT().Check(ctx, trackID).Return(nil)
				ur.EXPECT().GetByPlaylist(ctx, playlistID).Return(correctUsers, nil)
				pr.EXPECT().AddTrack(ctx, trackID, playlistID).Return(nil)
			},
		},
		{
			name:   "No Such Playlist",
			userID: correctUserID,
			mockBehavior: func(pr *playlistMocks.MockRepository, ur *userMocks.MockRepository,
				tr *trackMocks.MockRepository, playlistID, trackID, userID uint32) {

				pr.EXPECT().Check(ctx, playlistID).Return(errors.New(""))
			},
			expectError:      true,
			expectedErrorMsg: "can't find playlist",
		},
		{
			name:   "No Such Track",
			userID: correctUserID,
			mockBehavior: func(pr *playlistMocks.MockRepository, ur *userMocks.MockRepository,
				tr *trackMocks.MockRepository, playlistID, trackID, userID uint32) {

				pr.EXPECT().Check(ctx, playlistID).Return(nil)
				tr.EXPECT().Check(ctx, trackID).Return(errors.New(""))
			},
			expectError:      true,
			expectedErrorMsg: "can't find track",
		},
		{
			name:   "Users Issue",
			userID: correctUserID,
			mockBehavior: func(pr *playlistMocks.MockRepository, ur *userMocks.MockRepository,
				tr *trackMocks.MockRepository, playlistID, trackID, userID uint32) {

				pr.EXPECT().Check(ctx, playlistID).Return(nil)
				tr.EXPECT().Check(ctx, trackID).Return(nil)
				ur.EXPECT().GetByPlaylist(ctx, userID).Return(nil, errors.New(""))
			},
			expectError:      true,
			expectedErrorMsg: "can't get authors",
		},
		{
			name:   "Forbidden User",
			userID: uint32(2),
			mockBehavior: func(pr *playlistMocks.MockRepository, ur *userMocks.MockRepository,
				tr *trackMocks.MockRepository, playlistID, trackID, userID uint32) {

				pr.EXPECT().Check(ctx, playlistID).Return(nil)
				tr.EXPECT().Check(ctx, trackID).Return(nil)
				ur.EXPECT().GetByPlaylist(ctx, playlistID).Return(correctUsers, nil)
			},
			expectError:      true,
			expectedErrorMsg: "playlist can't be updated",
		},
		{
			name:   "Add Track Issue",
			userID: correctUserID,
			mockBehavior: func(pr *playlistMocks.MockRepository, ur *userMocks.MockRepository,
				tr *trackMocks.MockRepository, playlistID, trackID, userID uint32) {

				pr.EXPECT().Check(ctx, playlistID).Return(nil)
				tr.EXPECT().Check(ctx, trackID).Return(nil)
				ur.EXPECT().GetByPlaylist(ctx, playlistID).Return(correctUsers, nil)
				pr.EXPECT().AddTrack(ctx, trackID, playlistID).Return(errors.New(""))
			},
			expectError:      true,
			expectedErrorMsg: "can't add track into playlist",
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockBehavior(pr, ur, tr, correctPlaylistID, correctTrackID, tc.userID)

			err := u.AddTrack(ctx, correctTrackID, correctPlaylistID, tc.userID)

			if tc.expectError {
				assert.ErrorContains(t, err, tc.expectedErrorMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestPlaylistUsecase_DeleteTrack(t *testing.T) {
	type mockBehavior func(pr *playlistMocks.MockRepository, ur *userMocks.MockRepository,
		tr *trackMocks.MockRepository, playlistID, trackID, userID uint32)

	c := gomock.NewController(t)

	pr := playlistMocks.NewMockRepository(c)
	tr := trackMocks.NewMockRepository(c)
	ur := userMocks.NewMockRepository(c)
	cs := coverSaverMock{}

	u := NewUsecase(pr, tr, ur, cs)

	var correctUserID uint32 = 1
	var correctPlaylistID uint32 = 1
	var correctTrackID uint32 = 1

	correctUsers := []models.User{
		{
			ID: 1,
		},
	}

	testTable := []struct {
		name             string
		userID           uint32
		mockBehavior     mockBehavior
		expectError      bool
		expectedErrorMsg string
	}{
		{
			name:   "Common",
			userID: correctUserID,
			mockBehavior: func(pr *playlistMocks.MockRepository, ur *userMocks.MockRepository,
				tr *trackMocks.MockRepository, playlistID, trackID, userID uint32) {

				pr.EXPECT().Check(ctx, playlistID).Return(nil)
				tr.EXPECT().Check(ctx, trackID).Return(nil)
				ur.EXPECT().GetByPlaylist(ctx, playlistID).Return(correctUsers, nil)
				pr.EXPECT().DeleteTrack(ctx, trackID, playlistID).Return(nil)
			},
		},
		{
			name:   "No Such Playlist",
			userID: correctUserID,
			mockBehavior: func(pr *playlistMocks.MockRepository, ur *userMocks.MockRepository,
				tr *trackMocks.MockRepository, playlistID, trackID, userID uint32) {

				pr.EXPECT().Check(ctx, playlistID).Return(errors.New(""))
			},
			expectError:      true,
			expectedErrorMsg: "can't find playlist",
		},
		{
			name:   "No Such Track",
			userID: correctUserID,
			mockBehavior: func(pr *playlistMocks.MockRepository, ur *userMocks.MockRepository,
				tr *trackMocks.MockRepository, playlistID, trackID, userID uint32) {

				pr.EXPECT().Check(ctx, playlistID).Return(nil)
				tr.EXPECT().Check(ctx, trackID).Return(errors.New(""))
			},
			expectError:      true,
			expectedErrorMsg: "can't find track",
		},
		{
			name:   "Users Issue",
			userID: correctUserID,
			mockBehavior: func(pr *playlistMocks.MockRepository, ur *userMocks.MockRepository,
				tr *trackMocks.MockRepository, playlistID, trackID, userID uint32) {

				pr.EXPECT().Check(ctx, playlistID).Return(nil)
				tr.EXPECT().Check(ctx, trackID).Return(nil)
				ur.EXPECT().GetByPlaylist(ctx, userID).Return(nil, errors.New(""))
			},
			expectError:      true,
			expectedErrorMsg: "can't get authors",
		},
		{
			name:   "Forbidden User",
			userID: uint32(2),
			mockBehavior: func(pr *playlistMocks.MockRepository, ur *userMocks.MockRepository,
				tr *trackMocks.MockRepository, playlistID, trackID, userID uint32) {

				pr.EXPECT().Check(ctx, playlistID).Return(nil)
				tr.EXPECT().Check(ctx, trackID).Return(nil)
				ur.EXPECT().GetByPlaylist(ctx, playlistID).Return(correctUsers, nil)
			},
			expectError:      true,
			expectedErrorMsg: "playlist can't be updated",
		},
		{
			name:   "Delete Track Issue",
			userID: correctUserID,
			mockBehavior: func(pr *playlistMocks.MockRepository, ur *userMocks.MockRepository,
				tr *trackMocks.MockRepository, playlistID, trackID, userID uint32) {

				pr.EXPECT().Check(ctx, playlistID).Return(nil)
				tr.EXPECT().Check(ctx, trackID).Return(nil)
				ur.EXPECT().GetByPlaylist(ctx, playlistID).Return(correctUsers, nil)
				pr.EXPECT().DeleteTrack(ctx, trackID, playlistID).Return(errors.New(""))
			},
			expectError:      true,
			expectedErrorMsg: "can't delete track of playlist",
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockBehavior(pr, ur, tr, correctPlaylistID, correctTrackID, tc.userID)

			err := u.DeleteTrack(ctx, correctTrackID, correctPlaylistID, tc.userID)

			if tc.expectError {
				assert.ErrorContains(t, err, tc.expectedErrorMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestPlaylistUsecase_UpdateInfoAndMembers(t *testing.T) {
	type mockBehavior func(pr *playlistMocks.MockRepository, ur *userMocks.MockRepository,
		playlist models.Playlist, userID uint32)

	c := gomock.NewController(t)

	pr := playlistMocks.NewMockRepository(c)
	tr := trackMocks.NewMockRepository(c)
	ur := userMocks.NewMockRepository(c)
	cs := coverSaverMock{}

	u := NewUsecase(pr, tr, ur, cs)

	var correctUserID uint32 = 1
	var newUserID uint32 = 2

	var correctPlaylistID uint32 = 1
	oldPlaylist := &models.Playlist{
		ID:   correctPlaylistID,
		Name: "Жара",
	}
	oldAuthors := []models.User{
		{
			ID: correctUserID,
		},
	}

	correctUpdatedPlaylist := models.Playlist{
		ID:   correctPlaylistID,
		Name: "Прохлада",
	}
	newAuthorsID := []uint32{correctUserID, newUserID}

	testTable := []struct {
		name             string
		updatedPlaylist  models.Playlist
		newUsersID       []uint32
		userID           uint32
		mockBehavior     mockBehavior
		expectError      bool
		expectedErrorMsg string
	}{
		{
			name:            "Common",
			updatedPlaylist: correctUpdatedPlaylist,
			newUsersID:      newAuthorsID,
			userID:          correctUserID,
			mockBehavior: func(pr *playlistMocks.MockRepository, ur *userMocks.MockRepository,
				playlist models.Playlist, userID uint32) {

				pr.EXPECT().GetByID(ctx, playlist.ID).Return(oldPlaylist, nil)
				ur.EXPECT().GetByPlaylist(ctx, playlist.ID).Return(oldAuthors, nil)
				ur.EXPECT().GetByPlaylist(ctx, playlist.ID).Return(oldAuthors, nil)
				pr.EXPECT().UpdateWithMembers(ctx, playlist, []uint32{newUserID}).Return(nil)
			},
		},
		{
			name:            "No Such Playlist",
			updatedPlaylist: correctUpdatedPlaylist,
			newUsersID:      newAuthorsID,
			userID:          correctUserID,
			mockBehavior: func(pr *playlistMocks.MockRepository, ur *userMocks.MockRepository,
				playlist models.Playlist, userID uint32) {

				pr.EXPECT().GetByID(ctx, playlist.ID).Return(nil, errors.New(""))
			},
			expectError:      true,
			expectedErrorMsg: "can't find playlist",
		},
		{
			name:            "Users Issue",
			updatedPlaylist: correctUpdatedPlaylist,
			newUsersID:      newAuthorsID,
			userID:          correctUserID,
			mockBehavior: func(pr *playlistMocks.MockRepository, ur *userMocks.MockRepository,
				playlist models.Playlist, userID uint32) {

				pr.EXPECT().GetByID(ctx, playlist.ID).Return(oldPlaylist, nil)
				ur.EXPECT().GetByPlaylist(ctx, playlist.ID).Return(nil, errors.New(""))
			},
			expectError:      true,
			expectedErrorMsg: "can't get authors",
		},
		{
			name:            "Update Issue",
			updatedPlaylist: correctUpdatedPlaylist,
			newUsersID:      newAuthorsID,
			userID:          correctUserID,
			mockBehavior: func(pr *playlistMocks.MockRepository, ur *userMocks.MockRepository,
				playlist models.Playlist, userID uint32) {

				pr.EXPECT().GetByID(ctx, playlist.ID).Return(oldPlaylist, nil)
				ur.EXPECT().GetByPlaylist(ctx, playlist.ID).Return(oldAuthors, nil)
				ur.EXPECT().GetByPlaylist(ctx, playlist.ID).Return(oldAuthors, nil)
				pr.EXPECT().UpdateWithMembers(ctx, playlist, []uint32{newUserID}).Return(errors.New(""))
			},
			expectError:      true,
			expectedErrorMsg: "can't update playlist",
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockBehavior(pr, ur, tc.updatedPlaylist, tc.userID)

			err := u.UpdateInfoAndMembers(ctx, tc.updatedPlaylist, tc.newUsersID, tc.userID)

			if tc.expectError {
				assert.ErrorContains(t, err, tc.expectedErrorMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
