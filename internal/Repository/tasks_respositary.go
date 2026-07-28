package repository

import (
	"context"
	"time"

	"github.com/Davethompson01/rialo_hub_backend/config"
	"github.com/Davethompson01/rialo_hub_backend/internal/models"
)

func CreateTasks(api *config.ApiConfig, task models.Task) error {
	query := `
		INSERT INTO tasks(
			user_id,
			title,
			description,
			reward,
			status,
			deadline
		)
		VALUES ($1, $2, $3, $4, $5, $6) retur
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := api.DB.ExecContext(
		ctx,
		query,
		task.UserID,
		task.Title,
		task.Description,
		task.Reward,
		task.Status,
		task.Deadline,
	)

	return err
}

func TaskApplication(api *config.ApiConfig, application models.TaskApplication) error {
	query := `
		INSERT INTO task_application(
			task_id,
			employee_id,
			employer_id,
			skills,
			status,
			applied_at
		)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := api.DB.ExecContext(
		ctx,
		query,
		application.Task_id,
		application.Employee_id,
		application.Employer_id,
		application.Skills,
		application.Status,
		application.AppliedAt,
	)

	return err
}

// func ApplicationResponse(api *config.ApiConfig, response models.ApplicationResponse) error {
// 	query := `
// 		UPDATE task_application
// 		SET
// 			status = $1,
// 			responded_at = $2
// 		WHERE application_id = $3
// 	`

// 	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
// 	defer cancel()

// 	_, err := api.DB.ExecContext(
// 		ctx,
// 		query,
// 		response.Status,
// 		response.RespondedAt,
// 		response.ApplicationID,
// 	)
// 	return err
// }

func UpdateApplicationStatus(api *config.ApiConfig, applicationID int, status string) error {
	query := `
		UPDATE task_application
		SET status = $1
		WHERE id = $2
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := api.DB.ExecContext(ctx, query, status, applicationID)
	return err
}

func IsTaskAlreadyAssigned(api *config.ApiConfig, taskID int) (bool, error) {
	query := `
		SELECT EXISTS (
			SELECT 1
			FROM task_application
			WHERE task_id = $1
			AND status = 'accepted'
		)
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var exists bool

	err := api.DB.QueryRowContext(ctx, query, taskID).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func IsTaskOwner(api *config.ApiConfig, taskID, userID int) (bool, error) {

	query := `
		SELECT EXISTS(
			SELECT 1
			FROM tasks
			WHERE id = $1
			AND user_id = $2
		)
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var exists bool

	err := api.DB.QueryRowContext(ctx, query, taskID, userID).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func IsApplicationOwner(api *config.ApiConfig, applicationID, employeeID int) (bool, error) {

	query := `
		SELECT EXISTS(
			SELECT 1
			FROM task_application
			WHERE id = $1
			AND employee_id = $2
		)
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var exists bool

	err := api.DB.QueryRowContext(ctx, query, applicationID, employeeID).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func DeleteApplication(api *config.ApiConfig, applicationID int) error {

	query := `
		DELETE FROM task_application
		WHERE id = $1
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := api.DB.ExecContext(ctx, query, applicationID)
	return err
}

func UpdateTaskStatus(api *config.ApiConfig, taskID int, status string) error {

	query := `
		UPDATE tasks
		SET status = $1
		WHERE id = $2
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := api.DB.ExecContext(ctx, query, status, taskID)
	return err
}

func RejectPendingApplications(api *config.ApiConfig, taskID int) error {

	query := `
		UPDATE task_application
		SET status = 'rejected'
		WHERE task_id = $1
		AND status = 'pending'
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := api.DB.ExecContext(ctx, query, taskID)
	return err
}

func DeleteTask(api *config.ApiConfig, taskID int) error {

	query := `
		DELETE FROM tasks
		WHERE id = $1
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := api.DB.ExecContext(ctx, query, taskID)
	return err
}

func GetTaskApplications(api *config.ApiConfig, taskID int) ([]models.ApplicationResponse, error) {

	query := `
		SELECT
			ta.application_id,
			ta.task_id,
			u.user_id,
			u.username,
			u.avatar,
			u.reputation,
			ta.status
		FROM task_application ta
		JOIN users u
			ON ta.employee_id = u.user_id
		WHERE ta.task_id = $1
		ORDER BY ta.applied_at ASC
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := api.DB.QueryContext(ctx, query, taskID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var applications []models.ApplicationResponse

	for rows.Next() {
		var application models.ApplicationResponse

		err := rows.Scan(
			&application.ID,
			&application.TaskID,
			&application.ApplicantID,
			&application.Username,
			&application.Avatar,
			&application.Reputation,
			&application.Status,
		)
		if err != nil {
			return nil, err
		}

		applications = append(applications, application)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return applications, nil
}

func GetMyApplications(api *config.ApiConfig, employeeID int) ([]models.ApplicationResponse, error) {

	query := `
		SELECT
			ta.application_id,
			ta.task_id,
			u.user_id,
			u.username,
			u.avatar,
			u.reputation,
			ta.status
		FROM task_application ta
		JOIN users u
			ON ta.employee_id = u.user_id
		WHERE ta.employee_id = $1
		ORDER BY ta.applied_at DESC
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := api.DB.QueryContext(ctx, query, employeeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var applications []models.ApplicationResponse

	for rows.Next() {
		var application models.ApplicationResponse

		if err := rows.Scan(
			&application.ID,
			&application.TaskID,
			&application.ApplicantID,
			&application.Username,
			&application.Avatar,
			&application.Reputation,
			&application.Status,
		); err != nil {
			return nil, err
		}

		applications = append(applications, application)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return applications, nil
}
