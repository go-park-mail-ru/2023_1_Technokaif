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
	albumMocks "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/album/mocks"
)

var ctx = context.Background()

const albumTable = "Albums"
const trackTable = "Track"
const artistsAlbumsTable = "Artists_Albums"
const likedAlbumsTable = "Liked_albums"

var errPqInternal = errors.New("postgres is dead")

func TestAlbumRepositoryPostgreSQL_Check(t *testing.T) {
	// Init
	type mockBehavior func(albumID uint32)

	dbMock, sqlxMock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("%v", err)
	}
	defer dbMock.Close()

	c := gomock.NewController(t)

	tablesMock := albumMocks.NewMockTables(c)

	repo := NewPostgreSQL(sqlx.NewDb(dbMock, "postgres"), tablesMock)

	// Test filling
	const defaultAlbumToCheckID uint32 = 1

	testTable := []struct {
		name           string
		albumToCheckID uint32
		mockBehavior   mockBehavior
		expectError    bool
		expectedError  error
	}{
		{
			name:           "Common",
			albumToCheckID: defaultAlbumToCheckID,
			mockBehavior: func(artistID uint32) {
				tablesMock.EXPECT().Albums().Return(albumTable)

				row := sqlxMock.NewRows([]string{"exists"}).AddRow(true)
				sqlxMock.ExpectQuery("SELECT EXISTS").
					WithArgs(artistID).
					WillReturnRows(row)
			},
		},
		{
			name:           "No Such Album",
			albumToCheckID: defaultAlbumToCheckID,
			mockBehavior: func(artistID uint32) {
				tablesMock.EXPECT().Albums().Return(albumTable)

				row := sqlxMock.NewRows([]string{"exists"}).AddRow(false)
				sqlxMock.ExpectQuery("SELECT EXISTS").
					WithArgs(artistID).
					WillReturnRows(row)
			},
			expectError:   true,
			expectedError: &models.NoSuchAlbumError{AlbumID: defaultAlbumToCheckID},
		},
		{
			name:           "Internal PostgreSQL Error",
			albumToCheckID: defaultAlbumToCheckID,
			mockBehavior: func(artistID uint32) {
				tablesMock.EXPECT().Albums().Return(albumTable)

				sqlxMock.ExpectQuery("SELECT EXISTS").
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
			tc.mockBehavior(tc.albumToCheckID)

			err := repo.Check(ctx, tc.albumToCheckID)

			// Test
			if tc.expectError {
				assert.ErrorContains(t, err, tc.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestAlbumRepositoryPostgreSQL_Insert(t *testing.T) {
	// Init
	type mockBehavior func(album models.Album, artistsID []uint32, id uint32)

	dbMock, sqlxMock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("%v", err)
	}
	defer dbMock.Close()

	c := gomock.NewController(t)

	tablesMock := albumMocks.NewMockTables(c)

	repo := NewPostgreSQL(sqlx.NewDb(dbMock, "postgres"), tablesMock)

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

			id, err := repo.Insert(ctx, tc.album, tc.artistsID)

			// Test
			if tc.expectError {
				assert.ErrorIs(t, err, tc.expectedError)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedID, id)
			}
		})
	}
}

