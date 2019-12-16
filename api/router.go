package api

import (
	"github.com/gin-gonic/gin"

	"seaman/api/health"
)

func Router() *gin.Engine {
	r := gin.Default()
	r.GET("/health", health.Health)
	return r
}
