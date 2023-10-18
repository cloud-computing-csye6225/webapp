package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"webapp/services"
)

func CheckDB(services services.APIServices) gin.HandlerFunc {
	return func(c *gin.Context) {
		err := services.Database.Ping()
		if err != nil {
			c.AbortWithStatusJSON(http.StatusServiceUnavailable, gin.H{"error": "Service unavailable"})
			return
		}
	}
}
