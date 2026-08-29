package services

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/Davethompson01/rialo_hub_backend/config"
	repository "github.com/Davethompson01/rialo_hub_backend/internal/Repository"
	"github.com/Davethompson01/rialo_hub_backend/internal/models"
	"github.com/Davethompson01/rialo_hub_backend/internal/validation"
)

func AcceptEmployee(api *config.ApiConfig, applicationID, taskID, employerID int) (string, error) {
	ownsTask, err := repository.IsTaskOwner(api, taskID, employerID)
	if err != nil {
		return "", err
	}

	// if !ownsTask {
	// 	return "", errors.New("you are not authorized")
	// }
	if !ownsTask {
		return "", fmt.Errorf("you are not authorized %v", ownsTask)
	}

	assigned, err := repository.IsTaskAlreadyAssigned(api, taskID)
	if err != nil {
		return "", err
	}

	if assigned {
		return "", errors.New("this task already has an accepted applicant")
	}

	if err := repository.UpdateApplicationStatus(api, applicationID, "Accepted"); err != nil {
		return "", err
	}

	UpdateTask := repository.UpdateTaskStatus(api, taskID, "Accepted")
	if UpdateTask != nil {
		return "", UpdateTask
	}

	return "Applicant accepted successfully", nil
}

func CreateTasks(api *config.ApiConfig, task models.Task) (models.TaskResponse, error) {
	if err := validation.ValidateTasks(task); err != nil {
		return models.TaskResponse{}, err
	}

	createdTask, err := repository.CreateTasks(api, task)
	if err != nil {
		return models.TaskResponse{}, err
	}

	profileModels, err := repository.SelectUserByID(api, task.UserID)
	if err != nil {
		return models.TaskResponse{}, err
	}

	return models.TaskResponse{
		Task: createdTask,
		User: profileModels,
	}, nil
}

func RejectEmployee(api *config.ApiConfig, applicationID, taskID, employerID int) (string, error) {
	ownsTask, err := repository.IsTaskOwner(api, taskID, employerID)
	if err != nil {
		return "", err
	}

	if !ownsTask {
		return "", errors.New("you are not authorized")
	}

	exists, err := repository.IsApplicationForTask(api, applicationID, taskID)
	if err != nil {
		return "", err
	}

	if !exists {
		return "", errors.New("application does not belong to this task")
	}
	checkTasksExist := repository.CheckTasksExist(api, taskID)
	if !checkTasksExist {
		return "", fmt.Errorf("Task Not Found %v", checkTasksExist)
	}

	if err := repository.UpdateApplicationStatus(api, applicationID, "Rejected"); err != nil {
		return "", err
	}

	return "Applicant rejected successfully", nil
}

func GetTaskApplications(
	api *config.ApiConfig,
	taskID int,
	userID int,
) ([]models.ApplicationResponse, error) {

	ownerID, err := repository.GetTaskOwner(api, taskID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("task not found")
		}

		return nil, fmt.Errorf(
			"failed to get task owner: %w",
			err,
		)
	}

	if ownerID != userID {
		return nil, fmt.Errorf(
			"you are not authorized to view applications for this task",
		)
	}

	applications, err := repository.GetTaskApplications(
		api,
		taskID,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to get task applications: %w",
			err,
		)
	}

	return applications, nil
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

// func CloseTask(api *config.ApiConfig, taskID, employerID int) (string, error) {

// 	ownsTask, err := repository.IsTaskOwner(api, taskID, employerID)
// 	if err != nil {
// 		return "", err
// 	}

// 	if !ownsTask {
// 		return "", errors.New("you are not authorized")
// 	}
// 	if err := repository.UpdateTaskStatus(api, taskID, "closed"); err != nil {
// 		return "", err
// 	}

// 	if err := repository.RejectPendingApplications(api, taskID); err != nil {
// 		return "", err
// 	}

// 	return "Task closed successfully", nil
// }

func DeleteTask(api *config.ApiConfig, taskID, ownerID int) error {
	ownsTask, err := repository.IsTaskOwner(api, taskID, ownerID)
	if err != nil {
		return fmt.Errorf("failed to verify task ownership: %w", err)
	}

	if !ownsTask {
		return errors.New("you are not allowed to delete this task")
	}

	if err := repository.DeleteTaskWithApplications(api, taskID); err != nil {
		return fmt.Errorf("failed to delete task: %w", err)
	}

	return nil
}

