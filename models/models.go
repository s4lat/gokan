package models

type Person struct {
	Username      string  `json:"username"`
	FirstName     string  `json:"first_name"`
	LastName      string  `json:"last_name"`
	Email         string  `json:"email"`
	PasswordHash  string  `json:"password_hash"`
	Boards        []Board // LoadPersonBoards(), by board_id from contributor table
	AssignedTasks []Task  // LoadPersonAssignedTasks(), from executor_id in task
	ID            uint32  `json:"person_id"`
}

type Board struct {
	Name         string   `json:"board_name"`
	Contributors []Person // LoadBoardContributors() by person_id from contributor table
	Tasks        []Task   // LoadBoardTasksfrom board_id in task table
	Tags         []Tag    // from board_id in tag table
	ID           uint32   `json:"board_id"`
	OwnerID      uint32   `json:"owner_id"`
}

type Task struct {
	Name        string    `json:"task_name"`
	Description string    `json:"task_description"`
	Subtasks    []Subtask // LoadTaskSubtasks(), from parent_task_id in subtask table
	Tags        []Tag     // always load, from task_tag table
	ID          uint32    `json:"task_id"`
	BoardID     uint32    `json:"board_id"`
	AuthorID    uint32    `json:"author_id"`
	ExecutorID  uint32    `json:"executor_id"`
}

type Subtask struct {
	Name         string
	ID           uint32
	ParentTaskID uint32
}

type Tag struct {
	Name        string
	Description string
	ID          uint32
	BoardID     uint32
}
