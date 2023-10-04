package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"webapp/services"
)

func BasicAuth(services services.APIServices) gin.HandlerFunc {
	return func(context *gin.Context) {
		email, password, ok := context.Request.BasicAuth()
		fmt.Printf("In middleware email=%s password=%s ok=%v\n", email, password, ok)
		if !ok {
			context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid base64 encoding"})
		} else {
			account, err := services.AccountsService.GetAccountByEmail(email)
			if err != nil {
				fmt.Printf("Unable to get account with email %s, %s\n", email, err)
				context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unable to find username"})
			} else {
				success := services.AccountsService.CheckPasswordHash(password, account.Password)
				if success {
					context.Set("loggedInAccount", account)
					context.Next()
				} else {
					fmt.Println("Authentication failed")
					context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization failed"})
				}
			}
		}
	}
}
