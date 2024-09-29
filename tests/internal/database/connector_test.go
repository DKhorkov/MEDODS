package database__test

import (
	"database/sql"
	"testing"

	"github.com/DKhorkov/medods/internal/database"
	customerrors "github.com/DKhorkov/medods/internal/errors"
	testconfig "github.com/DKhorkov/medods/tests/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	_ "github.com/mattn/go-sqlite3"
)

func TestDatabaseConnect(t *testing.T) {
	testsConfig := testconfig.New()

	t.Run("successfully connect to database", func(t *testing.T) {
		connector := &database.CommonDBConnector{
			DSN:    testsConfig.Database.DSN,
			Driver: testsConfig.Database.Driver,
		}

		err := connector.Connect()
		require.NoError(t, err)
	})

	t.Run("connect to non existing database", func(t *testing.T) {
		connector := &database.CommonDBConnector{
			DSN:    "non existing database",
			Driver: "error",
		}

		err := connector.Connect()
		require.Error(t, err)
	})
}

func TestDatabaseGetTransaction(t *testing.T) {
	testsConfig := testconfig.New()

	t.Run("successfully return transaction", func(t *testing.T) {
		connector := &database.CommonDBConnector{
			DSN:    testsConfig.Database.DSN,
			Driver: testsConfig.Database.Driver,
		}

		if err := connector.Connect(); err != nil {
			t.Fatal(err)
		}

		transaction, err := connector.GetTransaction()
		require.NoError(t, err)
		assert.IsTypef(
			t,
			&sql.Tx{},
			transaction,
			"transaction type should be sql.Tx")
	})

	t.Run("get transaction from nil connection", func(t *testing.T) {
		connector := &database.CommonDBConnector{
			DSN:    "non existing database",
			Driver: "error",
		}

		transaction, err := connector.GetTransaction()
		require.Error(t, err)
		assert.IsTypef(
			t,
			customerrors.NilDBConnectionError{},
			err,
			"should be customerrors.NilDBConnectionError")
		assert.Nil(t, transaction)
	})
}

func TestDatabaseGetConnection(t *testing.T) {
	testsConfig := testconfig.New()

	t.Run("successfully return connection", func(t *testing.T) {
		connector := &database.CommonDBConnector{
			DSN:    testsConfig.Database.DSN,
			Driver: testsConfig.Database.Driver,
		}

		if err := connector.Connect(); err != nil {
			t.Fatal(err)
		}

		connection := connector.GetConnection()
		assert.NotNil(t, connection)
		assert.IsTypef(
			t,
			&sql.DB{},
			connection,
			"connection type should be sql.DB")
	})

	t.Run("connection was nil", func(t *testing.T) {
		connector := &database.CommonDBConnector{
			DSN:    testsConfig.Database.DSN,
			Driver: testsConfig.Database.Driver,
		}

		connection := connector.GetConnection()
		assert.NotNil(t, connection)
		assert.IsTypef(
			t,
			&sql.DB{},
			connection,
			"connection type should be sql.DB")
	})

	t.Run("connect to database is not possible", func(t *testing.T) {
		connector := &database.CommonDBConnector{
			DSN:    "non existing database",
			Driver: "error",
		}

		connection := connector.GetConnection()
		assert.Nil(t, connection)
	})
}
