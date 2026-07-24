package models

import "time"

type CommunityAlert struct {
	ID int

	TaskID int

	EmployerID int

	Role string

	CreatedAt time.Time
}
