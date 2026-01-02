package database_test

import (
	"database/sql"
	"fortyfour-backend/pkg/database"
	"testing"
)

func TestNewMySQLConnection(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		cfg     database.Config
		want    *sql.DB
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := database.NewMySQLConnection(tt.cfg)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("NewMySQLConnection() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("NewMySQLConnection() succeeded unexpectedly")
			}
			// TODO: update the condition below to compare got with tt.want.
			if true {
				t.Errorf("NewMySQLConnection() = %v, want %v", got, tt.want)
			}
		})
	}
}
