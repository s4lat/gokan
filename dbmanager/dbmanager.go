package dbmanager

import (
	"context"
	"fmt"
	// "github.com/jackc/pgerrcode"
	// "errors"
	// "github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	// "log"
)

type DBManager interface {
	CreatePerson(p Person) (Person, error)
	GetPersonByID(person_id uint32) (Person, error)
	GetPersonByEmail(email string) (Person, error)
	GetPersonByUsername(username string) (Person, error)
	RecreateAllTables() error
	IsTableExist(table_name string) (bool, error)
}

type PostgresDB struct {
	*pgxpool.Pool
}

// Creating returning new instance of PostgresDB with pool connected to dbURL
func NewPostgresDB(dbURL string) (PostgresDB, error) {
	dbPool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		return PostgresDB{}, fmt.Errorf("NewPostgreSQLManager(): %w", err)
	}
	return PostgresDB{dbPool}, nil
}

// Searching for person in DB by username, returning finded Person
func (pdb *PostgresDB) GetPersonByUsername(username string) (Person, error) {
	sql := "SELECT * FROM person WHERE username = $1;"

	var obtainedPerson Person
	err := pdb.QueryRow(context.Background(), sql, username).Scan(
		&obtainedPerson.ID,
		&obtainedPerson.Username,
		&obtainedPerson.FirstName,
		&obtainedPerson.LastName,
		&obtainedPerson.Email,
		&obtainedPerson.PasswordHash,
	)

	if err != nil {
		return Person{}, fmt.Errorf("GetPersonByUsername() -> %w", err)
	}
	return obtainedPerson, nil
}

// Searching for person in DB by email, returning finded Person
func (pdb *PostgresDB) GetPersonByEmail(email string) (Person, error) {
	sql := "SELECT * FROM person WHERE email = $1;"

	var obtainedPerson Person
	err := pdb.QueryRow(context.Background(), sql, email).Scan(
		&obtainedPerson.ID,
		&obtainedPerson.Username,
		&obtainedPerson.FirstName,
		&obtainedPerson.LastName,
		&obtainedPerson.Email,
		&obtainedPerson.PasswordHash,
	)

	if err != nil {
		return Person{}, fmt.Errorf("GetPersonByEmail() -> %w", err)
	}
	return obtainedPerson, nil
}

// Searching for person in DB by id, returning finded Person
func (pdb *PostgresDB) GetPersonByID(person_id uint32) (Person, error) {
	sql := "SELECT * FROM person WHERE person_id = $1;"

	var obtainedPerson Person
	err := pdb.QueryRow(context.Background(), sql, person_id).Scan(
		&obtainedPerson.ID,
		&obtainedPerson.Username,
		&obtainedPerson.FirstName,
		&obtainedPerson.LastName,
		&obtainedPerson.Email,
		&obtainedPerson.PasswordHash,
	)

	if err != nil {
		return Person{}, fmt.Errorf("GetPersonByID() -> %w", err)
	}
	return obtainedPerson, nil
}

// Creates new row in table 'person' with values from p fields
// Returning created Person
func (pdb *PostgresDB) CreatePerson(p Person) (Person, error) {
	sql := ("INSERT INTO " +
		"person (username, first_name, last_name, email, password_hash) " +
		"VALUES ($1, $2, $3, $4, $5)" +
		"RETURNING *;")

	var createdPerson Person
	err := pdb.QueryRow(context.Background(), sql,
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
		return Person{}, fmt.Errorf("CreatePerson() -> %w", err)
	}

	return createdPerson, nil
}

