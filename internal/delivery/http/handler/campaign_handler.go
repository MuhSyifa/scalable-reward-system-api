package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"scalable-reward-system/internal/domain"
	"scalable-reward-system/internal/utils"
)

type CampaignHandler struct {
	campaignRepo domain.CampaignRepository
	voucherRepo  domain.VoucherRepository
}

func NewCampaignHandler(cr domain.CampaignRepository, vr domain.VoucherRepository) *CampaignHandler {
	return &CampaignHandler{
		campaignRepo: cr,
		voucherRepo:  vr,
	}
}

func (h *CampaignHandler) ListActiveCampaigns(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "10")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, _ := strconv.Atoi(limitStr)
	offset, _ := strconv.Atoi(offsetStr)

	// Since we don't have Redis injected in this handler yet, we'll keep it simple for now, 
	// but normally you would do:
	// cacheKey := fmt.Sprintf("campaigns:active:%d:%d", limit, offset)
	// val, err := redis.Get(ctx, cacheKey).Result()
	// ... if hit return cache, else DB query

	campaigns, err := h.campaignRepo.ListActive(c.Request.Context(), limit, offset)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch campaigns")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Active campaigns retrieved", campaigns)
}

func (h *CampaignHandler) ListVouchers(c *gin.Context) {
	campaignIDStr := c.Param("id")
	campaignID, err := uuid.Parse(campaignIDStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid campaign ID")
		return
	}

	vouchers, err := h.voucherRepo.ListByCampaign(c.Request.Context(), campaignID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch vouchers")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Vouchers retrieved", vouchers)
}
