package models

import "time"

type Task struct {
	ID          int
	UserID      int
	Title       string
	Description string
	Reward      int
	Status      string
	Deadline    time.Time
	CreatedAt   time.Time
}

type TaskApplication struct {
	Task_id int
	Employee_id int    `json:"employee_id"`
	Employer_id int    `json:"employer_id"`
	Skills      string `json:"skills"`
	Status      string
	AppliedAt time.Time
}

type ApplicationResponse struct {
    ID          int
    TaskID      int
    ApplicantID int
    Username    string
    Avatar      string
    Reputation  int
    Status      string
}
