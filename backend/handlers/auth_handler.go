package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/leunameek/celestexmewave/internal/services"
	"github.com/leunameek/celestexmewave/internal/utils"
)

// RegisterRequest represents a registration request
type RegisterRequest struct {
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
	Password  string `json:"password" binding:"required"`
}

// LoginRequest represents a login request
type LoginRequest struct {
	EmailOrPhone string `json:"email" binding:"required"`
	Password     string `json:"password" binding:"required"`
}

// RefreshTokenRequest represents a refresh token request
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// PasswordResetRequest represents a password reset request
type PasswordResetRequest struct {
	Email string `json:"email"`
	Phone string `json:"phone"`
}

// VerifyResetCodeRequest represents a reset code verification request
type VerifyResetCodeRequest struct {
	Email       string `json:"email"`
	Phone       string `json:"phone"`
	ResetCode   string `json:"reset_code" binding:"required"`
	NewPassword string `json:"new_password" binding:"required"`
}

// Register registers a new user
func Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	user, err := services.RegisterUser(req.Email, req.Phone, req.FirstName, req.LastName, req.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Generate tokens for auto-login
	emailOrPhone := req.Email
	if emailOrPhone == "" {
		emailOrPhone = req.Phone
	}

	accessToken, err := utils.GenerateAccessToken(user.ID, emailOrPhone, user.FirstName, user.LastName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate access token"})
		return
	}

	refreshToken, err := utils.GenerateRefreshToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate refresh token"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":            user.ID,
		"email":         user.Email,
		"phone":         user.Phone,
		"first_name":    user.FirstName,
		"last_name":     user.LastName,
		"is_registered": user.IsRegistered,
		"created_at":    user.CreatedAt,
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"expires_in":    86400,
	})
}

// Login authenticates a user
func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	user, accessToken, refreshToken, err := services.LoginUser(req.EmailOrPhone, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"expires_in":    86400,
		"user": gin.H{
			"id":         user.ID,
			"email":      user.Email,
			"first_name": user.FirstName,
			"last_name":  user.LastName,
		},
	})
}

// RefreshToken refreshes the access token
func RefreshToken(c *gin.Context) {
	var req RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	accessToken, err := services.RefreshAccessToken(req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token": accessToken,
		"expires_in":   86400,
	})
}

// Logout logs out a user
func Logout(c *gin.Context) {
	// In a stateless JWT system, logout is handled on the client side
	// by removing the token from localStorage
	c.JSON(http.StatusOK, gin.H{"message": "logged out successfully"})
}

// RequestPasswordReset requests a password reset code
func RequestPasswordReset(c *gin.Context) {
	var req PasswordResetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	emailOrPhone := req.Email
	if emailOrPhone == "" {
		emailOrPhone = req.Phone
	}

	_, err := services.RequestPasswordReset(emailOrPhone)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "reset code sent to your email",
		"expires_in": 3600,
	})
}

// VerifyResetCode verifies a reset code and updates the password
func VerifyResetCode(c *gin.Context) {
	var req VerifyResetCodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	emailOrPhone := req.Email
	if emailOrPhone == "" {
		emailOrPhone = req.Phone
	}

	err := services.VerifyResetCode(emailOrPhone, req.ResetCode, req.NewPassword)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "password updated successfully",
		"redirect": "/login",
	})
}
