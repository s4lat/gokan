package dbmanager

type Person struct {
	ID            uint32
	Username      string
	FirstName     string
	LastName      string
	Email         string
	PWDH          string
	Boards        []Board // LoadPersonBoards(), by board_id from contributor table
	AssignedTasks []Task  // LoadPersonAssignedTasks(), from executor_id in task
}

type Board struct {
	ID           uint32
	Name         string
	OwnerID      uint32
	Contributors []Person // LoadBoardContributors() by person_id from contributor table
	Tasks        []Task   // LoadBoardTasksfrom board_id in task table
	Tags         []Tag    // from board_id in tag table
}

type Task struct {
	ID          uint32
	Name        string
	Description string
	BoardID     uint32
	Author      uint32
	ExecutorID  uint32
	Subtasks    []Subtask // LoadTaskSubtasks(), from parent_task_id in subtask table
	Tags        []Tag     // always load, from task_tag table
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
	BoardID    uint32
}
