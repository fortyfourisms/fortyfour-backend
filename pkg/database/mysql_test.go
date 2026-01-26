package database_test

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestNewMySQLConnection_MockOnly(t *testing.T) {
	// Test only the mock behavior
	db, mock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
	assert.NoError(t, err)
	defer db.Close()

	t.Run("Successful ping", func(t *testing.T) {
		mock.ExpectPing()
		err = db.Ping()
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Failed ping", func(t *testing.T) {
		mock.ExpectPing().WillReturnError(assert.AnError)
		err = db.Ping()
		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
