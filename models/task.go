package models

import (
	"context"
	"fmt"
)

// Task - task model struct.
type Task struct {
	Assignees   []TaskAssignee
	Author      TaskAuthor `json:"author"`
	Name        string     `json:"task_name"`
	Description string     `json:"task_description"`
	Subtasks    []Subtask
	Tags        []Tag
	ID          uint32 `json:"task_id"`
	BoardID     uint32 `json:"board_id"`
}

// Subtask - subtask model struct.
type Subtask struct {
	Name         string `json:"subtask_name"`
	ID           uint32 `json:"subtask_id"`
	ParentTaskID uint32 `json:"parent_task_id"`
}

// TaskAuthor - other name for SmallPerson struct, used for representing task author in Task struct.
type TaskAuthor SmallPerson

// TaskAssignee - other name for SmallPerson struct, used for representing task executor in Task struct.
type TaskAssignee SmallPerson

// TaskModel - struct that implements TaskManager interface for interacting with task table in db.
type TaskModel struct {
	DB DBConn
}

// Create - Creates new row in table 'task' with values from `t` fields,
// Returning created Task.
//
// Don't use directly, to create new task use BoardModel.AddTaskToBoard.
func (tm TaskModel) Create(t Task) (Task, error) {
	sql := ("WITH inserted_task AS (" +
		"INSERT INTO task " +
		"(task_name, task_description, board_id, author_id) " +
		"VALUES ($1, $2, $3, $4) RETURNING *) " +
		"SELECT inserted_task.*, username, first_name, last_name, email " +
		"FROM inserted_task JOIN person ON person_id = author_id;")

	var createdTask Task
	err := tm.DB.QueryRow(context.Background(), sql,
		t.Name,
		t.Description,
		t.BoardID,
		t.Author.ID,
	).Scan(
		&createdTask.ID,
		&createdTask.Name,
		&createdTask.Description,
		&createdTask.BoardID,
		&createdTask.Author.ID,
		&createdTask.Author.Username,
		&createdTask.Author.FirstName,
		&createdTask.Author.LastName,
		&createdTask.Author.Email,
	)

	if err != nil {
		return Task{}, fmt.Errorf("CreateTask -> %w", err)
	}

	return createdTask, nil
}

// GetByID - searching for task with task_id=taskID, returning Task.
func (tm TaskModel) GetByID(taskID uint32) (Task, error) {
	sql := ("SELECT task.*, " +
		"person.username, person.first_name, person.last_name, person.email " +
		"FROM task " +
		"JOIN person ON person_id = author_id " +
		"WHERE task_id = $1")

	var obtainedTask Task
	err := tm.DB.QueryRow(context.Background(), sql, taskID).Scan(
		&obtainedTask.ID,
		&obtainedTask.Name,
		&obtainedTask.Description,
		&obtainedTask.BoardID,
		&obtainedTask.Author.ID,
		&obtainedTask.Author.Username,
		&obtainedTask.Author.FirstName,
		&obtainedTask.Author.LastName,
		&obtainedTask.Author.Email,
	)

	if err != nil {
		return Task{}, fmt.Errorf("TaskModel.GetByID() -> %w", err)
	}

	obtainedTask, err = tm.loadEverything(obtainedTask)
	if err != nil {
		return Task{}, fmt.Errorf("TaskModel.GetByID() -> %w", err)
	}

	return obtainedTask, nil
}

// AssignPersonToTask - assigning task to person in assignee table.
func (tm TaskModel) AssignPersonToTask(person Person, task Task) (Task, error) {
	sql := "INSERT INTO assignee (ref_task_id, assignee_id) VALUES ($1, $2);"
	_, err := tm.DB.Exec(context.Background(), sql, task.ID, person.ID)
	if err != nil {
		return Task{}, fmt.Errorf("TaskModel.AssignTaskToPerson() -> %w", err)
	}

	task, err = tm.loadAssignees(task)
	if err != nil {
		return Task{}, fmt.Errorf("TaskModel.AssignTaskToPerson() -> %w", err)
	}
	return task, nil
}

