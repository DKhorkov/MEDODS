package mocks

import (
	"errors"
	"time"

	"github.com/DKhorkov/medods/internal/entities"
	customerrors "github.com/DKhorkov/medods/internal/errors"
)

type MockedAuthRepository struct {
	RefreshTokensStorage map[int]*entities.RefreshToken
}

func (repo *MockedAuthRepository) CreateRefreshToken(guid, value string, ttl time.Time) (int, error) {
	for _, refreshToken := range repo.RefreshTokensStorage {
		if refreshToken.Value == value {
			return 0, errors.New("refresh token already exists")
		}
	}

	refreshToken := &entities.RefreshToken{
		ID:        len(repo.RefreshTokensStorage) + 1,
		GUID:      guid,
		Value:     value,
		TTL:       ttl,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	repo.RefreshTokensStorage[refreshToken.ID] = refreshToken
	return refreshToken.ID, nil
}

func (repo *MockedAuthRepository) GetRefreshTokenByID(id int) (*entities.RefreshToken, error) {
	refreshToken := repo.RefreshTokensStorage[id]
	if refreshToken != nil {
		return refreshToken, nil
	}

	return nil, customerrors.RefreshTokenNotFoundError{}
}

func (repo *MockedAuthRepository) GetRefreshTokenByGUID(guid string) (*entities.RefreshToken, error) {
	for _, refreshToken := range repo.RefreshTokensStorage {
		if refreshToken.GUID == guid {
			return refreshToken, nil
		}
	}

	return nil, customerrors.RefreshTokenNotFoundError{}
}

func (repo *MockedAuthRepository) DeleteRefreshToken(token *entities.RefreshToken) error {
	refreshToken := repo.RefreshTokensStorage[token.ID]
	if refreshToken == nil {
		return customerrors.RefreshTokenNotFoundError{}
	}

	refreshToken.DeletedAt = time.Now().UTC()
	return nil
}
