package routes

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"io"
	"net/http"
	"time"
	"webapp/logger"
	"webapp/models"
	"webapp/services"
	"webapp/utils"
)

func AssignmentsPostHandler(services services.APIServices) gin.HandlerFunc {
	return func(c *gin.Context) {
		utils.StatIncrement("CreateAssignment", 1)

		// Read the Body content
		var body []byte
		if c.Request.Body != nil {
			body, _ = io.ReadAll(c.Request.Body)
		}
		// Restore the io.ReadCloser to its original state
		c.Request.Body = io.NopCloser(bytes.NewBuffer(body))

		isValid, validationErrors, err := utils.ValidateAssignmentInput(string(body))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if !isValid {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errors": validationErrors})
			return
		}

		var assignment models.Assignment

		if err := c.Bind(&assignment); err != nil {
			logger.Error("Failed to bind incoming payload with Gin", zap.Error(err))
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		//Check if deadline is in future
		if !isFutureTime(assignment.Deadline) {
			logger.Warn("Assignment deadline is in past")
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Assignment deadline cannot be in past"})
			return
		}

		// Recreating the model without timestamps
		assignment = models.Assignment{
			Name:          assignment.Name,
			Points:        assignment.Points,
			NumOfAttempts: assignment.NumOfAttempts,
			Deadline:      assignment.Deadline,
		}

		account, exists := c.Value("loggedInAccount").(models.Account)
		if exists {
			assignment.AccountID = account.ID
			err := services.AssignmentService.AddAssignment(&assignment)
			if err != nil {
				logger.Error("Failed to create an assignment", zap.Error(err))
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
				return
			}
			logger.Info("Assignment created successfully")
			c.JSON(http.StatusCreated, assignment)
		}
	}
}

func AssignmentGetHandler(services services.APIServices) gin.HandlerFunc {
	return func(c *gin.Context) {
		utils.StatIncrement("GetAssignments", 1)
		all, err := io.ReadAll(c.Request.Body)
		if err != nil {
			logger.Error("Error while reading the body", zap.Error(err))
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		} else {
			if len(all) > 0 {
				logger.Warn("Body is not allowed")
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "body is not allowed for GET request"})
				return
			}
		}

		assignments, err := services.AssignmentService.GetAssignment()
		if err != nil {
			logger.Error("Failed to get the assignments", zap.Error(err))
			c.AbortWithStatusJSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})
			return
		}
		logger.Info("Successfully got the assignments")
		c.JSON(http.StatusOK, assignments)
	}
}

func AssignmentGetByIDHandler(services services.APIServices) gin.HandlerFunc {
	return func(c *gin.Context) {
		utils.StatIncrement("GetAssignmentByID", 1)
		// Check for body
		all, err := io.ReadAll(c.Request.Body)
		if err != nil {
			logger.Error("Error while reading the body", zap.Error(err))
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		} else {
			if len(all) > 0 {
				logger.Warn("Body is not allowed")
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "body is not allowed for GET request"})
				return
			}
		}

		// Get params
		assignmentID := c.Param("id")

		// Get is UUID is invalid and sent bac request
		if !IsValidUUID(assignmentID) {
			logger.Warn("Assignment UUID is invalid", zap.Any("assignmentID", assignmentID))
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID"})
			return
		}

		// Get assignment
		assignment, err := services.AssignmentService.GetAssignmentByID(assignmentID)
		if err != nil {
			logger.Error("Failed to get the assignment", zap.Error(err))
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Assignment not found"})
			return
		}
		logger.Info("Successfully got the assignments")
		c.JSON(http.StatusOK, assignment)
	}
}

func AssignmentPutHandler(services services.APIServices) gin.HandlerFunc {
	return func(c *gin.Context) {
		utils.StatIncrement("UpdateAssignment", 1)

		// Read the Body content
		var body []byte
		if c.Request.Body != nil {
			body, _ = io.ReadAll(c.Request.Body)
		}

		// Restore the io.ReadCloser to its original state
		c.Request.Body = io.NopCloser(bytes.NewBuffer(body))

		isValid, validationErrors, err := utils.ValidateAssignmentInput(string(body))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if !isValid {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errors": validationErrors})
			return
		}

		// Get params
		assignmentID := c.Param("id")

		// Get is UUID is invalid and sent bac request
		if !IsValidUUID(assignmentID) {
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

		// Bind new data to the retrieved assignment
		if err := c.ShouldBindJSON(&assignment); err != nil {
			logger.Error("Failed to bind incoming payload with Gin", zap.Error(err))
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		//Checking authorization,
		account, exists := c.Value("loggedInAccount").(models.Account)
		if exists {
			if account.ID != assignment.AccountID {
				logger.Warn("Unauthorized access to the assignment")
				c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "unauthorized access"})
				return
			}
		}

		//Check if deadline is in future
		if !isFutureTime(assignment.Deadline) {
			logger.Warn("Assignment deadline is in past")
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Assignment deadline cannot be in past"})
			return
		}

		updatedAssignment, err := services.AssignmentService.UpdateAssignment(assignment)
		if err != nil {
			logger.Error("Failed to update an assignment", zap.Error(err))
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		logger.Info("Updated the assignment", zap.Any("oldAssignment", assignment), zap.Any("updatedAssignment", updatedAssignment))
		logger.Info("Assignment updated successfully")
		c.String(http.StatusNoContent, "")
	}
}

func AssignmentDeleteHandler(services services.APIServices) gin.HandlerFunc {
	return func(c *gin.Context) {
		utils.StatIncrement("DeleteAssignment", 1)

		// Check if body is empty
		all, err := io.ReadAll(c.Request.Body)
		if err != nil {
			logger.Error("Error while reading the body", zap.Error(err))
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		} else {
			if len(all) > 0 {
				logger.Warn("Body is not allowed")
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "body is not allowed for DELETE request"})
				return
			}
		}

		// Get params
		assignmentID := c.Param("id")

		// Get is UUID is invalid and sent bac request
		if !IsValidUUID(assignmentID) {
			logger.Warn("Assignment UUID is invalid", zap.Any("assignmentID", assignmentID))
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID"})
			return
		}

		// Get assignment from DB
		assignment, err := services.AssignmentService.GetAssignmentByID(assignmentID)
		if err != nil {
			logger.Error("Failed to get the assignment", zap.Error(err))
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Assignment not found"})
			return
		}

		//Checking authorization,
		account, exists := c.Value("loggedInAccount").(models.Account)
		if exists {
			if account.ID != assignment.AccountID {
				logger.Warn("Unauthorized access to the assignment")
				c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "unauthorized access"})
				return
			}
		}

		err = services.AssignmentService.DeleteAssignment(assignment)
		if err != nil {
			logger.Error("Failed to delete the assignment", zap.Error(err))
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		logger.Info("Successfully deleted the assignment")
		c.Status(http.StatusNoContent)
		return
	}
}

func AssignmentPatchHandler(services services.APIServices) gin.HandlerFunc {
	return func(c *gin.Context) {
		utils.StatIncrement("InvalidHandler", 1)

		logger.Warn("Method not allowed")
		c.JSON(http.StatusMethodNotAllowed, gin.H{"error": "Method not allowed"})
	}
}

func IsValidUUID(u string) bool {
	_, err := uuid.Parse(u)
	return err == nil
}

func isFutureTime(deadline time.Time) bool {
	currentTime := time.Now()
	return deadline.After(currentTime)
}
