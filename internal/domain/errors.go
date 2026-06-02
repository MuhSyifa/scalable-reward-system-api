package domain

import "errors"

var (
	ErrOutOfStock       = errors.New("voucher is out of stock")
	ErrInsufficientPts  = errors.New("insufficient points")
	ErrCampaignInactive = errors.New("campaign is not active")
	ErrNotFound         = errors.New("record not found")
)
