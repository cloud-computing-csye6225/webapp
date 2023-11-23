package routes

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"webapp/services"
	"webapp/utils"
)

func SubmissionsPostHandler(services services.APIServices) gin.HandlerFunc {
	return func(c *gin.Context) {
		utils.StatIncrement("CreateSubmission", 1)

		assignmentID := c.Param("id")

		c.JSON(http.StatusOK, gin.H{"id": assignmentID})
		return
	}
}
