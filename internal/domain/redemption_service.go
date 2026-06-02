package domain

import (
	"context"

	"github.com/google/uuid"
)

type RedemptionService interface {
	RedeemVoucher(ctx context.Context, userID, voucherID uuid.UUID) (*Redemption, error)
}
