package postgresql

import (
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
	trackMocks "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/track/mocks"
	commonTests "github.com/go-park-mail-ru/2023_1_Technokaif/internal/common/tests"
)

const trackTable = "Tracks"
const artistsTracksTable = "Artists_Tracks"
const likedTracksTable = "Liked_tracks"

var errPqInternal = errors.New("postgres is dead")

var defaultTrackAlbumID1 = uint32(1)
var defaultTrackAlbumID2 = uint32(2)
var defaultTracks = []models.Track{
	{
		ID:        1,
		Name:      "Lagg Out",
		AlbumID:   &defaultTrackAlbumID1,
		CoverSrc:  "/tracks/covers/laggout.png",
		RecordSrc: "/tracks/records/laggout.wav",
		Listens:   9999999,
	},
	{
		ID:        2,
		Name:      "Накануне",
		AlbumID:   &defaultTrackAlbumID2,
		CoverSrc:  "/tracks/covers/nakanune.png",
		RecordSrc: "/tracks/records/nakanune.wav",
		Listens:   10000000,
	},
}

func TestTrackRepositoryInsert(t *testing.T) {
	// Init
	type mockBehavior func(track models.Track, artistsID []uint32, id uint32)

	dbMock, sqlxMock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("%v", err)
	}
	defer dbMock.Close()

	c := gomock.NewController(t)

	l := commonTests.MockLogger(c)

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

