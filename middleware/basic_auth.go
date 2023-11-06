package middleware

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"webapp/logger"
	"webapp/services"
)

func BasicAuth(services services.APIServices) gin.HandlerFunc {
	return func(context *gin.Context) {
		email, password, ok := context.Request.BasicAuth()
		logger.Info("Authenticating user", zap.Any("user email", email))
		if !ok {
			logger.Error("Unable to authenticate user, failed to parse auth string")
			context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid base64 encoding"})
		} else {
			account, err := services.AccountsService.GetAccountByEmail(email)
			if err != nil {
				logger.Error("Unable to get account with email", zap.Error(err))
				context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unable to find username"})
			} else {
				success := services.AccountsService.CheckPasswordHash(password, account.Password)
				if success {
					context.Set("loggedInAccount", account)
					context.Next()
				} else {
					logger.Error("Authentication failed")
					context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization failed"})
				}
			}
		}
	}
}
