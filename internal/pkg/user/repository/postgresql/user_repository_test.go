package postgresql

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/models"
	"github.com/golang/mock/gomock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	userMocks "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/user/mocks"
)

var ctx = context.Background()

const userTable = "Users"
const userPlaylistTable = "Users_Playlists"

var errPqInternal = errors.New("postgres is dead")

func getCorrectUser(t *testing.T) *models.User {
	birthTime, err := time.Parse(time.RFC3339, "2003-08-23T00:00:00Z")
	require.NoError(t, err, "can't Parse birth date")

	birthDate := models.Date{Time: birthTime}

	return &models.User{
		ID:        1,
		Username:  "yarik_tri",
		Email:     "yarik1448kuzmin@gmail.com",
		FirstName: "Yaroslav",
		LastName:  "Kuzmin",
		BirthDate: birthDate,
		AvatarSrc: "/users/avatars/yarik_tri.png",
	}
}

func TestUserRepositoryPostgreSQL_Check(t *testing.T) {
	// Init
	type mockBehavior func(userID uint32)

	dbMock, sqlxMock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("%v", err)
	}
	defer dbMock.Close()

	c := gomock.NewController(t)

	tablesMock := userMocks.NewMockTables(c)

	repo := NewPostgreSQL(sqlx.NewDb(dbMock, "postgres"), tablesMock)

	// Test filling
	const defaultUserToCheckID uint32 = 1

	testTable := []struct {
		name          string
		userToCheckID uint32
		mockBehavior  mockBehavior
		expectError   bool
		expectedError error
	}{
		{
			name:          "Common",
			userToCheckID: defaultUserToCheckID,
			mockBehavior: func(userID uint32) {
				tablesMock.EXPECT().Users().Return(userTable)

				row := sqlxMock.NewRows([]string{"exists"}).AddRow(true)
				sqlxMock.ExpectQuery("SELECT EXISTS").
					WithArgs(userID).
					WillReturnRows(row)
			},
		},
		{
			name:          "No Such User",
			userToCheckID: defaultUserToCheckID,
			mockBehavior: func(userID uint32) {
				tablesMock.EXPECT().Users().Return(userTable)

				row := sqlxMock.NewRows([]string{"exists"}).AddRow(false)
				sqlxMock.ExpectQuery("SELECT EXISTS").
					WithArgs(userID).
					WillReturnRows(row)
			},
			expectError:   true,
			expectedError: &models.NoSuchUserError{UserID: defaultUserToCheckID},
		},
		{
			name:          "Internal PostgreSQL Error",
			userToCheckID: defaultUserToCheckID,
			mockBehavior: func(userID uint32) {
				tablesMock.EXPECT().Users().Return(userTable)

				sqlxMock.ExpectQuery("SELECT EXISTS").
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
			tc.mockBehavior(tc.userToCheckID)

			err := repo.Check(ctx, tc.userToCheckID)

			// Test
			if tc.expectError {
				assert.ErrorContains(t, err, tc.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestUserRepositoryPostgreSQL_GetByID(t *testing.T) {
	// Init
	type mockBehavior func(userID uint32, u models.User)

	dbMock, sqlxMock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("%v", err)
	}
	defer dbMock.Close()

	c := gomock.NewController(t)

	tablesMock := userMocks.NewMockTables(c)

	repo := NewPostgreSQL(sqlx.NewDb(dbMock, "postgres"), tablesMock)

	// Test filling
	const defaultUserToGetID uint32 = 1
	defaultUser := getCorrectUser(t)

	testTable := []struct {
		name          string
		userToGetID   uint32
		mockBehavior  mockBehavior
		expectedUser  models.User
		expectError   bool
		expectedError error
	}{
		{
			name:        "Common",
			userToGetID: defaultUserToGetID,
			mockBehavior: func(userID uint32, u models.User) {
				tablesMock.EXPECT().Users().Return(userTable)

				row := sqlxMock.NewRows(
					[]string{"id", "version", "username", "email", "password_hash", "salt",
						"first_name", "last_name", "birth_date", "avatar_src"}).
					AddRow(u.ID, u.Version, u.Username, u.Email, u.Password, u.Salt,
						u.FirstName, u.LastName, u.BirthDate.Time, u.AvatarSrc)
				sqlxMock.ExpectQuery("SELECT (.+) FROM " + userTable).
					WithArgs(userID).
					WillReturnRows(row)
			},
			expectedUser: *defaultUser,
		},
		{
			name:        "No Such User",
			userToGetID: defaultUserToGetID,
			mockBehavior: func(userID uint32, u models.User) {
				tablesMock.EXPECT().Users().Return(userTable)

				sqlxMock.ExpectQuery("SELECT (.+) FROM " + userTable).
					WithArgs(userID).
					WillReturnError(sql.ErrNoRows)
			},
			expectError:   true,
			expectedError: &models.NoSuchUserError{UserID: defaultUserToGetID},
		},
		{
			name:        "Internal PostgreSQL Error",
			userToGetID: defaultUserToGetID,
			mockBehavior: func(userID uint32, u models.User) {
				tablesMock.EXPECT().Users().Return(userTable)

				sqlxMock.ExpectQuery("SELECT (.+) FROM " + userTable).
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
			tc.mockBehavior(tc.userToGetID, tc.expectedUser)

			u, err := repo.GetByID(ctx, tc.userToGetID)

			// Test
			if tc.expectError {
				assert.ErrorContains(t, err, tc.expectedError.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedUser, *u)
			}
		})
	}
}

func TestUserRepositoryPostgreSQL_GetByUsername(t *testing.T) {
	// Init
	type mockBehavior func(username string, u *models.User)

	dbMock, sqlxMock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("%v", err)
	}
	defer dbMock.Close()

	c := gomock.NewController(t)

	tablesMock := userMocks.NewMockTables(c)

	repo := NewPostgreSQL(sqlx.NewDb(dbMock, "postgres"), tablesMock)

	// Test filling
	const defaultUsername = "yarik_tri"
	defaultUser := getCorrectUser(t)

	testTable := []struct {
		name          string
		username      string
		mockBehavior  mockBehavior
		expectedUser  *models.User
		expectError   bool
		expectedError error
	}{
		{
			name:     "Common",
			username: defaultUsername,
			mockBehavior: func(username string, u *models.User) {
				tablesMock.EXPECT().Users().Return(userTable)

				rows := sqlxMock.NewRows(
					[]string{"id", "version", "username", "email", "password_hash", "salt",
						"first_name", "last_name", "birth_date", "avatar_src"}).
					AddRow(u.ID, u.Version, u.Username, u.Email, u.Password, u.Salt,
						u.FirstName, u.LastName, u.BirthDate.Time, u.AvatarSrc)
				sqlxMock.ExpectQuery("SELECT (.+) FROM " + userTable).
					WithArgs(username).
					WillReturnRows(rows)
			},
			expectedUser: defaultUser,
		},
		{
			name:     "Internal PostgreSQL Error",
			username: defaultUsername,
			mockBehavior: func(username string, u *models.User) {
				tablesMock.EXPECT().Users().Return(userTable)

				sqlxMock.ExpectQuery("SELECT (.+) FROM " + userTable).
					WithArgs(username).
					WillReturnError(errPqInternal)
			},
			expectError:   true,
			expectedError: errPqInternal,
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			// Call mock
			tc.mockBehavior(tc.username, tc.expectedUser)

			u, err := repo.GetUserByUsername(ctx, tc.username)

			// Test
			if tc.expectError {
				assert.ErrorContains(t, err, tc.expectedError.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedUser, u)
			}
		})
	}
}

func TestUserRepositoryPostgreSQL_GetByPlaylist(t *testing.T) {
	// Init
	type mockBehavior func(playlistID uint32, u []models.User)

	dbMock, sqlxMock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("%v", err)
	}
	defer dbMock.Close()

	c := gomock.NewController(t)

	tablesMock := userMocks.NewMockTables(c)

	repo := NewPostgreSQL(sqlx.NewDb(dbMock, "postgres"), tablesMock)

	// Test filling
	const defaultPlaylistID uint32 = 1
	defaultUser := getCorrectUser(t)

	testTable := []struct {
		name          string
		playlistID    uint32
		mockBehavior  mockBehavior
		expectedUsers []models.User
		expectError   bool
		expectedError error
	}{
		{
			name:       "Common",
			playlistID: defaultPlaylistID,
			mockBehavior: func(playlistID uint32, u []models.User) {
				tablesMock.EXPECT().Users().Return(userTable)
				tablesMock.EXPECT().UsersPlaylists().Return(userPlaylistTable)

				rows := sqlxMock.NewRows(
					[]string{"id", "username", "email", "first_name",
						"last_name", "birth_date", "avatar_src"})
				for ind := range u {
					rows.AddRow(u[ind].ID, u[ind].Username, u[ind].Email, u[ind].FirstName,
						u[ind].LastName, u[ind].BirthDate.Time, u[ind].AvatarSrc)
				}
				sqlxMock.ExpectQuery("SELECT (.+) FROM " + userTable).
					WithArgs(playlistID).
					WillReturnRows(rows)
			},
			expectedUsers: []models.User{*defaultUser},
		},
		{
			name:       "Internal PostgreSQL Error",
			playlistID: defaultPlaylistID,
			mockBehavior: func(playlistID uint32, u []models.User) {
				tablesMock.EXPECT().Users().Return(userTable)
				tablesMock.EXPECT().UsersPlaylists().Return(userPlaylistTable)

				sqlxMock.ExpectQuery("SELECT (.+) FROM " + userTable).
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
			tc.mockBehavior(tc.playlistID, tc.expectedUsers)

			u, err := repo.GetByPlaylist(ctx, tc.playlistID)

			// Test
			if tc.expectError {
				assert.ErrorContains(t, err, tc.expectedError.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedUsers, u)
			}
		})
	}
}
