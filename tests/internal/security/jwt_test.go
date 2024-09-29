package security__test

import (
	"testing"
	"time"

	customerrors "github.com/DKhorkov/medods/internal/errors"

	"github.com/DKhorkov/medods/internal/security"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSecurityGenerateJWT(t *testing.T) {
	testCases := []struct {
		name          string
		data          security.JWTData
		message       string
		errorExpected bool
	}{
		{
			name: "generate valid token",
			data: security.JWTData{
				SecretKey: testsConfig.JWT.SecretKey,
				Algorithm: testsConfig.JWT.Algorithm,
				TTL:       testsConfig.JWT.RefreshTokenTTL,
				IP:        testsConfig.IP,
				Value:     testsConfig.RefreshToken.Value,
				GUID:      testsConfig.RefreshToken.GUID,
			},
			message:       "should return valid JWT token",
			errorExpected: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			token, err := security.GenerateJWT(tc.data)

			if tc.errorExpected {
				require.Error(t, err, tc.message)
				assert.Equal(
					t,
					"",
					token,
					tc.message)
			} else {
				require.NoError(t, err, tc.message)
				assert.NotEqual(
					t,
					"",
					token,
					tc.message)
			}
		})
	}
}

func TestSecurityParseJWT(t *testing.T) {
	testCases := []struct {
		name          string
		data          security.JWTData
		secretKey     string
		message       string
		errorExpected bool
		errorType     error
		expected      *security.JWTData
	}{
		{
			name: "parse JWT successfully",
			data: security.JWTData{
				SecretKey: testsConfig.JWT.SecretKey,
				Algorithm: testsConfig.JWT.Algorithm,
				TTL:       testsConfig.JWT.RefreshTokenTTL,
				IP:        testsConfig.IP,
				Value:     testsConfig.RefreshToken.Value,
				GUID:      testsConfig.RefreshToken.GUID,
			},
			message:       "should return valid JWT token",
			errorExpected: false,
			errorType:     nil,
			expected: &security.JWTData{
				IP:    testsConfig.IP,
				Value: testsConfig.RefreshToken.Value,
				GUID:  testsConfig.RefreshToken.GUID,
			},
			secretKey: testsConfig.JWT.SecretKey,
		},
		{
			name: "invalid secret key",
			data: security.JWTData{
				SecretKey: testsConfig.JWT.SecretKey,
				Algorithm: testsConfig.JWT.Algorithm,
				TTL:       testsConfig.JWT.RefreshTokenTTL,
				IP:        testsConfig.IP,
				Value:     testsConfig.RefreshToken.Value,
				GUID:      testsConfig.RefreshToken.GUID,
			},
			message:       "should raise an error due to invalid secret key",
			errorExpected: true,
			errorType:     customerrors.InvalidJWTError{},
			expected:      nil,
			secretKey:     "invalidSecretKey",
		},
		{
			name: "expired JWT",
			data: security.JWTData{
				SecretKey: testsConfig.JWT.SecretKey,
				Algorithm: testsConfig.JWT.Algorithm,
				TTL:       time.Minute * time.Duration(-1),
				IP:        testsConfig.IP,
				Value:     testsConfig.RefreshToken.Value,
				GUID:      testsConfig.RefreshToken.GUID,
			},
			message:       "should raise an error due to expired JWT",
			errorExpected: true,
			errorType:     customerrors.InvalidJWTError{},
			expected:      nil,
			secretKey:     testsConfig.JWT.SecretKey,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			token, err := security.GenerateJWT(tc.data)
			require.NoError(t, err, tc.message)

			parsedJWT, err := security.ParseJWT(token, tc.secretKey)
			if tc.errorExpected {
				require.Error(t, err, tc.message)
				assert.IsType(t, tc.errorType, err)
			} else {
				require.NoError(t, err, tc.message)
			}

			assert.Equal(
				t,
				tc.expected,
				parsedJWT,
				tc.message)
		})
	}
}