func TestAlbumRepositoryPostgreSQL_GetByID(t *testing.T) {
	// Init
	type mockBehavior func(albumID uint32, a models.Album)

	dbMock, sqlxMock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("%v", err)
	}
	defer dbMock.Close()

	c := gomock.NewController(t)

	tablesMock := albumMocks.NewMockTables(c)

	repo := NewPostgreSQL(sqlx.NewDb(dbMock, "postgres"), tablesMock)

	// Test filling
	const defaultAlbumToGetID uint32 = 1

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
		expectedAlbum models.Album
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
			expectedAlbum: defaultAlbum,
		},
		{
			name:         "No Such Album",
			albumToGetID: defaultAlbumToGetID,
			mockBehavior: func(albumID uint32, a models.Album) {
				tablesMock.EXPECT().Albums().Return(albumTable)

				sqlxMock.ExpectQuery("SELECT (.+) FROM " + albumTable).
					WithArgs(albumID).
					WillReturnError(sql.ErrNoRows)
			},
			expectError:   true,
			expectedError: &models.NoSuchAlbumError{AlbumID: defaultAlbumToGetID},
		},
		{
			name:         "Internal PostgreSQL Error",
			albumToGetID: defaultAlbumToGetID,
			mockBehavior: func(albumID uint32, a models.Album) {
				tablesMock.EXPECT().Albums().Return(albumTable)

				sqlxMock.ExpectQuery("SELECT (.+) FROM " + albumTable).
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
			tc.mockBehavior(tc.albumToGetID, tc.expectedAlbum)

			a, err := repo.GetByID(ctx, tc.albumToGetID)

			// Test
			if tc.expectError {
				assert.ErrorContains(t, err, tc.expectedError.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedAlbum, *a)
			}
		})
	}
}

