package services__test

import (
	"testing"
	"time"

	"github.com/DKhorkov/medods/internal/entities"
	mocks "github.com/DKhorkov/medods/internal/mocks/repositories"
	"github.com/DKhorkov/medods/internal/services"
	testconfig "github.com/DKhorkov/medods/tests/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var testsConfig = testconfig.New()

func TestServicesCreateRefreshToken(t *testing.T) {
	t.Run("create refresh token without old refresh token for provided GUID", func(t *testing.T) {
		authRepository := &mocks.MockedAuthRepository{RefreshTokensStorage: map[int]*entities.RefreshToken{}}
		authService := &services.CommonAuthService{AuthRepository: authRepository}

		previousRefreshTokensCount := len(authRepository.RefreshTokensStorage)
		refreshTokenID, err := authService.CreateRefreshToken(
			testsConfig.RefreshToken.GUID,
			testsConfig.RefreshToken.Value,
			time.Now().Add(time.Hour),
		)

		require.NoError(t, err)
		assert.Equal(
			t,
			previousRefreshTokensCount+1,
			refreshTokenID)
	})

	t.Run("create refresh token with old refresh token for provided GUID", func(t *testing.T) {
		oldRefreshToken := &entities.RefreshToken{
			ID:    1,
			Value: testsConfig.RefreshToken.Value,
			TTL:   time.Now().Add(time.Hour),
			GUID:  testsConfig.RefreshToken.GUID,
		}

		authRepository := &mocks.MockedAuthRepository{
			RefreshTokensStorage: map[int]*entities.RefreshToken{
				oldRefreshToken.ID: oldRefreshToken,
			},
		}

		authService := &services.CommonAuthService{AuthRepository: authRepository}

		previousRefreshTokensCount := len(authRepository.RefreshTokensStorage)
		refreshTokenID, err := authService.CreateRefreshToken(
			oldRefreshToken.GUID,
			"newTestValue",
			time.Now().Add(time.Hour),
		)

		require.NoError(t, err)
		assert.Equal(
			t,
			previousRefreshTokensCount+1,
			refreshTokenID)

		assert.True(t, oldRefreshToken.DeletedAt.Before(time.Now()))
	})
}
