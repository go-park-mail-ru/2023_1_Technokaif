package repository

import (
	"errors"
	"fmt"
	"testing"
	"time"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/logger"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/models"
	"github.com/stretchr/testify/assert"
)

func defaultUser() (models.User, error) {
	birthTime, err := time.Parse(time.RFC3339, "2003-08-23T00:00:00Z")
	if err != nil {
		return models.User{}, fmt.Errorf("can't Parse birth date: %v", err)
	}
	birthDate := models.Date{birthTime}

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
		BirhDate:  birthDate,
	}, nil
}

func userWithoutUsername() (models.User, error) {
	defaultUser, err := defaultUser()
	if err != nil {
		return models.User{}, err
	}

	userWithoutUsername := defaultUser
	userWithoutUsername.Username = ""

	return userWithoutUsername, nil
}

func TestAuthPostgresCreateUser(t *testing.T) {
	type mockBehavior func(u models.User, id int)

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("%v", err)
	}
	defer db.Close()

	l, err := logger.NewFLogger()
	if err != nil {
		t.Errorf("%v", err)
	}

	r := NewAuthPostgres(db, l)

	correctTestUser, err := defaultUser()
	if err != nil {
		t.Errorf("can't create default user: %v", err)
	}
	userWithoutUsername, err := userWithoutUsername()
	if err != nil {
		t.Errorf("can't create user without username: %v", err)
	}

	testTable := []struct {
		name          string
		userToCreate  models.User
		mockBehavior  mockBehavior
		expectedId    int
		expectError   bool
		expectedError string
	}{
		{
			name:         "Common",
			userToCreate: correctTestUser,
			mockBehavior: func(u models.User, id int) {
				// Just creates query row (expected result after exec)
				row := sqlmock.NewRows([]string{"id"}).AddRow(id)

				// Exec request and expect rows to be result
				mock.ExpectQuery("INSERT INTO "+UsersTable).
					WithArgs(u.Username, u.Email, u.Password, u.Salt,
						u.FirstName, u.LastName, u.Sex, u.BirhDate.Format(time.RFC3339)).
					WillReturnRows(row)
			},
			expectedId: 1,
		},
		{
			name:         "Empty required Fields",
			userToCreate: userWithoutUsername,
			mockBehavior: func(u models.User, id int) {
				mock.ExpectQuery("INSERT INTO "+UsersTable).
					WithArgs(u.Username, u.Email, u.Password, u.Salt,
						u.FirstName, u.LastName, u.Sex, u.BirhDate.Format(time.RFC3339)).
					WillReturnError(errors.New("required field is null"))
			},
			expectError:   true,
			expectedError: "required field is null",
		},
		{
			name:         "user already exists",
			userToCreate: correctTestUser,
			mockBehavior: func(u models.User, id int) {
				mock.ExpectQuery("INSERT INTO "+UsersTable).
					WithArgs(u.Username, u.Email, u.Password, u.Salt,
						u.FirstName, u.LastName, u.Sex, u.BirhDate.Format(time.RFC3339)).
					WillReturnError(errors.New("user already exists"))
			},
			expectError:   true,
			expectedError: "user already exists",
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockBehavior(tc.userToCreate, tc.expectedId)

			resultId, err := r.CreateUser(tc.userToCreate)
			if tc.expectError {
				assert.EqualError(t, err, tc.expectedError)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedId, resultId)
			}
		})
	}
}

