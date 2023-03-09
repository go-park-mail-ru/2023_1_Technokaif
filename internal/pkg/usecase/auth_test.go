package usecase

import (
	"fmt"
	mathRand "math/rand"
	cryptoRand "crypto/rand"
	"testing"
	"time"

	logMocks "github.com/go-park-mail-ru/2023_1_Technokaif/internal/logger/mocks"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/models"
	"github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/repository"
	mocks "github.com/go-park-mail-ru/2023_1_Technokaif/internal/pkg/repository/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestUsecaseAuthCreateUser(t *testing.T) {
	type mockBehavior func(a *mocks.MockAuth, u models.User)
	type result struct {
		Id  int
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
				BirhDate:  birthDate,
				Sex:       models.Male,
			},
			mockBehavior: func(a *mocks.MockAuth, u models.User) {
				// random salt, can't predict :(
				a.EXPECT().CreateUser(gomock.Any()).Return(1, nil)
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

			a := mocks.NewMockAuth(c)

			tc.mockBehavior(a, tc.user)

			r := &repository.Repository{
				Auth: a,
			}

			l := logMocks.NewMockLogger(c)
			l.EXPECT().Error(gomock.Any()).AnyTimes()
			l.EXPECT().Info(gomock.Any()).AnyTimes()

			u := NewUsecase(r, l)

			id, err := u.CreateUser(tc.user)

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

			a := mocks.NewMockAuth(c)
			r := &repository.Repository{
				Auth: a,
			}

			l := logMocks.NewMockLogger(c)
			l.EXPECT().Error(gomock.Any()).AnyTimes()
			l.EXPECT().Info(gomock.Any()).AnyTimes()

			u := NewUsecase(r, l)

			expectedUserID := uint(mathRand.Intn(10000))
			expectedUserVersion := uint(mathRand.Intn(10000))

			token, err := u.GenerateAccessToken(expectedUserID, expectedUserVersion)
			assert.NoError(t, err)

			gotUserID, gotUserVersion, err := u.CheckAccessToken(token)
			assert.NoError(t, err)
			assert.Equal(t, expectedUserID, gotUserID)
			assert.Equal(t, expectedUserVersion, gotUserVersion)		
		})
	}
}

func TestUsecaseCheckHash(t *testing.T) {
	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")

	randPass := func(n int) string {
		b := make([]rune, n)
		for i := range b {
			b[i] = letterRunes[mathRand.Intn(len(letterRunes))]
		}
		return string(b)
	}

	randSalt := func(n int) []byte {
		salt := make([]byte, n)
		cryptoRand.Read(salt)
		return salt
	}

	const iterations = 10
	const iterationsEq = 5
	for i := 0; i < iterations; i++ {
		t.Run(fmt.Sprintf("Success Hash test %d", i), func(t *testing.T) {

			pass := randPass(10)
			salt := randSalt(8)
			hash := hashPassword(pass, salt)

			for j := 0; j < iterationsEq; j++ {
				assert.Equal(t, hash, hashPassword(pass, salt))
			}		
		})
	}
}

