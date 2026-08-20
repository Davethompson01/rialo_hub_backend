package models

import "time"

type CreateOffers struct {
	TaskId     int
	EmployerID int
	UserID     int
	NewOffer   int
	Status     string
	CreatedAt  time.Time
}

type OfferNotification struct {
	TaskId     int
	UserID     int
	EmployeeID int
	NewOffer   int
}

type OfferResponse struct {
	OfferID     int
	TaskID      int
	ApplicantId int
	EmployerID  int
	Username    string
	Avatar      string
	Status      string
}
