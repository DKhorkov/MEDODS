package interfaces

import (
	"github.com/DKhorkov/medods/internal/entities"
)

type UseCases interface {
	CreateTokens(data entities.CreateTokensDTO) (*entities.Tokens, error)
	RefreshTokens(user entities.RefreshTokensDTO) (*entities.Tokens, error)
}