func TestAuthPostgresGetUserByUsername(t *testing.T) {
	type mockBehavior func(username string, u *models.User)

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("%v", err)
	}
	defer db.Close()

	l, err := logger.NewFLogger()
	if err != nil {
		t.Errorf("%v", err)
	}

	r := NewAuthPostgres(db, l)

	u, err := defaultUser()
	if err != nil {
		t.Errorf("can't create default user: %v", err)
	}

	testTable := []struct {
		name          string
		username      string
		mockBehavior  mockBehavior
		expectedUser  *models.User
		expectError   bool
		expectedError string
	}{
		{
			name:     "Common",
			username: "yarik_tri",
			mockBehavior: func(username string, u *models.User) {
				row := sqlmock.
					NewRows([]string{"id", "version", "username", "email", "password_hash",
						"salt", "first_name", "last_name", "sex", "birth_date"}).
					AddRow(u.ID, u.Version, u.Username, u.Email, u.Password, u.Salt,
						u.FirstName, u.LastName, u.Sex, u.BirhDate.Time)

				mock.ExpectQuery("SELECT (.+) FROM " + UsersTable).
					WithArgs(username).
					WillReturnRows(row)
			},
			expectedUser: &u,
		},
		{
			name:     "No such user",
			username: "yarik_dva",
			mockBehavior: func(username string, u *models.User) {
				mock.ExpectQuery("SELECT (.+) FROM " + UsersTable).
					WithArgs(username).
					WillReturnError(errors.New("no such user"))
			},
			expectError:   true,
			expectedError: "no such user",
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockBehavior(tc.username, tc.expectedUser)

			user, err := r.GetUserByUsername(tc.username)
			if tc.expectError {
				assert.Error(t, err, tc.expectedError)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedUser, user)
			}
		})
	}
}

func TestAuthPostgresGetUserByAuthData(t *testing.T) {
	type mockBehavior func(userID, userVersion uint, u *models.User)
	type authData struct {
		userID      uint
		userVersion uint
	}

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("%v", err)
	}
	defer db.Close()

	l, err := logger.NewFLogger()
	if err != nil {
		t.Errorf("%v", err)
	}

	r := NewAuthPostgres(db, l)

	u, err := defaultUser()
	if err != nil {
		t.Errorf("can't create default user: %v", err)
	}

	testTable := []struct {
		name          string
		authData      authData
		mockBehavior  mockBehavior
		expectedUser  *models.User
		expectError   bool
		expectedError string
	}{
		{
			name: "Common",
			authData: authData{
				userID:      1,
				userVersion: 1,
			},
			mockBehavior: func(userID, userVersion uint, u *models.User) {
				row := sqlmock.
					NewRows([]string{"id", "version", "username", "email", "password_hash",
						"salt", "first_name", "last_name", "sex", "birth_date"}).
					AddRow(u.ID, u.Version, u.Username, u.Email, u.Password, u.Salt,
						u.FirstName, u.LastName, u.Sex, u.BirhDate.Time)

				mock.ExpectQuery("SELECT (.+) FROM "+UsersTable).
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
			mockBehavior: func(userID, userVersion uint, u *models.User) {
				mock.ExpectQuery("SELECT (.+) FROM "+UsersTable).
					WithArgs(userID, userVersion).
					WillReturnError(errors.New("no such user"))
			},
			expectError:   true,
			expectedError: "no such user",
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockBehavior(tc.authData.userID, tc.authData.userVersion, tc.expectedUser)

			user, err := r.GetUserByAuthData(tc.authData.userID, tc.authData.userVersion)
			if tc.expectError {
				assert.EqualError(t, err, tc.expectedError)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedUser, user)
			}
		})
	}
}

func TestAuthPostgresIncreaseUserVersion(t *testing.T) {
	type mockBehavior func(userID uint)

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("%v", err)
	}
	defer db.Close()

	l, err := logger.NewFLogger()
	if err != nil {
		t.Errorf("%v", err)
	}

	r := NewAuthPostgres(db, l)

	testTable := []struct {
		name          string
		userID        uint
		mockBehavior  mockBehavior
		expectedId    uint
		expectError   bool
		expectedError string
	}{
		{
			name:   "Common",
			userID: 1,
			mockBehavior: func(userID uint) {
				row := sqlmock.NewRows([]string{"id"}).AddRow(userID)

				mock.ExpectQuery("UPDATE " + UsersTable).
					WithArgs(userID).
					WillReturnRows(row)
			},
			expectedId: 1,
		},
		{
			name:   "No such user",
			userID: 1,
			mockBehavior: func(userID uint) {
				mock.ExpectQuery("UPDATE " + UsersTable).
					WithArgs(userID).
					WillReturnError(errors.New("no such user"))
			},
			expectError:   true,
			expectedError: "no such user",
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockBehavior(tc.userID)

			err := r.IncreaseUserVersion(tc.userID)
			if tc.expectError {
				assert.EqualError(t, err, tc.expectedError)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
