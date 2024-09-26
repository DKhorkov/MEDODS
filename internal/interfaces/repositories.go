package interfaces

import (
	"github.com/DKhorkov/medods/internal/entities"
	"time"
)

type AuthRepository interface {
	CreateRefreshToken(GUID, value string, TTL time.Time) (int, error)
	GetRefreshTokenByValue(value string) (*entities.RefreshToken, error)
	GetRefreshTokenByGUID(GUID string) (*entities.RefreshToken, error)
	DeleteRefreshToken(token *entities.RefreshToken) error
}
