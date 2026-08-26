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
) (models.NegotiationResponse, error) {

	tx, err := api.DB.Begin()
	if err != nil {
		return models.NegotiationResponse{}, fmt.Errorf(
			"failed to begin transaction: %w",
			err,
		)
	}

	defer tx.Rollback()

	// ============================================================
	// 1. GET OR CREATE CONVERSATION
	// ============================================================

	var conversation models.SendMessage

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
		&conversation.ConversationID,
		&conversation.TaskId,
		&conversation.EmployerID,
		&conversation.ApplicantID,
		&conversation.CreatedAt,
	)

	// ------------------------------------------------------------
	// Conversation does not exist
	// ------------------------------------------------------------

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
			DO UPDATE SET
				task_id = EXCLUDED.task_id
			RETURNING
				conversation_id,
				task_id,
				employer_id,
				applicant_id,
				created_at
		`

		err = tx.QueryRow(
			query,
			negotiate.TaskId,
			negotiate.EmployerID,
			negotiate.ApplicantID,
			negotiate.CreatedAt,
		).Scan(
			&conversation.ConversationID,
			&conversation.TaskId,
			&conversation.EmployerID,
			&conversation.ApplicantID,
			&conversation.CreatedAt,
		)

		if err != nil {
			return models.NegotiationResponse{}, fmt.Errorf(
				"failed to create conversation: %w",
				err,
			)
		}

	} else if err != nil {

		return models.NegotiationResponse{}, fmt.Errorf(
			"failed to get conversation: %w",
			err,
		)
	}

	// ============================================================
	// 2. INSERT MESSAGE
	// ============================================================

	var savedMessage models.SendMessage

	err = tx.QueryRow(`
		INSERT INTO messages (
			conversation_id,
			sender_id,
			content,
			created_at
		)
		VALUES ($1, $2, $3, $4)
		RETURNING
			message_id,
			conversation_id,
			sender_id,
			content,
			created_at
	`,
		conversation.ConversationID,
		negotiate.CreatedBy,
		negotiate.Content,
		negotiate.CreatedAt,
	).Scan(
		&savedMessage.MessageID,
		&savedMessage.ConversationID,
		&savedMessage.SenderID,
		&savedMessage.Content,
		&savedMessage.CreatedAt,
	)

	if err != nil {
		return models.NegotiationResponse{}, fmt.Errorf(
			"failed to insert message: %w",
			err,
		)
	}

	// ============================================================
	// 3. CREATE OFFER
	// ============================================================

	var offerID int

	err = tx.QueryRow(`
		INSERT INTO offers (
			conversation_id,
			task_id,
			employer_id,
			applicant_id,
			amount,
			created_by,
			status,
			created_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING offer_id
	`,
		conversation.ConversationID,
		conversation.TaskId,
		conversation.EmployerID,
		conversation.ApplicantID,
		negotiate.Offer.NewOffer,
		negotiate.CreatedBy,
		negotiate.Status,
		negotiate.CreatedAt,
	).Scan(&offerID)

	if err != nil {
		return models.NegotiationResponse{}, fmt.Errorf(
			"failed to insert offer: %w",
			err,
		)
	}

	// ============================================================
	// DEBUG
	// ============================================================

	fmt.Printf(
		"DEBUG negotiation: conversation=%d task=%d employer=%d applicant=%d offer=%d message=%d\n",
		conversation.ConversationID,
		conversation.TaskId,
		conversation.EmployerID,
		conversation.ApplicantID,
		offerID,
		savedMessage.MessageID,
	)

	// ============================================================
	// 4. UPDATE NOTIFICATION COUNT
	// ============================================================

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
		return models.NegotiationResponse{}, fmt.Errorf(
			"failed to update notification count: %w",
			err,
		)
	}

	// ============================================================
	// 5. COMMIT TRANSACTION
	// ============================================================

	if err := tx.Commit(); err != nil {
		return models.NegotiationResponse{}, fmt.Errorf(
			"failed to commit negotiation: %w",
			err,
		)
	}

	// ============================================================
	// 6. SEND OFFER WEBSOCKET EVENT
	// ============================================================

	offerEvent := models.WebSocketEvent{
		Type: "new_offer",
		Data: models.OfferNotification{
			TaskId:         negotiate.TaskId,
			ConversationID: conversation.ConversationID,
			UserID:         negotiate.EmployerID,
			EmployeeID:     negotiate.ApplicantID,
			NewOffer:       negotiate.Offer.NewOffer,
			CreatedAt:      negotiate.CreatedAt,
		},
	}

	if err := api.Hub.SendToUser(
		negotiate.EmployerID,
		offerEvent,
	); err != nil {
		fmt.Printf(
			"failed to send websocket offer notification: %v\n",
			err,
		)
	}

	// ============================================================
	// 7. SEND MESSAGE WEBSOCKET EVENT
	// ============================================================

	messageEvent := models.WebSocketEvent{
		Type: "new_message",
		Data: savedMessage,
	}

	if err := api.Hub.SendToUser(
		negotiate.EmployerID,
		messageEvent,
	); err != nil {
		fmt.Printf(
			"failed to send websocket message: %v\n",
			err,
		)
	}

	// ============================================================
	// 8. RETURN RESPONSE
	// ============================================================

	return models.NegotiationResponse{
		TaskID:         conversation.TaskId,
		ConversationID: conversation.ConversationID,
		EmployerID:     conversation.EmployerID,
		ApplicantID:    conversation.ApplicantID,
		OfferID:        offerID,
		Amount:         negotiate.Offer.NewOffer,
		Status:         negotiate.Status,
		CreatedAt:      negotiate.CreatedAt,
	}, nil
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
