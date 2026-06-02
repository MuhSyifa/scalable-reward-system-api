package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"scalable-reward-system/internal/domain"
	"scalable-reward-system/internal/utils"
)

type RedemptionHandler struct {
	redemptionService domain.RedemptionService
}

func NewRedemptionHandler(rs domain.RedemptionService) *RedemptionHandler {
	return &RedemptionHandler{redemptionService: rs}
}

type RedeemRequest struct {
	VoucherID string `json:"voucher_id" binding:"required,uuid"`
}

func (h *RedemptionHandler) Redeem(c *gin.Context) {
	userIDStr, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}
	userID := userIDStr.(uuid.UUID)

	var req RedeemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request payload")
		return
	}

	voucherID, err := uuid.Parse(req.VoucherID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid voucher ID format")
		return
	}

	redemption, err := h.redemptionService.RedeemVoucher(c.Request.Context(), userID, voucherID)
	if err != nil {
		if err == domain.ErrOutOfStock || err == domain.ErrInsufficientPts {
			utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Voucher redeemed successfully", redemption)
}
