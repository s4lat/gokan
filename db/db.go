package db

import (
	"github.com/s4lat/gokan/models"
)

// DB - struct for interacting with database.
type DB struct {
	Person PersonManager
	Board  BoardManager
	Task   TaskManager
	System SystemManager
}

// PersonManager - interface for interacting with person table in db.
type PersonManager interface {
	Create(p models.Person) (models.Person, error)
	GetByID(personID uint32) (models.Person, error)
	GetByEmail(email string) (models.Person, error)
	GetByUsername(username string) (models.Person, error)
}

// TaskManager - interface for interacting with task table in db.
type TaskManager interface {
	Create(t models.Task) (models.Task, error)
}

// BoardManager - interface for interacting with board table in db.
type BoardManager interface {
	Create(b models.Board) (models.Board, error)
	GetByID(boardID uint32) (models.Board, error)
}

// SystemManager - interface for interacting with db structure.
type SystemManager interface {
	RecreateAllTables() error
	IsTableExist(tableName string) (bool, error)
}
