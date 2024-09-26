package usecases

import (
	"fmt"
	"github.com/DKhorkov/medods/internal/config"
	"github.com/DKhorkov/medods/internal/entities"
	"github.com/DKhorkov/medods/internal/interfaces"
	"github.com/DKhorkov/medods/internal/security"
	"github.com/DKhorkov/medods/internal/utils"
	"strconv"
	"time"
)

type CommonUseCases struct {
	AuthService interfaces.AuthService
	HashCost    int
	JWTConfig   config.JWTConfig
}

func (useCases *CommonUseCases) CreateTokens(data entities.CreateTokensDTO) (*entities.Tokens, error) {
	refreshTokenValue := fmt.Sprintf(
		"%s:%s",
		utils.GenerateRandomString(10),
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
	return &entities.Tokens{}, nil
}
