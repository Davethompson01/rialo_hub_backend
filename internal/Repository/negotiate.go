package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/Davethompson01/rialo_hub_backend/config"
	"github.com/Davethompson01/rialo_hub_backend/internal/models"
)

func CreateNegotiation(
	api *config.ApiConfig,
	negotiate models.SendMessage,
) (models.SendMessage, error) {

	tx, err := api.DB.Begin()
	if err != nil {
		return models.SendMessage{}, fmt.Errorf(
			"failed to begin transaction: %w",
			err,
		)
	}

	defer tx.Rollback()

	var message models.SendMessage

	// ------------------------------------------------
	// 1. Get existing conversation or create one
	// ------------------------------------------------

	query := `
		SELECT
			conversation_id,
			task_id,
			employer_id,
			applicant_id,
			created_at
		FROM conversations
		WHERE task_id = $1
		  AND employer_id = $2
		  AND applicant_id = $3
	`

	err = tx.QueryRow(
		query,
		negotiate.TaskId,
		negotiate.EmployerID,
		negotiate.ApplicantID,
	).Scan(
		&message.ConversationID,
		&message.TaskId,
		&message.EmployerID,
		&message.ApplicantID,
		&message.CreatedAt,
	)

	// Conversation doesn't exist
	if err == sql.ErrNoRows {

		query = `
			INSERT INTO conversations (
    task_id,
    employer_id,
    applicant_id,
    created_at
)
VALUES ($1, $2, $3, $4)
ON CONFLICT (task_id, employer_id, applicant_id)
DO UPDATE SET task_id = EXCLUDED.task_id
RETURNING conversation_id, task_id, employer_id, applicant_id, created_at;
		`

		err = tx.QueryRow(
			query,
			negotiate.TaskId,
			negotiate.EmployerID,
			negotiate.ApplicantID,
			negotiate.CreatedAt,
		).Scan(
			&message.ConversationID,
			&message.TaskId,
			&message.EmployerID,
			&message.ApplicantID,
			&message.CreatedAt,
		)

		if err != nil {
			return models.SendMessage{}, fmt.Errorf(
				"failed to create/get conversation: %w",
				err,
			)
		}

	} else if err != nil {

		return models.SendMessage{}, fmt.Errorf(
			"failed to get conversation: %w",
			err,
		)
	}

	// ------------------------------------------------
	// 2. Create offer
	// ------------------------------------------------

	_, err = tx.Exec(`
		INSERT INTO offers (
			task_id,
			employer_id,
			applicant_id,
			new_offer,
			created_by
			status,
			created_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`,
		negotiate.TaskId,
		negotiate.EmployerID,
		negotiate.ApplicantID,
		negotiate.Offer.NewOffer,
		negotiate.Status,
		negotiate.CreatedAt,
	)

	if err != nil {
		return models.SendMessage{}, fmt.Errorf(
			"failed to insert offer: %w",
			err,
		)
	}

	// ------------------------------------------------
	// 3. Update notification count
	// ------------------------------------------------

	_, err = tx.Exec(`
		INSERT INTO notifications (
			task_id,
			user_id,
			notifications_count
		)
		VALUES ($1, $2, 1)
		ON CONFLICT (task_id, user_id)
		DO UPDATE SET
			notifications_count =
				notifications.notifications_count + 1
	`,
		negotiate.TaskId,
		negotiate.EmployerID,
	)

	if err != nil {
		return models.SendMessage{}, fmt.Errorf(
			"failed to update notification count: %w",
			err,
		)
	}

	// ------------------------------------------------
	// 4. Commit transaction
	// ------------------------------------------------

	if err := tx.Commit(); err != nil {
		return models.SendMessage{}, fmt.Errorf(
			"failed to commit negotiation: %w",
			err,
		)
	}

	// ------------------------------------------------
	// 5. Send WebSocket event
	// ------------------------------------------------

	event := models.WebSocketEvent{
		Type: "new_offer",
		Data: models.OfferNotification{
			TaskId:         negotiate.TaskId,
			ConversationID: message.ConversationID,
			UserID:         negotiate.EmployerID,
			EmployeeID:     negotiate.ApplicantID,
			NewOffer:       negotiate.Offer.NewOffer,
			CreatedAt:      negotiate.CreatedAt,
		},
	}

	if err := api.Hub.SendToUser(
		negotiate.EmployerID,
		event,
	); err != nil {
		fmt.Printf(
			"failed to send websocket notification: %v\n",
			err,
		)
	}

	return message, nil
}

func GetAllApplicantOffer(
	api *config.ApiConfig,
	userID int,
) ([]models.OfferResponse, error) {

	query := `
		SELECT
			o.offer_id,
			o.task_id,
			o.conversation_id,
			o.employer_id,
			o.applicant_id,
			o.created_by,
			o.amount,
			o.status,
			u.username,
			COALESCE(u.profile_pics, '') AS avatar,
			o.created_at
		FROM offers o
		JOIN users u
			ON o.applicant_id = u.user_id
		WHERE o.employer_id = $1
		ORDER BY o.created_at DESC
	`

	ctx, cancel := context.WithTimeout(
		context.Background(),
		3*time.Second,
	)
	defer cancel()

	rows, err := api.DB.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to get applicant offers: %w",
			err,
		)
	}

	defer rows.Close()

	offers := make([]models.OfferResponse, 0)

	for rows.Next() {

		var offer models.OfferResponse

		err := rows.Scan(
			&offer.OfferID,
			&offer.TaskID,
			&offer.ConversationID,
			&offer.EmployerID,
			&offer.ApplicantID,
			&offer.CreatedBy,
			&offer.Amount,
			&offer.Status,
			&offer.Username,
			&offer.Avatar,
			&offer.CreatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf(
				"failed to scan applicant offer: %w",
				err,
			)
		}

		offers = append(offers, offer)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return offers, nil
}

func GetApplicationOffers(
	api *config.ApiConfig,
	userID int,
) ([]models.OfferResponse, error) {

	query := `
		SELECT
			o.offer_id,
			o.task_id,
			o.conversation_id,
			o.employer_id,
			o.applicant_id,
			o.created_by,
			o.amount,
			o.status,
			u.username,
			COALESCE(u.profile_pics, '') AS avatar,
			o.created_at
		FROM offers o
		JOIN users u
			ON o.employer_id = u.user_id
		WHERE o.applicant_id = $1
		ORDER BY o.created_at DESC
	`

	ctx, cancel := context.WithTimeout(
		context.Background(),
		3*time.Second,
	)
	defer cancel()

	rows, err := api.DB.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to get application offers: %w",
			err,
		)
	}

	defer rows.Close()

	offers := make([]models.OfferResponse, 0)

	for rows.Next() {

		var offer models.OfferResponse

		err := rows.Scan(
			&offer.OfferID,
			&offer.TaskID,
			&offer.ConversationID,
			&offer.EmployerID,
			&offer.ApplicantID,
			&offer.CreatedBy,
			&offer.Amount,
			&offer.Status,
			&offer.Username,
			&offer.Avatar,
			&offer.CreatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf(
				"failed to scan application offer: %w",
				err,
			)
		}

		offers = append(offers, offer)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return offers, nil
}

func OfferOwner(api *config.ApiConfig, userID int) (bool, error) {

	query :=
		`SELECT EXISTS (
		SELECT 1 FROM offers WHERE offer_id = $1 AND employer_id = $2
		)`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var exists bool

	err := api.DB.QueryRowContext(ctx, query, userID).Scan(&exists)

	if err != nil {
		return false, err
	}
	return exists, nil

}
