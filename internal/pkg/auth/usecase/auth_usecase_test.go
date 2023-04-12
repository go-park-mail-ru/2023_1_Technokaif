package usecase

import (
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	commonTests "github.com/go-park-mail-ru/2023_1_Technokaif/internal/common/tests"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/models"
	authMocks "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/auth/mocks"
	userMocks "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/user/mocks"
)

func TestUsecaseAuthCreateUser(t *testing.T) {
	// Init
	type mockBehavior func(r *userMocks.MockRepository, u models.User)
	type result struct {
		Id  uint32
		Err error
	}

	c := gomock.NewController(t)

	authMocksRepo := authMocks.NewMockRepository(c)
	userMocksRepo := userMocks.NewMockRepository(c)

	l := commonTests.MockLogger(c)

	u := NewUsecase(authMocksRepo, userMocksRepo, l)

	birthTime, err := time.Parse(time.RFC3339, "2003-08-23T00:00:00Z")
	require.NoError(t, err, "can't Parse birth date")
	
	birthDate := models.Date{Time: birthTime}

	testTable := []struct {
		name         string
		user         models.User
		mockBehavior mockBehavior
		expected     result
	}{
		{
			name: "Common",
			user: models.User{
				Username:  "yarik_tri",
				Email:     "yarik1448kuzmin@gmail.com",
				Password:  "Love1234",
				FirstName: "Yaroslav",
				LastName:  "Kuzmin",
				BirthDate: birthDate,
				Sex:       models.Male,
			},
			mockBehavior: func(r *userMocks.MockRepository, u models.User) {
				// random salt, can't predict :(
				r.EXPECT().CreateUser(gomock.Any()).Return(uint32(1), nil)
			},
			expected: result{
				Id:  1,
				Err: nil,
			},
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockBehavior(userMocksRepo, tc.user)
			id, err := u.SignUpUser(tc.user)

			assert.Equal(t, tc.expected.Id, id)
			assert.Equal(t, tc.expected.Err, err)
		})
	}
}
