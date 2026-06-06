package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/aridevk/dark-kitchen/packages/go/common/response"
)

func RequireRole(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole := GetUserRole(c)

		if userRole == "" {
			response.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", "Unauthorized")
			c.Abort()
			return
		}

		for _, role := range allowedRoles {
			if userRole == role {
				c.Next()
				return
			}
		}

		response.Error(c, http.StatusForbidden, "FORBIDDEN", "Forbidden")
		c.Abort()
	}
}
