package postgresql

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/golang/mock/gomock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"

	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/models"
	playlistMocks "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/playlist/mocks"
)

var ctx = context.Background()

const playlistTable = "Playlists"
const likedPlaylistsTable = "Liked_playlists"
const usersPlaylistsTable = "Users_Playlists"
const playlistsTracksTable = "Playlists_Tracks"

var errPqInternal = errors.New("postgres is dead")

func TestPlaylistRepositoryPostgreSQL_Check(t *testing.T) {
	// Init
	type mockBehavior func(playlistID uint32)

	dbMock, sqlxMock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("%v", err)
	}
	defer dbMock.Close()

	c := gomock.NewController(t)

	tablesMock := playlistMocks.NewMockTables(c)

	repo := NewPostgreSQL(sqlx.NewDb(dbMock, "postgres"), tablesMock)

	// Test filling
	const defaultPlaylistToCheckID uint32 = 1

	testTable := []struct {
		name              string
		playlistToCheckID uint32
		mockBehavior      mockBehavior
		expectError       bool
		expectedError     error
	}{
		{
			name:              "Common",
			playlistToCheckID: defaultPlaylistToCheckID,
			mockBehavior: func(playlistID uint32) {
				tablesMock.EXPECT().Playlists().Return(playlistTable)

				row := sqlxMock.NewRows([]string{"exists"}).AddRow(true)
				sqlxMock.ExpectQuery("SELECT EXISTS").
					WithArgs(playlistID).
					WillReturnRows(row)
			},
		},
		{
			name:              "No Such Playlist",
			playlistToCheckID: defaultPlaylistToCheckID,
			mockBehavior: func(playlistID uint32) {
				tablesMock.EXPECT().Playlists().Return(playlistTable)

				row := sqlxMock.NewRows([]string{"exists"}).AddRow(false)
				sqlxMock.ExpectQuery("SELECT EXISTS").
					WithArgs(playlistID).
					WillReturnRows(row)
			},
			expectError:   true,
			expectedError: &models.NoSuchPlaylistError{PlaylistID: defaultPlaylistToCheckID},
		},
		{
			name:              "Internal PostgreSQL Error",
			playlistToCheckID: defaultPlaylistToCheckID,
			mockBehavior: func(playlistID uint32) {
				tablesMock.EXPECT().Playlists().Return(playlistTable)

				sqlxMock.ExpectQuery("SELECT EXISTS").
					WithArgs(playlistID).
					WillReturnError(errPqInternal)
			},
			expectError:   true,
			expectedError: errPqInternal,
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			// Call mock
			tc.mockBehavior(tc.playlistToCheckID)

			err := repo.Check(ctx, tc.playlistToCheckID)

			// Test
			if tc.expectError {
				assert.ErrorContains(t, err, tc.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestPlaylistRepositoryPostgreSQL_Insert(t *testing.T) {
	// Init
	type mockBehavior func(p models.Playlist, usersID []uint32, id uint32)

	dbMock, sqlxMock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("%v", err)
	}
	defer dbMock.Close()

	c := gomock.NewController(t)

	tablesMock := playlistMocks.NewMockTables(c)

	repo := NewPostgreSQL(sqlx.NewDb(dbMock, "postgres"), tablesMock)

	defaultUsersIDToInsert := []uint32{1, 2, 3}

	description := "Ожидайте 3 июня"
	defaultPlaylistToInsert := models.Playlist{
		Name:        "Музыка для эпичной защиты",
		Description: &description,
		CoverSrc:    "/playlists/covers/1.png",
	}

	testTable := []struct {
		name          string
		playlist      models.Playlist
		usersID       []uint32
		mockBehavior  mockBehavior
		expectedID    uint32
		expectError   bool
		expectedError error
	}{
		{
			name:     "Common",
			playlist: defaultPlaylistToInsert,
			usersID:  defaultUsersIDToInsert,
			mockBehavior: func(p models.Playlist, usersID []uint32, id uint32) {
				tablesMock.EXPECT().Playlists().Return(playlistTable)
				tablesMock.EXPECT().UsersPlaylists().Return(usersPlaylistsTable)

				sqlxMock.ExpectBegin()

				row := sqlxMock.NewRows([]string{"id"}).AddRow(id)
				sqlxMock.ExpectQuery("INSERT INTO "+playlistTable).
					WithArgs(p.Name, p.Description, p.CoverSrc).
					WillReturnRows(row)

				for _, userID := range usersID {
					sqlxMock.ExpectExec("INSERT INTO "+usersPlaylistsTable).
						WithArgs(userID, id).
						WillReturnResult(driver.ResultNoRows)
				}

				sqlxMock.ExpectCommit()
			},
			expectedID: 1,
		},
		{
			name:     "Insert Users to Playlist Issue",
			playlist: defaultPlaylistToInsert,
			usersID:  defaultUsersIDToInsert,
			mockBehavior: func(p models.Playlist, usersID []uint32, id uint32) {
				tablesMock.EXPECT().Playlists().Return(playlistTable)
				tablesMock.EXPECT().UsersPlaylists().Return(usersPlaylistsTable)

				sqlxMock.ExpectBegin()

				row := sqlxMock.NewRows([]string{"id"}).AddRow(id)
				sqlxMock.ExpectQuery("INSERT INTO "+playlistTable).
					WithArgs(p.Name, p.Description, p.CoverSrc).
					WillReturnRows(row)

				sqlxMock.ExpectExec("INSERT INTO "+usersPlaylistsTable).
					WithArgs(usersID[0], id).
					WillReturnError(errPqInternal)

				sqlxMock.ExpectRollback()
			},
			expectedID:    1,
			expectError:   true,
			expectedError: errPqInternal,
		},
		{
			name:     "Insert Playlist Issue",
			playlist: defaultPlaylistToInsert,
			usersID:  defaultUsersIDToInsert,
			mockBehavior: func(p models.Playlist, usersID []uint32, id uint32) {
				tablesMock.EXPECT().Playlists().Return(playlistTable)

				sqlxMock.ExpectBegin()

				sqlxMock.ExpectQuery("INSERT INTO "+playlistTable).
					WithArgs(p.Name, p.Description, p.CoverSrc).
					WillReturnError(errPqInternal)

				sqlxMock.ExpectRollback()
			},
			expectError:   true,
			expectedError: errPqInternal,
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			// Call mock
			tc.mockBehavior(tc.playlist, tc.usersID, tc.expectedID)

			id, err := repo.Insert(ctx, tc.playlist, tc.usersID)

			// Test
			if tc.expectError {
				assert.ErrorIs(t, err, tc.expectedError)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, id, tc.expectedID)
			}
		})
	}
}

func TestPlaylistRepositoryPostgreSQL_GetByID(t *testing.T) {
	// Init
	type mockBehavior func(playlistID uint32, p models.Playlist)

	dbMock, sqlxMock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("%v", err)
	}
	defer dbMock.Close()

	c := gomock.NewController(t)

	tablesMock := playlistMocks.NewMockTables(c)

	repo := NewPostgreSQL(sqlx.NewDb(dbMock, "postgres"), tablesMock)

	// Test filling
	const defaultPlaylistToGetID uint32 = 1

	description := "Ожидайте 3 июня"
	defaultPlaylist := models.Playlist{
		ID:          defaultPlaylistToGetID,
		Name:        "Музыка для эпичной защиты",
		Description: &description,
		CoverSrc:    "/playlists/covers/epic.png",
	}

	testTable := []struct {
		name             string
		playlistToGetID  uint32
		mockBehavior     mockBehavior
		expectedPlaylist models.Playlist
		expectError      bool
		expectedError    error
	}{
		{
			name:            "Common",
			playlistToGetID: defaultPlaylistToGetID,
			mockBehavior: func(playlistID uint32, p models.Playlist) {
				tablesMock.EXPECT().Playlists().Return(playlistTable)

				row := sqlxMock.NewRows([]string{"id", "name", "description", "cover_src"}).
					AddRow(p.ID, p.Name, p.Description, p.CoverSrc)
				sqlxMock.ExpectQuery("SELECT (.+) FROM " + playlistTable).
					WithArgs(playlistID).
					WillReturnRows(row)
			},
			expectedPlaylist: defaultPlaylist,
		},
		{
			name:            "No Such Playlist",
			playlistToGetID: defaultPlaylistToGetID,
			mockBehavior: func(playlistID uint32, p models.Playlist) {
				tablesMock.EXPECT().Playlists().Return(playlistTable)

				sqlxMock.ExpectQuery("SELECT (.+) FROM " + playlistTable).
					WithArgs(playlistID).
					WillReturnError(sql.ErrNoRows)
			},
			expectError:   true,
			expectedError: &models.NoSuchPlaylistError{PlaylistID: defaultPlaylistToGetID},
		},
		{
			name:            "Internal PostgreSQL Error",
			playlistToGetID: defaultPlaylistToGetID,
			mockBehavior: func(playlistID uint32, p models.Playlist) {
				tablesMock.EXPECT().Playlists().Return(playlistTable)

				sqlxMock.ExpectQuery("SELECT (.+) FROM " + playlistTable).
					WithArgs(playlistID).
					WillReturnError(errPqInternal)
			},
			expectError:   true,
			expectedError: errPqInternal,
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			// Call mock
			tc.mockBehavior(tc.playlistToGetID, tc.expectedPlaylist)

			a, err := repo.GetByID(ctx, tc.playlistToGetID)

			// Test
			if tc.expectError {
				assert.ErrorContains(t, err, tc.expectedError.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedPlaylist, *a)
			}
		})
	}
}

