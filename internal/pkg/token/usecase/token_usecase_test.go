package usecase

import (
	"fmt"
	mathRand "math/rand"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	logMocks "github.com/go-park-mail-ru/2023_1_Technokaif/pkg/logger/mocks"
)

func TestUsecaseTokenGenerateAndCheckAccessToken(t *testing.T) {
	const iterations = 100

	c := gomock.NewController(t)

	l := logMocks.NewMockLogger(c)
	l.EXPECT().Error(gomock.Any()).AnyTimes()
	l.EXPECT().Info(gomock.Any()).AnyTimes()
	l.EXPECT().Errorf(gomock.Any(), gomock.Any()).AnyTimes()
	l.EXPECT().Infof(gomock.Any(), gomock.Any()).AnyTimes()
	u := NewUsecase(l)

	for i := 0; i < iterations; i++ {
		t.Run(fmt.Sprintf("Success Token test %d", i), func(t *testing.T) {
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

func TestUsecaseTokenGenerateAndCheckCSRFToken(t *testing.T) {
	const iterations = 100

	c := gomock.NewController(t)

	l := logMocks.NewMockLogger(c)
	l.EXPECT().Error(gomock.Any()).AnyTimes()
	l.EXPECT().Info(gomock.Any()).AnyTimes()
	l.EXPECT().Errorf(gomock.Any(), gomock.Any()).AnyTimes()
	l.EXPECT().Infof(gomock.Any(), gomock.Any()).AnyTimes()
	u := NewUsecase(l)

	for i := 0; i < iterations; i++ {
		t.Run(fmt.Sprintf("Success Token test %d", i), func(t *testing.T) {
			expectedUserID := uint32(mathRand.Intn(10000))

			token, err := u.GenerateCSRFToken(expectedUserID)
			assert.NoError(t, err)

			gotUserID, err := u.CheckCSRFToken(token)
			assert.NoError(t, err)
			assert.Equal(t, expectedUserID, gotUserID)
		})
	}
}
