package models

import (
	"context"
	"fmt"
)

// Board - board model struct.
type Board struct {
	Owner        BoardOwner    `json:"owner"`
	Name         string        `json:"board_name"`
	Contributors []Contributor // LoadBoardContributors() by person_id from contributor table
	Tasks        []Task
	Tags         []Tag
	ID           uint32 `json:"board_id"`
}

// BoardOwner - other name for SmallPerson struct, used for representing board owner in Board struct.
type BoardOwner SmallPerson

// Contributor - struct that used to represent contributors(persons) in Board.Contributors field.
type Contributor struct {
	Username  string
	FirstName string
	LastName  string
	Email     string
	ID        uint32
}

// BoardModel - struct that implements BoardManager interface for interacting with board table in db.
type BoardModel struct {
	DB DBConn
}

// Create - Creates new row in table 'board'.
// Returning created Board.
func (bm BoardModel) Create(board Board) (Board, error) {
	sql := ("WITH inserted_board AS ( " +
		"INSERT INTO board (board_name, owner_id) " +
		"VALUES ($1, $2) RETURNING *) " +
		"SELECT inserted_board.*, username, first_name, last_name, email " +
		"FROM inserted_board JOIN person ON person_id = owner_id;")

	var createdBoard Board
	err := bm.DB.QueryRow(context.Background(), sql,
		board.Name,
		board.Owner.ID,
	).Scan(
		&createdBoard.ID,
		&createdBoard.Name,
		&createdBoard.Owner.ID,
		&createdBoard.Owner.Username,
		&createdBoard.Owner.FirstName,
		&createdBoard.Owner.LastName,
		&createdBoard.Owner.Email,
	)

	if err != nil {
		return Board{}, fmt.Errorf("BoardModel.Create() -> %w", err)
	}

	return createdBoard, nil
}

// GetByID - searching for board in DB by ID, returning finded Board.
func (bm BoardModel) GetByID(boardID uint32) (Board, error) {
	sql := ("SELECT board.*, username, first_name, last_name, email " +
		"FROM board JOIN person ON person_id = owner_id " +
		"WHERE board_id = $1")

	var obtainedBoard Board
	err := bm.DB.QueryRow(context.Background(), sql, boardID).Scan(
		&obtainedBoard.ID,
		&obtainedBoard.Name,
		&obtainedBoard.Owner.ID,
		&obtainedBoard.Owner.Username,
		&obtainedBoard.Owner.FirstName,
		&obtainedBoard.Owner.LastName,
		&obtainedBoard.Owner.Email,
	)

	if err != nil {
		return Board{}, fmt.Errorf("BoardModel.GetByID() -> %w", err)
	}

	obtainedBoard, err = bm.loadEverything(obtainedBoard)
	if err != nil {
		return Board{}, fmt.Errorf("BoardModel.GetByID() -> %w", err)
	}

	return obtainedBoard, nil
}

// AddPersonToBoard - adds row in contributor table with values (person.ID, board.ID).
func (bm BoardModel) AddPersonToBoard(person Person, board Board) (Board, error) {
	sql := "INSERT INTO contributor (person_id, board_id) VALUES ($1, $2);"
	_, err := bm.DB.Exec(context.Background(), sql, person.ID, board.ID)

	if err != nil {
		return Board{}, fmt.Errorf("BoardModel.AddPersonToBoard() -> %w", err)
	}

	board, err = bm.loadContributors(board)
	if err != nil {
		return Board{}, fmt.Errorf("BoardModel.AddPersonToBoard() -> %w", err)
	}

	return board, nil
}

// AddTaskToBoard - add task to table 'task' in db with board_id = board.ID.
func (bm BoardModel) AddTaskToBoard(task Task, board Board) (Board, error) {
	task.BoardID = board.ID
	_, err := TaskModel(bm).Create(task)
	if err != nil {
		return Board{}, fmt.Errorf("BoardModel.AddTaskToBoard() -> %w", err)
	}

	board, err = bm.loadTasks(board)
	if err != nil {
		return Board{}, fmt.Errorf("BoardModel.GetByID() -> %w", err)
	}

	return board, nil
}

