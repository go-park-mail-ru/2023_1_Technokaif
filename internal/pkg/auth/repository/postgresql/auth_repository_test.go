package postgresql

import (
	"database/sql"
	"fmt"
	"testing"
	"time"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/golang/mock/gomock"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	commonTests "github.com/go-park-mail-ru/2023_1_Technokaif/internal/common/tests"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/models"
	authMocks "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/auth/mocks"
)

const usersTable = "Users"

func defaultUser() (models.User, error) {
	birthTime, err := time.Parse(time.RFC3339, "2003-08-23T00:00:00Z")
	if err != nil {
		return models.User{}, fmt.Errorf("can't Parse birth date: %v", err)
	}
	birthDate := models.Date{Time: birthTime}

	return models.User{
		ID:        1,
		Version:   1,
		Username:  "yarik_tri",
		Email:     "yarik1448kuzmin@gmail.com",
		Password:  "hash_password",
		Salt:      "salt",
		FirstName: "Yaroslav",
		LastName:  "Kuzmin",
		Sex:       models.Male,
		BirthDate: birthDate,
		AvatarSrc: "/yarik_champion.png",
	}, nil
}

var errPqInternal = errors.New("postgres is dead")

func TestAuthPostgresGetUserByAuthData(t *testing.T) {
	// Init
	type mockBehavior func(userID, userVersion uint32, u *models.User)

	type authData struct {
		userID      uint32
		userVersion uint32
	}

	dbMock, sqlMock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("%v", err)
	}
	defer dbMock.Close()

	c := gomock.NewController(t)

	l := commonTests.MockLogger(c)

	tablesMock := authMocks.NewMockTables(c)

	repo := NewPostgreSQL(sqlx.NewDb(dbMock, "postgres"), tablesMock, l)

	u, err := defaultUser()
	require.NoError(t, err, "can't create default user")

	testTable := []struct {
		name          string
		authData      authData
		mockBehavior  mockBehavior
		expectedUser  *models.User
		expectError   bool
		expectedError error
	}{
		{
			name: "Common",
			authData: authData{
				userID:      1,
				userVersion: 1,
			},
			mockBehavior: func(userID, userVersion uint32, u *models.User) {
				tablesMock.EXPECT().Users().Return(usersTable)

				row := sqlmock.
					NewRows([]string{"id", "version", "username", "email", "password_hash",
						"salt", "first_name", "last_name", "sex", "birth_date", "avatar_src"}).
					AddRow(u.ID, u.Version, u.Username, u.Email, u.Password, u.Salt,
						u.FirstName, u.LastName, u.Sex, u.BirthDate.Time, u.AvatarSrc)

				sqlMock.ExpectQuery("SELECT (.+) FROM "+usersTable).
					WithArgs(userID, userVersion).
					WillReturnRows(row)
			},
			expectedUser: &u,
		},
		{
			name: "No such user",
			authData: authData{
				userID:      1,
				userVersion: 2,
			},
			mockBehavior: func(userID, userVersion uint32, u *models.User) {
				tablesMock.EXPECT().Users().Return(usersTable)

				sqlMock.ExpectQuery("SELECT (.+) FROM "+usersTable).
					WithArgs(userID, userVersion).
					WillReturnError(sql.ErrNoRows)
			},
			expectError:   true,
			expectedError: &models.NoSuchUserError{UserID: 1},
		},
		{
			name: "Internal postgres error",
			authData: authData{
				userID:      1,
				userVersion: 2,
			},
			mockBehavior: func(userID, userVersion uint32, u *models.User) {
				tablesMock.EXPECT().Users().Return(usersTable)

				sqlMock.ExpectQuery("SELECT (.+) FROM "+usersTable).
					WithArgs(userID, userVersion).
					WillReturnError(errPqInternal)
			},
			expectError:   true,
			expectedError: errPqInternal,
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			// Call mock
			tc.mockBehavior(tc.authData.userID, tc.authData.userVersion, tc.expectedUser)

			user, err := repo.GetUserByAuthData(tc.authData.userID, tc.authData.userVersion)

			// Test
			if tc.expectError {
				assert.ErrorContains(t, err, tc.expectedError.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedUser, user)
			}
		})
	}
}

func TestAuthPostgresIncreaseUserVersion(t *testing.T) {
	// Init
	type mockBehavior func(userID uint32)

	dbMock, sqlxMock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("%v", err)
	}
	defer dbMock.Close()

	c := gomock.NewController(t)

	l := commonTests.MockLogger(c)

	tablesMock := authMocks.NewMockTables(c)

	repo := NewPostgreSQL(sqlx.NewDb(dbMock, "postgres"), tablesMock, l)

	testTable := []struct {
		name          string
		userID        uint32
		mockBehavior  mockBehavior
		expectedId    uint32
		expectError   bool
		expectedError error
	}{
		{
			name:   "Common",
			userID: 1,
			mockBehavior: func(userID uint32) {
				tablesMock.EXPECT().Users().Return(usersTable)

				row := sqlmock.NewRows([]string{"id"}).AddRow(userID)

				sqlxMock.ExpectQuery("UPDATE " + usersTable).
					WithArgs(userID).
					WillReturnRows(row)
			},
			expectedId: 1,
		},
		{
			name:   "No such user",
			userID: 1,
			mockBehavior: func(userID uint32) {
				tablesMock.EXPECT().Users().Return(usersTable)

				sqlxMock.ExpectQuery("UPDATE " + usersTable).
					WithArgs(userID).
					WillReturnError(sql.ErrNoRows)
			},
			expectError:   true,
			expectedError: &models.NoSuchUserError{UserID: 1},
		},
		{
			name:   "Internal postgres error",
			userID: 1,
			mockBehavior: func(userID uint32) {
				tablesMock.EXPECT().Users().Return(usersTable)

				sqlxMock.ExpectQuery("UPDATE " + usersTable).
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
			tc.mockBehavior(tc.userID)

			// Test
			err := repo.IncreaseUserVersion(tc.userID)
			if tc.expectError {
				assert.ErrorContains(t, err, tc.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
