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

	tables := []string{"person", "board", "task", "subtask", "tag", "task_tag", "contributor"}
	for _, table := range tables {
		if isExist, err := dbManager.IsTableExist(table); err != nil {
			t.Error(err)
		} else if !isExist {
			t.Errorf("Table '%s' does not created", table)
		}
	}
}

func TestIsTableExist(t *testing.T) {
	err := dbManager.RecreateAllTables()
	if err != nil {
		t.Error(err)
	}

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
}
