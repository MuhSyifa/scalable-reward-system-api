package service

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"

	"scalable-reward-system/internal/domain"
)

type redemptionService struct {
	db          *sqlx.DB
	redis       *redis.Client
	userRepo    domain.UserRepository
	voucherRepo domain.VoucherRepository
	txRepo      domain.TransactionRepository
}

func NewRedemptionService(db *sqlx.DB, redisClient *redis.Client, userRepo domain.UserRepository, voucherRepo domain.VoucherRepository, txRepo domain.TransactionRepository) domain.RedemptionService {
	return &redemptionService{
		db:          db,
		redis:       redisClient,
		userRepo:    userRepo,
		voucherRepo: voucherRepo,
		txRepo:      txRepo,
	}
}

func (s *redemptionService) RedeemVoucher(ctx context.Context, userID, voucherID uuid.UUID) (*domain.Redemption, error) {
	// 1. Distributed Lock to prevent concurrent double-redemption by the same user
	lockKey := "lock:redeem:" + userID.String()
	acquired, err := s.redis.SetNX(ctx, lockKey, 1, 5*time.Second).Result()
	if err != nil || !acquired {
		return nil, errors.New("request already in progress, please wait")
	}
	defer s.redis.Del(ctx, lockKey)

	// Fetch user & voucher to validate before starting DB tx (fail fast)
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	voucher, err := s.voucherRepo.GetByID(ctx, voucherID)
	if err != nil {
		return nil, err
	}

	if !voucher.IsActive || voucher.Stock <= 0 {
		return nil, domain.ErrOutOfStock
	}

	if user.TotalPoints < voucher.PointsRequired {
		return nil, domain.ErrInsufficientPts
	}

	var redemption *domain.Redemption

	// 2. Start atomic Database Transaction
	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		}
	}()

	// Lock the user row to prevent race conditions on balance update
	var currentPoints int
	err = tx.GetContext(ctx, &currentPoints, "SELECT total_points FROM users WHERE id = $1 FOR UPDATE", userID)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	if currentPoints < voucher.PointsRequired {
		tx.Rollback()
		return nil, domain.ErrInsufficientPts
	}

	// Lock the voucher row and decrement stock
	var currentStock int
	err = tx.GetContext(ctx, &currentStock, "SELECT stock FROM vouchers WHERE id = $1 FOR UPDATE", voucherID)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	if currentStock <= 0 {
		tx.Rollback()
		return nil, domain.ErrOutOfStock
	}

	// Update queries
	_, err = tx.ExecContext(ctx, "UPDATE vouchers SET stock = stock - 1, total_redeemed = total_redeemed + 1 WHERE id = $1", voucherID)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	_, err = tx.ExecContext(ctx, "UPDATE users SET total_points = total_points - $1 WHERE id = $2", voucher.PointsRequired, userID)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// Insert redemption record
	redemptionID := uuid.New()
	redemption = &domain.Redemption{
		ID:         redemptionID,
		UserID:     userID,
		VoucherID:  voucherID,
		PointsUsed: voucher.PointsRequired,
		Status:     "success",
		CreatedAt:  time.Now(),
	}

	_, err = tx.NamedExecContext(ctx, "INSERT INTO redemptions (id, user_id, voucher_id, points_used, status, created_at) VALUES (:id, :user_id, :voucher_id, :points_used, :status, :created_at)", redemption)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// Insert point transaction record
	pointTxID := uuid.New()
	refID := redemptionID.String()
	pointTx := &domain.PointTransaction{
		ID:              pointTxID,
		UserID:          userID,
		Amount:          -voucher.PointsRequired, // deduction
		TransactionType: domain.TxTypeRedeem,
		ReferenceID:     &refID,
		Status:          "success",
		CreatedAt:       time.Now(),
	}

	_, err = tx.NamedExecContext(ctx, "INSERT INTO point_transactions (id, user_id, amount, transaction_type, reference_id, status, created_at) VALUES (:id, :user_id, :amount, :transaction_type, :reference_id, :status, :created_at)", pointTx)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return nil, err
	}

	log.Info().Msgf("User %s successfully redeemed voucher %s", userID.String(), voucherID.String())

	return redemption, nil
}
