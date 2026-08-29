package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"time"

	"github.com/Davethompson01/rialo_hub_backend/config"
	"github.com/Davethompson01/rialo_hub_backend/internal/models"
)

func CreateTasks(api *config.ApiConfig, task models.Task) (models.Task, error) {
	query := `
		INSERT INTO tasks(
			user_id,
			title,
			description,
			reward,
			status,
			deadline
		)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING
			task_id,
			user_id,
			title,
			description,
			reward,
			status,
			deadline,
			created_at
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var createdTask models.Task

	err := api.DB.QueryRowContext(
		ctx,
		query,
		task.UserID,
		task.Title,
		task.Description,
		task.Reward,
		task.Status,
		task.Deadline,
	).Scan(
		&createdTask.ID,
		&createdTask.UserID,
		&createdTask.Title,
		&createdTask.Description,
		&createdTask.Reward,
		&createdTask.Status,
		&createdTask.Deadline,
		&createdTask.CreatedAt,
	)
	if err != nil {
		return models.Task{}, err
	}

	return createdTask, nil
}

func TaskApplication(api *config.ApiConfig, task models.TaskApplication) (models.TaskApplication, error) {
	query := `
        INSERT INTO task_application (
            task_id,
            employee_id,
            employer_id,
            skills,
            status
        )
        VALUES ($1, $2, $3, $4, $5)
        RETURNING
            task_id,
            employee_id,
            employer_id,
            skills,
            status
    `

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var application models.TaskApplication

	err := api.DB.QueryRowContext(
		ctx,
		query,
		task.Task_id,
		task.Employee_id,
		task.Employer_id,
		task.Skills,
		task.Status,
	).Scan(
		&application.Task_id,
		&application.Employee_id,
		&application.Employer_id,
		&application.Skills,
		&application.Status,
	)
	if err != nil {
		return models.TaskApplication{}, err
	}

	return application, nil
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
		WHERE task_application_id = $2
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
			AND status = 'Accepted'
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

func IsTaskOwner(api *config.ApiConfig, taskID, employerID int) (bool, error) {
	query := `
		SELECT EXISTS(
			SELECT 1
			FROM tasks
			WHERE task_id = $1
			AND user_id = $2
		)
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var exists bool

	err := api.DB.QueryRowContext(ctx, query, taskID, employerID).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func IsApplicationForTask(
	api *config.ApiConfig,
	applicationID int,
	taskID int,
) (bool, error) {

	query := `
		SELECT EXISTS (
			SELECT 1
			FROM task_application
			WHERE task_application_id = $1
			  AND task_id = $2
		)
	`

	var exists bool

	err := api.DB.QueryRow(
		query,
		applicationID,
		taskID,
	).Scan(&exists)

	if err != nil {
		return false, err
	}

	return exists, nil
}

func IsValidApplication(
	api *config.ApiConfig,
	applicationID, taskID, applicantID int,
) (bool, error) {

	query := `
		SELECT EXISTS (
			SELECT 1
			FROM task_application
			WHERE task_application_id = $1
			  AND task_id = $2
			  AND employee_id = $3
		)
	`

	ctx, cancel := context.WithTimeout(
		context.Background(),
		3*time.Second,
	)
	defer cancel()

	var exists bool

	err := api.DB.QueryRowContext(
		ctx,
		query,
		applicationID,
		taskID,
		applicantID,
	).Scan(&exists)

	if err != nil {
		return false, err
	}

	return exists, nil
}
func IsApplicationOwner(
	api *config.ApiConfig,
	applicationID, employeeID int,
) (bool, error) {

	query := `
		SELECT EXISTS (
			SELECT 1
			FROM task_application
			WHERE task_application_id = $1
			  AND employee_id = $2
		)
	`

	ctx, cancel := context.WithTimeout(
		context.Background(),
		3*time.Second,
	)
	defer cancel()

	var exists bool

	err := api.DB.QueryRowContext(
		ctx,
		query,
		applicationID,
		employeeID,
	).Scan(&exists)

	if err != nil {
		return false, err
	}

	return exists, nil
}

func HasEmployeeApplied(api *config.ApiConfig, taskID, employeeID int) (bool, error) {
	query := `
		SELECT EXISTS(
			SELECT 1
			FROM task_application
			WHERE task_id = $1
			  AND employee_id = $2
		)
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var exists bool

	err := api.DB.QueryRowContext(ctx, query, taskID, employeeID).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func DeleteApplication(api *config.ApiConfig, applicationID int) error {

	query := `
		DELETE FROM task_application
		WHERE task_application_id = $1
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
		WHERE task_id = $2
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
		WHERE task_id = $1
	`

	ctx, cancel := context.WithTimeout(
		context.Background(),
		3*time.Second,
	)
	defer cancel()

	result, err := api.DB.ExecContext(ctx, query, taskID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("task not found")
	}

	return nil
}

func DeleteTaskWithApplications(api *config.ApiConfig, taskID int) error {
	tx, err := api.DB.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer tx.Rollback()

	// Delete applications first because they reference the task.
	_, err = tx.Exec(`
		DELETE FROM task_application
		WHERE task_id = $1
	`, taskID)
	if err != nil {
		return fmt.Errorf("failed to delete task applications: %w", err)
	}

	// Then delete the task.
	result, err := tx.Exec(`
		DELETE FROM tasks
		WHERE task_id = $1
	`, taskID)
	if err != nil {
		return fmt.Errorf("failed to delete task: %w", err)
	}

	// Optional but useful: make sure a task was actually deleted.
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check deleted task: %w", err)
	}

	if rowsAffected == 0 {
		return errors.New("task not found")
	}

	// Everything succeeded.
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func GetTaskApplications(
	api *config.ApiConfig,
	taskID int,
) ([]models.ApplicationResponse, error) {

	query := `
		SELECT
			ta.task_application_id,
			ta.task_id,
			u.user_id,
			u.username,
			COALESCE(u.profile_pics, '') AS profile_pics,
			u.reputation,
			ta.status
		FROM task_application ta
		JOIN users u
			ON ta.employee_id = u.user_id
		WHERE ta.task_id = $1
		ORDER BY ta.applied_at ASC
	`

	ctx, cancel := context.WithTimeout(
		context.Background(),
		3*time.Second,
	)
	defer cancel()

	rows, err := api.DB.QueryContext(
		ctx,
		query,
		taskID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	applications := make([]models.ApplicationResponse, 0)

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
			ta.task_application_id,
			ta.task_id,
			u.user_id,
			u.username,
			COALESCE(u.profile_pics, '') AS profile_pics,
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

func CheckTasksExist(api *config.ApiConfig, tasks_id int) bool {
	var exists bool

	query := `SELECT EXISTS(
		SELECT 1 FROM tasks where task_id = $1
	)`
	err := api.DB.QueryRow(query, tasks_id).Scan(&exists)
	if err != nil {
		return false
	}

	return exists
}

func TaskFeeds(
	api *config.ApiConfig,
	userID int,
) ([]models.Taskfeed, error) {

	query := `
		SELECT
			tk.task_id,
			tk.user_id,
			tk.user_id AS employer_id,
			u.username,
			COALESCE(u.profile_pics, '') AS profile_pics,
			tk.title,
			tk.description,
			tk.reward,
			u.role,
			tk.status,
			tk.deadline,

			COUNT(ta.task_application_id) AS applicant_count,

			EXISTS (
				SELECT 1
				FROM task_application my_app
				WHERE my_app.task_id = tk.task_id
				  AND my_app.employee_id = $1
			) AS is_applied,

			(
				SELECT my_app.task_application_id
				FROM task_application my_app
				WHERE my_app.task_id = tk.task_id
				  AND my_app.employee_id = $1
				LIMIT 1
			) AS application_id

		FROM tasks tk

		JOIN users u
			ON tk.user_id = u.user_id

		LEFT JOIN task_application ta
			ON tk.task_id = ta.task_id

		GROUP BY
			tk.task_id,
			tk.user_id,
			u.username,
			u.profile_pics,
			tk.title,
			tk.description,
			tk.reward,
			u.role,
			tk.status,
			tk.deadline,
			tk.created_at

		ORDER BY tk.created_at DESC;
	`

	ctx, cancel := context.WithTimeout(
		context.Background(),
		3*time.Second,
	)
	defer cancel()

	rows, err := api.DB.QueryContext(
		ctx,
		query,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []models.Taskfeed

	for rows.Next() {

		var task models.Taskfeed

		err := rows.Scan(
			&task.ID,
			&task.UserID,
			&task.EmployerID,
			&task.Username,
			&task.ProfilePicture,
			&task.Title,
			&task.Description,
			&task.Reward,
			&task.Role,
			&task.Status,
			&task.Deadline,
			&task.ApplicantCount,
			&task.IsApplied,
			&task.ApplicationID,
		)

		if err != nil {
			return nil, err
		}

		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}

func GetTaskOwner(api *config.ApiConfig, taskID int) (int, error) {
	query := `
		SELECT user_id
		FROM tasks
		WHERE task_id = $1
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var employerID int

	err := api.DB.QueryRowContext(
		ctx,
		query,
		taskID,
	).Scan(&employerID)

	if err != nil {
		return 0, err
	}

	return employerID, nil
}

func GetTaskEmployer(
	api *config.ApiConfig,
	taskID int,
) (int, error) {

	query := `
		SELECT user_id
		FROM tasks
		WHERE task_id = $1
	`

	var employerID int

	err := api.DB.QueryRow(
		query,
		taskID,
	).Scan(&employerID)

	if err != nil {
		if err == sql.ErrNoRows {
			return 0, fmt.Errorf("task not found")
		}

		return 0, err
	}

	return employerID, nil
}
func IsTaskEmployer(
	api *config.ApiConfig,
	taskID int,
	employerID int,
) (bool, error) {

	query := `
		SELECT EXISTS (
			SELECT 1
			FROM tasks
			WHERE task_id = $1
			  AND user_id = $2
		)
	`

	var exists bool

	err := api.DB.QueryRow(
		query,
		taskID,
		employerID,
	).Scan(&exists)

	if err != nil {
		return false, err
	}

	return exists, nil
}

func GetConversation(
	api *config.ApiConfig,
	taskID int,
	conversationID int,
) (models.Conversation, error) {

	query := `
		SELECT
			conversation_id,
			task_id,
			employer_id,
			applicant_id,
			created_at
		FROM conversations
		WHERE conversation_id = $1
		  AND task_id = $2
	`

	var conversation models.Conversation

	err := api.DB.QueryRow(
		query,
		conversationID,
		taskID,
	).Scan(
		&conversation.ConversationID,
		&conversation.TaskID,
		&conversation.EmployerID,
		&conversation.ApplicantID,
		&conversation.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return models.Conversation{}, fmt.Errorf(
				"conversation not found",
			)
		}

		return models.Conversation{}, err
	}

	return conversation, nil
}

func GetConversationMessages(
	api *config.ApiConfig,
	conversationID int,
) ([]models.SendMessage, error) {

	query := `
		SELECT
			message_id,
			conversation_id,
			sender_id,
			content,
			created_at
		FROM messages
		WHERE conversation_id = $1
		ORDER BY created_at ASC
	`

	rows, err := api.DB.Query(
		query,
		conversationID,
	)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var messages []models.SendMessage

	for rows.Next() {

		var message models.SendMessage

		err := rows.Scan(
			&message.MessageID,
			&message.ConversationID,
			&message.SenderID,
			&message.Content,
			&message.CreatedAt,
		)

		if err != nil {
			return nil, err
		}

		messages = append(messages, message)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return messages, nil
}
