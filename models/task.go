package models

import (
	"context"
	"fmt"
)

// Task - task model struct.
type Task struct {
	Name        string    `json:"task_name"`
	Description string    `json:"task_description"`
	Subtasks    []Subtask // always load, from parent_task_id in subtask table
	Tags        []Tag     // always load, from task_tag table
	ID          uint32    `json:"task_id"`
	BoardID     uint32    `json:"board_id"`
	AuthorID    uint32    `json:"author_id"`
	ExecutorID  uint32    `json:"executor_id"`
}

// TaskModel - struct that implements TaskManager interface for interacting with task table in db.
type TaskModel struct {
	DB DBConn
}

// Create - Creates new row in table 'task' with values from `t` fields,
// Returning created Task.
func (tm TaskModel) Create(t Task) (Task, error) {
	sql := ("INSERT INTO task " +
		"(task_name, task_description, board_id, author_id, executor_id) " +
		"VALUES ($1, $2, $3, $4, $5) RETURNING *;")

	var createdTask Task
	err := tm.DB.QueryRow(context.Background(), sql,
		t.Name,
		t.Description,
		t.BoardID,
		t.AuthorID,
		t.ExecutorID,
	).Scan(
		&createdTask.ID,
		&createdTask.Name,
		&createdTask.Description,
		&createdTask.BoardID,
		&createdTask.AuthorID,
		&createdTask.ExecutorID,
	)

	if err != nil {
		return Task{}, fmt.Errorf("CreateTask -> %w", err)
	}

	return createdTask, nil
}

// GetByID - searching for task with task_id=TaskID, returning Task.
func (tm TaskModel) GetByID(taskID uint32) (Task, error) {
	sql := "SELECT * FROM task WHERE task_id = $1;"

	var obtainedTask Task
	err := tm.DB.QueryRow(context.Background(), sql, taskID).Scan(
		&obtainedTask.ID,
		&obtainedTask.Name,
		&obtainedTask.Description,
		&obtainedTask.BoardID,
		&obtainedTask.AuthorID,
		&obtainedTask.ExecutorID,
	)

	if err != nil {
		return Task{}, fmt.Errorf("PersonModel.GetByID() -> %w", err)
	}

	obtainedTask, err = tm.loadTags(obtainedTask)
	if err != nil {
		return Task{}, fmt.Errorf("TaskModel.AddTagToTask() -> %w", err)
	}

	return obtainedTask, nil
}

// AddTagToTask - add tag to task in task_tag table.
func (tm TaskModel) AddTagToTask(tag Tag, task Task) (Task, error) {
	sql := "insert into task_tag (ref_tag_id, ref_task_id) values ($1, $2);"
	_, err := tm.DB.Exec(context.Background(), sql, tag.ID, task.ID)
	if err != nil {
		return Task{}, fmt.Errorf("TaskModel.AddTagToTask() -> %w", err)
	}

	task, err = tm.loadTags(task)
	if err != nil {
		return Task{}, fmt.Errorf("TaskModel.AddTagToTask() -> %w", err)
	}
	return task, nil
}

// loadTags - loading tags in Task.Tags slice.
func (tm TaskModel) loadTags(task Task) (Task, error) {
	sql := ("SELECT tag.* FROM task " +
		"JOIN task_tag ON task_id = ref_task_id " +
		"JOIN tag ON tag_id = ref_tag_id " +
		"WHERE task_id = $1")

	rows, _ := tm.DB.Query(context.Background(), sql, task.ID)
	defer rows.Close()

	var tags []Tag
	for rows.Next() {
		var tag Tag
		err := rows.Scan(&tag.ID, &tag.Name, &tag.Description, &tag.BoardID)
		if err != nil {
			return Task{}, fmt.Errorf("TaskModel.loadTags -> %w", err)
		}
		tags = append(tags, tag)
	}

	if err := rows.Err(); err != nil {
		return Task{}, fmt.Errorf("TaskModel.loadTags -> %w", err)
	}

	task.Tags = tags
	return task, nil
}
