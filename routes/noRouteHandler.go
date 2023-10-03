package routes

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func NoRouteHandler(context *gin.Context) {
	context.Header("Cache-Control", "no-store, no-cache, must-revalidate;")
	context.String(http.StatusNotFound, "")
}
