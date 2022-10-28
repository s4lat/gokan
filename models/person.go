package models

import (
	"context"
	"fmt"
)

// Person - person model struct.
type Person struct {
	Username      string       `json:"username"`
	FirstName     string       `json:"first_name"`
	LastName      string       `json:"last_name"`
	Email         string       `json:"email"`
	PasswordHash  string       `json:"password_hash"`
	Boards        []SmallBoard // LoadPersonBoards(), by board_id from contributor table
	AssignedTasks []Task       // LoadPersonAssignedTasks(), from executor_id in task
	ID            uint32       `json:"person_id"`
}

// SmallPerson - is a struct, that used to save person data in some other structs, when
// we don't need to save all person information like password, board, assigned tasks and other.
type SmallPerson struct {
	Username  string `json:"username"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	ID        uint32 `json:"person_id"`
}

// Small - return SmallPerson representation of Person.
func (p *Person) Small() SmallPerson {
	return SmallPerson{Username: p.Username, FirstName: p.FirstName,
		LastName: p.LastName, Email: p.Email, ID: p.ID}
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

	obtainedPerson, err = pm.loadEverything(obtainedPerson)
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

	obtainedPerson, err = pm.loadEverything(obtainedPerson)
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

	obtainedPerson, err = pm.loadEverything(obtainedPerson)
	if err != nil {
		return Person{}, fmt.Errorf("PersonModel.GetByUsername() -> %w", err)
	}

	return obtainedPerson, nil
}

// loadEverything - combines loadAssignedTasks, loadBoards  in one method.
func (pm PersonModel) loadEverything(person Person) (Person, error) {
	person, err := pm.loadAssignedTasks(person)
	if err != nil {
		return Person{}, fmt.Errorf("PersonModel.loadEverything() -> %w", err)
	}

	person, err = pm.loadBoards(person)
	if err != nil {
		return Person{}, fmt.Errorf("PersonModel.loadEverything() -> %w", err)
	}

	return person, nil
}

// loadAssignedTasks - loading assigned to person tasks in Person.AssignedTasks slice.
func (pm PersonModel) loadAssignedTasks(person Person) (Person, error) {
	sql := "SELECT ref_task_id FROM assignee WHERE assignee_id = $1"

	rows, _ := pm.DB.Query(context.Background(), sql, person.ID)

	localTaskModel := TaskModel(pm)
	defer rows.Close()
	var assignedTasks []Task
	for rows.Next() {
		var taskID uint32
		err := rows.Scan(&taskID)
		if err != nil {
			return Person{}, fmt.Errorf("PersonModel.loadAssignedTasks() -> %w", err)
		}

		task, err := localTaskModel.GetByID(taskID)
		if err != nil {
			return Person{}, fmt.Errorf("PersonModel.loadAssignedTasks() -> %w", err)
		}

		assignedTasks = append(assignedTasks, task)
	}

	if err := rows.Err(); err != nil {
		return Person{}, fmt.Errorf("PersonModel.loadAssignedTasks() -> %w", err)
	}

	person.AssignedTasks = assignedTasks
	return person, nil
}

// loadBoards - loads owned and contributed by person, boards.
func (pm PersonModel) loadBoards(person Person) (Person, error) {
	sql := ("SELECT board.*, username, first_name, last_name, email " +
		"FROM board JOIN person ON person_id = owner_id " +
		"WHERE owner_id = $1 " +
		"UNION " +
		"SELECT board.*, username, first_name, last_name, email " +
		"FROM contributor " +
		"JOIN board ON contributor.board_id = board.board_id " +
		"JOIN person ON board.owner_id = person.person_id " +
		"WHERE contributor.person_id = $1")

	rows, _ := pm.DB.Query(context.Background(), sql, person.ID)
	var boards []SmallBoard
	for rows.Next() {
		var board SmallBoard
		err := rows.Scan(&board.ID, &board.Name,
			&board.Owner.ID, &board.Owner.Username, &board.Owner.FirstName,
			&board.Owner.LastName, &board.Owner.Email)

		if err != nil {
			return Person{}, fmt.Errorf("PersonModel.loadBoards() -> %w", err)
		}

		boards = append(boards, board)
	}

	if err := rows.Err(); err != nil {
		return Person{}, fmt.Errorf("PersonModel.loadBoards() -> %w", err)
	}

	person.Boards = boards
	return person, nil
}
