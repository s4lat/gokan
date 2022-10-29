//nolint:lll
package db

import (
	"context"

	"github.com/s4lat/gokan/models"
)

// DB - struct for interacting with database.
type DB struct {
	System SystemManager
	Person PersonManager
	Board  BoardManager
	Task   TaskManager
	Tag    TagManager
}

// NewDB - returning new initilized DB.
func NewDB(dbConn models.DBConn) DB {
	return DB{
		System: models.SystemModel{DB: dbConn},
		Person: models.PersonModel{DB: dbConn},
		Board:  models.BoardModel{DB: dbConn},
		Task:   models.TaskModel{DB: dbConn},
		Tag:    models.TagModel{DB: dbConn},
	}
}

// SystemManager - interface for interacting with db structure.
type SystemManager interface {
	RecreateAllTables(ctx context.Context) error
	IsTableExist(ctx context.Context, tableName string) (bool, error)
}

// PersonManager - interface for interacting with person table in db.
type PersonManager interface {
	Create(ctx context.Context, person models.Person) (models.Person, error)
	DeleteByID(ctx context.Context, personID uint32) error
	GetByID(ctx context.Context, personID uint32) (models.Person, error)
	GetByEmail(ctx context.Context, email string) (models.Person, error)
	GetByUsername(ctx context.Context, username string) (models.Person, error)
}

// BoardManager - interface for interacting with board table in db.
type BoardManager interface {
	Create(ctx context.Context, board models.Board) (models.Board, error)
	DeleteByID(ctx context.Context, boardID uint32) error
	GetByID(ctx context.Context, boardID uint32) (models.Board, error)
	AddContributorToBoard(ctx context.Context, contrib models.Contributor, board models.Board) (models.Board, error)
	RemoveContributorFromBoard(ctx context.Context, contrib models.Contributor, board models.Board) (models.Board, error)
	AddTaskToBoard(ctx context.Context, task models.Task, board models.Board) (models.Board, error)
	RemoveTaskFromBoard(ctx context.Context, task models.Task, board models.Board) (models.Board, error)
	AddTagToBoard(ctx context.Context, tag models.Tag, board models.Board) (models.Board, error)
	RemoveTagFromBoard(ctx context.Context, tag models.Tag, board models.Board) (models.Board, error)
}

// TaskManager - interface for interacting with task table in db.
type TaskManager interface {
	Create(ctx context.Context, task models.Task) (models.Task, error)
	DeleteByID(ctx context.Context, taskID uint32) error
	GetByID(ctx context.Context, taskID uint32) (models.Task, error)
	AddTagToTask(ctx context.Context, tag models.Tag, task models.Task) (models.Task, error)
	RemoveTagFromTask(ctx context.Context, tag models.Tag, task models.Task) (models.Task, error)
	AddAssigneeToTask(ctx context.Context, assignee models.TaskAssignee, task models.Task) (models.Task, error)
	RemoveAssignFromTask(ctx context.Context, person models.TaskAssignee, task models.Task) (models.Task, error)
	AddSubtaskToTask(ctx context.Context, subtask models.Subtask, task models.Task) (models.Task, error)
	RemoveSubtaskFromTask(ctx context.Context, subtask models.Subtask, task models.Task) (models.Task, error)
}

// TagManager - interface for interacting with tag table in db.
type TagManager interface {
	Create(ctx context.Context, tag models.Tag) (models.Tag, error)
	DeleteByID(ctx context.Context, tagID uint32) error
	GetByID(ctx context.Context, tagID uint32) (models.Tag, error)
}
