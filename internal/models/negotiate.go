package models

import "time"

type SendMessage struct {
	NegotiationID int
	TaskId        int
	EmployerID    int
	ApplicantID   int
	Message       string
	Offer         CreateOffers
	Status        string
	CreatedAt     time.Time
}
