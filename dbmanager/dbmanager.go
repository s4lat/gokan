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
	CreatePerson(person *Person) error
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

// Creates new row in table 'person' with values from p
// Values from created row pass to p by address
func (pm *PostgresDB) CreatePerson(p *Person) error {
	sql := ("INSERT INTO " +
		"person (username, first_name, last_name, email, password_hash) " +
		"VALUES ($1, $2, $3, $4, $5)" +
		"RETURNING *")

	var tmpPerson Person
	err := pm.QueryRow(context.Background(), sql,
		p.Username,
		p.FirstName,
		p.LastName,
		p.Email,
		p.PasswordHash,
	).Scan(
		&tmpPerson.ID,
		&tmpPerson.Username,
		&tmpPerson.FirstName,
		&tmpPerson.LastName,
		&tmpPerson.Email,
		&tmpPerson.PasswordHash,
	)

	if err != nil {
		return fmt.Errorf("CreatePerson() -> %w", err)
	}

	p.ID = tmpPerson.ID
	p.Username = tmpPerson.Username
	p.FirstName = tmpPerson.FirstName
	p.LastName = tmpPerson.LastName
	p.Email = tmpPerson.Email
	p.PasswordHash = tmpPerson.PasswordHash
	return nil
}

// Delete previously created and create all new tables required by the GoKan
func (pm *PostgresDB) RecreateAllTables() error {
	err := pm.dropAllTables()
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
		if _, err := pm.Exec(context.Background(), sql); err != nil {
			return fmt.Errorf("RecreateAllTables() -> %w", err)
		}
	}

	return nil
}

// Returning true if table exist in 'public' scheme, else false
func (pm *PostgresDB) IsTableExist(table_name string) (bool, error) {
	const sql = ("" +
		"SELECT EXISTS (" +
		"SELECT FROM pg_tables " +
		"WHERE schemaname = 'public' " +
		"AND tablename = $1)")

	row := pm.QueryRow(context.Background(), sql, table_name)

	var isExist bool
	if err := row.Scan(&isExist); err != nil {
		return false, fmt.Errorf("IsTableExist() -> %w", err)
	}

	return isExist, nil
}

// Drops public scheme with all tables
func (pm *PostgresDB) dropAllTables() error {
	const (
		sql1 = "DROP SCHEMA public CASCADE"
		sql2 = "CREATE SCHEMA public;"
	)

	if _, err := pm.Exec(context.Background(), sql1); err != nil {
		return fmt.Errorf("dropAllTables() -> %w", err)
	}

	if _, err := pm.Exec(context.Background(), sql2); err != nil {
		return fmt.Errorf("dropAllTables() -> %w", err)
	}

	return nil
}
