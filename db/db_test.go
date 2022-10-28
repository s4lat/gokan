//nolint:gocognit
package db

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"

	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/s4lat/gokan/models"
)

var DBURL = "postgres://user:password@localhost:5432/test"
var db DB

type MockedData struct {
	Persons  []models.Person  `json:"persons"`
	Boards   []models.Board   `json:"boards"`
	Tasks    []models.Task    `json:"tasks"`
	Tags     []models.Tag     `json:"tags"`
	Subtasks []models.Subtask `json:"subtasks"`
	TaskTag  []struct {
		TaskID uint32 `json:"ref_task_id"`
		TagID  uint32 `json:"ref_tag_id"`
	} `json:"task_tag"`

	Assignees []struct {
		TaskID     uint32 `json:"ref_task_id"`
		AssigneeID uint32 `json:"assignee_id"`
	} `json:"assignees"`
}

func LoadMockData() (MockedData, error) {
	jsonData, err := os.ReadFile("mock_db.json")
	if err != nil {
		return MockedData{}, fmt.Errorf("LoadMockData() -> %w", err)
	}

	var mockedData MockedData
	if err := json.Unmarshal(jsonData, &mockedData); err != nil {
		return MockedData{}, fmt.Errorf("LoadMockData() -> %w", err)
	}
	return mockedData, nil
}

func (md *MockedData) CreateMockedPersons() error {
	for _, person := range md.Persons {
		_, err := db.Person.Create(person)
		if err != nil {
			return fmt.Errorf("CreateMockedPersons -> %w", err)
		}
	}
	return nil
}

func (md *MockedData) CreateMockedBoards() error {
	for _, board := range md.Boards {
		_, err := db.Board.Create(board)
		if err != nil {
			return fmt.Errorf("CreateMockedBoards -> %w", err)
		}
	}
	return nil
}

func (md *MockedData) CreateMockedTasks() error {
	for _, task := range md.Tasks {
		_, err := db.Task.Create(task)
		if err != nil {
			return fmt.Errorf("CreateMockedTasks -> %w", err)
		}
	}
	return nil
}

func (md *MockedData) CreateMockedTags() error {
	for _, tag := range md.Tags {
		_, err := db.Tag.Create(tag)
		if err != nil {
			return fmt.Errorf("CreateMockedTags -> %w", err)
		}
	}
	return nil
}

func TestMain(m *testing.M) {
	dbPool, err := pgxpool.New(context.Background(), DBURL)
	if err != nil {
		log.Fatal(err)
	}

	dbConn := models.DBConn(dbPool)
	db = NewDB(dbConn)
	m.Run()
}

func TestSystemRecreateAllTables(t *testing.T) {
	if err := db.System.RecreateAllTables(); err != nil {
		t.Fatal(err)
	}

	tables := []string{"person", "board", "task", "subtask", "tag", "task_tag", "contributor"}
	for _, table := range tables {
		if isExist, err := db.System.IsTableExist(table); err != nil {
			t.Error(err)
		} else if !isExist {
			t.Errorf("Table '%s' does not created", table)
		}
	}
}

func TestSystemIsTableExist(t *testing.T) {
	if err := db.System.RecreateAllTables(); err != nil {
		t.Fatal(err)
	}

	if isExist, err := db.System.IsTableExist("kek"); err != nil {
		t.Error(err)
	} else if isExist {
		t.Error("Table 'kek' exist but it's not")
	}

	if isExist, err := db.System.IsTableExist("person"); err != nil {
		t.Error(err)
	} else if !isExist {
		t.Error("Table 'person' does not exist but it exist")
	}
}

func TestPersonCreate(t *testing.T) {
	if err := db.System.RecreateAllTables(); err != nil {
		t.Fatal(err)
	}

	mockedData, err := LoadMockData()
	if err != nil {
		t.Fatal(err)
	}

	cmpIgnore := cmpopts.IgnoreFields(models.Person{}, "Boards", "AssignedTasks")
	for _, mockedPerson := range mockedData.Persons {
		createdPerson, err := db.Person.Create(mockedPerson)
		if err != nil {
			t.Error(err)
		}

		if !cmp.Equal(createdPerson, mockedPerson, cmpIgnore) {
			t.Errorf("Created person not equal to mocked: \n\t%v \n\t%v",
				createdPerson, mockedPerson)
		}

		t.Logf("Created: %v", createdPerson)
	}

	if _, err := db.Person.Create(mockedData.Persons[0]); err == nil {
		t.Error("Person.Create() does't throw error when creating rows with same UNIQUE fields")
	}
}