func TestAlbumRepositoryPostgreSQL_DeleteByID(t *testing.T) {
	// Init
	type mockBehavior func(albumID uint32)

	dbMock, sqlxMock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("%v", err)
	}
	defer dbMock.Close()

	c := gomock.NewController(t)

	tablesMock := albumMocks.NewMockTables(c)

	repo := NewPostgreSQL(sqlx.NewDb(dbMock, "postgres"), tablesMock)

	// Test filling
	const defaultAlbumToDeleteID uint32 = 1

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
					WillReturnResult(driver.RowsAffected(1))
			},
		},
		{
			name:            "No Such Album",
			albumToDeleteID: defaultAlbumToDeleteID,
			mockBehavior: func(albumID uint32) {
				tablesMock.EXPECT().Albums().Return(albumTable)

				sqlxMock.ExpectExec("DELETE FROM " + albumTable).
					WithArgs(albumID).
					WillReturnResult(driver.RowsAffected(0))
			},
			expectError:   true,
			expectedError: &models.NoSuchAlbumError{AlbumID: defaultAlbumToDeleteID},
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

			err := repo.DeleteByID(ctx, tc.albumToDeleteID)

			// Test
			if tc.expectError {
				assert.ErrorContains(t, err, tc.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestAlbumRepositoryPostgreSQL_GetFeed(t *testing.T) {
	// Init
	type mockBehavior func(albums []models.Album)

	dbMock, sqlxMock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("%v", err)
	}
	defer dbMock.Close()

	c := gomock.NewController(t)

	tablesMock := albumMocks.NewMockTables(c)

	repo := NewPostgreSQL(sqlx.NewDb(dbMock, "postgres"), tablesMock)

	// Test filling
	descriptionID1 := "Антиутопия"
	descriptionID2 := "Крутой альбом от крутого дуета"
	defaultAlbums := []models.Album{
		{
			ID:          1,
			Name:        "Горгород",
			Description: &descriptionID1,
			CoverSrc:    "/albums/covers/gorgorod.png",
		},
		{
			ID:          2,
			Name:        "Стыд или Слава",
			Description: &descriptionID2,
			CoverSrc:    "/albums/covers/shameorglory.png",
		},
	}

	testTable := []struct {
		name           string
		mockBehavior   mockBehavior
		expectedAlbums []models.Album
		expectError    bool
		expectedError  error
	}{
		{
			name: "Common",
			mockBehavior: func(a []models.Album) {
				tablesMock.EXPECT().Albums().Return(albumTable)

				rows := sqlxMock.NewRows([]string{"id", "name", "description", "cover_src"})
				for ind := range a {
					rows.AddRow(a[ind].ID, a[ind].Name, a[ind].Description, a[ind].CoverSrc)
				}
				sqlxMock.ExpectQuery("SELECT (.+) FROM " + albumTable).
					WillReturnRows(rows)
			},
			expectedAlbums: defaultAlbums,
		},
		{
			name: "Internal PostgreSQL Error",
			mockBehavior: func(a []models.Album) {
				tablesMock.EXPECT().Albums().Return(albumTable)

				sqlxMock.ExpectQuery("SELECT (.+) FROM " + albumTable).
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
			tc.mockBehavior(tc.expectedAlbums)

			a, err := repo.GetFeed(ctx, feedAmountLimit)

			// Test
			if tc.expectError {
				assert.ErrorAs(t, err, &tc.expectedError)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedAlbums, a)
			}
		})
	}
}

func TestAlbumRepositoryPostgreSQL_GetByArtist(t *testing.T) {
	// Init
	type mockBehavior func(artistID uint32, albums []models.Album)

	dbMock, sqlxMock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("%v", err)
	}
	defer dbMock.Close()

	c := gomock.NewController(t)

	tablesMock := albumMocks.NewMockTables(c)

	repo := NewPostgreSQL(sqlx.NewDb(dbMock, "postgres"), tablesMock)

	// Test filling
	const defaultArtistID uint32 = 1

	descriptionID1 := "Антиутопия"
	descriptionID2 := "Грайм из Лондона"
	defaultAlbums := []models.Album{
		{
			ID:          1,
			Name:        "Горгород",
			Description: &descriptionID1,
			CoverSrc:    "/albums/covers/gorgorod.png",
		},
		{
			ID:          2,
			Name:        "Mixxxtape II",
			Description: &descriptionID2,
			CoverSrc:    "/albums/covers/mixxxtapeii.png",
		},
	}

	testTable := []struct {
		name           string
		artistID       uint32
		mockBehavior   mockBehavior
		expectedAlbums []models.Album
		expectError    bool
		expectedError  error
	}{
		{
			name:     "Common",
			artistID: defaultArtistID,
			mockBehavior: func(artistID uint32, a []models.Album) {
				tablesMock.EXPECT().Albums().Return(albumTable)
				tablesMock.EXPECT().ArtistsAlbums().Return(artistsAlbumsTable)

				rows := sqlxMock.NewRows([]string{"id", "name", "description", "cover_src"})
				for ind := range a {
					rows.AddRow(a[ind].ID, a[ind].Name, a[ind].Description, a[ind].CoverSrc)
				}
				sqlxMock.ExpectQuery(fmt.Sprintf("SELECT (.+) FROM %s a INNER JOIN %s",
					albumTable, artistsAlbumsTable)).
					WithArgs(artistID).
					WillReturnRows(rows)
			},
			expectedAlbums: defaultAlbums,
		},
		{
			name:     "No Such Artist",
			artistID: defaultArtistID,
			mockBehavior: func(artistID uint32, a []models.Album) {
				tablesMock.EXPECT().Albums().Return(albumTable)
				tablesMock.EXPECT().ArtistsAlbums().Return(artistsAlbumsTable)

				sqlxMock.ExpectQuery(fmt.Sprintf("SELECT (.+) FROM %s a INNER JOIN %s",
					albumTable, artistsAlbumsTable)).
					WithArgs(artistID).
					WillReturnError(sql.ErrNoRows)
			},
			expectError:   true,
			expectedError: &models.NoSuchArtistError{ArtistID: defaultArtistID},
		},
		{
			name:     "Internal PostgreSQL Error",
			artistID: defaultArtistID,
			mockBehavior: func(artistID uint32, a []models.Album) {
				tablesMock.EXPECT().Albums().Return(albumTable)
				tablesMock.EXPECT().ArtistsAlbums().Return(artistsAlbumsTable)

				sqlxMock.ExpectQuery(fmt.Sprintf("SELECT (.+) FROM %s a INNER JOIN %s",
					albumTable, artistsAlbumsTable)).
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
			tc.mockBehavior(tc.artistID, tc.expectedAlbums)

			a, err := repo.GetByArtist(ctx, tc.artistID)

			// Test
			if tc.expectError {
				assert.ErrorContains(t, err, tc.expectedError.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedAlbums, a)
			}
		})
	}
}

