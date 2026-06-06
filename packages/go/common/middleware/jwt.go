package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	commonAuth "github.com/aridevk/dark-kitchen/packages/go/common/auth"
	"github.com/aridevk/dark-kitchen/packages/go/common/response"
)

const UserIDKey = "user_id"
const UserEmailKey = "user_email"
const UserRoleKey = "user_role"

func JWT(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			response.Error(c, http.StatusUnauthorized, "MISSING_AUTH_HEADER", "Missing Authorization header")
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)

		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			response.Error(c, http.StatusUnauthorized, "INVALID_AUTH_HEADER", "Invalid Authorization header")
			c.Abort()
			return
		}

		tokenString := strings.TrimSpace(parts[1])

		if tokenString == "" {
			response.Error(c, http.StatusUnauthorized, "EMPTY_TOKEN", "Empty token")
			c.Abort()
			return
		}

		claims, err := commonAuth.ValidateToken(tokenString, secret)
		if err != nil {
			response.Error(c, http.StatusUnauthorized, "INVALID_TOKEN", "Invalid or expired token")
			c.Abort()
			return
		}

		c.Set(UserIDKey, claims.UserID)
		c.Set(UserEmailKey, claims.Email)
		c.Set(UserRoleKey, claims.Role)

		c.Next()
	}
}

func GetUserID(c *gin.Context) uint {
	value, exists := c.Get(UserIDKey)
	if !exists {
		return 0
	}

	userID, ok := value.(uint)
	if !ok {
		return 0
	}

	return userID
}

func GetUserEmail(c *gin.Context) string {
	value, exists := c.Get(UserEmailKey)
	if !exists {
		return ""
	}

	email, ok := value.(string)
	if !ok {
		return ""
	}

	return email
}

func GetUserRole(c *gin.Context) string {
	value, exists := c.Get(UserRoleKey)
	if !exists {
		return ""
	}

	role, ok := value.(string)
	if !ok {
		return ""
	}

	return role
}
