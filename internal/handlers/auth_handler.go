package handlers

import (
	"errors"
	"net/http"

	"storage-manager-api/internal/services"
	"storage-manager-api/internal/utils"

	"github.com/gin-gonic/gin"
)

// AuthRequest merepresentasikan payload JSON untuk register dan login.
type AuthRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// AuthHandler menangani request HTTP terkait autentikasi.
type AuthHandler struct {
	authService services.AuthService
}

// NewAuthHandler membuat instance baru AuthHandler.
func NewAuthHandler(authService services.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// Register menangani endpoint POST /api/v1/auth/register.
func (h *AuthHandler) Register(c *gin.Context) {
	var req AuthRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Input tidak valid")
		return
	}

	user, err := h.authService.Register(req.Email, req.Password)
	if err != nil {
		if errors.Is(err, services.ErrInvalidInput) || errors.Is(err, services.ErrEmailAlreadyExists) {
			utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "Gagal mendaftarkan user")
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "User berhasil didaftarkan", gin.H{
		"id":    user.ID,
		"email": user.Email,
	})
}

// Login menangani endpoint POST /api/v1/auth/login.
func (h *AuthHandler) Login(c *gin.Context) {
	var req AuthRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Input tidak valid")
		return
	}

	token, err := h.authService.Login(req.Email, req.Password)
	if err != nil {
		if errors.Is(err, services.ErrInvalidCredentials) {
			utils.ErrorResponse(c, http.StatusUnauthorized, "Email atau password salah")
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "Gagal memproses login")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Login berhasil", gin.H{
		"token": token,
	})
}
