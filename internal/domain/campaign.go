package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Campaign struct {
	ID          uuid.UUID `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
	StartDate   time.Time `json:"start_date" db:"start_date"`
	EndDate     time.Time `json:"end_date" db:"end_date"`
	IsActive    bool      `json:"is_active" db:"is_active"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

type CampaignRepository interface {
	Create(ctx context.Context, campaign *Campaign) error
	GetByID(ctx context.Context, id uuid.UUID) (*Campaign, error)
	ListActive(ctx context.Context, limit, offset int) ([]Campaign, error)
	Update(ctx context.Context, campaign *Campaign) error
}
