package dbmanager

// Rewrite to not gorm
type Person struct {
	ID        uint32
	Username  string
	FirstName string
	LastName  string
	Email     string
	PWDH      string
}

type Board struct {
	ID      uint32
	Name    string
	OwnerID uint32
}

type Task struct {
	ID          uint32
	Name        string
	Description string
	BoardID     uint32
	CreatorID   uint32
	ExecutorID  uint32
}

type Subtask struct {
	ID           uint32
	Name         string
	ParentTaskID uint32
}

type Tag struct {
	ID         uint32
	Name       string
	Descripton string
}
