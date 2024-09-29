package repositories__test

import (
	"testing"
	"time"

	"github.com/DKhorkov/medods/internal/entities"
	customerrors "github.com/DKhorkov/medods/internal/errors"
	testconfig "github.com/DKhorkov/medods/tests/config"

	"github.com/DKhorkov/medods/internal/database"
	"github.com/DKhorkov/medods/internal/repositories"
	testlifespan "github.com/DKhorkov/medods/tests/internal/repositories/lifespan"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	_ "github.com/mattn/go-sqlite3"
)

var testsConfig = testconfig.New()

func TestRepositoriesCreateRefreshToken(t *testing.T) {
	var ttl = time.Now().Add(time.Hour)

	t.Run("successfully create refreshToken", func(t *testing.T) {
		connection := testlifespan.StartUp(t)
		defer testlifespan.TearDown(t, connection)

		authRepository := repositories.CommonAuthRepository{
			DBConnector: &database.CommonDBConnector{
				Connection: connection,
			},
		}

		// Error and zero userID due to returning nil ID after register.
		// SQLite inner realization without AUTO_INCREMENT for SERIAL PRIMARY KEY
		refreshTokenID, err := authRepository.CreateRefreshToken(
			testsConfig.RefreshToken.GUID,
			testsConfig.RefreshToken.Value,
			ttl,
		)

		require.Error(t, err)
		assert.Equal(t, 0, refreshTokenID)

		var refreshTokensCount int
		err = connection.QueryRow(
			`
				SELECT COUNT(*)
				FROM refresh_tokens
			`,
		).Scan(&refreshTokensCount)
		require.NoError(t, err)
		assert.Equal(t, 1, refreshTokensCount)
	})

	t.Run("create refreshToken failure due to existence of refreshToken with same value", func(t *testing.T) {
		connection := testlifespan.StartUp(t)
		defer testlifespan.TearDown(t, connection)

		_, err := connection.Exec(
			`
			INSERT INTO refresh_tokens (guid, value, ttl) 
			VALUES ($1, $2, $3)
			RETURNING refresh_tokens.id
		`,
			testsConfig.RefreshToken.GUID,
			testsConfig.RefreshToken.Value,
			ttl,
		)

		if err != nil {
			t.Fatalf("failed to insert user: %v", err)
		}

		authRepository := repositories.CommonAuthRepository{
			DBConnector: &database.CommonDBConnector{
				Connection: connection,
			},
		}

		refreshTokenID, err := authRepository.CreateRefreshToken(
			testsConfig.RefreshToken.GUID,
			testsConfig.RefreshToken.Value,
			ttl,
		)

		require.Error(t, err)
		assert.Equal(t, 0, refreshTokenID)
	})
}

