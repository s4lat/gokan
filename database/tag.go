package database

import (
	"context"
	"fmt"
)

// Tag - tag model struct.
type Tag struct {
	Name        string `json:"tag_name"`
	Description string `json:"tag_description"`
	ID          uint32 `json:"tag_id"`
	BoardID     uint32 `json:"board_id"`
}

// TagModel - struct that implements TagManager interface for interacting with tag table in db.
type TagModel struct {
	DB DBConn
}

// Create - Creates new row in table 'tag'.
// Returning created Person.
//
// Don't use directly, to create new tag use BoardModel.AddTagToBoard.
func (tm TagModel) Create(ctx context.Context, tag Tag) (Tag, error) {
	sql := ("INSERT INTO " +
		"tag (tag_name, tag_description, board_id) " +
		"VALUES ($1, $2, $3)" +
		"RETURNING *;")

	var createdTag Tag
	err := tm.DB.QueryRow(ctx, sql,
		tag.Name,
		tag.Description,
		tag.BoardID,
	).Scan(
		&createdTag.ID,
		&createdTag.Name,
		&createdTag.Description,
		&createdTag.BoardID,
	)

	if err != nil {
		return Tag{}, fmt.Errorf("TagModel.Create() -> %w", err)
	}
	return createdTag, nil
}

// DeleteByID - deletes row from table 'tag'.
func (tm TagModel) DeleteByID(ctx context.Context, tagID uint32) error {
	sql := "DELETE FROM tag WHERE tag_id = $1;"
	_, err := tm.DB.Exec(ctx, sql, tagID)
	if err != nil {
		return fmt.Errorf("TagModel.DeleteByID() -> %w", err)
	}
	return nil
}

// GetByID - searching for tag in DB by ID, returning finded Tag.
func (tm TagModel) GetByID(ctx context.Context, tagID uint32) (Tag, error) {
	sql := "SELECT * FROM tag WHERE tag_id = $1;"

	var obtainedTag Tag
	err := tm.DB.QueryRow(ctx, sql, tagID).Scan(
		&obtainedTag.ID,
		&obtainedTag.Name,
		&obtainedTag.Description,
		&obtainedTag.BoardID,
	)

	if err != nil {
		return Tag{}, fmt.Errorf("BoardModel.GetByID() -> %w", err)
	}
	return obtainedTag, nil
}