func TestPlaylistRepositoryPostgreSQL_DeleteByID(t *testing.T) {
	// Init
	type mockBehavior func(playlistID uint32)

	dbMock, sqlxMock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("%v", err)
	}
	defer dbMock.Close()

	c := gomock.NewController(t)

	tablesMock := playlistMocks.NewMockTables(c)

	repo := NewPostgreSQL(sqlx.NewDb(dbMock, "postgres"), tablesMock)

	// Test filling
	const defaultPlaylistToDeleteID uint32 = 1

	testTable := []struct {
		name               string
		playlistToDeleteID uint32
		mockBehavior       mockBehavior
		expectError        bool
		expectedError      error
	}{
		{
			name:               "Common",
			playlistToDeleteID: defaultPlaylistToDeleteID,
			mockBehavior: func(playlistID uint32) {
				tablesMock.EXPECT().Playlists().Return(playlistTable)

				sqlxMock.ExpectExec("DELETE FROM " + playlistTable).
					WithArgs(playlistID).
					WillReturnResult(driver.RowsAffected(1))
			},
		},
		{
			name:               "No Such Playlist",
			playlistToDeleteID: defaultPlaylistToDeleteID,
			mockBehavior: func(playlistID uint32) {
				tablesMock.EXPECT().Playlists().Return(playlistTable)

				sqlxMock.ExpectExec("DELETE FROM " + playlistTable).
					WithArgs(playlistID).
					WillReturnResult(driver.RowsAffected(0))
			},
			expectError:   true,
			expectedError: &models.NoSuchPlaylistError{PlaylistID: defaultPlaylistToDeleteID},
		},
		{
			name:               "Internal PostgreSQL Error",
			playlistToDeleteID: defaultPlaylistToDeleteID,
			mockBehavior: func(playlistID uint32) {
				tablesMock.EXPECT().Playlists().Return(playlistTable)

				sqlxMock.ExpectExec("DELETE FROM " + playlistTable).
					WithArgs(playlistID).
					WillReturnError(errPqInternal)
			},
			expectError:   true,
			expectedError: errPqInternal,
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			// Call mock
			tc.mockBehavior(tc.playlistToDeleteID)

			err := repo.DeleteByID(ctx, tc.playlistToDeleteID)

			// Test
			if tc.expectError {
				assert.ErrorContains(t, err, tc.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestPlaylistRepositoryPostgreSQL_AddTrack(t *testing.T) {
	// Init
	type mockBehavior func(trackID, playlistID uint32)

	dbMock, sqlxMock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("%v", err)
	}
	defer dbMock.Close()

	c := gomock.NewController(t)

	tablesMock := playlistMocks.NewMockTables(c)

	repo := NewPostgreSQL(sqlx.NewDb(dbMock, "postgres"), tablesMock)

	var defaultTrackToInsertID uint32 = 1
	var defaultPlaylistID uint32 = 1

	testTable := []struct {
		name          string
		playlistID    uint32
		trackID       uint32
		mockBehavior  mockBehavior
		expectError   bool
		expectedError error
	}{
		{
			name:       "Common",
			playlistID: defaultPlaylistID,
			trackID:    defaultTrackToInsertID,
			mockBehavior: func(trackID, playlistID uint32) {
				tablesMock.EXPECT().PlaylistsTracks().Return(playlistsTracksTable)

				sqlxMock.ExpectExec("INSERT INTO "+playlistsTracksTable).
					WithArgs(trackID, playlistID).
					WillReturnResult(driver.ResultNoRows)
			},
		},
		{
			name:       "No Such Playlist Issue",
			playlistID: defaultPlaylistID,
			trackID:    defaultTrackToInsertID,
			mockBehavior: func(trackID, playlistID uint32) {
				tablesMock.EXPECT().PlaylistsTracks().Return(playlistsTracksTable)

				sqlxMock.ExpectExec("INSERT INTO "+playlistTable).
					WithArgs(trackID, playlistID).
					WillReturnError(sql.ErrNoRows)
			},
			expectError:   true,
			expectedError: &models.NoSuchPlaylistError{PlaylistID: defaultPlaylistID},
		},
		{
			name:       "Insert Track Into Playlist Issue",
			playlistID: defaultPlaylistID,
			trackID:    defaultTrackToInsertID,
			mockBehavior: func(trackID, playlistID uint32) {
				tablesMock.EXPECT().PlaylistsTracks().Return(playlistsTracksTable)

				sqlxMock.ExpectExec("INSERT INTO "+playlistsTracksTable).
					WithArgs(trackID, playlistID).
					WillReturnError(errPqInternal)
			},
			expectError:   true,
			expectedError: errPqInternal,
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			// Call mock
			tc.mockBehavior(tc.trackID, tc.playlistID)

			err := repo.AddTrack(ctx, tc.trackID, tc.playlistID)

			// Test
			if tc.expectError {
				assert.ErrorContains(t, err, tc.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestPlaylistRepositoryPostgreSQL_GetFeed(t *testing.T) {
	// Init
	type mockBehavior func(playlists []models.Playlist)

	dbMock, sqlxMock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("%v", err)
	}
	defer dbMock.Close()

	c := gomock.NewController(t)

	tablesMock := playlistMocks.NewMockTables(c)

	repo := NewPostgreSQL(sqlx.NewDb(dbMock, "postgres"), tablesMock)

	// Test filling
	descriptionID1 := "Ожидайте 3 июня"
	descriptionID2 := "Если вдруг решил отдохнуть"
	defaultPlaylists := []models.Playlist{
		{
			ID:          1,
			Name:        "Музыка для эпичной защиты",
			Description: &descriptionID1,
			CoverSrc:    "/playlists/covers/epic.png",
		},
		{
			ID:          2,
			Name:        "Для чилла",
			Description: &descriptionID2,
			CoverSrc:    "/playlists/covers/chill.png",
		},
	}

	testTable := []struct {
		name              string
		mockBehavior      mockBehavior
		expectedPlaylists []models.Playlist
		expectError       bool
		expectedError     error
	}{
		{
			name: "Common",
			mockBehavior: func(p []models.Playlist) {
				tablesMock.EXPECT().Playlists().Return(playlistTable)

				rows := sqlxMock.NewRows([]string{"id", "name", "description", "cover_src"})
				for ind := range p {
					rows.AddRow(p[ind].ID, p[ind].Name, p[ind].Description, p[ind].CoverSrc)
				}
				sqlxMock.ExpectQuery("SELECT (.+) FROM " + playlistTable).
					WillReturnRows(rows)
			},
			expectedPlaylists: defaultPlaylists,
		},
		{
			name: "Internal PostgreSQL Error",
			mockBehavior: func(a []models.Playlist) {
				tablesMock.EXPECT().Playlists().Return(playlistTable)

				sqlxMock.ExpectQuery("SELECT (.+) FROM " + playlistTable).
					WillReturnError(errPqInternal)
			},
			expectError:   true,
			expectedError: errPqInternal,
		},
	}

	var feedAmountLimit uint32 = 100
	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			// Call mock
			tc.mockBehavior(tc.expectedPlaylists)

			a, err := repo.GetFeed(ctx, feedAmountLimit)

			// Test
			if tc.expectError {
				assert.ErrorAs(t, err, &tc.expectedError)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedPlaylists, a)
			}
		})
	}
}

func TestPlaylistRepositoryPostgreSQL_GetByUser(t *testing.T) {
	// Init
	type mockBehavior func(userID uint32, playlists []models.Playlist)

	dbMock, sqlxMock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("%v", err)
	}
	defer dbMock.Close()

	c := gomock.NewController(t)

	tablesMock := playlistMocks.NewMockTables(c)

	repo := NewPostgreSQL(sqlx.NewDb(dbMock, "postgres"), tablesMock)

	// Test filling
	const defaultUserID uint32 = 1

	descriptionID1 := "Ожидайте 3 июня"
	descriptionID2 := "Если вдруг решил отдохнуть"
	defaultPlaylists := []models.Playlist{
		{
			ID:          1,
			Name:        "Музыка для эпичной защиты",
			Description: &descriptionID1,
			CoverSrc:    "/playlists/covers/epic.png",
		},
		{
			ID:          2,
			Name:        "Для чилла",
			Description: &descriptionID2,
			CoverSrc:    "/playlists/covers/chill.png",
		},
	}

	testTable := []struct {
		name              string
		userID            uint32
		mockBehavior      mockBehavior
		expectedPlaylists []models.Playlist
		expectError       bool
		expectedError     error
	}{
		{
			name:   "Common",
			userID: defaultUserID,
			mockBehavior: func(userID uint32, p []models.Playlist) {
				tablesMock.EXPECT().Playlists().Return(playlistTable)
				tablesMock.EXPECT().UsersPlaylists().Return(usersPlaylistsTable)

				rows := sqlxMock.NewRows([]string{"id", "name", "description", "cover_src"})
				for ind := range p {
					rows.AddRow(p[ind].ID, p[ind].Name, p[ind].Description, p[ind].CoverSrc)
				}
				sqlxMock.ExpectQuery(fmt.Sprintf("SELECT (.+) FROM %s p INNER JOIN %s",
					playlistTable, usersPlaylistsTable)).
					WithArgs(userID).
					WillReturnRows(rows)
			},
			expectedPlaylists: defaultPlaylists,
		},
		{
			name:   "No Such User",
			userID: defaultUserID,
			mockBehavior: func(userID uint32, playlists []models.Playlist) {
				tablesMock.EXPECT().Playlists().Return(playlistTable)
				tablesMock.EXPECT().UsersPlaylists().Return(usersPlaylistsTable)

				sqlxMock.ExpectQuery(fmt.Sprintf("SELECT (.+) FROM %s p INNER JOIN %s",
					playlistTable, usersPlaylistsTable)).
					WithArgs(userID).
					WillReturnError(sql.ErrNoRows)
			},
			expectError:   true,
			expectedError: &models.NoSuchUserError{UserID: defaultUserID},
		},
		{
			name:   "Internal PostgreSQL Error",
			userID: defaultUserID,
			mockBehavior: func(userID uint32, playlists []models.Playlist) {
				tablesMock.EXPECT().Playlists().Return(playlistTable)
				tablesMock.EXPECT().UsersPlaylists().Return(usersPlaylistsTable)

				sqlxMock.ExpectQuery(fmt.Sprintf("SELECT (.+) FROM %s p INNER JOIN %s",
					playlistTable, usersPlaylistsTable)).
					WithArgs(userID).
					WillReturnError(errPqInternal)
			},
			expectError:   true,
			expectedError: errPqInternal,
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			// Call mock
			tc.mockBehavior(tc.userID, tc.expectedPlaylists)

			a, err := repo.GetByUser(ctx, tc.userID)

			// Test
			if tc.expectError {
				assert.ErrorContains(t, err, tc.expectedError.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedPlaylists, a)
			}
		})
	}
}

func TestPlaylistRepositoryPostgreSQL_GetLikedByUser(t *testing.T) {
	// Init
	type mockBehavior func(userID uint32, playlists []models.Playlist)

	dbMock, sqlxMock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("%v", err)
	}
	defer dbMock.Close()

	c := gomock.NewController(t)

	tablesMock := playlistMocks.NewMockTables(c)

	repo := NewPostgreSQL(sqlx.NewDb(dbMock, "postgres"), tablesMock)

	// Test filling
	const defaultUserID uint32 = 1

	descriptionID1 := "Ожидайте 3 июня"
	descriptionID2 := "Если вдруг решил отдохнуть"
	defaultPlaylists := []models.Playlist{
		{
			ID:          1,
			Name:        "Музыка для эпичной защиты",
			Description: &descriptionID1,
			CoverSrc:    "/playlists/covers/epic.png",
		},
		{
			ID:          2,
			Name:        "Для чилла",
			Description: &descriptionID2,
			CoverSrc:    "/playlists/covers/chill.png",
		},
	}

	testTable := []struct {
		name              string
		userID            uint32
		mockBehavior      mockBehavior
		expectedPlaylists []models.Playlist
		expectError       bool
		expectedError     error
	}{
		{
			name:   "Common",
			userID: defaultUserID,
			mockBehavior: func(userID uint32, p []models.Playlist) {
				tablesMock.EXPECT().Playlists().Return(playlistTable)
				tablesMock.EXPECT().LikedPlaylists().Return(likedPlaylistsTable)

				rows := sqlxMock.NewRows([]string{"id", "name", "description", "cover_src"})
				for ind := range p {
					rows.AddRow(p[ind].ID, p[ind].Name, p[ind].Description, p[ind].CoverSrc)
				}
				sqlxMock.ExpectQuery(fmt.Sprintf("SELECT (.+) FROM %s p INNER JOIN %s",
					playlistTable, likedPlaylistsTable)).
					WithArgs(userID).
					WillReturnRows(rows)
			},
			expectedPlaylists: defaultPlaylists,
		},
		{
			name:   "No Such User",
			userID: defaultUserID,
			mockBehavior: func(userID uint32, p []models.Playlist) {
				tablesMock.EXPECT().Playlists().Return(playlistTable)
				tablesMock.EXPECT().LikedPlaylists().Return(likedPlaylistsTable)

				sqlxMock.ExpectQuery(fmt.Sprintf("SELECT (.+) FROM %s p INNER JOIN %s",
					playlistTable, likedPlaylistsTable)).
					WithArgs(userID).
					WillReturnError(sql.ErrNoRows)
			},
			expectError:   true,
			expectedError: &models.NoSuchUserError{UserID: defaultUserID},
		},
		{
			name:   "Internal PostgreSQL Error",
			userID: defaultUserID,
			mockBehavior: func(userID uint32, p []models.Playlist) {
				tablesMock.EXPECT().Playlists().Return(playlistTable)
				tablesMock.EXPECT().LikedPlaylists().Return(likedPlaylistsTable)

				sqlxMock.ExpectQuery(fmt.Sprintf("SELECT (.+) FROM %s p INNER JOIN %s",
					playlistTable, likedPlaylistsTable)).
					WithArgs(userID).
					WillReturnError(errPqInternal)
			},
			expectError:   true,
			expectedError: errPqInternal,
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			// Call mock
			tc.mockBehavior(tc.userID, tc.expectedPlaylists)

			a, err := repo.GetLikedByUser(ctx, tc.userID)

			// Test
			if tc.expectError {
				assert.ErrorContains(t, err, tc.expectedError.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedPlaylists, a)
			}
		})
	}
}

