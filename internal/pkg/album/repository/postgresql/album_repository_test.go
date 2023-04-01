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
	albumMocks "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/album/mocks"
	logMocks "github.com/go-park-mail-ru/2023_1_Technokaif/pkg/logger/mocks"
)

const albumTable = "Albums"
const artistsAlbumsTable = "Artists_Albums"
const likedAlbumsTable = "Liked_albums"

var errPqInternal = errors.New("postgres is dead")

func TestAlbumRepositoryInsert(t *testing.T) {
	// Init
	type mockBehavior func(album models.Album, artistsID []uint32, id uint32)

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

	tablesMock := albumMocks.NewMockTables(c)

	repo := NewPostgreSQL(sqlx.NewDb(dbMock, "postgres"), tablesMock, l)

	defaultAlbumToIsert := models.Album{
		Name:     "Горгород",
		CoverSrc: "/albums/covers/gorgorod.png",
	}

	defaultArtistsIDToInsert := []uint32{1, 2, 3}

	testTable := []struct {
		name          string
		album         models.Album
		artistsID     []uint32
		mockBehavior  mockBehavior
		expectedID    uint32
		expectError   bool
		expectedError error
	}{
		{
			name:      "Common",
			album:     defaultAlbumToIsert,
			artistsID: defaultArtistsIDToInsert,
			mockBehavior: func(a models.Album, artistsID []uint32, id uint32) {
				tablesMock.EXPECT().Albums().Return(albumTable)
				tablesMock.EXPECT().ArtistsAlbums().Return(artistsAlbumsTable)

				sqlxMock.ExpectBegin()

				row := sqlxMock.NewRows([]string{"id"}).AddRow(id)
				sqlxMock.ExpectQuery("INSERT INTO "+albumTable).
					WithArgs(a.Name, a.Description, a.CoverSrc).
					WillReturnRows(row)

				for _, artistID := range artistsID {
					sqlxMock.ExpectExec("INSERT INTO "+artistsAlbumsTable).
						WithArgs(artistID, id).
						WillReturnResult(driver.ResultNoRows)
				}

				sqlxMock.ExpectCommit()
			},
			expectedID: 1,
		},
		{
			name:      "Insert Artists to Album Issue",
			album:     defaultAlbumToIsert,
			artistsID: defaultArtistsIDToInsert,
			mockBehavior: func(a models.Album, artistsID []uint32, id uint32) {
				tablesMock.EXPECT().Albums().Return(albumTable)
				tablesMock.EXPECT().ArtistsAlbums().Return(artistsAlbumsTable)

				sqlxMock.ExpectBegin()

				row := sqlxMock.NewRows([]string{"id"}).AddRow(id)
				sqlxMock.ExpectQuery("INSERT INTO "+albumTable).
					WithArgs(a.Name, a.Description, a.CoverSrc).
					WillReturnRows(row)

				sqlxMock.ExpectExec("INSERT INTO "+artistsAlbumsTable).
					WithArgs(artistsID[0], id).
					WillReturnError(errPqInternal)

				sqlxMock.ExpectRollback()
			},
			expectedID:    1,
			expectError:   true,
			expectedError: errPqInternal,
		},
		{
			name:      "Insert Album Issue",
			album:     defaultAlbumToIsert,
			artistsID: defaultArtistsIDToInsert,
			mockBehavior: func(a models.Album, artistsID []uint32, id uint32) {
				tablesMock.EXPECT().Albums().Return(albumTable)

				sqlxMock.ExpectBegin()

				sqlxMock.ExpectQuery("INSERT INTO "+albumTable).
					WithArgs(a.Name, a.Description, a.CoverSrc).
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
			tc.mockBehavior(tc.album, tc.artistsID, tc.expectedID)

			id, err := repo.Insert(tc.album, tc.artistsID)

			// Test
			if tc.expectError {
				assert.ErrorIs(t, err, tc.expectedError)
			} else {
				assert.Equal(t, tc.expectedID, id)
				assert.NoError(t, err)
			}
		})
	}
}

