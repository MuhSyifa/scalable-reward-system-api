package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"scalable-reward-system/internal/domain"
)

type transactionRepository struct {
	db *sqlx.DB
}

func NewTransactionRepository(db *sqlx.DB) domain.TransactionRepository {
	return &transactionRepository{db: db}
}

func (r *transactionRepository) CreatePointTransaction(ctx context.Context, tx *domain.PointTransaction) error {
	query := `
		INSERT INTO point_transactions (id, user_id, amount, transaction_type, reference_id, status, created_at)
		VALUES (:id, :user_id, :amount, :transaction_type, :reference_id, :status, :created_at)
	`
	_, err := r.db.NamedExecContext(ctx, query, tx)
	return err
}

func (r *transactionRepository) CreateRedemption(ctx context.Context, red *domain.Redemption) error {
	query := `
		INSERT INTO redemptions (id, user_id, voucher_id, points_used, status, created_at)
		VALUES (:id, :user_id, :voucher_id, :points_used, :status, :created_at)
	`
	_, err := r.db.NamedExecContext(ctx, query, red)
	return err
}

func (r *transactionRepository) ListUserTransactions(ctx context.Context, userID uuid.UUID, limit, offset int) ([]domain.PointTransaction, error) {
	var txs []domain.PointTransaction
	query := `SELECT * FROM point_transactions WHERE user_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3`
	err := r.db.SelectContext(ctx, &txs, query, userID, limit, offset)
	return txs, err
}

// Transactional execute function for complex operations
func (r *transactionRepository) ExecTx(ctx context.Context, fn func(tx *sqlx.Tx) error) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}

	err = fn(tx)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return rbErr
		}
		return err
	}
	return tx.Commit()
}