// Delete previously created and create all new tables required by the GoKan
func (pdb *PostgresDB) RecreateAllTables() error {
	err := pdb.dropAllTables()
	if err != nil {
		return fmt.Errorf("RecreateAllTables() -> %w", err)
	}

	const (
		createPersonTableSQL = ("" +
			"CREATE TABLE person (" +
			"person_id serial PRIMARY KEY," +
			"username VARCHAR UNIQUE NOT NULL," +
			"first_name VARCHAR NOT NULL," +
			"last_name VARCHAR NOT NULL," +
			"email VARCHAR UNIQUE NOT NULL," +
			"password_hash VARCHAR NOT NULL" +
			");")

		createBoardTableSQL = ("" +
			"CREATE TABLE board (" +
			"board_id serial PRIMARY KEY," +
			"board_name VARCHAR NOT NULL," +
			"owner_id INTEGER REFERENCES person (person_id) ON DELETE SET NULL" +
			");")

		createTaskTableSQL = ("" +
			"CREATE TABLE task (" +
			"task_id serial PRIMARY KEY," +
			"task_name VARCHAR NOT NULL," +
			"task_description VARCHAR," +
			"board_id INTEGER REFERENCES board (board_id) ON DELETE CASCADE," +
			"author_id INTEGER REFERENCES person (person_id) ON DELETE SET NULL," +
			"executor_id INTEGER REFERENCES person (person_id) ON DELETE SET NULL" +
			");")

		createSubtaskTableSQL = ("" +
			"CREATE TABLE subtask (" +
			"subtask_id serial PRIMARY KEY," +
			"subtask_name VARCHAR NOT NULL," +
			"parent_task_id INTEGER REFERENCES task (task_id) ON DELETE CASCADE" +
			");")

		createTagTableSQL = ("" +
			"CREATE TABLE tag (" +
			"tag_id serial PRIMARY KEY," +
			"tag_name VARCHAR NOT NULL," +
			"description VARCHAR NOT NULL," +
			"board_id INTEGER REFERENCES board (board_id) ON DELETE CASCADE" +
			");")

		createTaskTagTableSQL = ("" +
			"CREATE TABLE task_tag (" +
			"task_id INTEGER REFERENCES task (task_id) ON DELETE CASCADE," +
			"tag_id INTEGER REFERENCES tag (tag_id) ON DELETE CASCADE," +
			"CONSTRAINT task_tag_pkey PRIMARY KEY (task_id, tag_id)" +
			");")

		createContributorTableSQL = ("" +
			"CREATE TABLE contributor (" +
			"person_id INTEGER REFERENCES person (person_id) ON DELETE CASCADE," +
			"board_id INTEGER REFERENCES board (board_id) ON DELETE CASCADE," +
			"CONSTRAINT contributor_pkey PRIMARY KEY (person_id, board_id)" +
			");")
	)

	sql_strings := []string{
		createPersonTableSQL,
		createBoardTableSQL,
		createTaskTableSQL,
		createSubtaskTableSQL,
		createTagTableSQL,
		createTaskTagTableSQL,
		createContributorTableSQL,
	}

	for _, sql := range sql_strings {
		if _, err := pdb.Exec(context.Background(), sql); err != nil {
			return fmt.Errorf("RecreateAllTables() -> %w", err)
		}
	}

	return nil
}

// Returning true if table exist in 'public' scheme, else false
func (pdb *PostgresDB) IsTableExist(table_name string) (bool, error) {
	const sql = ("" +
		"SELECT EXISTS (" +
		"SELECT FROM pg_tables " +
		"WHERE schemaname = 'public' " +
		"AND tablename = $1);")

	row := pdb.QueryRow(context.Background(), sql, table_name)

	var isExist bool
	if err := row.Scan(&isExist); err != nil {
		return false, fmt.Errorf("IsTableExist() -> %w", err)
	}

	return isExist, nil
}

// Drops public scheme with all tables
func (pdb *PostgresDB) dropAllTables() error {
	const (
		sql1 = "DROP SCHEMA public CASCADE"
		sql2 = "CREATE SCHEMA public;"
	)

	if _, err := pdb.Exec(context.Background(), sql1); err != nil {
		return fmt.Errorf("dropAllTables() -> %w", err)
	}

	if _, err := pdb.Exec(context.Background(), sql2); err != nil {
		return fmt.Errorf("dropAllTables() -> %w", err)
	}

	return nil
}
