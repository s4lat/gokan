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

	Contributors []struct {
		BoardID  uint32 `json:"board_id"`
		PersonID uint32 `json:"person_id"`
	} `json:"contributors"`
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
		_, err := db.Person.Create(context.Background(), person)
		if err != nil {
			return fmt.Errorf("CreateMockedPersons -> %w", err)
		}
	}
	return nil
}

func (md *MockedData) CreateMockedBoards() error {
	for _, board := range md.Boards {
		_, err := db.Board.Create(context.Background(), board)
		if err != nil {
			return fmt.Errorf("CreateMockedBoards -> %w", err)
		}
	}
	return nil
}

func (md *MockedData) CreateMockedTasks() error {
	for _, task := range md.Tasks {
		_, err := db.Task.Create(context.Background(), task)
		if err != nil {
			return fmt.Errorf("CreateMockedTasks -> %w", err)
		}
	}
	return nil
}

func (md *MockedData) CreateMockedTags() error {
	for _, tag := range md.Tags {
		_, err := db.Tag.Create(context.Background(), tag)
		if err != nil {
			return fmt.Errorf("CreateMockedTags -> %w", err)
		}
	}
	return nil
}

func TestMain(m *testing.M) {
	ctx := context.Background()
	dbPool, err := pgxpool.New(ctx, DBURL)
	if err != nil {
		log.Fatal(err)
	}

	dbConn := models.DBConn(dbPool)
	db = NewDB(dbConn)
	m.Run()
}

func TestSystemRecreateAllTables(t *testing.T) {
	ctx := context.Background()
	if err := db.System.RecreateAllTables(ctx); err != nil {
		t.Fatal(err)
	}

	tables := []string{"person", "board", "task", "subtask", "tag", "task_tag", "contributor"}
	for _, table := range tables {
		if isExist, err := db.System.IsTableExist(ctx, table); err != nil {
			t.Error(err)
		} else if !isExist {
			t.Errorf("Table '%s' does not created", table)
		}
	}
}

func TestSystemIsTableExist(t *testing.T) {
	ctx := context.Background()
	if err := db.System.RecreateAllTables(ctx); err != nil {
		t.Fatal(err)
	}

	if isExist, err := db.System.IsTableExist(ctx, "kek"); err != nil {
		t.Error(err)
	} else if isExist {
		t.Error("Table 'kek' exist but it's not")
	}

	if isExist, err := db.System.IsTableExist(ctx, "person"); err != nil {
		t.Error(err)
	} else if !isExist {
		t.Error("Table 'person' does not exist but it exist")
	}
}

func TestPersonCreate(t *testing.T) {
	ctx := context.Background()
	if err := db.System.RecreateAllTables(ctx); err != nil {
		t.Fatal(err)
	}

	mockedData, err := LoadMockData()
	if err != nil {
		t.Fatal(err)
	}

	cmpIgnore := cmpopts.IgnoreFields(models.Person{}, "Boards", "AssignedTasks")
	for _, mockedPerson := range mockedData.Persons {
		createdPerson, err := db.Person.Create(ctx, mockedPerson)
		if err != nil {
			t.Error(err)
		}

		if !cmp.Equal(createdPerson, mockedPerson, cmpIgnore) {
			t.Errorf("Created person not equal to mocked: \n\t%v \n\t%v",
				createdPerson, mockedPerson)
		}

		t.Logf("Created: %v", createdPerson)
	}

	if _, err := db.Person.Create(ctx, mockedData.Persons[0]); err == nil {
		t.Error("Person.Create() does't throw error when creating rows with same UNIQUE fields")
	}
}