func TestPlaylistRepositoryPostgreSQL_Like(t *testing.T) {
	// Init
	type mockBehavior func(playlistID uint32, userID uint32)

	type LikeInfo struct {
		playlistID uint32
		userID     uint32
	}

	dbMock, sqlxMock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("%v", err)
	}
	defer dbMock.Close()

	c := gomock.NewController(t)

	tablesMock := playlistMocks.NewMockTables(c)

	repo := NewPostgreSQL(sqlx.NewDb(dbMock, "postgres"), tablesMock)

	// Test filling
	const defaultPlaylistToLikeID uint32 = 1
	const defaultLikedUserID uint32 = 1

	defaultLikeInfo := LikeInfo{
		playlistID: defaultPlaylistToLikeID,
		userID:     defaultLikedUserID,
	}

	testTable := []struct {
		name           string
		likeInfo       LikeInfo
		mockBehavior   mockBehavior
		expectInserted bool
		expectError    bool
		expectedError  error
	}{
		{
			name:     "Common",
			likeInfo: defaultLikeInfo,
			mockBehavior: func(playlistID uint32, userID uint32) {
				tablesMock.EXPECT().LikedPlaylists().Return(likedPlaylistsTable)

				sqlxMock.ExpectExec("INSERT INTO "+likedPlaylistsTable).
					WithArgs(playlistID, userID).
					WillReturnResult(driver.ResultNoRows)
			},
			expectInserted: true,
		},
		{
			name:     "No Such Playlist",
			likeInfo: defaultLikeInfo,
			mockBehavior: func(playlistID uint32, userID uint32) {
				tablesMock.EXPECT().LikedPlaylists().Return(likedPlaylistsTable)

				sqlxMock.ExpectExec("INSERT INTO "+likedPlaylistsTable).
					WithArgs(playlistID, userID).
					WillReturnError(sql.ErrNoRows)
			},
			expectError:   true,
			expectedError: &models.NoSuchPlaylistError{PlaylistID: defaultPlaylistToLikeID},
		},
		{
			name:     "Internal PostgreSQL Error",
			likeInfo: defaultLikeInfo,
			mockBehavior: func(playlistID uint32, userID uint32) {
				tablesMock.EXPECT().LikedPlaylists().Return(likedPlaylistsTable)

				sqlxMock.ExpectExec("INSERT INTO "+likedPlaylistsTable).
					WithArgs(playlistID, userID).
					WillReturnError(errPqInternal)
			},
			expectError:   true,
			expectedError: errPqInternal,
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			// Call mock
			tc.mockBehavior(tc.likeInfo.playlistID, tc.likeInfo.userID)

			inserted, err := repo.InsertLike(ctx, tc.likeInfo.playlistID, tc.likeInfo.userID)

			// Test
			if tc.expectError {
				assert.ErrorContains(t, err, tc.expectedError.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, inserted, tc.expectInserted)
			}
		})
	}
}