func TestAlbumRepositoryPostgreSQL_GetByTrack(t *testing.T) {
	// Init
	type mockBehavior func(trackID uint32, album models.Album)

	dbMock, sqlxMock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("%v", err)
	}
	defer dbMock.Close()

	c := gomock.NewController(t)

	tablesMock := albumMocks.NewMockTables(c)

	repo := NewPostgreSQL(sqlx.NewDb(dbMock, "postgres"), tablesMock)

	// Test filling
	const defaultTrackID uint32 = 1

	description := "Антиутопия"
	defaultAlbum := models.Album{
		ID:          1,
		Name:        "Горгород",
		Description: &description,
		CoverSrc:    "/albums/covers/gorgorod.png",
	}

	testTable := []struct {
		name          string
		trackID       uint32
		mockBehavior  mockBehavior
		expectedAlbum models.Album
		expectError   bool
		expectedError error
	}{
		{
			name:    "Common",
			trackID: defaultTrackID,
			mockBehavior: func(trackID uint32, a models.Album) {
				tablesMock.EXPECT().Albums().Return(albumTable)
				tablesMock.EXPECT().Tracks().Return(trackTable)

				row := sqlxMock.NewRows([]string{"id", "name", "description", "cover_src"}).
					AddRow(a.ID, a.Name, a.Description, a.CoverSrc)
				sqlxMock.ExpectQuery(fmt.Sprintf("SELECT (.+) FROM %s a INNER JOIN %s",
					albumTable, trackTable)).
					WithArgs(trackID).
					WillReturnRows(row)
			},
			expectedAlbum: defaultAlbum,
		},
		{
			name:    "No Such Track",
			trackID: defaultTrackID,
			mockBehavior: func(trackID uint32, a models.Album) {
				tablesMock.EXPECT().Albums().Return(albumTable)
				tablesMock.EXPECT().Tracks().Return(trackTable)

				sqlxMock.ExpectQuery(fmt.Sprintf("SELECT (.+) FROM %s a INNER JOIN %s",
					albumTable, trackTable)).
					WithArgs(trackID).
					WillReturnError(sql.ErrNoRows)
			},
			expectError:   true,
			expectedError: &models.NoSuchTrackError{TrackID: defaultTrackID},
		},
		{
			name:    "Internal PostgreSQL Error",
			trackID: defaultTrackID,
			mockBehavior: func(trackID uint32, a models.Album) {
				tablesMock.EXPECT().Albums().Return(albumTable)
				tablesMock.EXPECT().Tracks().Return(trackTable)

				sqlxMock.ExpectQuery(fmt.Sprintf("SELECT (.+) FROM %s a INNER JOIN %s",
					albumTable, trackTable)).
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
			tc.mockBehavior(tc.trackID, tc.expectedAlbum)

			a, err := repo.GetByTrack(ctx, tc.trackID)

			// Test
			if tc.expectError {
				assert.ErrorContains(t, err, tc.expectedError.Error())
			} else {
				assert.Equal(t, tc.expectedAlbum, *a)
				assert.NoError(t, err)
			}
		})
	}
}

