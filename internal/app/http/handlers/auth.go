package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/I-Van-Radkov/messenger/internal/dto"
	"github.com/I-Van-Radkov/messenger/internal/services/auth"
	"github.com/gin-gonic/gin"
)

type AuthProvider interface {
	Register(ctx context.Context, email, username, password string) (string, error)
	Login(ctx context.Context, email, password string) (string, error)
}

type AuthHandlers struct {
	authService AuthProvider
}

func NewAuthHandlers(authService AuthProvider) *AuthHandlers {
	return &AuthHandlers{
		authService: authService,
	}
}

func (h *AuthHandlers) RegisterHandler(c *gin.Context) {
	ctx := c.Request.Context()

	var req dto.RegisterRequest
	if err := json.NewDecoder(c.Request.Body).Decode(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}
	defer c.Request.Body.Close()

	if req.Email == "" || req.Password == "" || req.Username == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "email, username and password are required",
		})
		return
	}

	tokenString, err := h.authService.Register(ctx, req.Email, req.Username, req.Password) // возвращает jwt token и ошибку
	if err != nil {
		switch {
		case errors.Is(err, auth.ErrEmailAlreadyExists):
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		case errors.Is(err, auth.ErrUsernameAlreadyExists):
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "Bearer",
		Value:    tokenString,
		Path:     "/",
		HttpOnly: true,
		Expires:  time.Now().Add(24 * time.Hour),
	})
	c.JSON(http.StatusCreated, gin.H{"message": "user registered successfully"})
}

func (h *AuthHandlers) LoginHandler(c *gin.Context) {
	ctx := c.Request.Context()

	var req dto.LoginRequest
	if err := json.NewDecoder(c.Request.Body).Decode(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}
	defer c.Request.Body.Close()

	if req.Email == "" || req.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "email and password are required",
		})
		return
	}

	tokenString, err := h.authService.Login(ctx, req.Email, req.Password)
	if err != nil {
		switch {
		case errors.Is(err, auth.ErrUserNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case errors.Is(err, auth.ErrInvalidPassword):
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "Bearer",
		Value:    tokenString,
		Path:     "/",
		HttpOnly: true,
		Expires:  time.Now().Add(24 * time.Hour),
	})
	c.JSON(http.StatusOK, gin.H{"message": "login successful"})
}

func (h *AuthHandlers) LogoutHandler(c *gin.Context) {
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "Bearer",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Expires:  time.Unix(0, 0),
	})
	c.JSON(http.StatusOK, gin.H{"message": "logged out successfully"})
}
