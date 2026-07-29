package models

import (
	"time"
)

type Task struct {
	ID          int
	UserID      int
	Title       string `validate:"required,min=8"`
	Description string `validate:"required,min=8"`
	Reward      int    `validate:"required"`
	Status      string
	Deadline    time.Time `validate:"required"`
	CreatedAt   time.Time
}

type TaskApplication struct {
	ID          int
	Task_id     int
	Employee_id int    `json:"employee_id"`
	Employer_id int    `json:"employer_id"`
	Skills      string `json:"skills"`
	Status      string
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

type UpdateApplicationStatus struct {
	ApplicationID int
	Status        string
}

type Profile struct {
	UserID           int
	Profile_pics     string
	Discord_username string
	Username         string
	Role             string
}

type TaskResponse struct {
	Task Task    `json:"task"`
	User Profile `json:"user"`
}
