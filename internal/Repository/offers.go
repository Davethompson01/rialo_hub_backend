package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/Davethompson01/rialo_hub_backend/config"
)

func AcceptOffer(
	api *config.ApiConfig,
	offerID int,
	conversationID int,
) error {

	query := `
		UPDATE offers
		SET status = 'accepted',
		    updated_at = NOW()
		WHERE offer_id = $1
		  AND conversation_id = $2
		  AND status = 'pending'
	`

	ctx, cancel := context.WithTimeout(
		context.Background(),
		3*time.Second,
	)
	defer cancel()

	result, err := api.DB.ExecContext(
		ctx,
		query,
		offerID,
		conversationID,
	)

	if err != nil {
		return fmt.Errorf(
			"failed to accept offer: %w",
			err,
		)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return fmt.Errorf(
			"offer not found or already processed",
		)
	}

	return nil
}
func RejectOffer(
	api *config.ApiConfig,
	offerID int,
	conversationID int,
) error {

	query := `
		UPDATE offers
		SET status = 'rejected',
		    updated_at = NOW()
		WHERE offer_id = $1
		  AND conversation_id = $2
		  AND status = 'pending'
	`

	ctx, cancel := context.WithTimeout(
		context.Background(),
		3*time.Second,
	)
	defer cancel()

	result, err := api.DB.ExecContext(
		ctx,
		query,
		offerID,
		conversationID,
	)

	if err != nil {
		return fmt.Errorf(
			"failed to reject offer: %w",
			err,
		)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return fmt.Errorf(
			"offer not found or already processed",
		)
	}

	return nil
}
