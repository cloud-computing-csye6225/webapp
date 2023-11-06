package middleware

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"webapp/logger"
	"webapp/services"
)

func CheckDB(services services.APIServices) gin.HandlerFunc {
	return func(c *gin.Context) {
		err := services.Database.Ping()
		if err != nil {
			logger.Error("Unable to establish connection with DB, aborting the request", zap.Error(err))
			c.AbortWithStatusJSON(http.StatusServiceUnavailable, gin.H{"error": "Service unavailable"})
			return
		}
	}
}
