package service

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"

	"scalable-reward-system/internal/config"
	"scalable-reward-system/internal/domain"
	"scalable-reward-system/internal/utils"
)

type userService struct {
	userRepo domain.UserRepository
	cfg      *config.Config
}

func NewUserService(userRepo domain.UserRepository, cfg *config.Config) domain.UserService {
	return &userService{
		userRepo: userRepo,
		cfg:      cfg,
	}
}

func (s *userService) Register(ctx context.Context, email, password string) (*domain.User, error) {
	// Check if user exists
	existing, _ := s.userRepo.GetByEmail(ctx, email)
	if existing != nil {
		return nil, errors.New("email already in use")
	}

	hash, err := utils.HashPassword(password)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	user := &domain.User{
		ID:           uuid.New(),
		Email:        email,
		PasswordHash: hash,
		Role:         domain.RoleUser, // default to user
		TotalPoints:  0,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userService) Login(ctx context.Context, email, password string) (string, string, error) {
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return "", "", errors.New("invalid email or password")
	}

	if !utils.CheckPasswordHash(password, user.PasswordHash) {
		return "", "", errors.New("invalid email or password")
	}

	accessToken, err := utils.GenerateToken(user.ID, string(user.Role), s.cfg.JWTSecret, 15*time.Minute)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := utils.GenerateToken(user.ID, string(user.Role), s.cfg.JWTSecret, 7*24*time.Hour)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (s *userService) GetProfile(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	return s.userRepo.GetByID(ctx, id)
}
