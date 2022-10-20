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
	RecreateAllTables() error
	IsTableExist(table_name string) (bool, error)
}

type PostgresDB struct {
	*pgxpool.Pool
}

func NewPostgresDB(dbURL string) (PostgresDB, error) {
	dbPool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		return PostgresDB{}, fmt.Errorf("NewPostgreSQLManager(): %w", err)
	}
	return PostgresDB{dbPool}, nil
}

func (pm *PostgresDB) RecreateAllTables() error {
	err := pm.dropAllTables()
	if err != nil {
		return fmt.Errorf("RecreateAllTables() -> %w", err)
	}

	const (
		createPersonTableSQL = ("" +
			"CREATE TABLE person (" +
			"person_id serial PRIMARY KEY," +
			"username VARCHAR NOT NULL," +
			"first_name VARCHAR NOT NULL," +
			"last_name VARCHAR NOT NULL," +
			"email VARCHAR NOT NULL," +
			"password VARCHAR NOT NULL" +
			");")
	)

	if _, err := pm.Exec(context.Background(), createPersonTableSQL); err != nil {
		return fmt.Errorf("RecreateAllTables() -> %w", err)
	}

	return nil
}

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

func (pm *PostgresDB) dropAllTables() error {
	const (
		sql1 = "DROP SCHEMA public CASCADE"
		sql2 = "CREATE SCHEMA public;"
	)

	if _, err := pm.Exec(context.Background(), sql1); err != nil {
		return fmt.Errorf("DropAllTables() -> %w", err)
	}

	if _, err := pm.Exec(context.Background(), sql2); err != nil {
		return fmt.Errorf("DropAllTables() -> %w", err)
	}

	return nil
}
