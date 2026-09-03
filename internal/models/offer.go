package models

import "time"

type CreateOffers struct {
	OfferID        int
	TaskId         int
	EmployerID     int
	UserID         int
	NewOffer       int
	Status         string
	ConversationID int
	CreatedAt      time.Time
}

type OfferNotification struct {
	TaskId         int
	ConversationID int
	UserID         int
	EmployeeID     int
	NewOffer       int
	CreatedAt      time.Time
}

type OfferResponse struct {
	OfferID       int `json:"offer_id"`
	ApplicationID int `json:"application_id"`
	TaskID         int       `json:"task_id"`
	ConversationID int       `json:"conversation_id"`
	EmployerID     int       `json:"employer_id"`
	ApplicantID    int       `json:"applicant_id"`
	CreatedBy      int       `json:"created_by"`
	Amount         float64   `json:"amount"`
	Status         string    `json:"status"`
	Username       string    `json:"username"`
	Avatar         string    `json:"avatar"`
	CreatedAt      time.Time `json:"created_at"`
}

type OfferActionRequest struct {
	OfferID        int `json:"offer_id"`
	ApplicationID  int `json:"application_id"`
	TaskID         int `json:"task_id"`
	ConversationID int `json:"conversation_id"`
}
