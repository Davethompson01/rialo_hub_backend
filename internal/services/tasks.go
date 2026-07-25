package services

import (
	"errors"

	"github.com/Davethompson01/rialo_hub_backend/config"
	repositary "github.com/Davethompson01/rialo_hub_backend/internal/Repositary"
	"github.com/Davethompson01/rialo_hub_backend/internal/models"
	"github.com/Davethompson01/rialo_hub_backend/internal/validation"
)

func AcceptEmployee(api *config.ApiConfig, applicationID, taskID int) (string, error) {
	assigned, err := repositary.IsTaskAlreadyAssigned(api, taskID)
	if err != nil {
		return "", err
	}

	if assigned {
		return "", errors.New("this task already has an accepted applicant")
	}

	if err := repositary.UpdateApplicationStatus(api, applicationID, "accepted"); err != nil {
		return "", err
	}

	return "Applicant accepted successfully", nil
}

func CreateTasks(api *config.ApiConfig, task models.Task) (string, error) {
	if err := validation.ValidateTasks(task); err != nil {
		return err.Error(), err
	}

	CreateTasks := repositary.CreateTasks(api, task)
	if CreateTasks != nil {
		return CreateTasks.Error(), nil
	}

	return "Tasks created Successfully", nil
}

func RejectEmployee(api *config.ApiConfig, applicationID int) (string, error) {

	if err := repositary.UpdateApplicationStatus(api, applicationID, "rejected"); err != nil {
		return "", err
	}

	return "Applicant rejected successfully", nil
}


// func CancelApplication(api *config.ApiConfig, applicationID, employeeID int) (string, error) {

// 	ownsApplication, err := repositary.IsApplicationOwner(api, applicationID, employeeID)
// 	if err != nil {
// 		return "", err
// 	}

// 	if !ownsApplication {
// 		return "", errors.New("you are not allowed to cancel this application")
// 	}

// 	if err := repository.DeleteApplication(api, applicationID); err != nil {
// 		return "", err
// 	}

// 	return "Application cancelled successfully", nil
// }