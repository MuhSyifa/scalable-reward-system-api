package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"scalable-reward-system/internal/domain"
)

type campaignRepository struct {
	db *sqlx.DB
}

func NewCampaignRepository(db *sqlx.DB) domain.CampaignRepository {
	return &campaignRepository{db: db}
}

func (r *campaignRepository) Create(ctx context.Context, campaign *domain.Campaign) error {
	query := `
		INSERT INTO campaigns (id, name, description, start_date, end_date, is_active, created_at, updated_at)
		VALUES (:id, :name, :description, :start_date, :end_date, :is_active, :created_at, :updated_at)
	`
	_, err := r.db.NamedExecContext(ctx, query, campaign)
	return err
}

func (r *campaignRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Campaign, error) {
	var campaign domain.Campaign
	query := `SELECT * FROM campaigns WHERE id = $1`
	err := r.db.GetContext(ctx, &campaign, query, id)
	if err != nil {
		return nil, err
	}
	return &campaign, nil
}

func (r *campaignRepository) ListActive(ctx context.Context, limit, offset int) ([]domain.Campaign, error) {
	var campaigns []domain.Campaign
	query := `
		SELECT * FROM campaigns 
		WHERE is_active = true AND start_date <= NOW() AND end_date >= NOW()
		ORDER BY start_date DESC
		LIMIT $1 OFFSET $2
	`
	err := r.db.SelectContext(ctx, &campaigns, query, limit, offset)
	return campaigns, err
}

func (r *campaignRepository) Update(ctx context.Context, campaign *domain.Campaign) error {
	query := `
		UPDATE campaigns 
		SET name = :name, description = :description, start_date = :start_date, 
		    end_date = :end_date, is_active = :is_active, updated_at = :updated_at
		WHERE id = :id
	`
	_, err := r.db.NamedExecContext(ctx, query, campaign)
	return err
}
