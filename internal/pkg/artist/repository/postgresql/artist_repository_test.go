package postgresql

import (
	"database/sql/driver"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/golang/mock/gomock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"

	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/models"
	artistMocks "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/artist/mocks"
	logMocks "github.com/go-park-mail-ru/2023_1_Technokaif/pkg/logger/mocks"
)

const artistTable = "Artists"
const likedArtistsTable = "Liked_artists"

var errPqInternal = errors.New("postgres is dead")

func TestTrackRepositoryInsert(t *testing.T) {
	// Init
	type mockBehavior func(a models.Artist, id uint32)

	dbMock, sqlxMock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("%v", err)
	}
	defer dbMock.Close()

	c := gomock.NewController(t)

	l := logMocks.NewMockLogger(c)
	l.EXPECT().Error(gomock.Any()).AnyTimes()
	l.EXPECT().Info(gomock.Any()).AnyTimes()
	l.EXPECT().Errorf(gomock.Any(), gomock.Any()).AnyTimes()
	l.EXPECT().Infof(gomock.Any(), gomock.Any()).AnyTimes()

	tablesMock := artistMocks.NewMockTables(c)

	repo := NewPostgreSQL(sqlx.NewDb(dbMock, "postgres"), tablesMock, l)

	// Test filling
	defaultUserOfArtistID := uint32(1)

	defaultArtistToIsert := models.Artist{
		Name:      "Oxxxymiron",
		UserID:    &defaultUserOfArtistID,
		AvatarSrc: "/artists/avatars/oxxxymiron.png",
	}

	testTable := []struct {
		name          string
		artist        models.Artist
		mockBehavior  mockBehavior
		expectedID    uint32
		expectError   bool
		expectedError error
	}{
		{
			name:   "Common",
			artist: defaultArtistToIsert,
			mockBehavior: func(a models.Artist, id uint32) {
				tablesMock.EXPECT().Artists().Return(artistTable)

				row := sqlxMock.NewRows([]string{"id"}).AddRow(id)
				sqlxMock.ExpectQuery("INSERT INTO "+artistTable).
					WithArgs(a.UserID, a.Name, a.AvatarSrc).
					WillReturnRows(row)
			},
			expectedID: 1,
		},
		{
			name:   "Insert Artists Issue",
			artist: defaultArtistToIsert,
			mockBehavior: func(a models.Artist, id uint32) {
				tablesMock.EXPECT().Artists().Return(artistTable)

				sqlxMock.ExpectQuery("INSERT INTO "+artistTable).
					WithArgs(a.UserID, a.Name, a.AvatarSrc).
					WillReturnError(errPqInternal)
			},
			expectedID:    1,
			expectError:   true,
			expectedError: errPqInternal,
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			// Call mock
			tc.mockBehavior(tc.artist, tc.expectedID)

			id, err := repo.Insert(tc.artist)

			// Test
			if tc.expectError {
				assert.ErrorIs(t, err, tc.expectedError)
			} else {
				assert.Equal(t, id, tc.expectedID)
				assert.NoError(t, err)
			}
		})
	}
}

func TestArtistRepositoryDeleteByID(t *testing.T) {
	// Init
	type mockBehavior func(artistID uint32)

	dbMock, sqlxMock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("%v", err)
	}
	defer dbMock.Close()

	c := gomock.NewController(t)

	l := logMocks.NewMockLogger(c)
	l.EXPECT().Error(gomock.Any()).AnyTimes()
	l.EXPECT().Info(gomock.Any()).AnyTimes()
	l.EXPECT().Errorf(gomock.Any(), gomock.Any()).AnyTimes()
	l.EXPECT().Infof(gomock.Any(), gomock.Any()).AnyTimes()

	tablesMock := artistMocks.NewMockTables(c)

	repo := NewPostgreSQL(sqlx.NewDb(dbMock, "postgres"), tablesMock, l)

	// Test filling
	defaultArtistToDeleteID := uint32(1)

	testTable := []struct {
		name             string
		artistToDeleteID uint32
		mockBehavior     mockBehavior
		expectError      bool
		expectedError    error
	}{
		{
			name:             "Common",
			artistToDeleteID: defaultArtistToDeleteID,
			mockBehavior: func(artistID uint32) {
				tablesMock.EXPECT().Artists().Return(artistTable)

				sqlxMock.ExpectExec("DELETE FROM " + artistTable).
					WithArgs(artistID).
					WillReturnResult(driver.ResultNoRows)
			},
		},
		{
			name:             "Internal PostgreSQL Error",
			artistToDeleteID: defaultArtistToDeleteID,
			mockBehavior: func(artistID uint32) {
				tablesMock.EXPECT().Artists().Return(artistTable)

				sqlxMock.ExpectExec("DELETE FROM " + artistTable).
					WithArgs(artistID).
					WillReturnError(errPqInternal)
			},
			expectError:   true,
			expectedError: errPqInternal,
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			// Call mock
			tc.mockBehavior(tc.artistToDeleteID)

			err := repo.DeleteByID(tc.artistToDeleteID)

			// Test
			if tc.expectError {
				assert.ErrorIs(t, err, tc.expectedError)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestArtistRepositoryLike(t *testing.T) {
	// Init
	type mockBehavior func(artistID uint32, userID uint32)

	type LikeInfo struct {
		artistID uint32
		userID   uint32
	}

	dbMock, sqlxMock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("%v", err)
	}
	defer dbMock.Close()

	c := gomock.NewController(t)

	l := logMocks.NewMockLogger(c)
	l.EXPECT().Error(gomock.Any()).AnyTimes()
	l.EXPECT().Info(gomock.Any()).AnyTimes()
	l.EXPECT().Errorf(gomock.Any(), gomock.Any()).AnyTimes()
	l.EXPECT().Infof(gomock.Any(), gomock.Any()).AnyTimes()

	tablesMock := artistMocks.NewMockTables(c)

	repo := NewPostgreSQL(sqlx.NewDb(dbMock, "postgres"), tablesMock, l)

	// Test filling
	defaultArtistToLikeID := uint32(1)
	defaultLikedUserID := uint32(1)

	defaultLikeInfo := LikeInfo{
		artistID: defaultArtistToLikeID,
		userID:   defaultLikedUserID,
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
			mockBehavior: func(artistID uint32, userID uint32) {
				tablesMock.EXPECT().LikedArtists().Return(likedArtistsTable)

				sqlxMock.ExpectExec("INSERT INTO "+likedArtistsTable).
					WithArgs(artistID, userID).
					WillReturnResult(driver.ResultNoRows)
			},
			expectInserted: true,
		},
		{
			name:     "Internal PostgreSQL Error",
			likeInfo: defaultLikeInfo,
			mockBehavior: func(artistID uint32, userID uint32) {
				tablesMock.EXPECT().LikedArtists().Return(likedArtistsTable)

				sqlxMock.ExpectExec("INSERT INTO "+likedArtistsTable).
					WithArgs(artistID, userID).
					WillReturnError(errPqInternal)
			},
			expectError:   true,
			expectedError: errPqInternal,
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			// Call mock
			tc.mockBehavior(tc.likeInfo.artistID, tc.likeInfo.userID)

			success, err := repo.InsertLike(tc.likeInfo.artistID, tc.likeInfo.userID)

			// Test
			if tc.expectError {
				assert.ErrorIs(t, err, tc.expectedError)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, success, tc.expectInserted)
			}
		})
	}
}
