package repositories

import (
	"strings"
	"time"

	"github.com/DKhorkov/medods/internal/database"
	"github.com/DKhorkov/medods/internal/entities"
	customerrors "github.com/DKhorkov/medods/internal/errors"
	"github.com/DKhorkov/medods/internal/interfaces"
)

type CommonAuthRepository struct {
	DBConnector interfaces.DBConnector
}

func (repo *CommonAuthRepository) CreateRefreshToken(guid, value string, ttl time.Time) (int, error) {
	var refreshTokenID int
	connection := repo.DBConnector.GetConnection()
	err := connection.QueryRow(
		`
			INSERT INTO refresh_tokens (guid, value, ttl) 
			VALUES ($1, $2, $3)
			RETURNING refresh_tokens.id
		`,
		guid,
		value,
		ttl,
	).Scan(&refreshTokenID)

	if err != nil {
		return 0, err
	}

	return refreshTokenID, nil
}

func (repo *CommonAuthRepository) GetRefreshTokenByID(id int) (*entities.RefreshToken, error) {
	refreshToken := &entities.RefreshToken{}
	columns := database.GetEntityColumns(refreshToken)
	connection := repo.DBConnector.GetConnection()
	err := connection.QueryRow(
		`
			SELECT *
			FROM refresh_tokens AS rt
			WHERE rt.id = $1 
			  AND rt.ttl > CURRENT_TIMESTAMP
			  AND rt.deleted_at IS NULL
		`,
		id,
	).Scan(columns...)

	if err != nil && !strings.Contains(err.Error(), "storing driver.Value type <nil> into type *time.Time") {
		return nil, customerrors.RefreshTokenNotFoundError{}
	}

	return refreshToken, nil
}

func (repo *CommonAuthRepository) GetRefreshTokenByGUID(guid string) (*entities.RefreshToken, error) {
	refreshToken := &entities.RefreshToken{}
	columns := database.GetEntityColumns(refreshToken)
	connection := repo.DBConnector.GetConnection()
	err := connection.QueryRow(
		`
			SELECT *
			FROM refresh_tokens AS rt
			WHERE rt.guid = $1 
			  AND rt.ttl > CURRENT_TIMESTAMP
			  AND rt.deleted_at IS NULL
		`,
		guid,
	).Scan(columns...)

	if err != nil && !strings.Contains(err.Error(), "storing driver.Value type <nil> into type *time.Time") {
		return nil, customerrors.RefreshTokenNotFoundError{}
	}

	return refreshToken, nil
}

func (repo *CommonAuthRepository) DeleteRefreshToken(token *entities.RefreshToken) error {
	connection := repo.DBConnector.GetConnection()
	err := connection.QueryRow(
		`
			UPDATE refresh_tokens
			SET deleted_at = CURRENT_TIMESTAMP
			WHERE id = $1
		`,
		token.ID,
	).Err()

	return err
}
