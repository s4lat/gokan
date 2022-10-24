package dbmanager

import (
	// "errors"
	// "context"
	// "github.com/jackc/pgx/v5/pgxpool"
	"encoding/json"
	"fmt"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/s4lat/gokan/models"
	"github.com/s4lat/gokan/postgresdb"
	"os"
	"testing"
)

var (
	DB_URL      = "postgres://user:password@localhost:5432/test"
	postgres, _ = postgresdb.NewPostgresDB(DB_URL)
	dbManager   = DBManager(&postgres)
)

type MockedData struct {
	Persons []models.Person `json:"persons"`
	Boards  []models.Board  `json:"boards"`
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

func CreateMockedPersons(mockData MockedData) error {
	for _, person := range mockData.Persons {
		_, err := dbManager.CreatePerson(person)
		if err != nil {
			return fmt.Errorf("CreateMockedPersons -> %w", err)
		}
	}
	return nil
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

	cmpIgnore := cmpopts.IgnoreFields(models.Person{}, "Boards", "AssignedTasks")
	for _, person := range mockData.Persons {
		createdPerson, err := dbManager.CreatePerson(person)
		if err != nil {
			t.Error(err)
		}

		if !cmp.Equal(createdPerson, person, cmpIgnore) {
			t.Errorf("Created person not equal to mocked: \n\t%v \n\t%v",
				createdPerson, person)
		}

		t.Logf("Created: %v", createdPerson)
	}

	if _, err := dbManager.CreatePerson(mockData.Persons[0]); err == nil {
		t.Error("CreatePerson() does't throw error when creating rows with same UNIQUE fields")
	}
}

func TestGetPersonByID(t *testing.T) {
	if err := dbManager.RecreateAllTables(); err != nil {
		t.Fatal(err)
	}

	mockData, err := LoadMockData()
	if err != nil {
		t.Fatal(err)
	}

	if err := CreateMockedPersons(mockData); err != nil {
		t.Fatal(err)
	}

	if _, err := dbManager.GetPersonByID(3030); err == nil {
		t.Error("Searching for non-existent ID not throwing error")
	}

	cmpIgnore := cmpopts.IgnoreFields(models.Person{}, "Boards", "AssignedTasks")
	for _, person := range mockData.Persons {
		obtainedPerson, err := dbManager.GetPersonByID(person.ID)
		if err != nil {
			t.Error(err)
		}

		t.Logf("Obtained: %v", obtainedPerson)

		if !cmp.Equal(obtainedPerson, person, cmpIgnore) {
			t.Errorf("Obtained person not equal to mocked: \n\t%v \n\t%v",
				obtainedPerson, person)
		}
	}
}

func TestGetPersonByEmail(t *testing.T) {
	if err := dbManager.RecreateAllTables(); err != nil {
		t.Fatal(err)
	}

	mockData, err := LoadMockData()
	if err != nil {
		t.Fatal(err)
	}

	if err := CreateMockedPersons(mockData); err != nil {
		t.Fatal(err)
	}

	if _, err := dbManager.GetPersonByEmail("aaa321@mail.ru"); err == nil {
		t.Error("Searching for non-existent email not throwing error")
	}

	cmpIgnore := cmpopts.IgnoreFields(models.Person{}, "Boards", "AssignedTasks")
	for _, person := range mockData.Persons {
		obtainedPerson, err := dbManager.GetPersonByEmail(person.Email)
		if err != nil {
			t.Error(err)
		}

		t.Logf("Obtained: %v", obtainedPerson)

		if !cmp.Equal(obtainedPerson, person, cmpIgnore) {
			t.Errorf("Obtained person not equal to mocked: \n\t%v \n\t%v",
				obtainedPerson, person)
		}
	}
}

func TestGetPersonByUsername(t *testing.T) {
	if err := dbManager.RecreateAllTables(); err != nil {
		t.Fatal(err)
	}

	mockData, err := LoadMockData()
	if err != nil {
		t.Fatal(err)
	}

	if err := CreateMockedPersons(mockData); err != nil {
		t.Fatal(err)
	}

	if _, err := dbManager.GetPersonByUsername("akaksda"); err == nil {
		t.Error("Searching for non-existent username not throwing error")
	}

	cmpIgnore := cmpopts.IgnoreFields(models.Person{}, "Boards", "AssignedTasks")
	for _, person := range mockData.Persons {
		obtainedPerson, err := dbManager.GetPersonByUsername(person.Username)
		if err != nil {
			t.Error(err)
		}

		t.Logf("Obtained: %v", obtainedPerson)

		if !cmp.Equal(obtainedPerson, person, cmpIgnore) {
			t.Errorf("Obtained person not equal to mocked: \n\t%v \n\t%v",
				obtainedPerson, person)
		}
	}
}

func TestCreateBoard(t *testing.T) {
	if err := dbManager.RecreateAllTables(); err != nil {
		t.Fatal(err)
	}

	mockData, err := LoadMockData()
	if err != nil {
		t.Fatal(err)
	}

	if err := CreateMockedPersons(mockData); err != nil {
		t.Fatal(err)
	}

	cmpIgnore := cmpopts.IgnoreFields(models.Board{}, "Contributors", "Tasks", "Tags")
	for _, board := range mockData.Boards {
		t.Logf("%v", board)
		createdBoard, err := dbManager.CreateBoard(board)
		if err != nil {
			t.Error(err)
		}

		if !cmp.Equal(createdBoard, board, cmpIgnore) {
			t.Errorf("Created board not equal to mocked: \n\t%v \n\t%v",
				createdBoard, board)
		}

		t.Logf("Created: %v", createdBoard)
	}

	badBoard := models.Board{Name: "badBoard", OwnerID: 1337}
	if _, err := dbManager.CreateBoard(badBoard); err == nil {
		t.Error("CreateBoard() does't throw error when creating rows with non-existent owner_id")
	}
}
