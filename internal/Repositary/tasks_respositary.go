package repositary

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
			deadline,
			created_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
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
		task.CreatedAt,
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




