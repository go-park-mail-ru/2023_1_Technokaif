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

	commonTests "github.com/go-park-mail-ru/2023_1_Technokaif/internal/common/tests"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/models"
	artistMocks "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/artist/mocks"
)

const artistTable = "Artists"
const likedArtistsTable = "Liked_artists"
const artistsAlbumsTable = "Artists_Albums"
const artistsTracksTable = "Artists_Tracks"

var errPqInternal = errors.New("postgres is dead")

func TestArtistRepositoryInsert(t *testing.T) {
	// Init
	type mockBehavior func(a models.Artist, id uint32)

	dbMock, sqlxMock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("%v", err)
	}
	defer dbMock.Close()

	c := gomock.NewController(t)

	l := commonTests.MockLogger(c)

	tablesMock := artistMocks.NewMockTables(c)

	repo := NewPostgreSQL(sqlx.NewDb(dbMock, "postgres"), tablesMock, l)

	// Test filling
	var defaultUserOfArtistID uint32 = 1

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
				assert.NoError(t, err)
				assert.Equal(t, id, tc.expectedID)
			}
		})
	}
}

func TestArtistRepositoryGetByID(t *testing.T) {
	// Init
	type mockBehavior func(artistID uint32, a models.Artist)

	dbMock, sqlxMock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("%v", err)
	}
	defer dbMock.Close()

	c := gomock.NewController(t)

	l := commonTests.MockLogger(c)

	tablesMock := artistMocks.NewMockTables(c)

	repo := NewPostgreSQL(sqlx.NewDb(dbMock, "postgres"), tablesMock, l)

	// Test filling
	const defaultArtistToGetID uint32 = 1

	defaultArtist := models.Artist{
		ID:        defaultArtistToGetID,
		Name:      "Oxxxymiron",
		AvatarSrc: "/artists/avatars/oxxxymiron.png",
	}

	testTable := []struct {
		name           string
		artistToGetID  uint32
		mockBehavior   mockBehavior
		expectedArtist models.Artist
		expectError    bool
		expectedError  error
	}{
		{
			name:          "Common",
			artistToGetID: defaultArtistToGetID,
			mockBehavior: func(artistID uint32, a models.Artist) {
				tablesMock.EXPECT().Artists().Return(artistTable)

				row := sqlxMock.NewRows([]string{"id", "name", "avatar_src"}).
					AddRow(a.ID, a.Name, a.AvatarSrc)
				sqlxMock.ExpectQuery("SELECT (.+) FROM " + artistTable).
					WithArgs(artistID).
					WillReturnRows(row)
			},
			expectedArtist: defaultArtist,
		},
		{
			name:          "No Such Artist",
			artistToGetID: defaultArtistToGetID,
			mockBehavior: func(artistID uint32, a models.Artist) {
				tablesMock.EXPECT().Artists().Return(artistTable)

				sqlxMock.ExpectQuery("SELECT (.+) FROM " + artistTable).
					WithArgs(artistID).
					WillReturnError(sql.ErrNoRows)
			},
			expectError:   true,
			expectedError: &models.NoSuchArtistError{ArtistID: defaultArtistToGetID},
		},
		{
			name:          "Internal PostgreSQL Error",
			artistToGetID: defaultArtistToGetID,
			mockBehavior: func(artistID uint32, a models.Artist) {
				tablesMock.EXPECT().Artists().Return(artistTable)

				sqlxMock.ExpectQuery("SELECT (.+) FROM " + artistTable).
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
			tc.mockBehavior(tc.artistToGetID, tc.expectedArtist)

			a, err := repo.GetByID(tc.artistToGetID)

			// Test
			if tc.expectError {
				assert.ErrorContains(t, err, tc.expectedError.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedArtist, *a)
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

	l := commonTests.MockLogger(c)

	tablesMock := artistMocks.NewMockTables(c)

	repo := NewPostgreSQL(sqlx.NewDb(dbMock, "postgres"), tablesMock, l)

	// Test filling
	const defaultArtistToDeleteID uint32 = 1

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
					WillReturnResult(driver.RowsAffected(1))
			},
		},
		{
			name:             "No Such Artist",
			artistToDeleteID: defaultArtistToDeleteID,
			mockBehavior: func(albumID uint32) {
				tablesMock.EXPECT().Artists().Return(artistTable)

				sqlxMock.ExpectExec("DELETE FROM " + artistTable).
					WithArgs(albumID).
					WillReturnResult(driver.RowsAffected(0))
			},
			expectError:   true,
			expectedError: &models.NoSuchArtistError{ArtistID: defaultArtistToDeleteID},
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
				assert.ErrorContains(t, err, tc.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestArtistRepositoryGetFeed(t *testing.T) {
	// Init
	type mockBehavior func(artists []models.Artist)

	dbMock, sqlxMock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("%v", err)
	}
	defer dbMock.Close()

	c := gomock.NewController(t)

	l := commonTests.MockLogger(c)

	tablesMock := artistMocks.NewMockTables(c)

	repo := NewPostgreSQL(sqlx.NewDb(dbMock, "postgres"), tablesMock, l)

	// Test filling
	defaultArtists := []models.Artist{
		{
			ID:        1,
			Name:      "Oxxxymiron",
			AvatarSrc: "/artists/avatars/oxxxymiron.png",
		},
		{
			ID:        2,
			Name:      "SALUKI",
			AvatarSrc: "/artists/avatars/saluki.png",
		},
	}

	testTable := []struct {
		name            string
		mockBehavior    mockBehavior
		expectedArtists []models.Artist
		expectError     bool
		expectedError   error
	}{
		{
			name: "Common",
			mockBehavior: func(a []models.Artist) {
				tablesMock.EXPECT().Artists().Return(artistTable)

				rows := sqlxMock.NewRows([]string{"id", "name", "avatar_src"}).
					AddRow(a[0].ID, a[0].Name, a[0].AvatarSrc).
					AddRow(a[1].ID, a[1].Name, a[1].AvatarSrc)
				sqlxMock.ExpectQuery("SELECT (.+) FROM " + artistTable).
					WillReturnRows(rows)
			},
			expectedArtists: defaultArtists,
		},
		{
			name: "Internal PostgreSQL Error",
			mockBehavior: func(a []models.Artist) {
				tablesMock.EXPECT().Artists().Return(artistTable)

				sqlxMock.ExpectQuery("SELECT (.+) FROM " + artistTable).
					WillReturnError(errPqInternal)
			},
			expectError:   true,
			expectedError: errPqInternal,
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			// Call mock
			tc.mockBehavior(tc.expectedArtists)

			a, err := repo.GetFeed()

			// Test
			if tc.expectError {
				assert.ErrorAs(t, err, &tc.expectedError)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedArtists, a)
			}
		})
	}
}

func TestArtistRepositoryGetByAlbum(t *testing.T) {
	// Init
	type mockBehavior func(albumID uint32, artists []models.Artist)

	dbMock, sqlxMock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("%v", err)
	}
	defer dbMock.Close()

	c := gomock.NewController(t)

	l := commonTests.MockLogger(c)

	tablesMock := artistMocks.NewMockTables(c)

	repo := NewPostgreSQL(sqlx.NewDb(dbMock, "postgres"), tablesMock, l)

	// Test filling
	const defaultAlbumID uint32 = 1

	defaultArtists := []models.Artist{
		{
			ID:        1,
			Name:      "Oxxxymiron",
			AvatarSrc: "/artists/avatars/oxxxymiron.png",
		},
		{
			ID:        2,
			Name:      "SALUKI",
			AvatarSrc: "/artists/avatars/saluki.png",
		},
	}

	testTable := []struct {
		name            string
		albumID         uint32
		mockBehavior    mockBehavior
		expectedArtists []models.Artist
		expectError     bool
		expectedError   error
	}{
		{
			name:    "Common",
			albumID: defaultAlbumID,
			mockBehavior: func(albumID uint32, a []models.Artist) {
				tablesMock.EXPECT().Artists().Return(artistTable)
				tablesMock.EXPECT().ArtistsAlbums().Return(artistsAlbumsTable)

				rows := sqlxMock.NewRows([]string{"id", "name", "avatar_src"}).
					AddRow(a[0].ID, a[0].Name, a[0].AvatarSrc).
					AddRow(a[1].ID, a[1].Name, a[1].AvatarSrc)
				sqlxMock.ExpectQuery(fmt.Sprintf("SELECT (.+) FROM %s a INNER JOIN %s",
					artistTable, artistsAlbumsTable)).
					WithArgs(albumID).
					WillReturnRows(rows)
			},
			expectedArtists: defaultArtists,
		},
		{
			name:    "No Such Album",
			albumID: defaultAlbumID,
			mockBehavior: func(albumID uint32, a []models.Artist) {
				tablesMock.EXPECT().Artists().Return(artistTable)
				tablesMock.EXPECT().ArtistsAlbums().Return(artistsAlbumsTable)

				sqlxMock.ExpectQuery(fmt.Sprintf("SELECT (.+) FROM %s a INNER JOIN %s",
					artistTable, artistsAlbumsTable)).
					WithArgs(albumID).
					WillReturnError(sql.ErrNoRows)
			},
			expectError:   true,
			expectedError: &models.NoSuchAlbumError{AlbumID: defaultAlbumID},
		},
		{
			name:    "Internal PostgreSQL Error",
			albumID: defaultAlbumID,
			mockBehavior: func(albumID uint32, a []models.Artist) {
				tablesMock.EXPECT().Artists().Return(artistTable)
				tablesMock.EXPECT().ArtistsAlbums().Return(artistsAlbumsTable)

				sqlxMock.ExpectQuery(fmt.Sprintf("SELECT (.+) FROM %s a INNER JOIN %s",
					artistTable, artistsAlbumsTable)).
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
			tc.mockBehavior(tc.albumID, tc.expectedArtists)

			a, err := repo.GetByAlbum(tc.albumID)

			// Test
			if tc.expectError {
				assert.ErrorContains(t, err, tc.expectedError.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedArtists, a)
			}
		})
	}
}