func TestPersonGetByID(t *testing.T) {
	ctx := context.Background()
	if err := db.System.RecreateAllTables(ctx); err != nil {
		t.Fatal(err)
	}

	mockedData, err := LoadMockData()
	if err != nil {
		t.Fatal(err)
	}

	if err := mockedData.CreateMockedPersons(); err != nil {
		t.Fatal(err)
	}

	if _, err := db.Person.GetByID(ctx, 3030); err == nil {
		t.Error("Searching for user with non-existent ID not throwing error")
	}

	cmpIgnore := cmpopts.IgnoreFields(models.Person{}, "Boards")
	for _, mockedPerson := range mockedData.Persons {
		obtainedPerson, err := db.Person.GetByID(ctx, mockedPerson.ID)
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
	ctx := context.Background()
	if err := db.System.RecreateAllTables(ctx); err != nil {
		t.Fatal(err)
	}

	mockedData, err := LoadMockData()
	if err != nil {
		t.Fatal(err)
	}

	if err := mockedData.CreateMockedPersons(); err != nil {
		t.Fatal(err)
	}

	if _, err := db.Person.GetByEmail(ctx, "aaa321@mail.ru"); err == nil {
		t.Error("Searching for user with non-existent email not throwing error")
	}

	cmpIgnore := cmpopts.IgnoreFields(models.Person{}, "Boards", "AssignedTasks")
	for _, mockedPerson := range mockedData.Persons {
		obtainedPerson, err := db.Person.GetByEmail(ctx, mockedPerson.Email)
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
	ctx := context.Background()
	if err := db.System.RecreateAllTables(ctx); err != nil {
		t.Fatal(err)
	}

	mockedData, err := LoadMockData()
	if err != nil {
		t.Fatal(err)
	}

	if err := mockedData.CreateMockedPersons(); err != nil {
		t.Fatal(err)
	}

	if _, err := db.Person.GetByUsername(ctx, "akaksda"); err == nil {
		t.Error("Searching for user with non-existent username not throwing error")
	}

	cmpIgnore := cmpopts.IgnoreFields(models.Person{}, "Boards", "AssignedTasks")
	for _, mockedPerson := range mockedData.Persons {
		obtainedPerson, err := db.Person.GetByUsername(ctx, mockedPerson.Username)
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
	ctx := context.Background()
	if err := db.System.RecreateAllTables(ctx); err != nil {
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
		createdBoard, err := db.Board.Create(ctx, board)
		if err != nil {
			t.Error(err)
		}

		if !cmp.Equal(createdBoard, board, cmpIgnore) {
			t.Errorf("Created board not equal to mocked: \n\t%v \n\t%v",
				createdBoard, board)
		}

		t.Logf("Created: %v", createdBoard)
	}

	badBoard := models.Board{Name: "badBoard", Owner: models.BoardOwner{ID: 1337}}
	if _, err := db.Board.Create(ctx, badBoard); err == nil {
		t.Error("BoardModel.Create() does't throw error when creating rows with non-existent owner_id")
	}
}

func TestBoardGetByID(t *testing.T) {
	ctx := context.Background()
	if err := db.System.RecreateAllTables(ctx); err != nil {
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

	if _, err := db.Board.GetByID(ctx, 1337); err == nil {
		t.Error("Searching for board with non-existent boardID not throwing error")
	}

	cmpIgnore := cmpopts.IgnoreFields(models.Board{}, "Contributors", "Tasks", "Tags")
	for _, mockedBoard := range mockedData.Boards {
		obtainedBoard, err := db.Board.GetByID(ctx, mockedBoard.ID)
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
	ctx := context.Background()
	if err := db.System.RecreateAllTables(ctx); err != nil {
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
		createdTask, err := db.Task.Create(ctx, mockedTask)
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
	ctx := context.Background()
	if err := db.System.RecreateAllTables(ctx); err != nil {
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

	if _, err := db.Task.GetByID(ctx, 1337); err == nil {
		t.Error("Searching for task with non-existent taskID not throwing error")
	}

	cmpIgnore := cmpopts.IgnoreFields(models.Task{}, "Subtasks", "Tags")
	for _, mockedTask := range mockedData.Tasks {
		obtainedTask, err := db.Task.GetByID(ctx, mockedTask.ID)
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
	ctx := context.Background()
	if err := db.System.RecreateAllTables(ctx); err != nil {
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
		createdTag, err := db.Tag.Create(ctx, mockedTag)
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
	ctx := context.Background()
	if err := db.System.RecreateAllTables(ctx); err != nil {
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

	if _, err := db.Tag.GetByID(ctx, 1337); err == nil {
		t.Error("Searching for tag with non-existent tagID not throwing error")
	}

	cmpIgnore := cmpopts.IgnoreFields(models.Task{}, "ID", "Assignees")
	for _, mockedTag := range mockedData.Tags {
		obtainedTag, err := db.Tag.GetByID(ctx, mockedTag.ID)
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
	ctx := context.Background()
	if err := db.System.RecreateAllTables(ctx); err != nil {
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
		task, err := db.Task.GetByID(ctx, taskTag.TaskID)
		if err != nil {
			t.Error(err)
		}
		tag, err := db.Tag.GetByID(ctx, taskTag.TagID)
		if err != nil {
			t.Error(err)
		}
		task, err = db.Task.AddTagToTask(ctx, tag, task)
		if err != nil {
			t.Error(err)
		}

		mockedTag, err := db.Tag.GetByID(ctx, taskTag.TagID)
		if err != nil {
			t.Error(err)
		}

		for _, tag := range task.Tags {
			if tag.ID == taskTag.TagID {
				if !cmp.Equal(tag, mockedTag) {
					t.Errorf("Obtained tag not equal to mocked: \n\t%v \n\t%v",
						tag, mockedTag)
				} else {
					t.Logf("Successfully added tag to task: %v - %v", task.ID, task.Tags)
					continue OuterFor
				}
			}
		}
		t.Errorf("Tag not added to task.Tags: \n\t%v\n\t%v", tag, task.Tags)
	}
}

func TestTaskAssignPersonToTask(t *testing.T) {
	ctx := context.Background()
	if err := db.System.RecreateAllTables(ctx); err != nil {
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
		task, err := db.Task.GetByID(ctx, mockedAssign.TaskID)
		if err != nil {
			t.Error(err)
		}
		person, err := db.Person.GetByID(ctx, mockedAssign.AssigneeID)
		if err != nil {
			t.Error(err)
		}

		task, err = db.Task.AssignPersonToTask(ctx, person, task)
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
	ctx := context.Background()
	if err := db.System.RecreateAllTables(ctx); err != nil {
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

OuterFor:
	for _, mockedSubtask := range mockedData.Subtasks {
		task, err := db.Task.GetByID(ctx, mockedSubtask.ParentTaskID)
		if err != nil {
			t.Error(err)
		}

		task, err = db.Task.AddSubtaskToTask(ctx, mockedSubtask, task)
		if err != nil {
			t.Error(err)
		}

		for _, subtask := range task.Subtasks {
			if subtask.ID == mockedSubtask.ID {
				if !cmp.Equal(subtask, mockedSubtask) {
					t.Errorf("Obtained subtask not equal to mocked: \n\t%v \n\t%v",
						subtask, mockedSubtask)
				} else {
					t.Logf("Successfully added subtask to task: %v - %v", task.ID, task.Subtasks)
					continue OuterFor
				}
			}
		}

		t.Errorf("Subtask not added to task.Subtasks: \n\t%v\n\t%v", mockedSubtask, task.Assignees)
	}
}

func TestBoardAddTagToBoard(t *testing.T) {
	ctx := context.Background()
	if err := db.System.RecreateAllTables(ctx); err != nil {
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

OuterFor:
	for _, mockedTag := range mockedData.Tags {
		board, err := db.Board.GetByID(ctx, mockedTag.BoardID)
		if err != nil {
			t.Error(err)
		}

		board, err = db.Board.AddTagToBoard(ctx, mockedTag, board)
		if err != nil {
			t.Error(err)
		}

		for _, tag := range board.Tags {
			if tag.ID == mockedTag.ID {
				if !cmp.Equal(tag, mockedTag) {
					t.Errorf("Added tag not equal to mocked: \n\t%v \n\t%v",
						tag, mockedTag)
				} else {
					t.Logf("Successfully added tag to board: %v - %v", board.ID, board.Tags)
					continue OuterFor
				}
			}
		}

		t.Errorf("Tag not added to board.Tags: \n\t%v\n\t%v", mockedTag, board.Tags)
	}
}

func TestBoardAddTaskToBoard(t *testing.T) {
	ctx := context.Background()
	if err := db.System.RecreateAllTables(ctx); err != nil {
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

OuterFor:
	for _, mockedTask := range mockedData.Tasks {
		board, err := db.Board.GetByID(ctx, mockedTask.BoardID)
		if err != nil {
			t.Error(err)
		}

		board, err = db.Board.AddTaskToBoard(ctx, mockedTask, board)
		if err != nil {
			t.Error(err)
		}

		for _, task := range board.Tasks {
			if task.ID == mockedTask.ID {
				if !cmp.Equal(task, mockedTask) {
					t.Errorf("Added task not equal to mocked: \n\t%v \n\t%v",
						task, mockedTask)
				} else {
					t.Logf("Successfully added task to board: %v - %v", board.ID, board.Tasks)
					continue OuterFor
				}
			}
		}

		t.Errorf("Task not added to board.Tasks: \n\t%v\n\t%v", mockedTask, board.Tasks)
	}
}

func TestBoardAddPersonToBoard(t *testing.T) {
	ctx := context.Background()
	if err := db.System.RecreateAllTables(ctx); err != nil {
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

OuterFor:
	for _, mockedContributor := range mockedData.Contributors {
		board, err := db.Board.GetByID(ctx, mockedContributor.BoardID)
		if err != nil {
			t.Error(err)
		}

		person, err := db.Person.GetByID(ctx, mockedContributor.PersonID)
		if err != nil {
			t.Error(err)
		}

		board, err = db.Board.AddPersonToBoard(ctx, person, board)
		if err != nil {
			t.Error(err)
		}

		for _, contributor := range board.Contributors {
			if contributor.ID == person.ID {
				if !person.IsContributor(contributor) {
					t.Errorf("Added contributor not equal to person: \n\t%v \n\t%v",
						contributor, person)
				} else {
					t.Logf("Successfully added contributor to board: %v - %v", board.ID, board.Contributors)
					continue OuterFor
				}
			}

		}

		t.Errorf("Contributor not added to board.Contributors: \n\t%v\n\t%v", person, board.Contributors)
	}
}

func TestPersonLoadAssignedTasks(t *testing.T) {
	ctx := context.Background()
	if err := db.System.RecreateAllTables(ctx); err != nil {
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

OuterFor:
	for _, assignRow := range mockedData.Assignees {
		mockedTask, err := db.Task.GetByID(ctx, assignRow.TaskID)
		if err != nil {
			t.Error(err)
		}

		person, err := db.Person.GetByID(ctx, assignRow.AssigneeID)
		if err != nil {
			t.Error(err)
		}

		mockedTask, err = db.Task.AssignPersonToTask(ctx, person, mockedTask)
		if err != nil {
			t.Error(err)
		}

		person, err = db.Person.GetByID(ctx, person.ID)
		if err != nil {
			t.Error(err)
		}
		for _, task := range person.AssignedTasks {
			if task.ID == assignRow.TaskID {
				if !cmp.Equal(task, mockedTask) {
					t.Errorf("Loaded assigned task not equal to mocked: \n\t%v \n\t%v",
						task, mockedTask)
				} else {
					t.Logf("Successfully loaded assigned task to person: %v - %v", person.ID, person.AssignedTasks)
					continue OuterFor
				}
			}
		}

		t.Errorf("Task not loaded to person.AssignedTasks: \n\t%v\n\t%v", mockedTask, person.AssignedTasks)
	}
}

func TestPersonLoadBoards(t *testing.T) {
	ctx := context.Background()
	if err := db.System.RecreateAllTables(ctx); err != nil {
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

	for _, mockedContributor := range mockedData.Contributors {
		board, err := db.Board.GetByID(ctx, mockedContributor.BoardID)
		if err != nil {
			t.Fatal(err)
		}

		person, err := db.Person.GetByID(ctx, mockedContributor.PersonID)
		if err != nil {
			t.Fatal(err)
		}

		board, err = db.Board.AddPersonToBoard(ctx, person, board)
		if err != nil {
			t.Fatal(err)
		}
	}

OuterFor:
	for _, mockedContributor := range mockedData.Contributors {
		person, err := db.Person.GetByID(ctx, mockedContributor.PersonID)
		if err != nil {
			t.Error(err)
		}

		mockedBoard, err := db.Board.GetByID(ctx, mockedContributor.BoardID)
		if err != nil {
			t.Error(err)
		}

		for _, board := range person.Boards {
			if board.ID == mockedBoard.ID {
				if !cmp.Equal(mockedBoard.Small(), board) {
					t.Errorf("Loaded board not equal to mocked: \n\t%v \n\t%v",
						board, mockedBoard.Small())
				} else {
					t.Logf("Successfully loaded board to person: %v - %v", person.ID, person.Boards)
					continue OuterFor
				}
			}
		}

		t.Errorf("Board not loaded to person.Boards: \n\t%v\n\t%v", mockedBoard.Small(), person.Boards)
	}
}

func TestBoardDeleteByID(t *testing.T) {
	ctx := context.Background()

	if err := db.System.RecreateAllTables(ctx); err != nil {
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

	for _, mockedBoard := range mockedData.Boards {
		err := db.Board.DeleteByID(ctx, mockedBoard.ID)
		if err != nil {
			t.Error(err)
		}
	}

	if err := db.Board.DeleteByID(ctx, 131); err != nil {
		t.Error("Board.DeleteByID() not throwing error when deleting non-existent board")
	}
}
