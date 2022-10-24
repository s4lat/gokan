package dbmanager

import (
	"github.com/s4lat/gokan/models"
)

// DBManager - interface for data base managing.
type DBManager interface {
	CreateBoard(b models.Board) (models.Board, error)
	CreatePerson(p models.Person) (models.Person, error)
	GetPersonByID(personID uint32) (models.Person, error)
	GetPersonByEmail(email string) (models.Person, error)
	GetPersonByUsername(username string) (models.Person, error)
	RecreateAllTables() error
	IsTableExist(tableName string) (bool, error)
}