func TestTrackRepositoryGetByID(t *testing.T) {
	// Init
	type mockBehavior func(trackID uint32, t models.Track)

	dbMock, sqlxMock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("%v", err)
	}
	defer dbMock.Close()

	c := gomock.NewController(t)

	l := commonTests.MockLogger(c)

	tablesMock := trackMocks.NewMockTables(c)

	repo := NewPostgreSQL(sqlx.NewDb(dbMock, "postgres"), tablesMock, l)

	// Test filling
	defaultTrackToGetID := uint32(1)
	defaultTrackAlbumID := uint32(1)

	defaultTrack := models.Track{
		ID:        defaultTrackToGetID,
		Name:      "Lagg Out",
		AlbumID:   &defaultTrackAlbumID,
		CoverSrc:  "/tracks/covers/laggout.png",
		RecordSrc: "/tracks/records/laggout.wav",
		Listens:   9999999,
	}

	testTable := []struct {
		name          string
		trackToGetID  uint32
		mockBehavior  mockBehavior
		expectedTrack models.Track
		expectError   bool
		expectedError error
	}{
		{
			name:         "Common",
			trackToGetID: defaultTrackToGetID,
			mockBehavior: func(trackID uint32, t models.Track) {
				tablesMock.EXPECT().Tracks().Return(trackTable)

				row := sqlxMock.NewRows([]string{"id", "name", "album_id", "cover_src", "record_src", "listens"}).
					AddRow(t.ID, t.Name, t.AlbumID, t.CoverSrc, t.RecordSrc, t.Listens)
				sqlxMock.ExpectQuery("SELECT (.+) FROM " + trackTable).
					WithArgs(trackID).
					WillReturnRows(row)
			},
			expectedTrack: defaultTrack,
		},
		{
			name:         "No Such Track",
			trackToGetID: defaultTrackToGetID,
			mockBehavior: func(trackID uint32, t models.Track) {
				tablesMock.EXPECT().Tracks().Return(trackTable)

				sqlxMock.ExpectQuery("SELECT (.+) FROM " + trackTable).
					WithArgs(trackID).
					WillReturnError(sql.ErrNoRows)
			},
			expectError:   true,
			expectedError: &models.NoSuchTrackError{TrackID: defaultTrackToGetID},
		},
		{
			name:         "Internal PostgreSQL Error",
			trackToGetID: defaultTrackToGetID,
			mockBehavior: func(trackID uint32, t models.Track) {
				tablesMock.EXPECT().Tracks().Return(trackTable)

				sqlxMock.ExpectQuery("SELECT (.+) FROM " + trackTable).
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
			tc.mockBehavior(tc.trackToGetID, tc.expectedTrack)

			tr, err := repo.GetByID(tc.trackToGetID)

			// Test
			if tc.expectError {
				assert.ErrorContains(t, err, tc.expectedError.Error())
			} else {
				assert.Equal(t, tc.expectedTrack, *tr)
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

	l := commonTests.MockLogger(c)

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
					WillReturnResult(driver.RowsAffected(1))
			},
		},
		{
			name:            "No Such Track",
			trackToDeleteID: defaultTrackToDeleteID,
			mockBehavior: func(trackID uint32) {
				tablesMock.EXPECT().Tracks().Return(trackTable)

				sqlxMock.ExpectExec("DELETE FROM " + trackTable).
					WithArgs(trackID).
					WillReturnResult(driver.RowsAffected(0))
			},
			expectError:   true,
			expectedError: &models.NoSuchTrackError{TrackID: defaultTrackToDeleteID},
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
				assert.ErrorContains(t, err, tc.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestTrackRepositoryGetFeed(t *testing.T) {
	// Init
	type mockBehavior func(tracks []models.Track)

	dbMock, sqlxMock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("%v", err)
	}
	defer dbMock.Close()

	c := gomock.NewController(t)

	l := commonTests.MockLogger(c)

	tablesMock := trackMocks.NewMockTables(c)

	repo := NewPostgreSQL(sqlx.NewDb(dbMock, "postgres"), tablesMock, l)

	// Test filling
	testTable := []struct {
		name           string
		mockBehavior   mockBehavior
		expectedTracks []models.Track
		expectError    bool
		expectedError  error
	}{
		{
			name: "Common",
			mockBehavior: func(t []models.Track) {
				tablesMock.EXPECT().Tracks().Return(trackTable)

				rows := sqlxMock.NewRows([]string{"id", "name", "album_id", "cover_src", "record_src", "listens"}).
					AddRow(t[0].ID, t[0].Name, t[0].AlbumID, t[0].CoverSrc, t[0].RecordSrc, t[0].Listens).
					AddRow(t[1].ID, t[1].Name, t[1].AlbumID, t[1].CoverSrc, t[1].RecordSrc, t[1].Listens)
				sqlxMock.ExpectQuery("SELECT (.+) FROM " + trackTable).
					WillReturnRows(rows)
			},
			expectedTracks: defaultTracks,
		},
		{
			name: "Internal PostgreSQL Error",
			mockBehavior: func(t []models.Track) {
				tablesMock.EXPECT().Tracks().Return(trackTable)

				sqlxMock.ExpectQuery("SELECT (.+) FROM " + trackTable).
					WillReturnError(errPqInternal)
			},
			expectError:   true,
			expectedError: errPqInternal,
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			// Call mock
			tc.mockBehavior(tc.expectedTracks)

			tr, err := repo.GetFeed()

			// Test
			if tc.expectError {
				assert.ErrorAs(t, err, &tc.expectedError)
			} else {
				assert.Equal(t, tc.expectedTracks, tr)
				assert.NoError(t, err)
			}
		})
	}
}

func TestTrackRepositoryGetByArtist(t *testing.T) {
	// Init
	type mockBehavior func(artistID uint32, tracks []models.Track)

	dbMock, sqlxMock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("%v", err)
	}
	defer dbMock.Close()

	c := gomock.NewController(t)

	l := commonTests.MockLogger(c)

	tablesMock := trackMocks.NewMockTables(c)

	repo := NewPostgreSQL(sqlx.NewDb(dbMock, "postgres"), tablesMock, l)

	// Test filling
	defaultArtistID := uint32(1)

	testTable := []struct {
		name           string
		artistID       uint32
		mockBehavior   mockBehavior
		expectedTracks []models.Track
		expectError    bool
		expectedError  error
	}{
		{
			name:     "Common",
			artistID: defaultArtistID,
			mockBehavior: func(artistID uint32, t []models.Track) {
				tablesMock.EXPECT().Tracks().Return(trackTable)
				tablesMock.EXPECT().ArtistsTracks().Return(artistsTracksTable)

				rows := sqlxMock.NewRows([]string{"id", "name", "album_id", "cover_src", "record_src", "listens"}).
					AddRow(t[0].ID, t[0].Name, t[0].AlbumID, t[0].CoverSrc, t[0].RecordSrc, t[0].Listens).
					AddRow(t[1].ID, t[1].Name, t[1].AlbumID, t[1].CoverSrc, t[1].RecordSrc, t[1].Listens)
				sqlxMock.ExpectQuery(fmt.Sprintf("SELECT (.+) FROM %s t INNER JOIN %s",
					trackTable, artistsTracksTable)).
					WithArgs(artistID).
					WillReturnRows(rows)
			},
			expectedTracks: defaultTracks,
		},
		{
			name:     "No Such Artist",
			artistID: defaultArtistID,
			mockBehavior: func(artistID uint32, t []models.Track) {
				tablesMock.EXPECT().Tracks().Return(trackTable)
				tablesMock.EXPECT().ArtistsTracks().Return(artistsTracksTable)

				sqlxMock.ExpectQuery(fmt.Sprintf("SELECT (.+) FROM %s t INNER JOIN %s",
					trackTable, artistsTracksTable)).
					WithArgs(artistID).
					WillReturnError(sql.ErrNoRows)
			},
			expectError:   true,
			expectedError: &models.NoSuchArtistError{ArtistID: defaultArtistID},
		},
		{
			name:     "Internal PostgreSQL Error",
			artistID: defaultArtistID,
			mockBehavior: func(artistID uint32, t []models.Track) {
				tablesMock.EXPECT().Tracks().Return(trackTable)
				tablesMock.EXPECT().ArtistsTracks().Return(artistsTracksTable)

				sqlxMock.ExpectQuery(fmt.Sprintf("SELECT (.+) FROM %s t INNER JOIN %s",
					trackTable, artistsTracksTable)).
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
			tc.mockBehavior(tc.artistID, tc.expectedTracks)

			a, err := repo.GetByArtist(tc.artistID)

			// Test
			if tc.expectError {
				assert.ErrorContains(t, err, tc.expectedError.Error())
			} else {
				assert.Equal(t, tc.expectedTracks, a)
				assert.NoError(t, err)
			}
		})
	}
}

func TestTrackRepositoryGetByAlbum(t *testing.T) {
	// Init
	type mockBehavior func(albumID uint32, tracks []models.Track)

	dbMock, sqlxMock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("%v", err)
	}
	defer dbMock.Close()

	c := gomock.NewController(t)

	l := commonTests.MockLogger(c)

	tablesMock := trackMocks.NewMockTables(c)

	repo := NewPostgreSQL(sqlx.NewDb(dbMock, "postgres"), tablesMock, l)

	// Test filling
	defaultAlbumID := uint32(1)

	defaultAlbumPosition := uint32(1)
	for ind := range defaultTracks {
		defaultTracks[ind].AlbumPosition = &defaultAlbumPosition
	}

	defer func() {
		for ind := range defaultTracks {
			defaultTracks[ind].AlbumPosition = nil
		}
	}()

	testTable := []struct {
		name           string
		albumID        uint32
		mockBehavior   mockBehavior
		expectedTracks []models.Track
		expectError    bool
		expectedError  error
	}{
		{
			name:    "Common",
			albumID: defaultAlbumID,
			mockBehavior: func(albumID uint32, t []models.Track) {
				tablesMock.EXPECT().Tracks().Return(trackTable)

				rows := sqlxMock.NewRows([]string{"id", "name", "album_id", "album_position", "cover_src", "record_src", "listens"}).
					AddRow(t[0].ID, t[0].Name, t[0].AlbumID, t[0].AlbumPosition, t[0].CoverSrc, t[0].RecordSrc, t[0].Listens).
					AddRow(t[1].ID, t[1].Name, t[1].AlbumID, t[1].AlbumPosition, t[1].CoverSrc, t[1].RecordSrc, t[1].Listens)
				sqlxMock.ExpectQuery(fmt.Sprintf("SELECT (.+) FROM %s", trackTable)).
					WithArgs(albumID).
					WillReturnRows(rows)
			},
			expectedTracks: defaultTracks,
		},
		{
			name:    "No Such Track",
			albumID: defaultAlbumID,
			mockBehavior: func(albumID uint32, t []models.Track) {
				tablesMock.EXPECT().Tracks().Return(trackTable)

				sqlxMock.ExpectQuery(fmt.Sprintf("SELECT (.+) FROM %s", trackTable)).
					WithArgs(albumID).
					WillReturnError(sql.ErrNoRows)
			},
			expectError:   true,
			expectedError: &models.NoSuchAlbumError{AlbumID: defaultAlbumID},
		},
		{
			name:    "Internal PostgreSQL Error",
			albumID: defaultAlbumID,
			mockBehavior: func(albumID uint32, t []models.Track) {
				tablesMock.EXPECT().Tracks().Return(trackTable)

				sqlxMock.ExpectQuery(fmt.Sprintf("SELECT (.+) FROM %s", trackTable)).
					WithArgs(albumID).
					WillReturnError(errPqInternal)
			},
			expectError:   true,
			expectedError: errPqInternal,
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			// Call mock
			tc.mockBehavior(tc.albumID, tc.expectedTracks)

			tr, err := repo.GetByAlbum(tc.albumID)

			// Test
			if tc.expectError {
				assert.ErrorContains(t, err, tc.expectedError.Error())
			} else {
				assert.Equal(t, tc.expectedTracks, tr)
				assert.NoError(t, err)
			}
		})
	}
}

