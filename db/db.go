package db

import (
	"github.com/s4lat/gokan/models"
)

type DB struct {
	Person PersonManager
	Board  BoardManager
	Task   TaskManager
	System SystemManager
}

type PersonManager interface {
	Create(p models.Person) (models.Person, error)
	GetByID(personID uint32) (models.Person, error)
	GetByEmail(email string) (models.Person, error)
	GetByUsername(username string) (models.Person, error)
}

type TaskManager interface {
	Create(t models.Task) (models.Task, error)
}

type BoardManager interface {
	Create(b models.Board) (models.Board, error)
	GetByID(boardID uint32) (models.Board, error)
}

type SystemManager interface {
	RecreateAllTables() error
	IsTableExist(tableName string) (bool, error)
}

// DBManager - interface for database managing.
// type DBManager interface {
// CreateTask(t models.Task) (models.Task, error)
// }
