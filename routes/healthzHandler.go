package routes

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"webapp/db"
)

func HealthzGetReqHandler(db db.Database) gin.HandlerFunc {
	return func(context *gin.Context) {
		context.Header("Cache-Control", "no-store, no-cache, must-revalidate;")
		if context.Request.Method == http.MethodGet {
			querystring := context.Request.URL.RawQuery
			all, err := io.ReadAll(context.Request.Body)
			if err != nil {
				fmt.Printf("Error while reading the body, %s\n", err)
			}
			if querystring != "" || len(all) > 0 {
				context.String(http.StatusBadRequest, "")
			} else {
				err := db.Ping()
				if err != nil {
					context.String(http.StatusServiceUnavailable, "")
					return
				}
				context.String(http.StatusOK, "")
			}
		} else {
			context.String(http.StatusMethodNotAllowed, "")
		}
	}
}
