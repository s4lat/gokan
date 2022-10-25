package models

import (
	"context"
	"fmt"
)

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

type PersonModel struct {
	DB DBConn
}

// CreatePerson - Creates new row in table 'person' with values from `p` fields,
// Returning created Person.
func (pm PersonModel) Create(p Person) (Person, error) {
	sql := ("INSERT INTO " +
		"person (username, first_name, last_name, email, password_hash) " +
		"VALUES ($1, $2, $3, $4, $5)" +
		"RETURNING *;")

	var createdPerson Person
	err := pm.DB.QueryRow(context.Background(), sql,
		p.Username,
		p.FirstName,
		p.LastName,
		p.Email,
		p.PasswordHash,
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
