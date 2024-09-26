package interfaces

import (
	"github.com/DKhorkov/medods/internal/entities"
	"time"
)

type AuthService interface {
	CreateRefreshToken(GUID, value string, TTL time.Time) (int, error)
	GetRefreshTokenByValue(value string) (*entities.RefreshToken, error)
}
