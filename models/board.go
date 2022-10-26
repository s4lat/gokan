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

// Create - Creates new row in table 'board' with values from `b` fields,
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
