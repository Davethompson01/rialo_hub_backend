package services

import (
	"fmt"

	"github.com/Davethompson01/rialo_hub_backend/config"
	repository "github.com/Davethompson01/rialo_hub_backend/internal/Repository"
	"github.com/Davethompson01/rialo_hub_backend/internal/models"
	"github.com/Davethompson01/rialo_hub_backend/internal/validation"
)

func CreateNegotiation(
	api *config.ApiConfig,
	negotiate models.SendMessage,
) (models.NegotiationResponse, error) {

	if err := validation.ValidateNegotiation(negotiate); err != nil {
		return models.NegotiationResponse{}, err
	}

	validApplication, err := repository.IsValidApplication(
		api,
		negotiate.ApplicationID,
		negotiate.TaskId,
		negotiate.ApplicantID,
	)
	if err != nil {
		return models.NegotiationResponse{}, err
	}

	if !validApplication {
		return models.NegotiationResponse{}, fmt.Errorf(
			"invalid application for this task and applicant %w,", validApplication,
		)
	}

	negotiation, err := repository.CreateNegotiation(api, negotiate)
	if err != nil {
		return models.NegotiationResponse{}, err
	}

	return negotiation, nil
}
func AcceptOffers(
	api *config.ApiConfig,
	applicationID int,
	taskID int,
	employerID int,
	offerID int,
	conversationID int,
) (string, error) {

	// First accept/assign the employee.
	accept, err := AcceptEmployee(
		api,
		applicationID,
		taskID,
		employerID,
	)

	if err != nil {
		return "", err
	}

	// Then accept the offer.
	err = repository.AcceptOffer(
		api,
		offerID,
		conversationID,
	)

	if err != nil {
		return "", err
	}

	return accept, nil
}

func RejectOffers(
	api *config.ApiConfig,
	applicationID int,
	taskID int,
	offerID int,
	employerID int,
	conversationID int,
) (string, error) {

	rejected, err := RejectEmployee(
		api,
		applicationID,
		taskID,
		employerID,
	)

	if err != nil {
		return "", err
	}

	err = repository.RejectOffer(
		api,
		offerID,
		conversationID,
	)

	if err != nil {
		return "", err
	}

	return rejected, nil
}
func SendMessage(api *config.ApiConfig, message models.SendMessage) (models.MessageResponse, error) {

	SendMessage, err := repository.SendMessage(api, message)
	if err != nil {
		return models.MessageResponse{}, err
	}
	return SendMessage, nil
}

func GetAllApplicantOffer(
	api *config.ApiConfig,
	userID int,
) ([]models.OfferResponse, error) {

	return repository.GetAllApplicantOffer(api, userID)
}
func GetApplicationOffers(
	api *config.ApiConfig,
	userID int,
) ([]models.OfferResponse, error) {

	return repository.GetApplicationOffers(api, userID)
}
func StartNegotiation(
	api *config.ApiConfig,
	negotiate models.SendMessage,
) (models.NegotiationResponse, error) {

	if err := validation.ValidateNegotiation(negotiate); err != nil {
		return models.NegotiationResponse{}, err
	}

	// Check that this application belongs to this applicant
	// and belongs to this task.
	validApplication, err := repository.IsValidApplication(
		api,
		negotiate.ApplicationID,
		negotiate.TaskId,
		negotiate.ApplicantID,
	)
	if err != nil {
		return models.NegotiationResponse{}, err
	}

	if !validApplication {
		return models.NegotiationResponse{}, fmt.Errorf(
			"invalid application for this task",
		)
	}

	// Get the employer from tasks.user_id
	employerID, err := repository.GetTaskEmployer(
		api,
		negotiate.TaskId,
	)
	if err != nil {
		return models.NegotiationResponse{}, err
	}

	// Set employer from database.
	negotiate.EmployerID = employerID

	return repository.CreateNegotiationTransaction(
		api,
		negotiate,
	)
}

func GetEmployerNegotiation(
	api *config.ApiConfig,
	taskID int,
	conversationID int,
	employerID int,
) (models.NegotiationResponse, error) {

	// Verify that this employer owns the task
	valid, err := repository.IsTaskEmployer(
		api,
		taskID,
		employerID,
	)

	if err != nil {
		return models.NegotiationResponse{}, err
	}

	if !valid {
		return models.NegotiationResponse{}, fmt.Errorf(
			"you are not the employer for this task",
		)
	}

	// Get the conversation
	conversation, err := repository.GetConversation(
		api,
		taskID,
		conversationID,
	)

	if err != nil {
		return models.NegotiationResponse{}, err
	}

	// Get all messages
	messages, err := repository.GetConversationMessages(
		api,
		conversationID,
	)

	if err != nil {
		return models.NegotiationResponse{}, err
	}

	return models.NegotiationResponse{
		TaskID:         conversation.TaskID,
		ConversationID: conversation.ConversationID,
		EmployerID:     conversation.EmployerID,
		ApplicantID:    conversation.ApplicantID,
		Messages:       messages,
	}, nil
}

