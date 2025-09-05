package handlers

import (
	"context"
	"database/sql"

	"github.com/I-Van-Radkov/messenger/internal/dto"
	"github.com/I-Van-Radkov/messenger/internal/models"
	"github.com/gin-gonic/gin"
)

type UserProvider interface {
	GetByUsername(ctx context.Context, username string) (*models.User, error)
	GetByID(ctx context.Context, id int64) (*models.User, error)
}

type UserHandlers struct {
	userService UserProvider
}

func NewUserHandlers(userService UserProvider) *UserHandlers {
	return &UserHandlers{
		userService: userService,
	}
}

func (h *UserHandlers) SearchHandler(c *gin.Context) {
	ctx := c.Request.Context()

	username := c.Query("username")

	user, err := h.userService.GetByUsername(ctx, username)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(400, gin.H{"error": "user not found"})
		} else {
			c.JSON(400, gin.H{"error": "internal server error"})
		}
		return
	}

	dtoUserSearch := dto.UserSearchResponse{
		ID:       user.ID,
		Username: user.Username,
	}

	c.JSON(200, dtoUserSearch)
}

func (h *UserHandlers) GetUserHandler(c *gin.Context) {
	ctx := c.Request.Context()

	userId := c.GetInt64("user_id")

	user, err := h.userService.GetByID(ctx, userId)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(400, gin.H{"error": "user not found"})
		} else {
			c.JSON(400, gin.H{"error": "internal server error"})
		}
		return
	}

	dtoUserProfile := dto.UserProfileResponse{
		ID:        user.ID,
		Email:     user.Email,
		Username:  user.Username,
		CreatedAt: user.CreatedAt,
	}

	c.JSON(200, dtoUserProfile)
}
