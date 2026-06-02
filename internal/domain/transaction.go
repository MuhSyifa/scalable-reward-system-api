package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type TransactionType string

const (
	TxTypeEarn   TransactionType = "earn"
	TxTypeRedeem TransactionType = "redeem"
	TxTypeExpire TransactionType = "expire"
)

type PointTransaction struct {
	ID              uuid.UUID       `json:"id" db:"id"`
	UserID          uuid.UUID       `json:"user_id" db:"user_id"`
	Amount          int             `json:"amount" db:"amount"`
	TransactionType TransactionType `json:"transaction_type" db:"transaction_type"`
	ReferenceID     *string         `json:"reference_id" db:"reference_id"` // E.g., redemption ID or order ID
	Status          string          `json:"status" db:"status"`
	CreatedAt       time.Time       `json:"created_at" db:"created_at"`
}

type Redemption struct {
	ID         uuid.UUID `json:"id" db:"id"`
	UserID     uuid.UUID `json:"user_id" db:"user_id"`
	VoucherID  uuid.UUID `json:"voucher_id" db:"voucher_id"`
	PointsUsed int       `json:"points_used" db:"points_used"`
	Status     string    `json:"status" db:"status"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
}

type TransactionRepository interface {
	CreatePointTransaction(ctx context.Context, tx *PointTransaction) error
	CreateRedemption(ctx context.Context, r *Redemption) error
	ListUserTransactions(ctx context.Context, userID uuid.UUID, limit, offset int) ([]PointTransaction, error)
}
