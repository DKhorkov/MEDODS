package usecases__test

import (
	"strconv"
	"testing"
	"time"

	"github.com/DKhorkov/hmtm-sso/pkg/logging"
	"github.com/DKhorkov/medods/internal/entities"
	customerrors "github.com/DKhorkov/medods/internal/errors"
	mocks "github.com/DKhorkov/medods/internal/mocks/repositories"
	"github.com/DKhorkov/medods/internal/security"
	"github.com/DKhorkov/medods/internal/services"
	"github.com/DKhorkov/medods/internal/usecases"
	testconfig "github.com/DKhorkov/medods/tests/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var testsConfig = testconfig.New()

func TestUseCasesCreateTokens(t *testing.T) {
	t.Run("create tokens successfully", func(t *testing.T) {
		authRepository := &mocks.MockedAuthRepository{RefreshTokensStorage: map[int]*entities.RefreshToken{}}
		usersRepository := &mocks.MockedUsersRepository{}
		authService := &services.CommonAuthService{AuthRepository: authRepository}
		usersService := &services.CommonUsersService{UsersRepository: usersRepository}
		useCases := &usecases.CommonUseCases{
			AuthService:  authService,
			UsersService: usersService,
			HashCost:     testsConfig.HashCost,
			JWTConfig:    testsConfig.JWT,
			SMTPConfig:   testsConfig.SMTP,
			Logger:       logging.GetInstance(testsConfig.Logging.Level, testsConfig.Logging.LogFilePath),
		}

		tokens, err := useCases.CreateTokens(
			entities.CreateTokensDTO{
				GUID: testsConfig.RefreshToken.GUID,
				IP:   testsConfig.IP,
			},
		)

		require.NoError(t, err)
		assert.NotEqual(
			t,
			"",
			tokens.AccessToken)
		assert.NotEqual(
			t,
			"",
			tokens.RefreshToken)
	})
}

func TestUseCasesRefreshTokens(t *testing.T) {
	t.Run("refresh tokens successfully", func(t *testing.T) {
		hashedRefreshTokenValue, err := security.HashRefreshToken(testsConfig.RefreshToken.Value, testsConfig.HashCost)
		if err != nil {
			t.Fatal(err)
		}

		dbRefreshToken := &entities.RefreshToken{
			ID:    1,
			Value: hashedRefreshTokenValue,
			TTL:   time.Now().Add(time.Hour),
			GUID:  testsConfig.RefreshToken.GUID,
		}

		authRepository := &mocks.MockedAuthRepository{
			RefreshTokensStorage: map[int]*entities.RefreshToken{
				dbRefreshToken.ID: dbRefreshToken,
			},
		}

		usersRepository := &mocks.MockedUsersRepository{}
		authService := &services.CommonAuthService{AuthRepository: authRepository}
		usersService := &services.CommonUsersService{UsersRepository: usersRepository}
		useCases := &usecases.CommonUseCases{
			AuthService:  authService,
			UsersService: usersService,
			HashCost:     testsConfig.HashCost,
			JWTConfig:    testsConfig.JWT,
			SMTPConfig:   testsConfig.SMTP,
			Logger:       logging.GetInstance(testsConfig.Logging.Level, testsConfig.Logging.LogFilePath),
		}

		refreshToken, err := security.GenerateJWT(
			security.JWTData{
				SecretKey: testsConfig.JWT.SecretKey,
				Algorithm: testsConfig.JWT.Algorithm,
				TTL:       testsConfig.JWT.RefreshTokenTTL,
				IP:        testsConfig.IP,
				Value:     testsConfig.RefreshToken.Value,
				GUID:      testsConfig.RefreshToken.GUID,
			},
		)

		if err != nil {
			t.Fatal(err)
		}

		accessToken, err := security.GenerateJWT(
			security.JWTData{
				SecretKey: testsConfig.JWT.SecretKey,
				Algorithm: testsConfig.JWT.Algorithm,
				TTL:       testsConfig.JWT.RefreshTokenTTL,
				IP:        testsConfig.IP,
				Value:     strconv.Itoa(dbRefreshToken.ID),
				GUID:      testsConfig.RefreshToken.GUID,
			},
		)

		if err != nil {
			t.Fatal(err)
		}

		tokens, err := useCases.RefreshTokens(
			entities.RefreshTokensDTO{
				Tokens: entities.Tokens{
					AccessToken:  accessToken,
					RefreshToken: refreshToken,
				},
				IP: testsConfig.IP,
			},
		)

		require.NoError(t, err)
		assert.NotEqual(
			t,
			"",
			tokens.AccessToken)
		assert.NotEqual(
			t,
			"",
			tokens.RefreshToken)
	})

	t.Run("refresh tokens with another IP", func(t *testing.T) {
		hashedRefreshTokenValue, err := security.HashRefreshToken(testsConfig.RefreshToken.Value, testsConfig.HashCost)
		if err != nil {
			t.Fatal(err)
		}

		dbRefreshToken := &entities.RefreshToken{
			ID:    1,
			Value: hashedRefreshTokenValue,
			TTL:   time.Now().Add(time.Hour),
			GUID:  testsConfig.RefreshToken.GUID,
		}

		authRepository := &mocks.MockedAuthRepository{
			RefreshTokensStorage: map[int]*entities.RefreshToken{
				dbRefreshToken.ID: dbRefreshToken,
			},
		}

		usersRepository := &mocks.MockedUsersRepository{}
		authService := &services.CommonAuthService{AuthRepository: authRepository}
		usersService := &services.CommonUsersService{UsersRepository: usersRepository}
		useCases := &usecases.CommonUseCases{
			AuthService:  authService,
			UsersService: usersService,
			HashCost:     testsConfig.HashCost,
			JWTConfig:    testsConfig.JWT,
			SMTPConfig:   testsConfig.SMTP,
			Logger:       logging.GetInstance(testsConfig.Logging.Level, testsConfig.Logging.LogFilePath),
		}

		refreshToken, err := security.GenerateJWT(
			security.JWTData{
				SecretKey: testsConfig.JWT.SecretKey,
				Algorithm: testsConfig.JWT.Algorithm,
				TTL:       testsConfig.JWT.RefreshTokenTTL,
				IP:        testsConfig.IP,
				Value:     testsConfig.RefreshToken.Value,
				GUID:      testsConfig.RefreshToken.GUID,
			},
		)

		if err != nil {
			t.Fatal(err)
		}

		accessToken, err := security.GenerateJWT(
			security.JWTData{
				SecretKey: testsConfig.JWT.SecretKey,
				Algorithm: testsConfig.JWT.Algorithm,
				TTL:       testsConfig.JWT.RefreshTokenTTL,
				IP:        testsConfig.IP,
				Value:     strconv.Itoa(dbRefreshToken.ID),
				GUID:      testsConfig.RefreshToken.GUID,
			},
		)

		if err != nil {
			t.Fatal(err)
		}

		tokens, err := useCases.RefreshTokens(
			entities.RefreshTokensDTO{
				Tokens: entities.Tokens{
					AccessToken:  accessToken,
					RefreshToken: refreshToken,
				},
				IP: "[::1]",
			},
		)

		require.Error(t, err)
		assert.IsType(t, customerrors.IPAddressDoesNotMatchWithTokensIPError{}, err)
		assert.Nil(t, tokens)
	})

	t.Run("access token does not belong to refresh token", func(t *testing.T) {
		hashedRefreshToken1Value, err := security.HashRefreshToken(testsConfig.RefreshToken.Value, testsConfig.HashCost)
		if err != nil {
			t.Fatal(err)
		}

		dbRefreshToken1 := &entities.RefreshToken{
			ID:    1,
			Value: hashedRefreshToken1Value,
			TTL:   time.Now().Add(time.Hour),
			GUID:  testsConfig.RefreshToken.GUID,
		}

		hashedRefreshToken2Value, err := security.HashRefreshToken("another value", testsConfig.HashCost)
		if err != nil {
			t.Fatal(err)
		}

		dbRefreshToken2 := &entities.RefreshToken{
			ID:    1,
			Value: hashedRefreshToken2Value,
			TTL:   time.Now().Add(time.Hour),
			GUID:  testsConfig.RefreshToken.GUID,
		}

		authRepository := &mocks.MockedAuthRepository{
			RefreshTokensStorage: map[int]*entities.RefreshToken{
				dbRefreshToken1.ID: dbRefreshToken1,
				dbRefreshToken2.ID: dbRefreshToken2,
			},
		}

		usersRepository := &mocks.MockedUsersRepository{}
		authService := &services.CommonAuthService{AuthRepository: authRepository}
		usersService := &services.CommonUsersService{UsersRepository: usersRepository}
		useCases := &usecases.CommonUseCases{
			AuthService:  authService,
			UsersService: usersService,
			HashCost:     testsConfig.HashCost,
			JWTConfig:    testsConfig.JWT,
			SMTPConfig:   testsConfig.SMTP,
			Logger:       logging.GetInstance(testsConfig.Logging.Level, testsConfig.Logging.LogFilePath),
		}

		refreshToken, err := security.GenerateJWT(
			security.JWTData{
				SecretKey: testsConfig.JWT.SecretKey,
				Algorithm: testsConfig.JWT.Algorithm,
				TTL:       testsConfig.JWT.RefreshTokenTTL,
				IP:        testsConfig.IP,
				Value:     testsConfig.RefreshToken.Value,
				GUID:      testsConfig.RefreshToken.GUID,
			},
		)

		if err != nil {
			t.Fatal(err)
		}

		accessToken, err := security.GenerateJWT(
			security.JWTData{
				SecretKey: testsConfig.JWT.SecretKey,
				Algorithm: testsConfig.JWT.Algorithm,
				TTL:       testsConfig.JWT.RefreshTokenTTL,
				IP:        testsConfig.IP,
				Value:     strconv.Itoa(dbRefreshToken2.ID),
				GUID:      testsConfig.RefreshToken.GUID,
			},
		)

		if err != nil {
			t.Fatal(err)
		}

		tokens, err := useCases.RefreshTokens(
			entities.RefreshTokensDTO{
				Tokens: entities.Tokens{
					AccessToken:  accessToken,
					RefreshToken: refreshToken,
				},
				IP: testsConfig.IP,
			},
		)

		require.Error(t, err)
		assert.IsType(t, customerrors.AccessTokenDoesNotBelongToRefreshTokenError{}, err)
		assert.Nil(t, tokens)
	})
}
