package postgresdb

import (
	"context"
	"fmt"

	// "github.com/jackc/pgerrcode"
	// "errors"
	// "github.com/jackc/pgx/v5/pgconn".
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/s4lat/gokan/models"
	// "log".
)

// PostgresDB - struct that implements DBManager interface using pgx module.
type PostgresDB struct {
	*pgxpool.Pool
}

// NewPostgresDB - returning new instance of PostgresDB with pgxpool connected to dbURL.
func NewPostgresDB(dbURL string) (PostgresDB, error) {
	dbPool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		return PostgresDB{}, fmt.Errorf("NewPostgreSQLManager(): %w", err)
	}
	return PostgresDB{dbPool}, nil
}

func (pdb *PostgresDB) GetBoardByID(boardID uint32) (models.Board, error) {
	sql := "SELECT * FROM board WHERE board_id = $1;"

	var obtainedBoard models.Board
	err := pdb.QueryRow(context.Background(), sql, boardID).Scan(
		&obtainedBoard.ID,
		&obtainedBoard.Name,
		&obtainedBoard.OwnerID,
	)

	if err != nil {
		return models.Board{}, fmt.Errorf("GetBoardByID() -> %w", err)
	}
	return obtainedBoard, nil
}

// GetPersonByUsername - searching for person in DB by username, returning finded Person.
func (pdb *PostgresDB) GetPersonByUsername(username string) (models.Person, error) {
	sql := "SELECT * FROM person WHERE username = $1;"

	var obtainedPerson models.Person
	err := pdb.QueryRow(context.Background(), sql, username).Scan(
		&obtainedPerson.ID,
		&obtainedPerson.Username,
		&obtainedPerson.FirstName,
		&obtainedPerson.LastName,
		&obtainedPerson.Email,
		&obtainedPerson.PasswordHash,
	)

	if err != nil {
		return models.Person{}, fmt.Errorf("GetPersonByUsername() -> %w", err)
	}
	return obtainedPerson, nil
}

// GetPersonByEmail - searching for person in DB by email, returning finded Person.
func (pdb *PostgresDB) GetPersonByEmail(email string) (models.Person, error) {
	sql := "SELECT * FROM person WHERE email = $1;"

	var obtainedPerson models.Person
	err := pdb.QueryRow(context.Background(), sql, email).Scan(
		&obtainedPerson.ID,
		&obtainedPerson.Username,
		&obtainedPerson.FirstName,
		&obtainedPerson.LastName,
		&obtainedPerson.Email,
		&obtainedPerson.PasswordHash,
	)

	if err != nil {
		return models.Person{}, fmt.Errorf("GetPersonByEmail() -> %w", err)
	}
	return obtainedPerson, nil
}

// GetPersonByID - searching for person in DB by id, returning finded Person.
func (pdb *PostgresDB) GetPersonByID(personID uint32) (models.Person, error) {
	sql := "SELECT * FROM person WHERE person_id = $1;"

	var obtainedPerson models.Person
	err := pdb.QueryRow(context.Background(), sql, personID).Scan(
		&obtainedPerson.ID,
		&obtainedPerson.Username,
		&obtainedPerson.FirstName,
		&obtainedPerson.LastName,
		&obtainedPerson.Email,
		&obtainedPerson.PasswordHash,
	)

	if err != nil {
		return models.Person{}, fmt.Errorf("GetPersonByID() -> %w", err)
	}
	return obtainedPerson, nil
}

// CreatePerson - Creates new row in table 'person' with values from `p` fields,
// Returning created Person.
func (pdb *PostgresDB) CreatePerson(p models.Person) (models.Person, error) {
	sql := ("INSERT INTO " +
		"person (username, first_name, last_name, email, password_hash) " +
		"VALUES ($1, $2, $3, $4, $5)" +
		"RETURNING *;")

	var createdPerson models.Person
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
		return models.Person{}, fmt.Errorf("CreatePerson() -> %w", err)
	}

	return createdPerson, nil
}

// CreateBoard - Creates new row in table 'board' with values from `b` fields,
// Returning created Board.
func (pdb *PostgresDB) CreateBoard(b models.Board) (models.Board, error) {
	sql := "INSERT INTO board (board_name, owner_id) VALUES ($1, $2) RETURNING *;"

	var createdBoard models.Board
	err := pdb.QueryRow(context.Background(), sql,
		b.Name,
		b.OwnerID,
	).Scan(
		&createdBoard.ID,
		&createdBoard.Name,
		&createdBoard.OwnerID,
	)

	if err != nil {
		return models.Board{}, fmt.Errorf("CreateBoard() -> %w", err)
	}

	return createdBoard, nil
}

// RecreateAllTables - drops previously created table and creates tables required by the GoKan.
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
			"owner_id INTEGER REFERENCES person (person_id) ON DELETE CASCADE NOT NULL" +
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

	sqlStrings := []string{
		createPersonTableSQL,
		createBoardTableSQL,
		createTaskTableSQL,
		createSubtaskTableSQL,
		createTagTableSQL,
		createTaskTagTableSQL,
		createContributorTableSQL,
	}

	for _, sql := range sqlStrings {
		if _, err := pdb.Exec(context.Background(), sql); err != nil {
			return fmt.Errorf("RecreateAllTables() -> %w", err)
		}
	}

	return nil
}

// IsTableExist - returning `true` if table exist in 'public' scheme, else `false`.
func (pdb *PostgresDB) IsTableExist(tableName string) (bool, error) {
	const sql = ("" +
		"SELECT EXISTS (" +
		"SELECT FROM pg_tables " +
		"WHERE schemaname = 'public' " +
		"AND tablename = $1);")

	row := pdb.QueryRow(context.Background(), sql, tableName)

	var isExist bool
	if err := row.Scan(&isExist); err != nil {
		return false, fmt.Errorf("IsTableExist() -> %w", err)
	}

	return isExist, nil
}

// dropAllTables - drops public scheme with all tables.
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
