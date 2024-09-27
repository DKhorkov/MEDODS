package interfaces

import (
	"time"

	"github.com/DKhorkov/medods/internal/entities"
)

type AuthRepository interface {
	CreateRefreshToken(guid, value string, ttl time.Time) (int, error)
	GetRefreshTokenByID(id int) (*entities.RefreshToken, error)
	GetRefreshTokenByGUID(guid string) (*entities.RefreshToken, error)
	DeleteRefreshToken(token *entities.RefreshToken) error
}

type UsersRepository interface {
	GetUserEmail(guid string) (string, error)
}