func TestTrackRepositoryGetLikedByUser(t *testing.T) {
	// Init
	type mockBehavior func(userID uint32, tracks []models.Track)

	dbMock, sqlxMock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("%v", err)
	}
	defer dbMock.Close()

	c := gomock.NewController(t)

	l := commonTests.MockLogger(c)
	tablesMock := trackMocks.NewMockTables(c)

	repo := NewPostgreSQL(sqlx.NewDb(dbMock, "postgres"), tablesMock, l)

	// Test filling
	defaultUserID := uint32(1)

	testTable := []struct {
		name           string
		userID         uint32
		mockBehavior   mockBehavior
		expectedTracks []models.Track
		expectError    bool
		expectedError  error
	}{
		{
			name:   "Common",
			userID: defaultUserID,
			mockBehavior: func(userID uint32, t []models.Track) {
				tablesMock.EXPECT().Tracks().Return(trackTable)
				tablesMock.EXPECT().LikedTracks().Return(likedTracksTable)

				rows := sqlxMock.NewRows([]string{"id", "name", "album_id", "cover_src", "record_src", "listens"}).
					AddRow(t[0].ID, t[0].Name, t[0].AlbumID, t[0].CoverSrc, t[0].RecordSrc, t[0].Listens).
					AddRow(t[1].ID, t[1].Name, t[1].AlbumID, t[1].CoverSrc, t[1].RecordSrc, t[1].Listens)
				sqlxMock.ExpectQuery(fmt.Sprintf("SELECT (.+) FROM %s t INNER JOIN %s",
					trackTable, likedTracksTable)).
					WithArgs(userID).
					WillReturnRows(rows)
			},
			expectedTracks: defaultTracks,
		},
		{
			name:   "No Such User",
			userID: defaultUserID,
			mockBehavior: func(userID uint32, t []models.Track) {
				tablesMock.EXPECT().Tracks().Return(trackTable)
				tablesMock.EXPECT().LikedTracks().Return(likedTracksTable)

				sqlxMock.ExpectQuery(fmt.Sprintf("SELECT (.+) FROM %s t INNER JOIN %s",
					trackTable, likedTracksTable)).
					WithArgs(userID).
					WillReturnError(sql.ErrNoRows)
			},
			expectError:   true,
			expectedError: &models.NoSuchUserError{UserID: defaultUserID},
		},
		{
			name:   "Internal PostgreSQL Error",
			userID: defaultUserID,
			mockBehavior: func(userID uint32, t []models.Track) {
				tablesMock.EXPECT().Tracks().Return(trackTable)
				tablesMock.EXPECT().LikedTracks().Return(likedTracksTable)

				sqlxMock.ExpectQuery(fmt.Sprintf("SELECT (.+) FROM %s t INNER JOIN %s",
					trackTable, likedTracksTable)).
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
			tc.mockBehavior(tc.userID, tc.expectedTracks)

			tr, err := repo.GetLikedByUser(tc.userID)

			// Test
			if tc.expectError {
				assert.ErrorContains(t, err, tc.expectedError.Error())
			} else {
				assert.Equal(t, tc.expectedTracks, tr)
				assert.NoError(t, err)
			}
		})
	}
}