func TestPlaylistRepositoryPostgreSQL_DeleteLike(t *testing.T) {
	// Init
	type mockBehavior func(playlistID uint32, userID uint32)

	type LikeInfo struct {
		playlistID uint32
		userID     uint32
	}

	dbMock, sqlxMock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("%v", err)
	}
	defer dbMock.Close()

	c := gomock.NewController(t)

	tablesMock := playlistMocks.NewMockTables(c)

	repo := NewPostgreSQL(sqlx.NewDb(dbMock, "postgres"), tablesMock)

	// Test filling
	defaultLikeInfo := LikeInfo{
		playlistID: 1,
		userID:     1,
	}

	testTable := []struct {
		name           string
		likeInfo       LikeInfo
		mockBehavior   mockBehavior
		expectInserted bool
		expectError    bool
		expectedError  error
	}{
		{
			name:     "Common",
			likeInfo: defaultLikeInfo,
			mockBehavior: func(playlistID uint32, userID uint32) {
				tablesMock.EXPECT().LikedPlaylists().Return(likedPlaylistsTable)

				sqlxMock.ExpectExec("DELETE FROM "+likedPlaylistsTable).
					WithArgs(playlistID, userID).
					WillReturnResult(driver.RowsAffected(1))
			},
			expectInserted: true,
		},
		{
			name:     "No Such Like",
			likeInfo: defaultLikeInfo,
			mockBehavior: func(playlistID, userID uint32) {
				tablesMock.EXPECT().LikedPlaylists().Return(likedPlaylistsTable)

				sqlxMock.ExpectExec("DELETE FROM "+likedPlaylistsTable).
					WithArgs(playlistID, userID).
					WillReturnResult(driver.RowsAffected(0))
			},
			expectInserted: false,
		},
		{
			name:     "Internal PostgreSQL Error",
			likeInfo: defaultLikeInfo,
			mockBehavior: func(playlistID uint32, userID uint32) {
				tablesMock.EXPECT().LikedPlaylists().Return(likedPlaylistsTable)

				sqlxMock.ExpectExec("DELETE FROM "+likedPlaylistsTable).
					WithArgs(playlistID, userID).
					WillReturnError(errPqInternal)
			},
			expectError:   true,
			expectedError: errPqInternal,
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			// Call mock
			tc.mockBehavior(tc.likeInfo.playlistID, tc.likeInfo.userID)

			inserted, err := repo.DeleteLike(ctx, tc.likeInfo.playlistID, tc.likeInfo.userID)

			// Test
			if tc.expectError {
				assert.ErrorIs(t, err, tc.expectedError)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, inserted, tc.expectInserted)
			}
		})
	}
}

