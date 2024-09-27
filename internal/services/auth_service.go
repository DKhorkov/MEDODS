package services

import (
	"time"

	"github.com/DKhorkov/medods/internal/entities"
	"github.com/DKhorkov/medods/internal/interfaces"
)

type CommonAuthService struct {
	AuthRepository interfaces.AuthRepository
}

func (service *CommonAuthService) CreateRefreshToken(guid, value string, ttl time.Time) (int, error) {
	if oldRefreshToken, err := service.AuthRepository.GetRefreshTokenByGUID(guid); err == nil {
		if err = service.AuthRepository.DeleteRefreshToken(oldRefreshToken); err != nil {
			return 0, err
		}
	}

	return service.AuthRepository.CreateRefreshToken(guid, value, ttl)
}

func (service *CommonAuthService) GetRefreshTokenByID(id int) (*entities.RefreshToken, error) {
	return service.AuthRepository.GetRefreshTokenByID(id)
}
