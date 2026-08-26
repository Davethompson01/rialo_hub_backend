package models

import "time"

type SendMessage struct {
	MessageID      int
	TaskId         int
	ConversationID int
	SenderID       int
	EmployerID     int
	ApplicationID  int
	ApplicantID    int
	RecieverID     int
	Content        string
	Offer          CreateOffers
	Status         string
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      time.Time
	CreatedBy      int
}

type NegotiationResponse struct {
	TaskID         int       `json:"task_id"`
	ConversationID int       `json:"conversation_id"`
	EmployerID     int       `json:"employer_id"`
	ApplicantID    int       `json:"applicant_id"`
	OfferID        int       `json:"offer_id"`
	Amount         int       `json:"amount"`
	Status         string    `json:"status"`
	CreatedAt      time.Time `json:"created_at"`
}


type MessageResponse struct {
    MessageID      int       `json:"message_id"`
    ConversationID int       `json:"conversation_id"`
    SenderID       int       `json:"sender_id"`
    Content        string    `json:"content"`
    CreatedAt      time.Time `json:"created_at"`
}


