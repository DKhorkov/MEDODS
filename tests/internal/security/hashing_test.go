package security__test

import (
	"testing"

	testconfig "github.com/DKhorkov/medods/tests/config"

	"github.com/DKhorkov/medods/internal/security"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var testsConfig = testconfig.New()

func TestSecurityHashRefreshToken(t *testing.T) {
	testCases := []struct {
		name          string
		hashCost      int
		refreshToken  string
		message       string
		errorExpected bool
	}{
		{
			name:          "refreshToken successfully hashed",
			hashCost:      testsConfig.HashCost,
			refreshToken:  "refreshToken",
			message:       "should return hash for refreshToken",
			errorExpected: false,
		},
		{
			name:          "too long refreshToken > 72 bytes",
			hashCost:      testsConfig.HashCost,
			refreshToken:  "tooLongRefreshTokenThatCanNotBeLessThanSeventyTwoBytesForSureAndThereCouldAlsoBeSomeStory",
			message:       "should return error",
			errorExpected: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			hashedRefreshToken, err := security.HashRefreshToken(tc.refreshToken, tc.hashCost)

			if tc.errorExpected {
				require.Error(t, err, tc.message)
				assert.Equal(
					t,
					"",
					hashedRefreshToken,
					tc.message)
			} else {
				require.NoError(t, err, tc.message)
				assert.NotEqual(
					t,
					"",
					hashedRefreshToken,
					tc.message)
			}
		})
	}
}

func TestSecurityValidateRefreshToken(t *testing.T) {
	refreshTokenToHash := "refreshToken"

	testCases := []struct {
		name         string
		expected     bool
		hashCost     int
		refreshToken string
		message      string
	}{
		{
			name:         "hashed refreshToken was created based on provided refreshToken",
			refreshToken: refreshTokenToHash,
			hashCost:     testsConfig.HashCost,
			expected:     true,
			message:      "should return true",
		},
		{
			name:         "hash refreshToken was not created based on provided refreshToken",
			refreshToken: "Incorrect refreshToken",
			hashCost:     testsConfig.HashCost,
			expected:     false,
			message:      "should return false",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			hashedRefreshToken, _ := security.HashRefreshToken(refreshTokenToHash, tc.hashCost)
			refreshTokenIsValid := security.ValidateRefreshToken(tc.refreshToken, hashedRefreshToken)

			assert.Equal(
				t,
				tc.expected,
				refreshTokenIsValid,
				tc.message)
		})
	}
}
