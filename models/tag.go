package models

import (
	"context"
	"fmt"
)

// Tag - tag model struct.
type Tag struct {
	Name        string `json:"tag_name"`
	Description string `json:"tag_description"`
	ID          uint32
	BoardID     uint32 `json:"board_id"`
}

// TagModel - struct that implements TagManager interface for interacting with tag table in db.
type TagModel struct {
	DB DBConn
}

// Create - Creates new row in table 'tag'.
// Returning created Person.
func (tm TagModel) Create(tag Tag) (Tag, error) {
	sql := ("INSERT INTO " +
		"tag (tag_name, tag_description, board_id) " +
		"VALUES ($1, $2, $3)" +
		"RETURNING *;")

	var createdTag Tag
	err := tm.DB.QueryRow(context.Background(), sql,
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