func TestAlbumRepositoryGetByID(t *testing.T) {
	// Init
	type mockBehavior func(albumID uint32, a models.Album)

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

	tablesMock := albumMocks.NewMockTables(c)

	repo := NewPostgreSQL(sqlx.NewDb(dbMock, "postgres"), tablesMock, l)

	// Test filling
	defaultAlbumToGetID := uint32(1)

	description := "Антиутопия"
	defaultAlbum := models.Album{
		ID:          defaultAlbumToGetID,
		Name:        "Горгород",
		Description: &description,
		CoverSrc:    "/albums/covers/gorgorod.png",
	}

	testTable := []struct {
		name          string
		albumToGetID  uint32
		mockBehavior  mockBehavior
		expectedAlbum *models.Album
		expectError   bool
		expectedError error
	}{
		{
			name:         "Common",
			albumToGetID: defaultAlbumToGetID,
			mockBehavior: func(albumID uint32, a models.Album) {
				tablesMock.EXPECT().Albums().Return(albumTable)

				row := sqlxMock.NewRows([]string{"id", "name", "description", "cover_src"}).
					AddRow(a.ID, a.Name, a.Description, a.CoverSrc)
				sqlxMock.ExpectQuery("SELECT (.+) FROM " + albumTable).
					WithArgs(albumID).
					WillReturnRows(row)
			},
			expectedAlbum: &defaultAlbum,
		},
		// {
		// 	name:         "Internal PostgreSQL Error",
		// 	albumToGetID: defaultAlbumToGetID,
		// 	mockBehavior: func(albumID uint32, a models.Album) {
		// 		tablesMock.EXPECT().Albums().Return(albumTable)

		// 		sqlxMock.ExpectQuery("SELECT (.+) FROM " + albumTable).
		// 			WithArgs(albumID).
		// 			WillReturnError(errPqInternal)
		// 	},
		// 	expectError:   true,
		// 	expectedError: errPqInternal,
		// },
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			// Call mock
			tc.mockBehavior(tc.albumToGetID, *tc.expectedAlbum)

			a, err := repo.GetByID(tc.albumToGetID)

			// Test
			if tc.expectError {
				assert.ErrorIs(t, err, tc.expectedError)
			} else {
				assert.Equal(t, tc.expectedAlbum, a)
				assert.NoError(t, err)
			}
		})
	}
}

func TestAlbumRepositoryDeleteByID(t *testing.T) {
	// Init
	type mockBehavior func(albumID uint32)

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

	tablesMock := albumMocks.NewMockTables(c)

	repo := NewPostgreSQL(sqlx.NewDb(dbMock, "postgres"), tablesMock, l)

	// Test filling
	defaultAlbumToDeleteID := uint32(1)

	testTable := []struct {
		name            string
		albumToDeleteID uint32
		mockBehavior    mockBehavior
		expectError     bool
		expectedError   error
	}{
		{
			name:            "Common",
			albumToDeleteID: defaultAlbumToDeleteID,
			mockBehavior: func(albumID uint32) {
				tablesMock.EXPECT().Albums().Return(albumTable)

				sqlxMock.ExpectExec("DELETE FROM " + albumTable).
					WithArgs(albumID).
					WillReturnResult(driver.ResultNoRows)
			},
		},
		{
			name:            "Internal PostgreSQL Error",
			albumToDeleteID: defaultAlbumToDeleteID,
			mockBehavior: func(albumID uint32) {
				tablesMock.EXPECT().Albums().Return(albumTable)

				sqlxMock.ExpectExec("DELETE FROM " + albumTable).
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
			tc.mockBehavior(tc.albumToDeleteID)

			err := repo.DeleteByID(tc.albumToDeleteID)

			// Test
			if tc.expectError {
				assert.ErrorIs(t, err, tc.expectedError)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestAlbumRepositoryLike(t *testing.T) {
	// Init
	type mockBehavior func(albumID uint32, userID uint32)

	type LikeInfo struct {
		albumID uint32
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

	tablesMock := albumMocks.NewMockTables(c)

	repo := NewPostgreSQL(sqlx.NewDb(dbMock, "postgres"), tablesMock, l)

	// Test filling
	defaultAlbumToLikeID := uint32(1)
	defaultLikedUserID := uint32(1)

	defaultLikeInfo := LikeInfo{
		albumID: defaultAlbumToLikeID,
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
				tablesMock.EXPECT().LikedAlbums().Return(likedAlbumsTable)

				sqlxMock.ExpectExec("INSERT INTO "+likedAlbumsTable).
					WithArgs(artistID, userID).
					WillReturnResult(driver.ResultNoRows)
			},
			expectInserted: true,
		},
		{
			name:     "Internal PostgreSQL Error",
			likeInfo: defaultLikeInfo,
			mockBehavior: func(artistID uint32, userID uint32) {
				tablesMock.EXPECT().LikedAlbums().Return(likedAlbumsTable)

				sqlxMock.ExpectExec("INSERT INTO "+likedAlbumsTable).
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
			tc.mockBehavior(tc.likeInfo.albumID, tc.likeInfo.userID)

			success, err := repo.InsertLike(tc.likeInfo.albumID, tc.likeInfo.userID)

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