func ApplyForTasks(api *config.ApiConfig, task models.TaskApplication) (models.TaskApplication, error) {
	// 1. Check if task exists
	taskExists := repository.CheckTasksExist(api, task.Task_id)

	if !taskExists {
		return models.TaskApplication{}, fmt.Errorf("task not found")
	}

	// 2. Get the task owner from the database
	taskOwnerID, err := repository.GetTaskOwner(api, task.Task_id)
	if err != nil {
		return models.TaskApplication{}, fmt.Errorf("failed to get task owner: %w", err)
	}

	// 3. Prevent the task owner from applying to their own task
	if taskOwnerID == task.Employee_id {
		return models.TaskApplication{}, fmt.Errorf("you created this task")
	}

	// 4. Check if task has already been assigned
	taskAlreadyAssigned, err := repository.IsTaskAlreadyAssigned(api, task.Task_id)
	if err != nil {
		return models.TaskApplication{}, fmt.Errorf(
			"failed to check task assignment: %w",
			err,
		)
	}

	if taskAlreadyAssigned {
		return models.TaskApplication{}, fmt.Errorf("task has already been assigned")
	}

	// 5. Check if employee has already applied
	alreadyApplied, err := repository.HasEmployeeApplied(
		api,
		task.Task_id,
		task.Employee_id,
	)
	if err != nil {
		return models.TaskApplication{}, fmt.Errorf(
			"failed to check previous application: %w",
			err,
		)
	}

	if alreadyApplied {
		return models.TaskApplication{}, fmt.Errorf("you already applied to this task")
	}

	// 6. Set the actual task owner from the database
	task.Employer_id = taskOwnerID

	// 7. Application should start as Pending
	task.Status = "Ongoing"

	// 8. Create application
	createdApplication, err := repository.TaskApplication(api, task)
	if err != nil {
		return models.TaskApplication{}, fmt.Errorf(
			"failed to apply for task: %w",
			err,
		)
	}

	return createdApplication, nil
}

// func ApplyForTasks(api *config.ApiConfig, task models.TaskApplication) (models.TaskApplication, error) {
// 	checkTasksExist := repository.CheckTasksExist(api, task.Task_id)
// 	fmt.Println(checkTasksExist, task.Task_id)
// 	if !checkTasksExist {
// 		return models.TaskApplication{}, fmt.Errorf("Task Not Found %v", checkTasksExist)
// 	}

// 	checkTaskOwner, err := repository.IsTaskOwner(api, task.Task_id, task.Employer_id)
// 	if err != nil {
// 		return models.TaskApplication{}, err
// 	}

// 	if checkTaskOwner {
// 		return models.TaskApplication{}, fmt.Errorf("you created this task")
// 	}

// 	checkIsTaskAlreadyAssigned, err := repository.IsTaskAlreadyAssigned(api, task.Task_id)
// 	if err != nil {
// 		return models.TaskApplication{}, err
// 	}
// 	applicantAppliedAlready, err := repository.HasEmployeeApplied(api, task.Task_id, task.Employee_id)
// 	if err != nil {
// 		return models.TaskApplication{}, err
// 	}
// 	if applicantAppliedAlready {
// 		return models.TaskApplication{}, fmt.Errorf("you applied to this task already")
// 	}

// 	if checkIsTaskAlreadyAssigned {
// 		return models.TaskApplication{}, fmt.Errorf("task has already been assigned")
// 	}

// 	createdApplication, err := repository.TaskApplication(api, task)
// 	if err != nil {
// 		return models.TaskApplication{}, fmt.Errorf("failed to apply for task: %w", err)
// 	}

// 	return createdApplication, nil

// }

// func TasksFeeds(api *config.ApiConfig) ([]models.Taskfeed, error) {
// 	getTask, err := repository.TaskFeeds(api, )
// 	if err != nil {
// 		return []models.Taskfeed{}, fmt.Errorf("failed to apply for task: %w", err)
// 	}
// 	return getTask, nil
// }
