package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Voucher struct {
	ID             uuid.UUID `json:"id" db:"id"`
	CampaignID     uuid.UUID `json:"campaign_id" db:"campaign_id"`
	Code           string    `json:"code" db:"code"`
	PointsRequired int       `json:"points_required" db:"points_required"`
	Stock          int       `json:"stock" db:"stock"`
	TotalRedeemed  int       `json:"total_redeemed" db:"total_redeemed"`
	IsActive       bool      `json:"is_active" db:"is_active"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`
}

type VoucherRepository interface {
	Create(ctx context.Context, voucher *Voucher) error
	GetByID(ctx context.Context, id uuid.UUID) (*Voucher, error)
	ListByCampaign(ctx context.Context, campaignID uuid.UUID) ([]Voucher, error)
	DecrementStock(ctx context.Context, id uuid.UUID) error
}
