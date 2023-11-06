package routes

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"io"
	"net/http"
	"webapp/db"
	"webapp/logger"
)

func HealthzGetReqHandler(db db.Database) gin.HandlerFunc {
	return func(context *gin.Context) {
		context.Header("Cache-Control", "no-store, no-cache, must-revalidate;")
		if context.Request.Method == http.MethodGet {
			querystring := context.Request.URL.RawQuery
			all, err := io.ReadAll(context.Request.Body)
			if err != nil {
				logger.Error("Error while reading the body", zap.Error(err))
			}
			if querystring != "" || len(all) > 0 {
				logger.Warn("Query parameters/body is not allowed for healthz")
				context.String(http.StatusBadRequest, "")
			} else {
				err := db.Ping()
				if err != nil {
					context.String(http.StatusServiceUnavailable, "")
					return
				}
				logger.Info("Webservice is healthy")
				context.String(http.StatusOK, "")
			}
		} else {
			logger.Warn("Invalid request URL")
			context.String(http.StatusMethodNotAllowed, "")
		}
	}
}
