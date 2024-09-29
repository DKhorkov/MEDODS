package database__test

import (
	"testing"

	"github.com/DKhorkov/medods/internal/database"
	"github.com/stretchr/testify/assert"
)

func TestDatabaseGetEntityColumns(t *testing.T) {
	t.Run("slice of correct len and capacity", func(t *testing.T) {
		testStruct := &struct {
			Column1 int
			Column2 string
		}{}

		columns := database.GetEntityColumns(testStruct)
		assert.Len(t, columns, 2)
		assert.IsTypef(
			t,
			[]interface{}{},
			columns,
			"should return a slice of []interface{}")
	})
}
