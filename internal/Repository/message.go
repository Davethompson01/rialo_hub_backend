package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/Davethompson01/rialo_hub_backend/config"
	"github.com/Davethompson01/rialo_hub_backend/internal/models"
)

func SendMessage(
	api *config.ApiConfig,
	message models.SendMessage,
) (models.MessageResponse, error) {

	ctx, cancel := context.WithTimeout(
		context.Background(),
		3*time.Second,
	)
	defer cancel()

	query := `
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
	`

	var savedMessage models.MessageResponse
	fmt.Println("ConversationID:", message.ConversationID)
	fmt.Println("SenderID:", message.SenderID)
	fmt.Println("Content:", message.Content)
	fmt.Println("CreatedAt:", message.CreatedAt)

	err := api.DB.QueryRowContext(
		ctx,
		query,
		message.ConversationID,
		message.SenderID,
		message.Content,
		message.CreatedAt,
	).Scan(
		&savedMessage.MessageID,
		&savedMessage.ConversationID,
		&savedMessage.SenderID,
		&savedMessage.Content,
		&savedMessage.CreatedAt,
	)

	if err != nil {
		return models.MessageResponse{}, fmt.Errorf(
			"failed to send message: %w",
			err,
		)
	}

	// Determine the other participant.
	var receiverID int

	err = api.DB.QueryRowContext(
		ctx,
		`
		SELECT
    CASE
        WHEN employer_id = $1 THEN applicant_id
        WHEN applicant_id = $1 THEN employer_id
        ELSE NULL
    END
FROM conversations
WHERE conversation_id = $2
		`,
		message.SenderID,
		message.ConversationID,
	).Scan(&receiverID)
	fmt.Println(message.ConversationID)
	if err != nil {
		return models.MessageResponse{}, fmt.Errorf(
			"failed to find message receiver: %w",
			err,
		)
	}

	event := models.WebSocketEvent{
		Type: "new_message",
		Data: savedMessage,
	}

	if err := api.Hub.SendToUser(
		receiverID,
		event,
	); err != nil {
		fmt.Printf(
			"failed to send websocket message: %v\n",
			err,
		)
	}

	return savedMessage, nil
}