func TestRepositoriesGetRefreshTokenByID(t *testing.T) {
	const refreshTokenID = 1

	t.Run("get existing refreshToken", func(t *testing.T) {
		connection := testlifespan.StartUp(t)
		defer testlifespan.TearDown(t, connection)

		testRefreshToken := &entities.RefreshToken{
			ID:    refreshTokenID,
			Value: testsConfig.RefreshToken.Value,
			TTL:   time.Now().Add(time.Hour),
			GUID:  testsConfig.RefreshToken.GUID,
		}

		_, err := connection.Exec(
			`
				INSERT INTO refresh_tokens (id, guid, value, ttl) 
				VALUES ($1, $2, $3, $4)
			`,
			testRefreshToken.ID,
			testRefreshToken.GUID,
			testRefreshToken.Value,
			testRefreshToken.TTL,
		)

		if err != nil {
			t.Fatalf("failed to insert refreshToken: %v", err)
		}

		authRepository := repositories.CommonAuthRepository{
			DBConnector: &database.CommonDBConnector{
				Connection: connection,
			},
		}

		refreshToken, err := authRepository.GetRefreshTokenByID(refreshTokenID)
		require.NoError(t, err)
		assert.NotNil(t, refreshToken)
		assert.Equal(
			t,
			time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC),
			refreshToken.DeletedAt)
	})

	t.Run("get non existing refreshToken", func(t *testing.T) {
		connection := testlifespan.StartUp(t)
		defer testlifespan.TearDown(t, connection)

		authRepository := repositories.CommonAuthRepository{
			DBConnector: &database.CommonDBConnector{
				Connection: connection,
			},
		}

		refreshToken, err := authRepository.GetRefreshTokenByID(refreshTokenID)
		require.Error(t, err)
		assert.IsType(t, customerrors.RefreshTokenNotFoundError{}, err)
		assert.Nil(t, refreshToken)
	})

	t.Run("get deleted refreshToken", func(t *testing.T) {
		connection := testlifespan.StartUp(t)
		defer testlifespan.TearDown(t, connection)

		_, err := connection.Exec(
			`
				INSERT INTO refresh_tokens (id, guid, value, ttl, deleted_at) 
				VALUES ($1, $2, $3, $4, $5)
			`,
			refreshTokenID,
			testsConfig.RefreshToken.GUID,
			testsConfig.RefreshToken.Value,
			time.Now().Add(time.Hour),
			time.Now().Add(time.Hour*time.Duration(-1)),
		)

		if err != nil {
			t.Fatalf("failed to insert refreshToken: %v", err)
		}

		authRepository := repositories.CommonAuthRepository{
			DBConnector: &database.CommonDBConnector{
				Connection: connection,
			},
		}

		refreshToken, err := authRepository.GetRefreshTokenByID(refreshTokenID)
		require.Error(t, err)
		assert.IsType(t, customerrors.RefreshTokenNotFoundError{}, err)
		assert.Nil(t, refreshToken)
	})

	t.Run("get expired refreshToken", func(t *testing.T) {
		connection := testlifespan.StartUp(t)
		defer testlifespan.TearDown(t, connection)

		_, err := connection.Exec(
			`
				INSERT INTO refresh_tokens (id, guid, value, ttl) 
				VALUES ($1, $2, $3, $4)
			`,
			refreshTokenID,
			testsConfig.RefreshToken.GUID,
			testsConfig.RefreshToken.Value,
			time.Now().Add(time.Hour*time.Duration(-1)),
		)

		if err != nil {
			t.Fatalf("failed to insert refreshToken: %v", err)
		}

		authRepository := repositories.CommonAuthRepository{
			DBConnector: &database.CommonDBConnector{
				Connection: connection,
			},
		}

		// No error due to in memory sqlite and CURRENT_TIMESTAMP returns 1 January 1970
		refreshToken, err := authRepository.GetRefreshTokenByID(refreshTokenID)
		require.NoError(t, err)
		assert.True(t, refreshToken.TTL.Before(time.Now()))
	})
}