func TestTrackRepositoryLike(t *testing.T) {
	// Init
	type mockBehavior func(trackID, userID uint32)

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

	l := commonTests.MockLogger(c)

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
			mockBehavior: func(artistID, userID uint32) {
				tablesMock.EXPECT().LikedTracks().Return(likedTracksTable)

				sqlxMock.ExpectExec("INSERT INTO "+likedTracksTable).
					WithArgs(artistID, userID).
					WillReturnResult(driver.ResultNoRows)
			},
			expectInserted: true,
		},
		{
			name:     "No Such Track",
			likeInfo: defaultLikeInfo,
			mockBehavior: func(artistID, userID uint32) {
				tablesMock.EXPECT().LikedTracks().Return(likedTracksTable)

				sqlxMock.ExpectExec("INSERT INTO "+likedTracksTable).
					WithArgs(artistID, userID).
					WillReturnError(sql.ErrNoRows)
			},
			expectError:   true,
			expectedError: &models.NoSuchTrackError{TrackID: defaultLikeInfo.trackID},
		},
		{
			name:     "Internal PostgreSQL Error",
			likeInfo: defaultLikeInfo,
			mockBehavior: func(artistID, userID uint32) {
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

			inserted, err := repo.InsertLike(tc.likeInfo.trackID, tc.likeInfo.userID)

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

func TestTrackRepositoryDeleteLike(t *testing.T) {
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

	l := commonTests.MockLogger(c)

	tablesMock := trackMocks.NewMockTables(c)

	repo := NewPostgreSQL(sqlx.NewDb(dbMock, "postgres"), tablesMock, l)

	// Test filling
	defaultLikeInfo := LikeInfo{
		trackID: 1,
		userID:  1,
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
			mockBehavior: func(trackID uint32, userID uint32) {
				tablesMock.EXPECT().LikedTracks().Return(likedTracksTable)

				sqlxMock.ExpectExec("DELETE FROM "+likedTracksTable).
					WithArgs(trackID, userID).
					WillReturnResult(driver.RowsAffected(1))
			},
			expectInserted: true,
		},
		{
			name:     "No Such Like",
			likeInfo: defaultLikeInfo,
			mockBehavior: func(trackID, userID uint32) {
				tablesMock.EXPECT().LikedTracks().Return(likedTracksTable)

				sqlxMock.ExpectExec("DELETE FROM "+likedTracksTable).
					WithArgs(trackID, userID).
					WillReturnResult(driver.RowsAffected(0))
			},
			expectInserted: false,
		},
		{
			name:     "Internal PostgreSQL Error",
			likeInfo: defaultLikeInfo,
			mockBehavior: func(trackID uint32, userID uint32) {
				tablesMock.EXPECT().LikedTracks().Return(likedTracksTable)

				sqlxMock.ExpectExec("DELETE FROM "+likedTracksTable).
					WithArgs(trackID, userID).
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

			inserted, err := repo.DeleteLike(tc.likeInfo.trackID, tc.likeInfo.userID)

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