func TestAlbumRepositoryPostgreSQL_GetLikedByUser(t *testing.T) {
	// Init
	type mockBehavior func(userID uint32, albums []models.Album)

	dbMock, sqlxMock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("%v", err)
	}
	defer dbMock.Close()

	c := gomock.NewController(t)

	tablesMock := albumMocks.NewMockTables(c)

	repo := NewPostgreSQL(sqlx.NewDb(dbMock, "postgres"), tablesMock)

	// Test filling
	const defaultUserID uint32 = 1

	descriptionID1 := "Антиутопия"
	descriptionID2 := "Грайм из Лондона"
	defaultAlbums := []models.Album{
		{
			ID:          1,
			Name:        "Горгород",
			Description: &descriptionID1,
			CoverSrc:    "/albums/covers/gorgorod.png",
		},
		{
			ID:          2,
			Name:        "Mixxxtape II",
			Description: &descriptionID2,
			CoverSrc:    "/albums/covers/mixxxtapeii.png",
		},
	}

	testTable := []struct {
		name           string
		userID         uint32
		mockBehavior   mockBehavior
		expectedAlbums []models.Album
		expectError    bool
		expectedError  error
	}{
		{
			name:   "Common",
			userID: defaultUserID,
			mockBehavior: func(userID uint32, a []models.Album) {
				tablesMock.EXPECT().Albums().Return(albumTable)
				tablesMock.EXPECT().LikedAlbums().Return(likedAlbumsTable)

				rows := sqlxMock.NewRows([]string{"id", "name", "description", "cover_src"})
				for ind := range a {
					rows.AddRow(a[ind].ID, a[ind].Name, a[ind].Description, a[ind].CoverSrc)
				}
				sqlxMock.ExpectQuery(fmt.Sprintf("SELECT (.+) FROM %s a INNER JOIN %s",
					albumTable, likedAlbumsTable)).
					WithArgs(userID).
					WillReturnRows(rows)
			},
			expectedAlbums: defaultAlbums,
		},
		{
			name:   "No Such User",
			userID: defaultUserID,
			mockBehavior: func(userID uint32, a []models.Album) {
				tablesMock.EXPECT().Albums().Return(albumTable)
				tablesMock.EXPECT().LikedAlbums().Return(likedAlbumsTable)

				sqlxMock.ExpectQuery(fmt.Sprintf("SELECT (.+) FROM %s a INNER JOIN %s",
					albumTable, likedAlbumsTable)).
					WithArgs(userID).
					WillReturnError(sql.ErrNoRows)
			},
			expectError:   true,
			expectedError: &models.NoSuchUserError{UserID: defaultUserID},
		},
		{
			name:   "Internal PostgreSQL Error",
			userID: defaultUserID,
			mockBehavior: func(userID uint32, a []models.Album) {
				tablesMock.EXPECT().Albums().Return(albumTable)
				tablesMock.EXPECT().LikedAlbums().Return(likedAlbumsTable)

				sqlxMock.ExpectQuery(fmt.Sprintf("SELECT (.+) FROM %s a INNER JOIN %s",
					albumTable, likedAlbumsTable)).
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
			tc.mockBehavior(tc.userID, tc.expectedAlbums)

			a, err := repo.GetLikedByUser(ctx, tc.userID)

			// Test
			if tc.expectError {
				assert.ErrorContains(t, err, tc.expectedError.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedAlbums, a)
			}
		})
	}
}

