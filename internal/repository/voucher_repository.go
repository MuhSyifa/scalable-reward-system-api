package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"scalable-reward-system/internal/domain"
)

type voucherRepository struct {
	db *sqlx.DB
}

func NewVoucherRepository(db *sqlx.DB) domain.VoucherRepository {
	return &voucherRepository{db: db}
}

func (r *voucherRepository) Create(ctx context.Context, voucher *domain.Voucher) error {
	query := `
		INSERT INTO vouchers (id, campaign_id, code, points_required, stock, total_redeemed, is_active, created_at, updated_at)
		VALUES (:id, :campaign_id, :code, :points_required, :stock, :total_redeemed, :is_active, :created_at, :updated_at)
	`
	_, err := r.db.NamedExecContext(ctx, query, voucher)
	return err
}

func (r *voucherRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Voucher, error) {
	var voucher domain.Voucher
	query := `SELECT * FROM vouchers WHERE id = $1`
	err := r.db.GetContext(ctx, &voucher, query, id)
	if err != nil {
		return nil, err
	}
	return &voucher, nil
}

func (r *voucherRepository) ListByCampaign(ctx context.Context, campaignID uuid.UUID) ([]domain.Voucher, error) {
	var vouchers []domain.Voucher
	query := `SELECT * FROM vouchers WHERE campaign_id = $1 AND is_active = true ORDER BY created_at DESC`
	err := r.db.SelectContext(ctx, &vouchers, query, campaignID)
	return vouchers, err
}

func (r *voucherRepository) DecrementStock(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE vouchers SET stock = stock - 1, total_redeemed = total_redeemed + 1, updated_at = NOW() WHERE id = $1 AND stock > 0`
	res, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return domain.ErrOutOfStock
	}
	return nil
}
