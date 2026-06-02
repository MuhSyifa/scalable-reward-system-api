package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"scalable-reward-system/internal/domain"
	"scalable-reward-system/internal/utils"
)

type AuthHandler struct {
	userService domain.UserService
}

func NewAuthHandler(us domain.UserService) *AuthHandler {
	return &AuthHandler{userService: us}
}

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request payload")
		return
	}

	user, err := h.userService.Register(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		utils.ErrorResponse(c, http.StatusConflict, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "User registered successfully", user)
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request payload")
		return
	}

	accessToken, refreshToken, err := h.userService.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Login successful", gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}
