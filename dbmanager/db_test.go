package dbmanager

import (
	// "errors"
	// "context"
	// "github.com/jackc/pgx/v5/pgxpool"
	"encoding/json"
	"fmt"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"os"
	"testing"
)

var (
	DB_URL      = "postgres://user:password@localhost:5432/test"
	postgres, _ = NewPostgresDB(DB_URL)
	dbManager   = DBManager(&postgres)
)

type MockedData struct {
	Persons []Person `json:"persons"`
}

func LoadMockData() (MockedData, error) {
	json_data, err := os.ReadFile("mock_db.json")
	if err != nil {
		return MockedData{}, fmt.Errorf("LoadMockData -> %w", err)
	}

	var mockData MockedData
	if err := json.Unmarshal(json_data, &mockData); err != nil {
		return MockedData{}, fmt.Errorf("LoadMockData -> %w", err)
	}
	return mockData, nil
}

func TestMain(m *testing.M) {
	m.Run()
}

func TestRecreateAllTables(t *testing.T) {
	if err := dbManager.RecreateAllTables(); err != nil {
		t.Fatal(err)
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
	if err := dbManager.RecreateAllTables(); err != nil {
		t.Fatal(err)
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

func TestCreatePerson(t *testing.T) {
	if err := dbManager.RecreateAllTables(); err != nil {
		t.Fatal(err)
	}

	mockData, err := LoadMockData()
	if err != nil {
		t.Fatal(err)
	}

	created_persons := make([]Person, len(mockData.Persons))
	copy(created_persons, mockData.Persons)

	cmpIgnore := cmpopts.IgnoreFields(Person{}, "ID", "Boards", "AssignedTasks")
	for i, _ := range created_persons {
		if err := dbManager.CreatePerson(&created_persons[i]); err != nil {
			t.Error(err)
		}

		if !cmp.Equal(created_persons[i], mockData.Persons[i], cmpIgnore) {
			t.Errorf("Created persons not equal: \n\t%+v \n\t%+v",
				created_persons[i], mockData.Persons[i])
		}
	}

	if err := dbManager.CreatePerson(&mockData.Persons[0]); err == nil {
		t.Error("CreatePerson() does't throw error when creating rows with same UNIQUE fields")
	}
}
