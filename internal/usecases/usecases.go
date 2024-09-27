package usecases

import (
	"fmt"
	"log/slog"
	"strconv"
	"time"

	"github.com/DKhorkov/hmtm-sso/pkg/logging"
	customerrors "github.com/DKhorkov/medods/internal/errors"

	"github.com/DKhorkov/medods/internal/config"
	"github.com/DKhorkov/medods/internal/entities"
	"github.com/DKhorkov/medods/internal/interfaces"
	"github.com/DKhorkov/medods/internal/security"
)

type CommonUseCases struct {
	AuthService  interfaces.AuthService
	UsersService interfaces.UsersService
	HashCost     int
	JWTConfig    config.JWTConfig
	SMTPConfig   config.SMTPConfig
	Logger       *slog.Logger
}

func (useCases *CommonUseCases) CreateTokens(data entities.CreateTokensDTO) (*entities.Tokens, error) {
	randomSeedLength := 10
	refreshTokenValue := fmt.Sprintf(
		"%s:%s",
		generateRandomString(randomSeedLength),
		data.GUID,
	)

	hashedRefreshTokenValue, err := security.HashRefreshToken(refreshTokenValue, useCases.HashCost)
	if err != nil {
		return nil, err
	}

	refreshTokenID, err := useCases.AuthService.CreateRefreshToken(
		data.GUID,
		hashedRefreshTokenValue,
		time.Now().Add(useCases.JWTConfig.RefreshTokenTTL),
	)

	if err != nil {
		return nil, err
	}

	refreshToken, err := security.GenerateJWT(
		security.JWTData{
			IP:        data.IP,
			GUID:      data.GUID,
			Value:     refreshTokenValue,
			SecretKey: useCases.JWTConfig.SecretKey,
			Algorithm: useCases.JWTConfig.Algorithm,
			TTL:       useCases.JWTConfig.RefreshTokenTTL,
		},
	)

	if err != nil {
		return nil, err
	}

	accessToken, err := security.GenerateJWT(
		security.JWTData{
			IP:        data.IP,
			GUID:      data.GUID,
			Value:     strconv.Itoa(refreshTokenID),
			SecretKey: useCases.JWTConfig.SecretKey,
			Algorithm: useCases.JWTConfig.Algorithm,
			TTL:       useCases.JWTConfig.AccessTokenTTL,
		},
	)

	if err != nil {
		return nil, err
	}

	return &entities.Tokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (useCases *CommonUseCases) RefreshTokens(data entities.RefreshTokensDTO) (*entities.Tokens, error) {
	accessTokenPayload, err := security.ParseJWT(data.Tokens.AccessToken, useCases.JWTConfig.SecretKey)
	if err != nil {
		return nil, err
	}

	refreshTokenPayload, err := security.ParseJWT(data.Tokens.RefreshToken, useCases.JWTConfig.SecretKey)
	if err != nil {
		return nil, err
	}

	if accessTokenPayload.IP != data.IP || refreshTokenPayload.IP != data.IP {
		go func() {
			email, err := useCases.UsersService.GetUserEmail(refreshTokenPayload.GUID)
			if err != nil {
				useCases.Logger.Error(
					"Failed to get user email",
					"Traceback",
					logging.GetLogTraceback(),
					"Error",
					err,
				)
			}

			sendEmail(
				"MEDODS warning\n",
				fmt.Sprintf("Someone tried to refresh your tokens from next IP - %s", data.IP),
				[]string{email},
				useCases.SMTPConfig,
				useCases.Logger,
			)
		}()

		return nil, customerrors.IPAddressDoesNotMatchWithTokensIPError{}
	}

	refreshTokenID, err := strconv.Atoi(accessTokenPayload.Value)
	if err != nil {
		return nil, err
	}

	dbRefreshToken, err := useCases.AuthService.GetRefreshTokenByID(refreshTokenID)
	if err != nil {
		return nil, err
	}

	if !security.ValidateRefreshToken(refreshTokenPayload.Value, dbRefreshToken.Value) {
		return nil, customerrors.AccessTokenDoesNotBelongToRefreshTokenError{}
	}

	return useCases.CreateTokens(
		entities.CreateTokensDTO{
			GUID: dbRefreshToken.GUID,
			IP:   data.IP,
		},
	)
}