func TestPlaylistRepositoryPostgreSQL_IsLiked(t *testing.T) {
	// Init
	type mockBehavior func(playlistID, userID uint32)

	dbMock, sqlxMock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("%v", err)
	}
	defer dbMock.Close()

	c := gomock.NewController(t)

	tablesMock := playlistMocks.NewMockTables(c)

	repo := NewPostgreSQL(sqlx.NewDb(dbMock, "postgres"), tablesMock)

	// Test filling
	const defaultPlaylistID uint32 = 1
	const defaultUserID uint32 = 1

	testTable := []struct {
		name          string
		playlistID    uint32
		userID        uint32
		mockBehavior  mockBehavior
		expectError   bool
		expectedError error
		isLiked       bool
	}{
		{
			name:       "Liked",
			playlistID: defaultPlaylistID,
			userID:     defaultUserID,
			mockBehavior: func(playlistID, userID uint32) {
				tablesMock.EXPECT().LikedPlaylists().Return(likedPlaylistsTable)

				row := sqlxMock.NewRows([]string{"exists"}).AddRow(true)
				sqlxMock.ExpectQuery("SELECT EXISTS").
					WithArgs(playlistID, userID).
					WillReturnRows(row)
			},
			isLiked: true,
		},
		{
			name:       "Isn't liked",
			playlistID: defaultPlaylistID,
			userID:     defaultUserID,
			mockBehavior: func(playlistID, userID uint32) {
				tablesMock.EXPECT().LikedPlaylists().Return(likedPlaylistsTable)

				row := sqlxMock.NewRows([]string{"exists"}).AddRow(false)
				sqlxMock.ExpectQuery("SELECT EXISTS").
					WithArgs(playlistID, userID).
					WillReturnRows(row)
			},
		},
		{
			name:       "Internal PostgreSQL Error",
			playlistID: defaultPlaylistID,
			userID:     defaultUserID,
			mockBehavior: func(playlistID, userID uint32) {
				tablesMock.EXPECT().LikedPlaylists().Return(likedPlaylistsTable)

				sqlxMock.ExpectQuery("SELECT EXISTS").
					WithArgs(playlistID, userID).
					WillReturnError(errPqInternal)
			},
			expectError:   true,
			expectedError: errPqInternal,
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			// Call mock
			tc.mockBehavior(tc.playlistID, tc.userID)

			isLiked, err := repo.IsLiked(ctx, tc.playlistID, tc.userID)

			// Test
			if tc.expectError {
				assert.ErrorContains(t, err, tc.expectedError.Error())
			} else {
				assert.Equal(t, tc.isLiked, isLiked)
				assert.NoError(t, err)
			}
		})
	}
}
