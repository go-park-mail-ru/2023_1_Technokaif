package usecase

import (
	"fmt"
	mathRand "math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUsecaseToken_GenerateAndCheckAccessToken(t *testing.T) {
	const iterations = 100

	u := NewUsecase()

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

func TestUsecaseToken_GenerateAndCheckCSRFToken(t *testing.T) {
	const iterations = 100

	u := NewUsecase()

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