func TestAlbumRepositoryPostgreSQL_Like(t *testing.T) {
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

	tablesMock := albumMocks.NewMockTables(c)

	repo := NewPostgreSQL(sqlx.NewDb(dbMock, "postgres"), tablesMock)

	// Test filling
	const defaultAlbumToLikeID uint32 = 1
	const defaultLikedUserID uint32 = 1

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
			mockBehavior: func(albumID uint32, userID uint32) {
				tablesMock.EXPECT().LikedAlbums().Return(likedAlbumsTable)

				sqlxMock.ExpectExec("INSERT INTO "+likedAlbumsTable).
					WithArgs(albumID, userID).
					WillReturnResult(driver.ResultNoRows)
			},
			expectInserted: true,
		},
		{
			name:     "No Such Album",
			likeInfo: defaultLikeInfo,
			mockBehavior: func(albumID uint32, userID uint32) {
				tablesMock.EXPECT().LikedAlbums().Return(likedAlbumsTable)

				sqlxMock.ExpectExec("INSERT INTO "+likedAlbumsTable).
					WithArgs(albumID, userID).
					WillReturnError(sql.ErrNoRows)
			},
			expectError:   true,
			expectedError: &models.NoSuchAlbumError{AlbumID: defaultAlbumToLikeID},
		},
		{
			name:     "Internal PostgreSQL Error",
			likeInfo: defaultLikeInfo,
			mockBehavior: func(albumID uint32, userID uint32) {
				tablesMock.EXPECT().LikedAlbums().Return(likedAlbumsTable)

				sqlxMock.ExpectExec("INSERT INTO "+likedAlbumsTable).
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
			tc.mockBehavior(tc.likeInfo.albumID, tc.likeInfo.userID)

			inserted, err := repo.InsertLike(ctx, tc.likeInfo.albumID, tc.likeInfo.userID)

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

func TestAlbumRepositoryPostgreSQL_DeleteLike(t *testing.T) {
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

	tablesMock := albumMocks.NewMockTables(c)

	repo := NewPostgreSQL(sqlx.NewDb(dbMock, "postgres"), tablesMock)

	// Test filling
	defaultLikeInfo := LikeInfo{
		albumID: 1,
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
			mockBehavior: func(albumID uint32, userID uint32) {
				tablesMock.EXPECT().LikedAlbums().Return(likedAlbumsTable)

				sqlxMock.ExpectExec("DELETE FROM "+likedAlbumsTable).
					WithArgs(albumID, userID).
					WillReturnResult(driver.RowsAffected(1))
			},
			expectInserted: true,
		},
		{
			name:     "No Such Like",
			likeInfo: defaultLikeInfo,
			mockBehavior: func(albumID, userID uint32) {
				tablesMock.EXPECT().LikedAlbums().Return(likedAlbumsTable)

				sqlxMock.ExpectExec("DELETE FROM "+likedAlbumsTable).
					WithArgs(albumID, userID).
					WillReturnResult(driver.RowsAffected(0))
			},
			expectInserted: false,
		},
		{
			name:     "Internal PostgreSQL Error",
			likeInfo: defaultLikeInfo,
			mockBehavior: func(albumID uint32, userID uint32) {
				tablesMock.EXPECT().LikedAlbums().Return(likedAlbumsTable)

				sqlxMock.ExpectExec("DELETE FROM "+likedAlbumsTable).
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
			tc.mockBehavior(tc.likeInfo.albumID, tc.likeInfo.userID)

			inserted, err := repo.DeleteLike(ctx, tc.likeInfo.albumID, tc.likeInfo.userID)

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

func TestAlbumRepositoryPostgreSQL_IsLiked(t *testing.T) {
	// Init
	type mockBehavior func(trackID, userID uint32)

	dbMock, sqlxMock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("%v", err)
	}
	defer dbMock.Close()

	c := gomock.NewController(t)

	tablesMock := albumMocks.NewMockTables(c)

	repo := NewPostgreSQL(sqlx.NewDb(dbMock, "postgres"), tablesMock)

	// Test filling
	const defaultAlbumID uint32 = 1
	const defaultUserID uint32 = 1

	testTable := []struct {
		name          string
		albumID       uint32
		userID        uint32
		mockBehavior  mockBehavior
		expectError   bool
		expectedError error
		isLiked       bool
	}{
		{
			name:    "Liked",
			albumID: defaultAlbumID,
			userID:  defaultUserID,
			mockBehavior: func(albumID, userID uint32) {
				tablesMock.EXPECT().LikedAlbums().Return(likedAlbumsTable)

				row := sqlxMock.NewRows([]string{"exists"}).AddRow(true)
				sqlxMock.ExpectQuery("SELECT EXISTS").
					WithArgs(albumID, userID).
					WillReturnRows(row)
			},
			isLiked: true,
		},
		{
			name:    "Isn't liked",
			albumID: defaultAlbumID,
			userID:  defaultUserID,
			mockBehavior: func(albumID, userID uint32) {
				tablesMock.EXPECT().LikedAlbums().Return(likedAlbumsTable)

				row := sqlxMock.NewRows([]string{"exists"}).AddRow(false)
				sqlxMock.ExpectQuery("SELECT EXISTS").
					WithArgs(albumID, userID).
					WillReturnRows(row)
			},
		},
		{
			name:    "Internal PostgreSQL Error",
			albumID: defaultAlbumID,
			userID:  defaultUserID,
			mockBehavior: func(albumID, userID uint32) {
				tablesMock.EXPECT().LikedAlbums().Return(likedAlbumsTable)

				sqlxMock.ExpectQuery("SELECT EXISTS").
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
			tc.mockBehavior(tc.albumID, tc.userID)

			isLiked, err := repo.IsLiked(ctx, tc.albumID, tc.userID)

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