func TestArtistRepositoryGetByTrack(t *testing.T) {
	// Init
	type mockBehavior func(trackID uint32, artist []models.Artist)

	dbMock, sqlxMock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("%v", err)
	}
	defer dbMock.Close()

	c := gomock.NewController(t)

	l := commonTests.MockLogger(c)

	tablesMock := artistMocks.NewMockTables(c)

	repo := NewPostgreSQL(sqlx.NewDb(dbMock, "postgres"), tablesMock, l)

	// Test filling
	const defaultTrackID uint32 = 1

	defaultArtists := []models.Artist{
		{
			ID:        1,
			Name:      "Oxxxymiron",
			AvatarSrc: "/artists/avatars/oxxxymiron.png",
		},
		{
			ID:        2,
			Name:      "SALUKI",
			AvatarSrc: "/artists/avatars/saluki.png",
		},
	}

	testTable := []struct {
		name            string
		trackID         uint32
		mockBehavior    mockBehavior
		expectedArtists []models.Artist
		expectError     bool
		expectedError   error
	}{
		{
			name:    "Common",
			trackID: defaultTrackID,
			mockBehavior: func(trackID uint32, a []models.Artist) {
				tablesMock.EXPECT().Artists().Return(artistTable)
				tablesMock.EXPECT().ArtistsTracks().Return(artistsTracksTable)

				row := sqlxMock.NewRows([]string{"id", "name", "avatar_src"}).
					AddRow(a[0].ID, a[0].Name, a[0].AvatarSrc).
					AddRow(a[1].ID, a[1].Name, a[1].AvatarSrc)
				sqlxMock.ExpectQuery(fmt.Sprintf("SELECT (.+) FROM %s a INNER JOIN %s",
					artistTable, artistsTracksTable)).
					WithArgs(trackID).
					WillReturnRows(row)
			},
			expectedArtists: defaultArtists,
		},
		{
			name:    "No Such Track",
			trackID: defaultTrackID,
			mockBehavior: func(trackID uint32, a []models.Artist) {
				tablesMock.EXPECT().Artists().Return(artistTable)
				tablesMock.EXPECT().ArtistsTracks().Return(artistsTracksTable)

				sqlxMock.ExpectQuery(fmt.Sprintf("SELECT (.+) FROM %s a INNER JOIN %s",
					artistTable, artistsTracksTable)).
					WithArgs(trackID).
					WillReturnError(sql.ErrNoRows)
			},
			expectError:   true,
			expectedError: &models.NoSuchTrackError{TrackID: defaultTrackID},
		},
		{
			name:    "Internal PostgreSQL Error",
			trackID: defaultTrackID,
			mockBehavior: func(trackID uint32, a []models.Artist) {
				tablesMock.EXPECT().Artists().Return(artistTable)
				tablesMock.EXPECT().ArtistsTracks().Return(artistsTracksTable)

				sqlxMock.ExpectQuery(fmt.Sprintf("SELECT (.+) FROM %s a INNER JOIN %s",
					artistTable, artistsTracksTable)).
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
			tc.mockBehavior(tc.trackID, tc.expectedArtists)

			a, err := repo.GetByTrack(tc.trackID)

			// Test
			if tc.expectError {
				assert.ErrorContains(t, err, tc.expectedError.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedArtists, a)
			}
		})
	}
}

