package models

import "time"

type Notification struct {
	ID int

	UserID int

	Type string

	Message string

	IsRead bool

	CreatedAt time.Time
}
