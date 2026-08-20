package services

import (
	"github.com/Davethompson01/rialo_hub_backend/config"
	repository "github.com/Davethompson01/rialo_hub_backend/internal/Repository"
	"github.com/Davethompson01/rialo_hub_backend/internal/models"
	"github.com/Davethompson01/rialo_hub_backend/internal/validation"
)

func CreateNegotiation(api *config.ApiConfig, negotiate models.SendMessage) (models.SendMessage, error) {
	if err := validation.ValidateNegotiation(negotiate); err != nil {
		return models.SendMessage{}, err
	}

	CreateNegotiation, err := repository.CreateNegotiation(api, negotiate)
	if err != nil {
		return models.SendMessage{}, err
	}

	return CreateNegotiation, nil
}

func AcceptNegotiation(api *config.ApiConfig, applicationID, taskID, employerID int) (string, error) {
	accept, err := AcceptEmployee(api, applicationID, taskID, employerID)
	if err != nil {
		return "", nil
	}
	return accept, nil
}

func RejecteNegotiation(api *config.ApiConfig, applicationID, taskID, employerID int) (string, error) {
	rejected, err := RejectEmployee(api, applicationID, taskID, employerID)
	if err != nil {
		return "", nil
	}
	return rejected, nil
}
