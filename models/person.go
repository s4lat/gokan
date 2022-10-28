package models

import (
	"context"
	"fmt"
)

// Person - person model struct.
type Person struct {
	Username      string  `json:"username"`
	FirstName     string  `json:"first_name"`
	LastName      string  `json:"last_name"`
	Email         string  `json:"email"`
	PasswordHash  string  `json:"password_hash"`
	Boards        []Board // LoadPersonBoards(), by board_id from contributor table
	AssignedTasks []Task  // LoadPersonAssignedTasks(), from executor_id in task
	ID            uint32  `json:"person_id"`
}

// IsContributor - checks if p of type Person represent same row from db as contrib of type Contributor.
func (p *Person) IsContributor(contrib Contributor) bool {
	isEqual := true

	switch {
	case contrib.Username != p.Username:
		isEqual = false
	case contrib.FirstName != p.FirstName:
		isEqual = false
	case contrib.LastName != p.LastName:
		isEqual = false
	case contrib.Email != p.Email:
		isEqual = false
	case contrib.ID != p.ID:
		isEqual = false
	}

	return isEqual
}

// PersonModel - struct that implements PersonManager interface for interacting with person table in db.
type PersonModel struct {
	DB DBConn
}

// Create - Creates new row in table 'person'.
// Returning created Person.
func (pm PersonModel) Create(person Person) (Person, error) {
	sql := ("INSERT INTO " +
		"person (username, first_name, last_name, email, password_hash) " +
		"VALUES ($1, $2, $3, $4, $5)" +
		"RETURNING *;")

	var createdPerson Person
	err := pm.DB.QueryRow(context.Background(), sql,
		person.Username,
		person.FirstName,
		person.LastName,
		person.Email,
		person.PasswordHash,
	).Scan(
		&createdPerson.ID,
		&createdPerson.Username,
		&createdPerson.FirstName,
		&createdPerson.LastName,
		&createdPerson.Email,
		&createdPerson.PasswordHash,
	)

	if err != nil {
		return Person{}, fmt.Errorf("PersonModel.Create() -> %w", err)
	}

	return createdPerson, nil
}

// GetByID - searching for person in DB by id, returning finded Person.
func (pm PersonModel) GetByID(personID uint32) (Person, error) {
	sql := "SELECT * FROM person WHERE person_id = $1;"

	var obtainedPerson Person
	err := pm.DB.QueryRow(context.Background(), sql, personID).Scan(
		&obtainedPerson.ID,
		&obtainedPerson.Username,
		&obtainedPerson.FirstName,
		&obtainedPerson.LastName,
		&obtainedPerson.Email,
		&obtainedPerson.PasswordHash,
	)

	if err != nil {
		return Person{}, fmt.Errorf("PersonModel.GetByID() -> %w", err)
	}
	return obtainedPerson, nil
}

// GetByEmail - searching for person in DB by email, returning finded Person.
func (pm PersonModel) GetByEmail(email string) (Person, error) {
	sql := "SELECT * FROM person WHERE email = $1;"

	var obtainedPerson Person
	err := pm.DB.QueryRow(context.Background(), sql, email).Scan(
		&obtainedPerson.ID,
		&obtainedPerson.Username,
		&obtainedPerson.FirstName,
		&obtainedPerson.LastName,
		&obtainedPerson.Email,
		&obtainedPerson.PasswordHash,
	)

	if err != nil {
		return Person{}, fmt.Errorf("PersonModel.GetByEmail() -> %w", err)
	}
	return obtainedPerson, nil
}

// GetByUsername - searching for person in DB by username, returning finded Person.
func (pm PersonModel) GetByUsername(username string) (Person, error) {
	sql := "SELECT * FROM person WHERE username = $1;"

	var obtainedPerson Person
	err := pm.DB.QueryRow(context.Background(), sql, username).Scan(
		&obtainedPerson.ID,
		&obtainedPerson.Username,
		&obtainedPerson.FirstName,
		&obtainedPerson.LastName,
		&obtainedPerson.Email,
		&obtainedPerson.PasswordHash,
	)

	if err != nil {
		return Person{}, fmt.Errorf("PersonModel.GetByUsername() -> %w", err)
	}
	return obtainedPerson, nil
}
