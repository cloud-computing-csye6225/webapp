package middleware

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"time"
	"webapp/logger"
	"webapp/services"
)

func CheckDB(services services.APIServices) gin.HandlerFunc {
	return func(c *gin.Context) {
		logger.Info("Checking if DB is active")
		start := time.Now()
		err := services.Database.Ping()
		elapsed := time.Since(start)
		logger.Debug("Ping timing test in checkDB middleware", zap.Any("time", elapsed))
		if err != nil {
			logger.Error("Unable to establish connection with DB, aborting the request", zap.Error(err))
			c.AbortWithStatusJSON(http.StatusServiceUnavailable, gin.H{"error": "Service unavailable"})
			return
		}
	}
}
