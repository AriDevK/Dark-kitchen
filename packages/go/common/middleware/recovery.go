package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/aridevk/dark-kitchen/packages/go/common/response"
)

func Recovery(log *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				requestID := GetRequestID(c)

				log.Error("panic recovered",
					zap.Any("error", err),
					zap.String("request_id", requestID),
				)

				response.Error(
					c,
					http.StatusInternalServerError,
					"INTERNAL_SERVER_ERROR",
					"Internal server error",
				)

				c.Abort()
			}
		}()

		c.Next()
	}
}