func TestPersonGetByID(t *testing.T) {
	if err := db.System.RecreateAllTables(); err != nil {
		t.Fatal(err)
	}

	mockedData, err := LoadMockData()
	if err != nil {
		t.Fatal(err)
	}

	if err := mockedData.CreateMockedPersons(); err != nil {
		t.Fatal(err)
	}

	if _, err := db.Person.GetByID(3030); err == nil {
		t.Error("Searching for user with non-existent ID not throwing error")
	}

	cmpIgnore := cmpopts.IgnoreFields(models.Person{}, "Boards", "AssignedTasks")
	for _, mockedPerson := range mockedData.Persons {
		obtainedPerson, err := db.Person.GetByID(mockedPerson.ID)
		if err != nil {
			t.Error(err)
		}

		t.Logf("Obtained: %v", obtainedPerson)

		if !cmp.Equal(obtainedPerson, mockedPerson, cmpIgnore) {
			t.Errorf("Obtained person not equal to mocked: \n\t%v \n\t%v",
				obtainedPerson, mockedPerson)
		}
	}
}

func TestPersonGetByEmail(t *testing.T) {
	if err := db.System.RecreateAllTables(); err != nil {
		t.Fatal(err)
	}

	mockedData, err := LoadMockData()
	if err != nil {
		t.Fatal(err)
	}

	if err := mockedData.CreateMockedPersons(); err != nil {
		t.Fatal(err)
	}

	if _, err := db.Person.GetByEmail("aaa321@mail.ru"); err == nil {
		t.Error("Searching for user with non-existent email not throwing error")
	}

	cmpIgnore := cmpopts.IgnoreFields(models.Person{}, "Boards", "AssignedTasks")
	for _, mockedPerson := range mockedData.Persons {
		obtainedPerson, err := db.Person.GetByEmail(mockedPerson.Email)
		if err != nil {
			t.Error(err)
		}

		t.Logf("Obtained: %v", obtainedPerson)

		if !cmp.Equal(obtainedPerson, mockedPerson, cmpIgnore) {
			t.Errorf("Obtained person not equal to mocked: \n\t%v \n\t%v",
				obtainedPerson, mockedPerson)
		}
	}
}

func TestPersonGetByUsername(t *testing.T) {
	if err := db.System.RecreateAllTables(); err != nil {
		t.Fatal(err)
	}

	mockedData, err := LoadMockData()
	if err != nil {
		t.Fatal(err)
	}

	if err := mockedData.CreateMockedPersons(); err != nil {
		t.Fatal(err)
	}

	if _, err := db.Person.GetByUsername("akaksda"); err == nil {
		t.Error("Searching for user with non-existent username not throwing error")
	}

	cmpIgnore := cmpopts.IgnoreFields(models.Person{}, "Boards", "AssignedTasks")
	for _, mockedPerson := range mockedData.Persons {
		obtainedPerson, err := db.Person.GetByUsername(mockedPerson.Username)
		if err != nil {
			t.Error(err)
		}

		t.Logf("Obtained: %v", obtainedPerson)

		if !cmp.Equal(obtainedPerson, mockedPerson, cmpIgnore) {
			t.Errorf("Obtained person not equal to mocked: \n\t%v \n\t%v",
				obtainedPerson, mockedPerson)
		}
	}
}

