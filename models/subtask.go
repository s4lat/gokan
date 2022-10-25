package models

// Subtask - subtask model struct.
type Subtask struct {
	Name         string
	ID           uint32
	ParentTaskID uint32
}
