package middlwares

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func JWTAuthMiddleware(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := c.Cookie("Bearer")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "authorization token required",
			})
			return
		}
		fmt.Println(token)

		parsedToken, err := jwt.Parse(token, func(t *jwt.Token) (any, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signed method: %v", t.Header["alg"])
			}

			return []byte(jwtSecret), nil
		})
		if err != nil {
			if errors.Is(err, jwt.ErrTokenExpired) {
				c.JSON(http.StatusUnauthorized, gin.H{
					"error": "token expired",
				})
				return
			}
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "invalid token",
			})
			return
		}

		if !parsedToken.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "invalid token",
			})
			return
		}

		claims, ok := parsedToken.Claims.(jwt.MapClaims)
		if ok {
			if userIDFloat, ok := claims["user_id"].(float64); ok {
				userID := int64(userIDFloat)
				c.Set("user_id", userID)
			} else {
				c.JSON(http.StatusUnauthorized, gin.H{
					"error": "invalid user_id type in token",
				})
				return
			}
		}

		c.Next()
	}
}