func TestArtistRepositoryGetLikedByUser(t *testing.T) {
	// Init
	type mockBehavior func(userID uint32, artists []models.Artist)

	dbMock, sqlxMock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("%v", err)
	}
	defer dbMock.Close()

	c := gomock.NewController(t)

	l := commonTests.MockLogger(c)

	tablesMock := artistMocks.NewMockTables(c)

	repo := NewPostgreSQL(sqlx.NewDb(dbMock, "postgres"), tablesMock, l)

	// Test filling
	const defaultUserID uint32 = 1

	defaultArtists := []models.Artist{
		{
			ID:        1,
			Name:      "Oxxxymiron",
			AvatarSrc: "/artists/avatars/oxxxymiron.png",
		},
		{
			ID:        2,
			Name:      "SALUKI",
			AvatarSrc: "/artists/avatars/saluki.png",
		},
	}

	testTable := []struct {
		name            string
		userID          uint32
		mockBehavior    mockBehavior
		expectedArtists []models.Artist
		expectError     bool
		expectedError   error
	}{
		{
			name:   "Common",
			userID: defaultUserID,
			mockBehavior: func(userID uint32, a []models.Artist) {
				tablesMock.EXPECT().Artists().Return(artistTable)
				tablesMock.EXPECT().LikedArtists().Return(likedArtistsTable)

				rows := sqlxMock.NewRows([]string{"id", "name", "avatar_src"}).
					AddRow(a[0].ID, a[0].Name, a[0].AvatarSrc).
					AddRow(a[1].ID, a[1].Name, a[1].AvatarSrc)
				sqlxMock.ExpectQuery(fmt.Sprintf("SELECT (.+) FROM %s a INNER JOIN %s",
					artistTable, likedArtistsTable)).
					WithArgs(userID).
					WillReturnRows(rows)
			},
			expectedArtists: defaultArtists,
		},
		{
			name:   "No Such User",
			userID: defaultUserID,
			mockBehavior: func(userID uint32, a []models.Artist) {
				tablesMock.EXPECT().Artists().Return(artistTable)
				tablesMock.EXPECT().LikedArtists().Return(likedArtistsTable)

				sqlxMock.ExpectQuery(fmt.Sprintf("SELECT (.+) FROM %s a INNER JOIN %s",
					artistTable, likedArtistsTable)).
					WithArgs(userID).
					WillReturnError(sql.ErrNoRows)
			},
			expectError:   true,
			expectedError: &models.NoSuchUserError{UserID: defaultUserID},
		},
		{
			name:   "Internal PostgreSQL Error",
			userID: defaultUserID,
			mockBehavior: func(userID uint32, a []models.Artist) {
				tablesMock.EXPECT().Artists().Return(artistTable)
				tablesMock.EXPECT().LikedArtists().Return(likedArtistsTable)

				sqlxMock.ExpectQuery(fmt.Sprintf("SELECT (.+) FROM %s a INNER JOIN %s",
					artistTable, likedArtistsTable)).
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
			tc.mockBehavior(tc.userID, tc.expectedArtists)

			a, err := repo.GetLikedByUser(tc.userID)

			// Test
			if tc.expectError {
				assert.ErrorContains(t, err, tc.expectedError.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedArtists, a)
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

	l := commonTests.MockLogger(c)

	tablesMock := artistMocks.NewMockTables(c)

	repo := NewPostgreSQL(sqlx.NewDb(dbMock, "postgres"), tablesMock, l)

	// Test filling
	const defaultArtistToLikeID uint32 = 1
	const defaultLikedUserID uint32 = 1

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
			name:     "No Such Artist",
			likeInfo: defaultLikeInfo,
			mockBehavior: func(albumID uint32, userID uint32) {
				tablesMock.EXPECT().LikedArtists().Return(likedArtistsTable)

				sqlxMock.ExpectExec("INSERT INTO "+likedArtistsTable).
					WithArgs(albumID, userID).
					WillReturnError(sql.ErrNoRows)
			},
			expectError:   true,
			expectedError: &models.NoSuchArtistError{ArtistID: defaultArtistToLikeID},
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

			inserted, err := repo.InsertLike(tc.likeInfo.artistID, tc.likeInfo.userID)

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

func TestAlbumRepositoryDeleteLike(t *testing.T) {
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

	l := commonTests.MockLogger(c)

	tablesMock := artistMocks.NewMockTables(c)

	repo := NewPostgreSQL(sqlx.NewDb(dbMock, "postgres"), tablesMock, l)

	// Test filling
	defaultLikeInfo := LikeInfo{
		artistID: 1,
		userID:   1,
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
			mockBehavior: func(albumID uint32, userID uint32) {
				tablesMock.EXPECT().LikedArtists().Return(likedArtistsTable)

				sqlxMock.ExpectExec("DELETE FROM "+likedArtistsTable).
					WithArgs(albumID, userID).
					WillReturnResult(driver.RowsAffected(1))
			},
			expectInserted: true,
		},
		{
			name:     "No Such Like",
			likeInfo: defaultLikeInfo,
			mockBehavior: func(albumID, userID uint32) {
				tablesMock.EXPECT().LikedArtists().Return(likedArtistsTable)

				sqlxMock.ExpectExec("DELETE FROM "+likedArtistsTable).
					WithArgs(albumID, userID).
					WillReturnResult(driver.RowsAffected(0))
			},
			expectInserted: false,
		},
		{
			name:     "Internal PostgreSQL Error",
			likeInfo: defaultLikeInfo,
			mockBehavior: func(albumID uint32, userID uint32) {
				tablesMock.EXPECT().LikedArtists().Return(likedArtistsTable)

				sqlxMock.ExpectExec("DELETE FROM "+likedArtistsTable).
					WithArgs(albumID, userID).
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

			inserted, err := repo.DeleteLike(tc.likeInfo.artistID, tc.likeInfo.userID)

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
