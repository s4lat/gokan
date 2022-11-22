package database

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

// SmallBoard - is a struct, that used to save board data in some other structs, when
// we don't need to save all board information like contributors, tasks, tags.
type SmallBoard struct {
	Name  string
	Owner BoardOwner
	ID    uint32
}

// Small - return SmallBoard representation of Person.
func (b Board) Small() SmallBoard {
	return SmallBoard{Name: b.Name, Owner: b.Owner, ID: b.ID}
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
func (bm BoardModel) Create(ctx context.Context, board Board) (Board, error) {
	sql := ("WITH inserted_board AS ( " +
		"INSERT INTO board (board_name, owner_id) " +
		"VALUES ($1, $2) RETURNING *) " +
		"SELECT inserted_board.*, username, first_name, last_name, email " +
		"FROM inserted_board JOIN person ON person_id = owner_id;")

	var createdBoard Board
	err := bm.DB.QueryRow(ctx, sql,
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

// DeleteByID - deletes row from table 'board'.
func (bm BoardModel) DeleteByID(ctx context.Context, boardID uint32) error {
	sql := "DELETE FROM board WHERE board_id = $1;"
	_, err := bm.DB.Exec(ctx, sql, boardID)
	if err != nil {
		return fmt.Errorf("BoardModel.DeleteByID() -> %w", err)
	}
	return nil
}

// GetByID - searching for board in DB by ID, returning finded Board.
func (bm BoardModel) GetByID(ctx context.Context, boardID uint32) (Board, error) {
	sql := ("SELECT board.*, username, first_name, last_name, email " +
		"FROM board JOIN person ON person_id = owner_id " +
		"WHERE board_id = $1")

	var obtainedBoard Board
	err := bm.DB.QueryRow(ctx, sql, boardID).Scan(
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

	obtainedBoard, err = bm.loadEverything(ctx, obtainedBoard)
	if err != nil {
		return Board{}, fmt.Errorf("BoardModel.GetByID() -> %w", err)
	}

	return obtainedBoard, nil
}

// AddContributorToBoard - adds row in contributor table with values (person.ID, board.ID).
func (bm BoardModel) AddContributorToBoard(ctx context.Context, contrib Contributor, board Board) (Board, error) {
	if contrib.ID == board.Owner.ID {
		return Board{}, fmt.Errorf("BoardModel.AddPersonToBoard ->" +
			"person is board owner, no need to add in contributors")
	}

	sql := "INSERT INTO contributor (person_id, board_id) VALUES ($1, $2);"
	_, err := bm.DB.Exec(ctx, sql, contrib.ID, board.ID)

	if err != nil {
		return Board{}, fmt.Errorf("BoardModel.AddPersonToBoard() -> %w", err)
	}

	board, err = bm.loadContributors(ctx, board)
	if err != nil {
		return Board{}, fmt.Errorf("BoardModel.AddPersonToBoard() -> %w", err)
	}

	return board, nil
}

// RemoveContributorFromBoard - removes row in contributor table with values (person.ID, board.ID).
func (bm BoardModel) RemoveContributorFromBoard(ctx context.Context, contrib Contributor, board Board) (Board, error) {
	sql := "DELETE FROM contributor WHERE person_id = $1 AND board_id = $2"
	_, err := bm.DB.Exec(ctx, sql, contrib.ID, board.ID)
	if err != nil {
		return Board{}, fmt.Errorf("BoardModel.RemoveContributorFromBoard() -> %w", err)
	}

	board, err = bm.GetByID(ctx, board.ID)
	if err != nil {
		return Board{}, fmt.Errorf("BoardModel.RemoveContributorFromBoard() -> %w", err)
	}
	return board, nil
}

// AddTaskToBoard - add task to table 'task' in db with board_id = board.ID.
func (bm BoardModel) AddTaskToBoard(ctx context.Context, task Task, board Board) (Board, error) {
	task.BoardID = board.ID
	_, err := TaskModel(bm).Create(ctx, task)
	if err != nil {
		return Board{}, fmt.Errorf("BoardModel.AddTaskToBoard() -> %w", err)
	}

	board, err = bm.loadTasks(ctx, board)
	if err != nil {
		return Board{}, fmt.Errorf("BoardModel.GetByID() -> %w", err)
	}

	return board, nil
}

// RemoveTaskFromBoard - removes task from board.
func (bm BoardModel) RemoveTaskFromBoard(ctx context.Context, task Task, board Board) (Board, error) {
	if task.BoardID != board.ID {
		return Board{}, fmt.Errorf("BoardModel.RemoveTaskFromBoard() -> task.BoardID(%d) != board.ID(%d)",
			task.BoardID, board.ID)
	}

	err := TaskModel(bm).DeleteByID(ctx, task.ID)
	if err != nil {
		return Board{}, fmt.Errorf("BoardModel.RemoveTaskFromBoard() -> %w", err)
	}

	board, err = bm.GetByID(ctx, board.ID)
	if err != nil {
		return Board{}, fmt.Errorf("BoardModel.RemoveTaskFromBoard() -> %w", err)
	}
	return board, nil
}

// AddTagToBoard - add tag to table 'tag' in db with board_id = board.ID.
func (bm BoardModel) AddTagToBoard(ctx context.Context, tag Tag, board Board) (Board, error) {
	tag.BoardID = board.ID
	_, err := TagModel(bm).Create(ctx, tag)
	if err != nil {
		return Board{}, fmt.Errorf("BoardModel.AddTagToBoard() -> %w", err)
	}

	board, err = bm.loadTags(ctx, board)
	if err != nil {
		return Board{}, fmt.Errorf("BoardModel.AddTagToBoard() -> %w", err)
	}

	return board, nil
}

// RemoveTagFromBoard - removes task from board.
func (bm BoardModel) RemoveTagFromBoard(ctx context.Context, tag Tag, board Board) (Board, error) {
	if tag.BoardID != board.ID {
		return Board{}, fmt.Errorf("BoardModel.RemoveTagFromBoard() -> tag.BoardID(%d) != board.ID(%d)",
			tag.BoardID, board.ID)
	}

	err := TagModel(bm).DeleteByID(ctx, tag.ID)
	if err != nil {
		return Board{}, fmt.Errorf("BoardModel.RemoveTagFromBoard() -> %w", err)
	}

	board, err = bm.GetByID(ctx, board.ID)
	if err != nil {
		return Board{}, fmt.Errorf("BoardModel.RemoveTagFromBoard() -> %w", err)
	}
	return board, nil
}

// loadContributors - loading contributors in Board.Contributors slice.
func (bm BoardModel) loadContributors(ctx context.Context, board Board) (Board, error) {
	sql := ("SELECT contributor.person_id, " +
		"person.username, person.first_name, person.last_name, person.email " +
		"FROM contributor JOIN person ON person.person_id = contributor.person_id " +
		"WHERE board_id = $1")

	rows, _ := bm.DB.Query(ctx, sql, board.ID)
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
func (bm BoardModel) loadEverything(ctx context.Context, board Board) (Board, error) {
	board, err := bm.loadTags(ctx, board)
	if err != nil {
		return Board{}, fmt.Errorf("BoardModel.loadEverything() -> %w", err)
	}

	board, err = bm.loadTasks(ctx, board)
	if err != nil {
		return Board{}, fmt.Errorf("BoardModel.loadEverything() -> %w", err)
	}

	board, err = bm.loadContributors(ctx, board)
	if err != nil {
		return Board{}, fmt.Errorf("BoardModel.loadEverything() -> %w", err)
	}

	return board, nil
}

// loadTags - loading tags in Board.Tags slice.
func (bm BoardModel) loadTags(ctx context.Context, board Board) (Board, error) {
	sql := "SELECT * FROM tag WHERE board_id=$1"

	rows, _ := bm.DB.Query(ctx, sql, board.ID)
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
func (bm BoardModel) loadTasks(ctx context.Context, board Board) (Board, error) {
	sql := "SELECT task_id FROM task WHERE board_id=$1"

	rows, _ := bm.DB.Query(ctx, sql, board.ID)
	defer rows.Close()

	localTaskModel := TaskModel(bm)
	var tasks []Task
	for rows.Next() {
		var taskID uint32
		err := rows.Scan(&taskID)
		if err != nil {
			return Board{}, fmt.Errorf("BoardModel.loadTasks() -> %w", err)
		}

		task, err := localTaskModel.GetByID(ctx, taskID)
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