func TestBoardCreate(t *testing.T) {
	if err := db.System.RecreateAllTables(); err != nil {
		t.Fatal(err)
	}

	mockedData, err := LoadMockData()
	if err != nil {
		t.Fatal(err)
	}

	if err := mockedData.CreateMockedPersons(); err != nil {
		t.Fatal(err)
	}

	cmpIgnore := cmpopts.IgnoreFields(models.Board{}, "Contributors", "Tasks", "Tags")
	for _, board := range mockedData.Boards {
		t.Logf("%v", board)
		createdBoard, err := db.Board.Create(board)
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
	if _, err := db.Board.Create(badBoard); err == nil {
		t.Error("BoardModel.Create() does't throw error when creating rows with non-existent owner_id")
	}
}

func BoardGetByID(t *testing.T) {
	if err := db.System.RecreateAllTables(); err != nil {
		t.Fatal(err)
	}
	mockedData, err := LoadMockData()
	if err != nil {
		t.Fatal(err)
	}
	if err := mockedData.CreateMockedPersons(); err != nil {
		t.Fatal(err)
	}
	if err := mockedData.CreateMockedBoards(); err != nil {
		t.Fatal(err)
	}

	if _, err := db.Board.GetByID(1337); err == nil {
		t.Error("Searching for board with non-existent boardID not throwing error")
	}

	cmpIgnore := cmpopts.IgnoreFields(models.Board{}, "Contributors", "Tasks", "Tags")
	for _, mockedBoard := range mockedData.Boards {
		obtainedBoard, err := db.Board.GetByID(mockedBoard.ID)
		if err != nil {
			t.Error(err)
		}

		if !cmp.Equal(obtainedBoard, mockedBoard, cmpIgnore) {
			t.Errorf("Obtained board not equal to mocked: \n\t%v \n\t%v",
				obtainedBoard, mockedBoard)
		}

		t.Logf("Obtained: %v", obtainedBoard)
	}
}

func TestTaskCreate(t *testing.T) {
	if err := db.System.RecreateAllTables(); err != nil {
		t.Fatal(err)
	}
	mockedData, err := LoadMockData()
	if err != nil {
		t.Fatal(err)
	}
	if err := mockedData.CreateMockedPersons(); err != nil {
		t.Fatal(err)
	}
	if err := mockedData.CreateMockedBoards(); err != nil {
		t.Fatal(err)
	}

	cmpIgnore := cmpopts.IgnoreFields(models.Task{}, "Subtasks", "Tags", "Assignees")
	for _, mockedTask := range mockedData.Tasks {
		createdTask, err := db.Task.Create(mockedTask)
		if err != nil {
			t.Error(err)
		}

		if !cmp.Equal(createdTask, mockedTask, cmpIgnore) {
			t.Errorf("Created task not equal to mocked: \n\t%v \n\t%v",
				createdTask, mockedTask)
		}

		t.Logf("Created: %v", createdTask)
	}
}

func TestTaskGetByID(t *testing.T) {
	if err := db.System.RecreateAllTables(); err != nil {
		t.Fatal(err)
	}
	mockedData, err := LoadMockData()
	if err != nil {
		t.Fatal(err)
	}
	if err := mockedData.CreateMockedPersons(); err != nil {
		t.Fatal(err)
	}
	if err := mockedData.CreateMockedBoards(); err != nil {
		t.Fatal(err)
	}
	if err := mockedData.CreateMockedTasks(); err != nil {
		t.Fatal(err)
	}

	if _, err := db.Task.GetByID(1337); err == nil {
		t.Error("Searching for task with non-existent taskID not throwing error")
	}

	cmpIgnore := cmpopts.IgnoreFields(models.Task{}, "Subtasks", "Tags")
	for _, mockedTask := range mockedData.Tasks {
		obtainedTask, err := db.Task.GetByID(mockedTask.ID)
		if err != nil {
			t.Error(err)
		}

		if !cmp.Equal(obtainedTask, mockedTask, cmpIgnore) {
			t.Errorf("Obtained task not equal to mocked: \n\t%v \n\t%v",
				obtainedTask, mockedTask)
		}
		t.Logf("Obtained: %v", obtainedTask)
	}
}

func TestTagCreate(t *testing.T) {
	if err := db.System.RecreateAllTables(); err != nil {
		t.Fatal(err)
	}
	mockedData, err := LoadMockData()
	if err != nil {
		t.Fatal(err)
	}
	if err := mockedData.CreateMockedPersons(); err != nil {
		t.Fatal(err)
	}
	if err := mockedData.CreateMockedBoards(); err != nil {
		t.Fatal(err)
	}
	if err := mockedData.CreateMockedTasks(); err != nil {
		t.Fatal(err)
	}

	cmpIgnore := cmpopts.IgnoreFields(models.Tag{}, "ID")
	for _, mockedTag := range mockedData.Tags {
		createdTag, err := db.Tag.Create(mockedTag)
		if err != nil {
			t.Error(err)
		}

		if !cmp.Equal(mockedTag, createdTag, cmpIgnore) {
			t.Errorf("Created tag not equal to mocked: \n\t%v \n\t%v",
				createdTag, mockedTag)
		}
		t.Logf("Created: %v", createdTag)
	}
}

func TestTagGetByID(t *testing.T) {
	if err := db.System.RecreateAllTables(); err != nil {
		t.Fatal(err)
	}
	mockedData, err := LoadMockData()
	if err != nil {
		t.Fatal(err)
	}
	if err := mockedData.CreateMockedPersons(); err != nil {
		t.Fatal(err)
	}
	if err := mockedData.CreateMockedBoards(); err != nil {
		t.Fatal(err)
	}
	if err := mockedData.CreateMockedTasks(); err != nil {
		t.Fatal(err)
	}
	if err := mockedData.CreateMockedTags(); err != nil {
		t.Fatal(err)
	}

	if _, err := db.Tag.GetByID(1337); err == nil {
		t.Error("Searching for tag with non-existent tagID not throwing error")
	}

	cmpIgnore := cmpopts.IgnoreFields(models.Task{}, "ID", "Assignees")
	for _, mockedTag := range mockedData.Tags {
		obtainedTag, err := db.Tag.GetByID(mockedTag.ID)
		if err != nil {
			t.Error(err)
		}

		if !cmp.Equal(obtainedTag, mockedTag, cmpIgnore) {
			t.Errorf("Obtained task not equal to mocked: \n\t%v \n\t%v",
				obtainedTag, mockedTag)
		}
		t.Logf("Obtained: %v", obtainedTag)
	}
}

func TestTaskAddTagToTask(t *testing.T) {
	if err := db.System.RecreateAllTables(); err != nil {
		t.Fatal(err)
	}
	mockedData, err := LoadMockData()
	if err != nil {
		t.Fatal(err)
	}
	if err := mockedData.CreateMockedPersons(); err != nil {
		t.Fatal(err)
	}
	if err := mockedData.CreateMockedBoards(); err != nil {
		t.Fatal(err)
	}
	if err := mockedData.CreateMockedTasks(); err != nil {
		t.Fatal(err)
	}
	if err := mockedData.CreateMockedTags(); err != nil {
		t.Fatal(err)
	}

OuterFor:
	for _, taskTag := range mockedData.TaskTag {
		task, err := db.Task.GetByID(taskTag.TaskID)
		if err != nil {
			t.Error(err)
		}
		tag, err := db.Tag.GetByID(taskTag.TagID)
		if err != nil {
			t.Error(err)
		}
		task, err = db.Task.AddTagToTask(tag, task)
		if err != nil {
			t.Error(err)
		}

		mockedTag, err := db.Tag.GetByID(taskTag.TagID)
		if err != nil {
			t.Error(err)
		}

		for _, tag := range task.Tags {
			if tag.ID == taskTag.TagID {
				if !cmp.Equal(tag, mockedTag) {
					t.Errorf("Obtained tag not equal to mocked: \n\t%v \n\t%v",
						tag, mockedTag)
				}
				t.Logf("Successfully added tag to task: %v - %v", task.ID, task.Tags)
				continue OuterFor
			}
		}
		t.Errorf("Tag not added to task.Tags: \n\t%v\n\t%v", tag, task.Tags)
	}
}

func TestTaskAssignPersonToTask(t *testing.T) {
	if err := db.System.RecreateAllTables(); err != nil {
		t.Fatal(err)
	}
	mockedData, err := LoadMockData()
	if err != nil {
		t.Fatal(err)
	}
	if err := mockedData.CreateMockedPersons(); err != nil {
		t.Fatal(err)
	}
	if err := mockedData.CreateMockedBoards(); err != nil {
		t.Fatal(err)
	}
	if err := mockedData.CreateMockedTasks(); err != nil {
		t.Fatal(err)
	}
	if err := mockedData.CreateMockedTags(); err != nil {
		t.Fatal(err)
	}

OuterFor:
	for _, mockedAssign := range mockedData.Assignees {
		task, err := db.Task.GetByID(mockedAssign.TaskID)
		if err != nil {
			t.Error(err)
		}
		person, err := db.Person.GetByID(mockedAssign.AssigneeID)
		if err != nil {
			t.Error(err)
		}

		task, err = db.Task.AssignPersonToTask(person, task)
		if err != nil {
			t.Error(err)
		}

		for _, assignee := range task.Assignees {
			if assignee.ID == mockedAssign.AssigneeID {
				t.Logf("Successfully assigned task to person: %v - %v", task.ID, task.Assignees)
				continue OuterFor
			}
		}

		t.Errorf("Person not added to task.Assignees: \n\t%v\n\t%v", person, task.Assignees)
	}
}

func TestTaskAddSubtaskToTask(t *testing.T) {
	if err := db.System.RecreateAllTables(); err != nil {
		t.Fatal(err)
	}
	mockedData, err := LoadMockData()
	if err != nil {
		t.Fatal(err)
	}
	if err := mockedData.CreateMockedPersons(); err != nil {
		t.Fatal(err)
	}
	if err := mockedData.CreateMockedBoards(); err != nil {
		t.Fatal(err)
	}
	if err := mockedData.CreateMockedTasks(); err != nil {
		t.Fatal(err)
	}
	if err := mockedData.CreateMockedTags(); err != nil {
		t.Fatal(err)
	}

OuterFor:
	for _, mockedSubtask := range mockedData.Subtasks {
		task, err := db.Task.GetByID(mockedSubtask.ParentTaskID)
		if err != nil {
			t.Error(err)
		}

		task, err = db.Task.AddSubtaskToTask(mockedSubtask, task)
		if err != nil {
			t.Error(err)
		}

		for _, subtask := range task.Subtasks {
			if subtask.ID == mockedSubtask.ID && !cmp.Equal(subtask, mockedSubtask) {
				t.Errorf("Obtained subtask not equal to mocked: \n\t%v \n\t%v",
					subtask, mockedSubtask)
			}
			t.Logf("Successfully added subtask to task: %v - %v", task.ID, task.Subtasks)
			continue OuterFor
		}

		t.Errorf("Subtask not added to task.Subtasks: \n\t%v\n\t%v", mockedSubtask, task.Assignees)
	}
}
