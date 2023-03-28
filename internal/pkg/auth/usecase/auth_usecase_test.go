package usecase

import (
	"fmt"
	mathRand "math/rand"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	logMocks "github.com/go-park-mail-ru/2023_1_Technokaif/pkg/logger/mocks"
	authMocks "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/auth/mocks"
	userMocks "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/user/mocks"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/models"
)

func TestUsecaseAuthCreateUser(t *testing.T) {  // Cringe
	type mockBehavior func(r *userMocks.MockRepository, u models.User)
	type result struct {
		Id  uint32
		Err error
	}

	birthTime, err := time.Parse(time.RFC3339, "2003-08-23T00:00:00Z")
	if err != nil {
		t.Errorf("can't Parse birth date: %v", err)
	}
	birthDate := models.Date{birthTime}

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
			c := gomock.NewController(t)
			defer c.Finish()

			authMocksRepo := authMocks.NewMockRepository(c)
			userMocksRepo := userMocks.NewMockRepository(c)

			tc.mockBehavior(userMocksRepo, tc.user)

			l := logMocks.NewMockLogger(c)
			l.EXPECT().Error(gomock.Any()).AnyTimes()
			l.EXPECT().Info(gomock.Any()).AnyTimes()
			l.EXPECT().Errorf(gomock.Any(), gomock.Any()).AnyTimes()
			l.EXPECT().Infof(gomock.Any(), gomock.Any()).AnyTimes()

			u := NewUsecase(authMocksRepo, userMocksRepo, l)

			id, err := u.SignUpUser(tc.user)

			assert.Equal(t, tc.expected.Id, id)
			assert.Equal(t, tc.expected.Err, err)
		})
	}
}

func TestUsecaseAuthGenerateAndCheckToken(t *testing.T) {

	const iterations = 100

	for i := 0; i < iterations; i++ {
		t.Run(fmt.Sprintf("Success Token test %d", i), func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			authMocksRepo := authMocks.NewMockRepository(c)
			userMocksRepo := userMocks.NewMockRepository(c)

			l := logMocks.NewMockLogger(c)
			l.EXPECT().Error(gomock.Any()).AnyTimes()
			l.EXPECT().Info(gomock.Any()).AnyTimes()
			l.EXPECT().Errorf(gomock.Any(), gomock.Any()).AnyTimes()
			l.EXPECT().Infof(gomock.Any(), gomock.Any()).AnyTimes()

			u := NewUsecase(authMocksRepo, userMocksRepo, l)

			expectedUserID := uint32(mathRand.Intn(10000))
			expectedUserVersion := uint32(mathRand.Intn(10000))

			token, err := u.GenerateAccessToken(expectedUserID, expectedUserVersion)
			assert.NoError(t, err)

			gotUserID, gotUserVersion, err := u.CheckAccessToken(token)
			assert.NoError(t, err)
			assert.Equal(t, expectedUserID, gotUserID)
			assert.Equal(t, expectedUserVersion, gotUserVersion)
		})
	}
}
