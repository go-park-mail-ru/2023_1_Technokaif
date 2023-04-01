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
	trackMocks "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/track/mocks"
	logMocks "github.com/go-park-mail-ru/2023_1_Technokaif/pkg/logger/mocks"
)

const trackTable = "Tracks"
const artistsTracksTable = "Artists_Tracks"
const likedTracksTable = "Liked_tracks"

var errPqInternal = errors.New("postgres is dead")

func TestTrackRepositoryInsert(t *testing.T) {
	// Init
	type mockBehavior func(track models.Track, artistsID []uint32, id uint32)

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

	tablesMock := trackMocks.NewMockTables(c)

	repo := NewPostgreSQL(sqlx.NewDb(dbMock, "postgres"), tablesMock, l)

	// Test filling
	albumID := uint32(1)
	AlbumPosition := uint32(1)

	defaultTrackToIsert := models.Track{
		Name:          "LAGG OUT",
		AlbumID:       &albumID,
		AlbumPosition: &AlbumPosition,
		CoverSrc:      "/tracks/covers/laggout.png",
		RecordSrc:     "/tracks/records/laggout.png",
	}

	defaultArtistsIDToInsert := []uint32{1, 2, 3}

	testTable := []struct {
		name          string
		track         models.Track
		artistsID     []uint32
		mockBehavior  mockBehavior
		expectedID    uint32
		expectError   bool
		expectedError error
	}{
		{
			name:      "Common",
			track:     defaultTrackToIsert,
			artistsID: defaultArtistsIDToInsert,
			mockBehavior: func(t models.Track, artistsID []uint32, id uint32) {
				tablesMock.EXPECT().Tracks().Return(trackTable)
				tablesMock.EXPECT().ArtistsTracks().Return(artistsTracksTable)

				sqlxMock.ExpectBegin()

				row := sqlxMock.NewRows([]string{"id"}).AddRow(id)
				sqlxMock.ExpectQuery("INSERT INTO "+trackTable).
					WithArgs(t.Name, t.AlbumID, t.AlbumPosition, t.CoverSrc, t.RecordSrc).
					WillReturnRows(row)

				for _, artistID := range artistsID {
					sqlxMock.ExpectExec("INSERT INTO "+artistsTracksTable).
						WithArgs(artistID, id).
						WillReturnResult(driver.ResultNoRows)
				}

				sqlxMock.ExpectCommit()
			},
			expectedID: 1,
		},
		{
			name:      "Insert Artists to Track Issue",
			track:     defaultTrackToIsert,
			artistsID: defaultArtistsIDToInsert,
			mockBehavior: func(t models.Track, artistsID []uint32, id uint32) {
				tablesMock.EXPECT().Tracks().Return(trackTable)
				tablesMock.EXPECT().ArtistsTracks().Return(artistsTracksTable)

				sqlxMock.ExpectBegin()

				row := sqlxMock.NewRows([]string{"id"}).AddRow(id)
				sqlxMock.ExpectQuery("INSERT INTO "+trackTable).
					WithArgs(t.Name, t.AlbumID, t.AlbumPosition, t.CoverSrc, t.RecordSrc).
					WillReturnRows(row)

				sqlxMock.ExpectExec("INSERT INTO "+artistsTracksTable).
					WithArgs(artistsID[0], id).
					WillReturnError(errPqInternal)

				sqlxMock.ExpectRollback()
			},
			expectedID:    1,
			expectError:   true,
			expectedError: errPqInternal,
		},
		{
			name:      "Insert Track Issue",
			track:     defaultTrackToIsert,
			artistsID: defaultArtistsIDToInsert,
			mockBehavior: func(t models.Track, artistsID []uint32, id uint32) {
				tablesMock.EXPECT().Tracks().Return(trackTable)

				sqlxMock.ExpectBegin()

				sqlxMock.ExpectQuery("INSERT INTO "+trackTable).
					WithArgs(t.Name, t.AlbumID, t.AlbumPosition, t.CoverSrc, t.RecordSrc).
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
			tc.mockBehavior(tc.track, tc.artistsID, tc.expectedID)

			id, err := repo.Insert(tc.track, tc.artistsID)

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

func TestTrackRepositoryDeleteByID(t *testing.T) {
	// Init
	type mockBehavior func(trackID uint32)

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

	tablesMock := trackMocks.NewMockTables(c)

	repo := NewPostgreSQL(sqlx.NewDb(dbMock, "postgres"), tablesMock, l)

	// Test filling
	defaultTrackToDeleteID := uint32(1)

	testTable := []struct {
		name            string
		trackToDeleteID uint32
		mockBehavior    mockBehavior
		expectError     bool
		expectedError   error
	}{
		{
			name:            "Common",
			trackToDeleteID: defaultTrackToDeleteID,
			mockBehavior: func(trackID uint32) {
				tablesMock.EXPECT().Tracks().Return(trackTable)

				sqlxMock.ExpectExec("DELETE FROM " + trackTable).
					WithArgs(trackID).
					WillReturnResult(driver.ResultNoRows)
			},
		},
		{
			name:            "Internal PostgreSQL Error",
			trackToDeleteID: defaultTrackToDeleteID,
			mockBehavior: func(trackID uint32) {
				tablesMock.EXPECT().Tracks().Return(trackTable)

				sqlxMock.ExpectExec("DELETE FROM " + trackTable).
					WithArgs(trackID).
					WillReturnError(errPqInternal)
			},
			expectError:   true,
			expectedError: errPqInternal,
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			// Call mock
			tc.mockBehavior(tc.trackToDeleteID)

			err := repo.DeleteByID(tc.trackToDeleteID)

			// Test
			if tc.expectError {
				assert.ErrorIs(t, err, tc.expectedError)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestTrackRepositoryLike(t *testing.T) {
	// Init
	type mockBehavior func(trackID uint32, userID uint32)

	type LikeInfo struct {
		trackID uint32
		userID  uint32
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

	tablesMock := trackMocks.NewMockTables(c)

	repo := NewPostgreSQL(sqlx.NewDb(dbMock, "postgres"), tablesMock, l)

	// Test filling
	defaultTrackToLikeID := uint32(1)
	defaultLikedUserID := uint32(1)

	defaultLikeInfo := LikeInfo{
		trackID: defaultTrackToLikeID,
		userID:  defaultLikedUserID,
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
				tablesMock.EXPECT().LikedTracks().Return(likedTracksTable)

				sqlxMock.ExpectExec("INSERT INTO "+likedTracksTable).
					WithArgs(artistID, userID).
					WillReturnResult(driver.ResultNoRows)
			},
			expectInserted: true,
		},
		{
			name:     "Internal PostgreSQL Error",
			likeInfo: defaultLikeInfo,
			mockBehavior: func(artistID uint32, userID uint32) {
				tablesMock.EXPECT().LikedTracks().Return(likedTracksTable)

				sqlxMock.ExpectExec("INSERT INTO "+likedTracksTable).
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
			tc.mockBehavior(tc.likeInfo.trackID, tc.likeInfo.userID)

			success, err := repo.InsertLike(tc.likeInfo.trackID, tc.likeInfo.userID)

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
