package models

import "time"

type SendMessage struct {
	MessageID      int
	TaskId         int
	ConversationID int
	SenderID       int
	EmployerID     int
	ApplicationID int
	ApplicantID    int
	RecieverID int
	Content        string
	Offer          CreateOffers
	Status         string
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      time.Time
}