// AddTagToTask - add tag to task in task_tag table.
func (tm TaskModel) AddTagToTask(tag Tag, task Task) (Task, error) {
	sql := "INSERT INTO task_tag (ref_tag_id, ref_task_id) VALUES ($1, $2);"
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

// AddSubtaskToTask - add subtask to task in subtask table.
func (tm TaskModel) AddSubtaskToTask(subtask Subtask, task Task) (Task, error) {
	sql := ("INSERT INTO subtask (subtask_name, parent_task_id) " +
		"VALUES ($1, $2);")
	_, err := tm.DB.Exec(context.Background(), sql, subtask.Name, subtask.ParentTaskID)
	if err != nil {
		return Task{}, fmt.Errorf("TaskModel.AddSubtaskToTask() -> %w", err)
	}

	task, err = tm.loadSubtasks(task)
	if err != nil {
		return Task{}, fmt.Errorf("TaskModel.AddSubtaskToTask() -> %w", err)
	}

	return task, nil
}

// loadEverything - combines loadTags, loadSubtasks, loadAssignees in one method.
func (tm TaskModel) loadEverything(task Task) (Task, error) {
	task, err := tm.loadTags(task)
	if err != nil {
		return Task{}, fmt.Errorf("TaskModel.loadEverything() -> %w", err)
	}

	task, err = tm.loadSubtasks(task)
	if err != nil {
		return Task{}, fmt.Errorf("TaskModel.loadEverything() -> %w", err)
	}

	task, err = tm.loadAssignees(task)
	if err != nil {
		return Task{}, fmt.Errorf("TaskModel.loadEverything() -> %w", err)
	}

	task, err = tm.loadTags(task)
	if err != nil {
		return Task{}, fmt.Errorf("TaskModel.loadEverything() -> %w", err)
	}

	return task, nil
}

// loadSubtasks - loading subtasks to Task.Subtasks list.
func (tm TaskModel) loadSubtasks(task Task) (Task, error) {
	sql := ("SELECT subtask_id, subtask_name, parent_task_id " +
		"FROM task JOIN subtask ON parent_task_id = task_id " +
		"WHERE task_id = $1")

	rows, _ := tm.DB.Query(context.Background(), sql, task.ID)
	defer rows.Close()

	var subtasks []Subtask
	for rows.Next() {
		var subtask Subtask
		err := rows.Scan(&subtask.ID, &subtask.Name, &subtask.ParentTaskID)
		if err != nil {
			return Task{}, fmt.Errorf("TaskModel.loadSubtasks() -> %w", err)
		}

		subtasks = append(subtasks, subtask)
	}

	if err := rows.Err(); err != nil {
		return Task{}, fmt.Errorf("TaskModel.loadSubtasks() -> %w", err)
	}

	task.Subtasks = subtasks
	return task, nil
}

// loadAssignees - loading assigness to Task.Assigness list.
func (tm TaskModel) loadAssignees(task Task) (Task, error) {
	sql := ("SELECT assignee_id, " +
		"person.username, person.first_name, person.last_name, person.email " +
		"FROM assignee JOIN person ON person_id = assignee_id " +
		"WHERE ref_task_id = $1")

	rows, _ := tm.DB.Query(context.Background(), sql, task.ID)
	defer rows.Close()

	var assignees []TaskAssignee
	for rows.Next() {
		var assignee TaskAssignee
		err := rows.Scan(&assignee.ID, &assignee.Username, &assignee.FirstName,
			&assignee.LastName, &assignee.Email)
		if err != nil {
			return Task{}, fmt.Errorf("TaskModel.loadAssigness() -> %w", err)
		}

		assignees = append(assignees, assignee)
	}

	if err := rows.Err(); err != nil {
		return Task{}, fmt.Errorf("TaskModel.loadAssigness() -> %w", err)
	}

	task.Assignees = assignees
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
			return Task{}, fmt.Errorf("TaskModel.loadTags() -> %w", err)
		}
		tags = append(tags, tag)
	}

	if err := rows.Err(); err != nil {
		return Task{}, fmt.Errorf("TaskModel.loadTags() -> %w", err)
	}

	task.Tags = tags
	return task, nil
}
