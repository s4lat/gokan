package dbmanager

import (
	// "errors"
	// "context"
	// "github.com/jackc/pgx/v5/pgxpool"
	"testing"
)

var (
	DB_URL      = "postgres://user:password@localhost:5432/test"
	postgres, _ = NewPostgresDB(DB_URL)
	dbManager   = DBManager(&postgres)
)

func TestMain(m *testing.M) {
	m.Run()
}

func TestRecreateAllTables(t *testing.T) {
	err := dbManager.RecreateAllTables()
	if err != nil {
		t.Error(err)
	}

	if isExist, err := dbManager.IsTableExist("person"); err != nil {
		t.Error(err)
	} else if !isExist {
		t.Error("Table 'person' does not created")
	}
}

func TestIsTableExist(t *testing.T) {
	dbManager.RecreateAllTables()

	if isExist, err := dbManager.IsTableExist("kek"); err != nil {
		t.Error(err)
	} else if isExist {
		t.Error("Table 'kek' exist but it's not")
	}

	if isExist, err := dbManager.IsTableExist("person"); err != nil {
		t.Error(err)
	} else if !isExist {
		t.Error("Table 'person' does not exist but it exist")
	}

	// Checking by raw SQL
}
