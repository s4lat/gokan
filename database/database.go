//nolint:lll
package database

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
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
func NewDB(dbConn DBConn) DB {
	return DB{
		System: SystemModel{DB: dbConn},
		Person: PersonModel{DB: dbConn},
		Board:  BoardModel{DB: dbConn},
		Task:   TaskModel{DB: dbConn},
		Tag:    TagModel{DB: dbConn},
	}
}

// DBConn - interface for data models to interact with db.
type DBConn interface {
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
}

// SystemManager - interface for interacting with db structure.
type SystemManager interface {
	RecreateAllTables(ctx context.Context) error
	IsTableExist(ctx context.Context, tableName string) (bool, error)
}

// PersonManager - interface for interacting with person table in db.
type PersonManager interface {
	Create(ctx context.Context, person Person) (Person, error)
	DeleteByID(ctx context.Context, personID uint32) error
	GetByID(ctx context.Context, personID uint32) (Person, error)
	GetByEmail(ctx context.Context, email string) (Person, error)
	GetByUsername(ctx context.Context, username string) (Person, error)
}

// BoardManager - interface for interacting with board table in db.
type BoardManager interface {
	Create(ctx context.Context, board Board) (Board, error)
	DeleteByID(ctx context.Context, boardID uint32) error
	GetByID(ctx context.Context, boardID uint32) (Board, error)
	AddContributorToBoard(ctx context.Context, contrib Contributor, board Board) (Board, error)
	RemoveContributorFromBoard(ctx context.Context, contrib Contributor, board Board) (Board, error)
	AddTaskToBoard(ctx context.Context, task Task, board Board) (Board, error)
	RemoveTaskFromBoard(ctx context.Context, task Task, board Board) (Board, error)
	AddTagToBoard(ctx context.Context, tag Tag, board Board) (Board, error)
	RemoveTagFromBoard(ctx context.Context, tag Tag, board Board) (Board, error)
}

// TaskManager - interface for interacting with task table in db.
type TaskManager interface {
	Create(ctx context.Context, task Task) (Task, error)
	DeleteByID(ctx context.Context, taskID uint32) error
	GetByID(ctx context.Context, taskID uint32) (Task, error)
	AddTagToTask(ctx context.Context, tag Tag, task Task) (Task, error)
	RemoveTagFromTask(ctx context.Context, tag Tag, task Task) (Task, error)
	AddAssigneeToTask(ctx context.Context, assignee TaskAssignee, task Task) (Task, error)
	RemoveAssignFromTask(ctx context.Context, person TaskAssignee, task Task) (Task, error)
	AddSubtaskToTask(ctx context.Context, subtask Subtask, task Task) (Task, error)
	RemoveSubtaskFromTask(ctx context.Context, subtask Subtask, task Task) (Task, error)
}

// TagManager - interface for interacting with tag table in db.
type TagManager interface {
	Create(ctx context.Context, tag Tag) (Tag, error)
	DeleteByID(ctx context.Context, tagID uint32) error
	GetByID(ctx context.Context, tagID uint32) (Tag, error)
}
