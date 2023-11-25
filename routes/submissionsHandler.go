package routes

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"webapp/logger"
	"webapp/models"
	"webapp/services"
	"webapp/utils"
)

func SubmissionsPostHandler(services services.APIServices) gin.HandlerFunc {
	return func(c *gin.Context) {
		utils.StatIncrement("CreateSubmission", 1)

		var submission models.Submission

		if err := c.Bind(&submission); err != nil {
			logger.Error("Failed to bind incoming payload with Gin", zap.Error(err))
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		//Get assignment ID
		assignmentID := c.Param("id")

		// Get is UUID is invalid and sent bac request
		if !utils.IsValidUUID(assignmentID) {
			logger.Warn("Assignment UUID is invalid", zap.Any("assignmentID", assignmentID))
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID"})
			return
		}

		//Get assignment from DB
		assignment, err := services.AssignmentService.GetAssignmentByID(assignmentID)
		if err != nil {
			logger.Error("Failed to get the assignment", zap.Error(err))
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Assignment not found"})
			return
		}

		account, exists := c.Value("loggedInAccount").(models.Account)

		submission = models.Submission{
			SubmissionURL: submission.SubmissionURL,
			AssignmentID:  assignment.ID,
		}

		if exists {
			submission.AccountID = account.ID
			err := services.SubmissionsService.CreateSubmission(&submission)
			if err != nil {
				logger.Error("Failed to submit assignment", zap.Error(err))
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
				return
			}
			val, _ := services.SubmissionsService.GetSubmissionByID(submission.ID)
			logger.Info("Assignment submitted successfully", zap.Any("submission", val))
			c.JSON(http.StatusCreated, submission)
		}
		
		return
	}
}
