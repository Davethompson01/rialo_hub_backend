package repository

import (
	"context"
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
		return models.SendMessage{}, fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer tx.Rollback()

	// 1. Create negotiation/message
	query := `
		INSERT INTO conversations (
			task_id,
			created_at
		)
		VALUES ($1, $2)
	`

	var message models.SendMessage

	err = tx.QueryRow(
		query,
		negotiate.TaskId,
		negotiate.CreatedAt,
	).Scan(
		&message.TaskId,
	)

	if err != nil {
		return models.SendMessage{}, fmt.Errorf(
			"failed to create Conversation: %w",
			err,
		)
	}

	// 2. Create offer
	_, err = tx.Exec(`
		INSERT INTO offer (
			task_id,
			employer_id,
			applicant_id,
			new_offer,
			status,
			created_at
		)
		VALUES ($1, $2, $3, $4, $5, $6)
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

	// 3. Update/create notification
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

	// 4. Commit transaction
	if err := tx.Commit(); err != nil {
		return models.SendMessage{}, fmt.Errorf(
			"failed to commit negotiation: %w",
			err,
		)
	}

	// 5. Send real-time WebSocket notification
	event := models.WebSocketEvent{
		Type: "new_offer",
		Data: models.OfferNotification{
			TaskId:     negotiate.TaskId,
			UserID:     negotiate.EmployerID,
			EmployeeID: negotiate.ApplicantID,
			NewOffer:   negotiate.Offer.NewOffer,
		},
	}

	if err := api.Hub.SendToUser(
		negotiate.EmployerID,
		event,
	); err != nil {
		// Don't fail the negotiation because the user is offline
		fmt.Printf("failed to send websocket notification: %v\n", err)
	}

	return message, nil
}

func GetAllApplicantOffer(api *config.ApiConfig, UserID int) ([]models.OfferResponse, error) {
	query := `SELECT 
					ne.task_id,
					ne.employer_id,
					ne.applicant_id,
					ne.message,
					ne.new_offer,
					u.username,
			COALESCE(u.profile_pics, '') AS profile_pics,
			FROM negotiations ne JOIN users 
			ON ne.employer_id = u user_id
			WHERE ne.employer_id = $1 
	`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := api.DB.QueryContext(ctx, query, UserID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var offers []models.OfferResponse

	for rows.Next() {
		var offer models.OfferResponse
		if err := rows.Scan(
			&offer.OfferID,
			&offer.EmployerID,
			&offer.ApplicantId,
			&offer.Avatar,
			&offer.Status,
			&offer.TaskID,
			&offer.Username,
		); err != nil {
			return nil, err
		}

		offers = append(offers, offer)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return offers, nil
}

func GetApplicationOffers(api *config.ApiConfig, UserID int) ([]models.OfferResponse, error) {

	query := `SELECT 
					ne.task_id,
					ne.employer_id,
					ne.applicant_id,
					ne.message,
					ne.new_offer,
					u.username,
			COALESCE(u.profile_pics, '') AS profile_pics,
			FROM negotiations ne JOIN users 
			ON ne.employee_id = u user_id
			WHERE ne.employee_id = $1 
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := api.DB.QueryContext(ctx, query, UserID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var offers []models.OfferResponse

	for rows.Next() {
		var offer models.OfferResponse
		if err := rows.Scan(
			&offer.OfferID,
			&offer.EmployerID,
			&offer.ApplicantId,
			&offer.Avatar,
			&offer.Status,
			&offer.TaskID,
			&offer.Username,
		); err != nil {
			return nil, err
		}

		offers = append(offers, offer)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return offers, nil
}
