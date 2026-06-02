package worker

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

type PointExpirationWorker struct {
	db *sqlx.DB
}

func NewPointExpirationWorker(db *sqlx.DB) *PointExpirationWorker {
	return &PointExpirationWorker{db: db}
}

// Start runs the worker in the background
func (w *PointExpirationWorker) Start(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	go func() {
		for {
			select {
			case <-ticker.C:
				w.processExpiredPoints(ctx)
			case <-ctx.Done():
				ticker.Stop()
				log.Info().Msg("Point expiration worker stopped")
				return
			}
		}
	}()
}

// processExpiredPoints runs the logic to deduct expired points
func (w *PointExpirationWorker) processExpiredPoints(ctx context.Context) {
	// In a real system, you'd find unredeemed points older than X months
	// For this portfolio project, we simulate the query
	
	query := `
		-- This is a simulated query for portfolio purposes
		-- In reality, we would track expiration per point_transaction (e.g., FIFO method)
		UPDATE users 
		SET total_points = 0 
		WHERE id IN (
			SELECT user_id FROM point_transactions 
			WHERE created_at < NOW() - INTERVAL '1 year'
			GROUP BY user_id
		)
	`
	
	res, err := w.db.ExecContext(ctx, query)
	if err != nil {
		log.Error().Err(err).Msg("Failed to process expired points")
		return
	}
	
	affected, _ := res.RowsAffected()
	if affected > 0 {
		log.Info().Msgf("Processed expired points for %d users", affected)
	}
}
