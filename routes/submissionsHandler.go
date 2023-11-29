package routes

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"io"
	"net/http"
	"time"
	"webapp/logger"
	"webapp/models"
	"webapp/services"
	"webapp/utils"
)

func SubmissionsPostHandler(services services.APIServices) gin.HandlerFunc {
	return func(c *gin.Context) {
		utils.StatIncrement("CreateSubmission", 1)

		// Read the Body content
		var body []byte
		if c.Request.Body != nil {
			body, _ = io.ReadAll(c.Request.Body)
		}
		// Restore the io.ReadCloser to its original state
		c.Request.Body = io.NopCloser(bytes.NewBuffer(body))

		isValid, validationErrors, err := utils.ValidateSubmissionInput(string(body))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if !isValid {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errors": validationErrors})
			return
		}

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

			isAttemptValid, err := services.SubmissionsService.CheckSubmissionAttemptValidity(submission, assignment.NumOfAttempts)
			if err != nil {
				logger.Error("Failed to check attempt validity", zap.Error(err))
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
				return
			} else {
				isSubmissionOnTime := services.SubmissionsService.CheckForLateSubmission(assignment.Deadline)

				if !isSubmissionOnTime {
					c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Cannot submit assignment after the deadline"})
					return
				}
				if !isAttemptValid {
					c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Submission attempts exceeded for this assignment"})
					return
				}
			}

			err = services.SubmissionsService.CreateSubmission(&submission)
			if err != nil {
				logger.Error("Failed to submit assignment", zap.Error(err))
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
				return
			}

			messageToLambda := struct {
				SubmissionID   string
				AssignmentName string
				FirstName      string
				EmailID        string
				SubmissionUrl  string
				SubmissionTime time.Time
			}{
				SubmissionID:   submission.ID,
				AssignmentName: assignment.Name,
				FirstName:      account.FirstName,
				EmailID:        account.Email,
				SubmissionUrl:  submission.SubmissionURL,
				SubmissionTime: submission.SubmissionCreated,
			}
			marshal, err := json.Marshal(messageToLambda)
			if err != nil {
				return
			}

			services.AWSService.PublishSubmissionToSNS(string(marshal))
			c.JSON(http.StatusCreated, submission)
		}
		return
	}
}
