package services

import (
	"errors"

	"github.com/Davethompson01/rialo_hub_backend/config"
	repository "github.com/Davethompson01/rialo_hub_backend/internal/Repository"
	"github.com/Davethompson01/rialo_hub_backend/internal/models"
	"github.com/Davethompson01/rialo_hub_backend/internal/validation"
)

func AcceptEmployee(api *config.ApiConfig, applicationID, taskID int) (string, error) {
	assigned, err := repository.IsTaskAlreadyAssigned(api, taskID)
	if err != nil {
		return "", err
	}

	if assigned {
		return "", errors.New("this task already has an accepted applicant")
	}

	if err := repository.UpdateApplicationStatus(api, applicationID, "accepted"); err != nil {
		return "", err
	}

	return "Applicant accepted successfully", nil
}

func CreateTasks(api *config.ApiConfig, task models.Task) (string, error) {
	if err := validation.ValidateTasks(task); err != nil {
		return err.Error(), err
	}

	CreateTasks := repository.CreateTasks(api, task)
	if CreateTasks != nil {
		return CreateTasks.Error(), nil
	}

	return "Tasks created Successfully", nil
}

func RejectEmployee(api *config.ApiConfig, applicationID int) (string, error) {

	if err := repository.UpdateApplicationStatus(api, applicationID, "rejected"); err != nil {
		return "", err
	}

	return "Applicant rejected successfully", nil
}

func GetTaskApplications(api *config.ApiConfig, taskID, employerID int) ([]models.ApplicationResponse, error) {

	ownsTask, err := repository.IsTaskOwner(api, taskID, employerID)
	if err != nil {
		return nil, err
	}

	if !ownsTask {
		return nil, errors.New("you are not authorized")
	}

	return repository.GetTaskApplications(api, taskID)
}

func GetMyApplications(api *config.ApiConfig, employeeID int) ([]models.ApplicationResponse, error) {

	return repository.GetMyApplications(api, employeeID)
}

func CancelApplication(api *config.ApiConfig, applicationID, employeeID int) (string, error) {

	ownsApplication, err := repository.IsApplicationOwner(api, applicationID, employeeID)
	if err != nil {
		return "", err
	}

	if !ownsApplication {
		return "", errors.New("you are not allowed to cancel this application")
	}

	if err := repository.DeleteApplication(api, applicationID); err != nil {
		return "", err
	}

	return "Application cancelled successfully", nil
}

func CloseTask(api *config.ApiConfig, taskID int) (string, error) {

	if err := repository.UpdateTaskStatus(api, taskID, "closed"); err != nil {
		return "", err
	}

	if err := repository.RejectPendingApplications(api, taskID); err != nil {
		return "", err
	}

	return "Task closed successfully", nil
}

func DeleteTask(api *config.ApiConfig, taskID, ownerID int) (string, error) {

	ownsTask, err := repository.IsTaskOwner(api, taskID, ownerID)
	if err != nil {
		return "", err
	}

	if !ownsTask {
		return "", errors.New("you are not allowed to delete this task")
	}

	if err := repository.DeleteTask(api, taskID); err != nil {
		return "", err
	}

	return "Task deleted successfully", nil
}
