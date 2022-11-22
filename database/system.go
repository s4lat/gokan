package database

import (
	"context"
	"fmt"
)

// SystemModel - struct that implements SystemManager interface for interacting with database structure.
type SystemModel struct {
	DB DBConn
}

// RecreateAllTables - drops previously created table and creates tables required by the GoKan.
func (sm SystemModel) RecreateAllTables(ctx context.Context) error {
	err := sm.dropAllTables(ctx)
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
			"author_id INTEGER REFERENCES person (person_id) ON DELETE SET NULL NOT NULL" +
			");")

		createAssigneeSQL = ("" +
			"CREATE TABLE assignee (" +
			"ref_task_id INTEGER REFERENCES task (task_id) ON DELETE CASCADE NOT NULL," +
			"assignee_id INTEGER REFERENCES person (person_id) ON DELETE CASCADE NOT NULL," +
			"CONSTRAINT assignee_pkey PRIMARY KEY (ref_task_id, assignee_id)" +
			");")

		createSubtaskTableSQL = ("" +
			"CREATE TABLE subtask (" +
			"subtask_id serial PRIMARY KEY," +
			"subtask_name VARCHAR NOT NULL," +
			"parent_task_id INTEGER REFERENCES task (task_id) ON DELETE CASCADE NOT NULL" +
			");")

		createTagTableSQL = ("" +
			"CREATE TABLE tag (" +
			"tag_id serial PRIMARY KEY," +
			"tag_name VARCHAR NOT NULL," +
			"tag_description VARCHAR NOT NULL," +
			"board_id INTEGER REFERENCES board (board_id) ON DELETE CASCADE" +
			");")

		createTaskTagTableSQL = ("" +
			"CREATE TABLE task_tag (" +
			"ref_task_id INTEGER REFERENCES task (task_id) ON DELETE CASCADE," +
			"ref_tag_id INTEGER REFERENCES tag (tag_id) ON DELETE CASCADE," +
			"CONSTRAINT task_tag_pkey PRIMARY KEY (ref_task_id, ref_tag_id)" +
			");")

		createContributorTableSQL = ("" +
			"CREATE TABLE contributor (" +
			"person_id INTEGER REFERENCES person (person_id) ON DELETE CASCADE," +
			"board_id INTEGER REFERENCES board (board_id) ON DELETE CASCADE," +
			"CONSTRAINT contributor_pkey PRIMARY KEY (person_id, board_id)" +
			");")
	)

	nullPersonSQL := ("INSERT INTO " +
		"person (person_id, username, first_name, last_name, email, password_hash) " +
		"VALUES (0, 'null', 'null', 'null', 'null', 'null')")

	sqlStrings := []string{
		createPersonTableSQL,
		createBoardTableSQL,
		createTaskTableSQL,
		createAssigneeSQL,
		createSubtaskTableSQL,
		createTagTableSQL,
		createTaskTagTableSQL,
		createContributorTableSQL,
		nullPersonSQL,
	}

	for _, sql := range sqlStrings {
		if _, err := sm.DB.Exec(ctx, sql); err != nil {
			return fmt.Errorf("RecreateAllTables() -> %w", err)
		}
	}

	return nil
}

// IsTableExist - returning `true` if table exist in 'public' scheme, else `false`.
func (sm SystemModel) IsTableExist(ctx context.Context, tableName string) (bool, error) {
	const sql = ("" +
		"SELECT EXISTS (" +
		"SELECT FROM pg_tables " +
		"WHERE schemaname = 'public' " +
		"AND tablename = $1);")

	row := sm.DB.QueryRow(ctx, sql, tableName)

	var isExist bool
	if err := row.Scan(&isExist); err != nil {
		return false, fmt.Errorf("IsTableExist() -> %w", err)
	}

	return isExist, nil
}

// dropAllTables - drops public scheme with all tables.
func (sm SystemModel) dropAllTables(ctx context.Context) error {
	const (
		sql1 = "DROP SCHEMA IF EXISTS public CASCADE"
		sql2 = "CREATE SCHEMA public;"
	)

	if _, err := sm.DB.Exec(ctx, sql1); err != nil {
		return fmt.Errorf("dropAllTables() -> %w", err)
	}

	if _, err := sm.DB.Exec(ctx, sql2); err != nil {
		return fmt.Errorf("dropAllTables() -> %w", err)
	}

	return nil
}
