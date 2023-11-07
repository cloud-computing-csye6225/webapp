package routes

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"webapp/logger"
	"webapp/utils"
)

func NoRouteHandler(context *gin.Context) {
	utils.StatIncrement("InvalidRoutes", 1)
	context.Header("Cache-Control", "no-store, no-cache, must-revalidate;")
	logger.Warn("Invalid request")
	context.String(http.StatusNotFound, "")
}
