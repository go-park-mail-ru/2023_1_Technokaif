package repository

import (
	"errors"
	"testing"
	"time"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/logger"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestAuthPostgres_CreateUser(t *testing.T) {
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

	type mockBehavior func(u models.User, id int)

	birthTime, err := time.Parse(time.RFC3339, "2003-08-23T00:00:00Z")
	if err != nil {
		t.Errorf("can't Parse birth date: %v", err)
	}
	birthDate := models.Date{birthTime}

	correctTestUser := models.User{
		Username:  "yarik_tri",
		Password:  "HASHHASHASHHASHSALTHASH",
		Salt:      "SALT",
		Email:     "yarik1448kuzmin@gmail.com",
		FirstName: "Yaroslav",
		LastName:  "Kuzmin",
		BirhDate:  birthDate,
		Sex:       models.Male,
	}

	testTable := []struct {
		name         string
		userToCreate models.User
		mockBehavior mockBehavior
		id           int
		hasError     bool
	}{
		{
			name:         "Common",
			userToCreate: correctTestUser,
			mockBehavior: func(u models.User, id int) {
				// Just creates query row (expected result after exec)
				rows := sqlmock.NewRows([]string{"id"}).AddRow(id)

				// Exec request and expect rows to be result
				mock.ExpectQuery("INSERT INTO Users").
					WithArgs(u.Username, u.Email, u.Password, u.Salt,
						u.FirstName, u.LastName, u.Sex, u.BirhDate.Format(time.RFC3339)).
					WillReturnRows(rows)
			},
			id: 1,
		},
		{
			name: "Empty Not Null Fields",
			userToCreate: models.User{
				Username:  "",
				Password:  "HASHHASHASHHASHSALTHASH",
				Salt:      "SALT",
				Email:     "yarik1448kuzmin@gmail.com",
				FirstName: "Yaroslav",
				LastName:  "Kuzmin",
				BirhDate:  birthDate,
				Sex:       models.Male,
			},
			mockBehavior: func(u models.User, id int) {
				// Just creates query row (expected result after exec)
				rows := sqlmock.NewRows([]string{"id"}).AddRow(id).
					RowError(1, errors.New("null NOT NULL field"))

				// Exec request and expect rows to be result
				mock.ExpectQuery("INSERT INTO Users").
					WithArgs(u.Username, u.Email, u.Password, u.Salt,
						u.FirstName, u.LastName, u.Sex, u.BirhDate.Format(time.RFC3339)).
					WillReturnRows(rows)
			},
			hasError: true,
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockBehavior(tc.userToCreate, tc.id)

			resultId, err := r.CreateUser(tc.userToCreate)
			if tc.hasError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.id, resultId)
			}
		})
	}
}
