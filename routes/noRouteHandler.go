package routes

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"webapp/logger"
)

func NoRouteHandler(context *gin.Context) {
	context.Header("Cache-Control", "no-store, no-cache, must-revalidate;")
	logger.Warn("Invalid request")
	context.String(http.StatusNotFound, "")
}