func TestRepositoriesGetRefreshTokenByGUID(t *testing.T) {
	const refreshTokenID = 1

	t.Run("get existing refreshToken", func(t *testing.T) {
		connection := testlifespan.StartUp(t)
		defer testlifespan.TearDown(t, connection)

		testRefreshToken := &entities.RefreshToken{
			ID:    refreshTokenID,
			Value: testsConfig.RefreshToken.Value,
			TTL:   time.Now().Add(time.Hour),
			GUID:  testsConfig.RefreshToken.GUID,
		}

		_, err := connection.Exec(
			`
				INSERT INTO refresh_tokens (id, guid, value, ttl) 
				VALUES ($1, $2, $3, $4)
			`,
			testRefreshToken.ID,
			testRefreshToken.GUID,
			testRefreshToken.Value,
			testRefreshToken.TTL,
		)

		if err != nil {
			t.Fatalf("failed to insert refreshToken: %v", err)
		}

		authRepository := repositories.CommonAuthRepository{
			DBConnector: &database.CommonDBConnector{
				Connection: connection,
			},
		}

		refreshToken, err := authRepository.GetRefreshTokenByGUID(testsConfig.RefreshToken.GUID)
		require.NoError(t, err)
		assert.NotNil(t, refreshToken)
		assert.Equal(
			t,
			time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC),
			refreshToken.DeletedAt)
	})

	t.Run("get non existing refreshToken", func(t *testing.T) {
		connection := testlifespan.StartUp(t)
		defer testlifespan.TearDown(t, connection)

		authRepository := repositories.CommonAuthRepository{
			DBConnector: &database.CommonDBConnector{
				Connection: connection,
			},
		}

		refreshToken, err := authRepository.GetRefreshTokenByGUID(testsConfig.RefreshToken.GUID)
		require.Error(t, err)
		assert.IsType(t, customerrors.RefreshTokenNotFoundError{}, err)
		assert.Nil(t, refreshToken)
	})

	t.Run("get deleted refreshToken", func(t *testing.T) {
		connection := testlifespan.StartUp(t)
		defer testlifespan.TearDown(t, connection)

		_, err := connection.Exec(
			`
				INSERT INTO refresh_tokens (id, guid, value, ttl, deleted_at) 
				VALUES ($1, $2, $3, $4, $5)
			`,
			refreshTokenID,
			testsConfig.RefreshToken.GUID,
			testsConfig.RefreshToken.Value,
			time.Now().Add(time.Hour),
			time.Now().Add(time.Hour*time.Duration(-1)),
		)

		if err != nil {
			t.Fatalf("failed to insert refreshToken: %v", err)
		}

		authRepository := repositories.CommonAuthRepository{
			DBConnector: &database.CommonDBConnector{
				Connection: connection,
			},
		}

		refreshToken, err := authRepository.GetRefreshTokenByGUID(testsConfig.RefreshToken.GUID)
		require.Error(t, err)
		assert.IsType(t, customerrors.RefreshTokenNotFoundError{}, err)
		assert.Nil(t, refreshToken)
	})

	t.Run("get expired refreshToken", func(t *testing.T) {
		connection := testlifespan.StartUp(t)
		defer testlifespan.TearDown(t, connection)

		_, err := connection.Exec(
			`
				INSERT INTO refresh_tokens (id, guid, value, ttl) 
				VALUES ($1, $2, $3, $4)
			`,
			refreshTokenID,
			testsConfig.RefreshToken.GUID,
			testsConfig.RefreshToken.Value,
			time.Now().Add(time.Hour*time.Duration(-1)),
		)

		if err != nil {
			t.Fatalf("failed to insert refreshToken: %v", err)
		}

		authRepository := repositories.CommonAuthRepository{
			DBConnector: &database.CommonDBConnector{
				Connection: connection,
			},
		}

		// No error due to in memory sqlite and CURRENT_TIMESTAMP returns 1 January 1970
		refreshToken, err := authRepository.GetRefreshTokenByGUID(testsConfig.RefreshToken.GUID)
		require.NoError(t, err)
		assert.True(t, refreshToken.TTL.Before(time.Now()))
	})
}

func TestRepositoriesDeleteRefreshToken(t *testing.T) {
	t.Run("delete existing refreshToken", func(t *testing.T) {
		connection := testlifespan.StartUp(t)
		defer testlifespan.TearDown(t, connection)

		testRefreshToken := &entities.RefreshToken{
			ID:    1,
			Value: testsConfig.RefreshToken.Value,
			TTL:   time.Now().Add(time.Hour),
			GUID:  testsConfig.RefreshToken.GUID,
		}

		_, err := connection.Exec(
			`
				INSERT INTO refresh_tokens (id, guid, value, ttl) 
				VALUES ($1, $2, $3, $4)
			`,
			testRefreshToken.ID,
			testRefreshToken.GUID,
			testRefreshToken.Value,
			testRefreshToken.TTL,
		)

		if err != nil {
			t.Fatalf("failed to insert refreshToken: %v", err)
		}

		authRepository := repositories.CommonAuthRepository{
			DBConnector: &database.CommonDBConnector{
				Connection: connection,
			},
		}

		err = authRepository.DeleteRefreshToken(testRefreshToken)
		require.NoError(t, err)
	})

	t.Run("delete non existing refreshToken", func(t *testing.T) {
		connection := testlifespan.StartUp(t)
		defer testlifespan.TearDown(t, connection)

		testRefreshToken := &entities.RefreshToken{
			ID:    1,
			Value: testsConfig.RefreshToken.Value,
			TTL:   time.Now().Add(time.Hour),
			GUID:  testsConfig.RefreshToken.GUID,
		}

		authRepository := repositories.CommonAuthRepository{
			DBConnector: &database.CommonDBConnector{
				Connection: connection,
			},
		}

		// No error due to update stmt inside of DeleteRefreshToken method
		err := authRepository.DeleteRefreshToken(testRefreshToken)
		require.NoError(t, err)
	})
}
