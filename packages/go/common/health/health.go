package health

import (
	"github.com/gin-gonic/gin"

	"github.com/aridevk/dark-kitchen/packages/go/common/response"
)

type Status struct {
	Status  string `json:"status"`
	Service string `json:"service"`
}

func Handler(serviceName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		response.OK(c, Status{
			Status:  "UP",
			Service: serviceName,
		})
	}
}
