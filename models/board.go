package models

import (
	"context"
	"fmt"
)

// Board - board model struct.
type Board struct {
	Name         string   `json:"board_name"`
	Contributors []Person // LoadBoardContributors() by person_id from contributor table
	Tasks        []Task   // LoadBoardTasks from board_id in task table
	Tags         []Tag    // LoadBoardTags from board_id in tag table
	ID           uint32   `json:"board_id"`
	OwnerID      uint32   `json:"owner_id"`
}

// BoardModel - struct that implements BoardManager interface for interacting with board table in db.
type BoardModel struct {
	DB DBConn
}

// Create - Creates new row in table 'board'.
// Returning created Board.
func (bm BoardModel) Create(board Board) (Board, error) {
	sql := "INSERT INTO board (board_name, owner_id) VALUES ($1, $2) RETURNING *;"

	var createdBoard Board
	err := bm.DB.QueryRow(context.Background(), sql,
		board.Name,
		board.OwnerID,
	).Scan(
		&createdBoard.ID,
		&createdBoard.Name,
		&createdBoard.OwnerID,
	)

	if err != nil {
		return Board{}, fmt.Errorf("BoardModel.Create() -> %w", err)
	}

	return createdBoard, nil
}

// GetByID - searching for board in DB by ID, returning finded Board.
func (bm BoardModel) GetByID(boardID uint32) (Board, error) {
	sql := "SELECT * FROM board WHERE board_id = $1;"

	var obtainedBoard Board
	err := bm.DB.QueryRow(context.Background(), sql, boardID).Scan(
		&obtainedBoard.ID,
		&obtainedBoard.Name,
		&obtainedBoard.OwnerID,
	)

	if err != nil {
		return Board{}, fmt.Errorf("BoardModel.GetByID() -> %w", err)
	}
	return obtainedBoard, nil
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
