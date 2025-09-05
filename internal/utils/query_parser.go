package utils

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetPaginationParams(c *gin.Context) (int, int, error) {
	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil {
		return 0, 0, errors.New("invalid limit format")
	}

	offset, err := strconv.Atoi(c.DefaultQuery("offset", "0"))
	if err != nil {
		return 0, 0, errors.New("invalid offset format")
	}

	if limit < 1 || limit > 100 {
		limit = 100
	}

	if offset < 0 {
		offset = 0
	}

	return limit, offset, nil
}
