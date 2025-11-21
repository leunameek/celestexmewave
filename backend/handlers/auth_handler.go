package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/leunameek/celestexmewave/internal/services"
	"github.com/leunameek/celestexmewave/internal/utils"
)

// Peticion de registro, formato chill
type RegisterRequest struct {
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
	Password  string `json:"password" binding:"required"`
}

// Peticion de login, bien basica
type LoginRequest struct {
	EmailOrPhone string `json:"email" binding:"required"`
	Password     string `json:"password" binding:"required"`
}

// Peticion para refrescar token
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// Peticion pa resetear clave
type PasswordResetRequest struct {
	Email string `json:"email"`
	Phone string `json:"phone"`
}

// Peticion pa verificar codigo de reset
type VerifyResetCodeRequest struct {
	Email       string `json:"email"`
	Phone       string `json:"phone"`
	ResetCode   string `json:"reset_code" binding:"required"`
	NewPassword string `json:"new_password" binding:"required"`
}

// Registro de usuario nuevo
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

	// Generamos tokens pa que quede logueado de una
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

// Login para autenticar
func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	user, accessToken, refreshToken, err := services.LoginUser(req.EmailOrPhone, req.Password)
	if err != nil {
		if errors.Is(err, services.ErrInvalidCredentials) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Credenciales inválidas"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo iniciar sesión"})
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

// Refrescar el token de acceso
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

// Logout para limpiar sesion (jwt es stateless)
func Logout(c *gin.Context) {
	// En JWT stateless, el logout lo hace el cliente borrando el token
	c.JSON(http.StatusOK, gin.H{"message": "logged out successfully"})
}

// Pedir codigo de reset
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

// Verificar codigo de reset y cambiar clave
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
