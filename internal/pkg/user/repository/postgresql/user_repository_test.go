package postgresql

import (
	"context"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/models"
	"github.com/golang/mock/gomock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"

	userMocks "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/user/mocks"
)

var ctx = context.Background()

const userTable = "Users"

var errPqInternal = errors.New("postgres is dead")

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
