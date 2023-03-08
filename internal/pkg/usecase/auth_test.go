package usecase

import (
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
