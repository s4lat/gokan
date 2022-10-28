package db

import (
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
	RecreateAllTables() error
	IsTableExist(tableName string) (bool, error)
}

// PersonManager - interface for interacting with person table in db.
type PersonManager interface {
	Create(person models.Person) (models.Person, error)
	GetByID(personID uint32) (models.Person, error)
	GetByEmail(email string) (models.Person, error)
	GetByUsername(username string) (models.Person, error)
}

// BoardManager - interface for interacting with board table in db.
type BoardManager interface {
	Create(board models.Board) (models.Board, error)
	GetByID(boardID uint32) (models.Board, error)
}

// TaskManager - interface for interacting with task table in db.
type TaskManager interface {
	Create(task models.Task) (models.Task, error)
	GetByID(taskID uint32) (models.Task, error)
	AddTagToTask(tag models.Tag, task models.Task) (models.Task, error)
	AssignPersonToTask(person models.Person, task models.Task) (models.Task, error)
	AddSubtaskToTask(subtask models.Subtask, task models.Task) (models.Task, error)
}

// TagManager - interface for interacting with tag table in db.
type TagManager interface {
	Create(tag models.Tag) (models.Tag, error)
	GetByID(tagID uint32) (models.Tag, error)
}
