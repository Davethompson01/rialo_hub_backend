package services

import (
	"github.com/Davethompson01/rialo_hub_backend/config"
	repository "github.com/Davethompson01/rialo_hub_backend/internal/Repository"
	"github.com/Davethompson01/rialo_hub_backend/internal/models"
	"github.com/Davethompson01/rialo_hub_backend/internal/validation"
)

func CreateNegotiation(api *config.ApiConfig, negotiate models.SendMessage) (models.NegotiationResponse, error) {
	if err := validation.ValidateNegotiation(negotiate); err != nil {
		return models.NegotiationResponse{}, err
	}

	CreateNegotiation, err := repository.CreateNegotiation(api, negotiate)
	if err != nil {
		return models.NegotiationResponse{}, err
	}

	return CreateNegotiation, nil
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
