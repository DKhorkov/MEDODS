package interfaces

import (
	"time"

	"github.com/DKhorkov/medods/internal/entities"
)

type AuthService interface {
	CreateRefreshToken(guid, value string, ttl time.Time) (int, error)
	GetRefreshTokenByID(id int) (*entities.RefreshToken, error)
}

type UsersService interface {
	GetUserEmail(guid string) (string, error)
}
