package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/Davethompson01/rialo_hub_backend/config"
	"github.com/Davethompson01/rialo_hub_backend/internal/models"
)

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

	ctx, cancel := context.WithTimeout(
		context.Background(),
		3*time.Second,
	)
	defer cancel()

	rows, err := api.DB.QueryContext(
		ctx,
		query,
		conversationID,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	messages := make([]models.SendMessage, 0)

	for rows.Next() {

		var message models.SendMessage

		if err := rows.Scan(
			&message.MessageID,
			&message.ConversationID,
			&message.SenderID,
			&message.Content,
			&message.CreatedAt,
		); err != nil {
			return nil, err
		}

		messages = append(messages, message)
	}

	return messages, rows.Err()
}

func GetMessages(
	api *config.ApiConfig,
	conversationID int,
) ([]models.MessageResponse, error) {

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

	ctx, cancel := context.WithTimeout(
		context.Background(),
		3*time.Second,
	)
	defer cancel()

	rows, err := api.DB.QueryContext(
		ctx,
		query,
		conversationID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get messages: %w", err)
	}

	defer rows.Close()

	messages := make([]models.MessageResponse, 0)

	for rows.Next() {

		var message models.MessageResponse

		if err := rows.Scan(
			&message.MessageID,
			&message.ConversationID,
			&message.SenderID,
			&message.Content,
			&message.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf(
				"failed to scan message: %w",
				err,
			)
		}

		messages = append(messages, message)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf(
			"failed while reading messages: %w",
			err,
		)
	}

	return messages, nil
}
