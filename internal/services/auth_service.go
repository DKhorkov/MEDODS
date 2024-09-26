package services

import (
	"github.com/DKhorkov/medods/internal/entities"
	"github.com/DKhorkov/medods/internal/interfaces"
	"time"
)

type CommonAuthService struct {
	AuthRepository interfaces.AuthRepository
}

func (service *CommonAuthService) CreateRefreshToken(GUID, value string, TTL time.Time) (int, error) {
	if oldRefreshToken, err := service.AuthRepository.GetRefreshTokenByGUID(GUID); err == nil {
		if err = service.AuthRepository.DeleteRefreshToken(oldRefreshToken); err != nil {
			return 0, err
		}
	}

	return service.AuthRepository.CreateRefreshToken(GUID, value, TTL)
}

func (service *CommonAuthService) GetRefreshTokenByValue(value string) (*entities.RefreshToken, error) {
	return service.AuthRepository.GetRefreshTokenByValue(value)
}