// AddTagToBoard - add tag to table 'tag' in db with board_id = board.ID.
func (bm BoardModel) AddTagToBoard(tag Tag, board Board) (Board, error) {
	tag.BoardID = board.ID
	_, err := TagModel(bm).Create(tag)
	if err != nil {
		return Board{}, fmt.Errorf("BoardModel.AddTagToBoard() -> %w", err)
	}

	board, err = bm.loadTags(board)
	if err != nil {
		return Board{}, fmt.Errorf("BoardModel.AddTagToBoard() -> %w", err)
	}

	return board, nil
}

// loadContributors - loading contributors in Board.Contributors slice.
func (bm BoardModel) loadContributors(board Board) (Board, error) {
	sql := ("SELECT contributor.person_id, " +
		"person.username, person.first_name, person.last_name, person.email " +
		"FROM contributor JOIN person ON person.person_id = contributor.person_id " +
		"WHERE board_id = $1")

	rows, _ := bm.DB.Query(context.Background(), sql, board.ID)
	defer rows.Close()

	var contributors []Contributor
	for rows.Next() {
		var contributor Contributor
		err := rows.Scan(&contributor.ID, &contributor.Username,
			&contributor.FirstName, &contributor.LastName, &contributor.Email)
		if err != nil {
			return Board{}, fmt.Errorf("BoardModel.loadContributors() -> %w", err)
		}
		contributors = append(contributors, contributor)
	}

	if err := rows.Err(); err != nil {
		return Board{}, fmt.Errorf("BoardModel.loadContributors() -> %w", err)
	}

	board.Contributors = contributors
	return board, nil
}

// loadEverything - combines loadTags, loadTasks, loadContributors in one method.
func (bm BoardModel) loadEverything(board Board) (Board, error) {
	board, err := bm.loadTags(board)
	if err != nil {
		return Board{}, fmt.Errorf("BoardModel.loadEverything() -> %w", err)
	}

	board, err = bm.loadTasks(board)
	if err != nil {
		return Board{}, fmt.Errorf("BoardModel.loadEverything() -> %w", err)
	}

	board, err = bm.loadContributors(board)
	if err != nil {
		return Board{}, fmt.Errorf("BoardModel.loadEverything() -> %w", err)
	}

	return board, nil
}

// loadTags - loading tags in Board.Tags slice.
func (bm BoardModel) loadTags(board Board) (Board, error) {
	sql := "SELECT * FROM tag WHERE board_id=$1"

	rows, _ := bm.DB.Query(context.Background(), sql, board.ID)
	defer rows.Close()

	var tags []Tag
	for rows.Next() {
		var tag Tag
		err := rows.Scan(&tag.ID, &tag.Name, &tag.Description, &tag.BoardID)
		if err != nil {
			return Board{}, fmt.Errorf("BoardModel.loadTags() -> %w", err)
		}
		tags = append(tags, tag)
	}

	if err := rows.Err(); err != nil {
		return Board{}, fmt.Errorf("BoardModel.loadTags() -> %w", err)
	}

	board.Tags = tags
	return board, nil
}

// loadTasks - loading tasks in Board.Tasks slice.
func (bm BoardModel) loadTasks(board Board) (Board, error) {
	sql := "SELECT task_id FROM task WHERE board_id=$1"

	rows, _ := bm.DB.Query(context.Background(), sql, board.ID)
	defer rows.Close()

	localTaskModel := TaskModel(bm)
	var tasks []Task
	for rows.Next() {
		var taskID uint32
		err := rows.Scan(&taskID)
		if err != nil {
			return Board{}, fmt.Errorf("BoardModel.loadTasks() -> %w", err)
		}

		task, err := localTaskModel.GetByID(taskID)
		if err != nil {
			return Board{}, fmt.Errorf("BoardModel.loadTasks() -> %w", err)
		}

		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return Board{}, fmt.Errorf("BoardModel.loadTasks() -> %w", err)
	}

	board.Tasks = tasks
	return board, nil
}
