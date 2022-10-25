package models

type Subtask struct {
	Name         string
	ID           uint32
	ParentTaskID uint32
}
