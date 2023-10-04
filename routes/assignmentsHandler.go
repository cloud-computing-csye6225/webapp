package routes

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"io"
	"net/http"
	"time"
	"webapp/models"
	"webapp/services"
)

func AssignmentsPostHandler(services services.APIServices) gin.HandlerFunc {
	return func(c *gin.Context) {
		var assignment models.Assignment

		if err := c.Bind(&assignment); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		//Check if deadline is in future
		if !isFutureTime(assignment.Deadline) {
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
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
				return
			}
			c.JSON(http.StatusCreated, assignment)
		}
	}
}

func AssignmentGetHandler(services services.APIServices) gin.HandlerFunc {
	return func(c *gin.Context) {
		all, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		} else {
			fmt.Printf("In get handler, %v", all)
			if len(all) > 0 {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "body is not allowed for GET request"})
				return
			}
		}

		assignments, err := services.AssignmentService.GetAssignment()
		if err != nil {
			c.AbortWithStatusJSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, assignments)
	}
}

func AssignmentGetByIDHandler(services services.APIServices) gin.HandlerFunc {
	return func(c *gin.Context) {

		// Check for body
		all, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		} else {
			if len(all) > 0 {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "body is not allowed for GET request"})
				return
			}
		}

		// Get params
		assignmentID := c.Param("id")

		// Get is UUID is invalid and sent bac request
		if !IsValidUUID(assignmentID) {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID"})
			return
		}

		// Get assignment
		assignment, err := services.AssignmentService.GetAssignmentByID(assignmentID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Assignment not found"})
			return
		}
		c.JSON(http.StatusOK, assignment)
	}
}

func AssignmentPutHandler(services services.APIServices) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get params
		assignmentID := c.Param("id")

		// Get is UUID is invalid and sent bac request
		if !IsValidUUID(assignmentID) {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID"})
			return
		}

		//Get assignment from DB
		assignment, err := services.AssignmentService.GetAssignmentByID(assignmentID)
		if err != nil {
			fmt.Printf("Error getting the assignment with id %v, %v\n", assignmentID, err)
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Assignment not found"})
			return
		}

		// Bind new data to the retrieved assignment
		if err := c.ShouldBindJSON(&assignment); err != nil {
			fmt.Printf("Error binding the assignment, %v\n", err)
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		//Checking authorization,
		account, exists := c.Value("loggedInAccount").(models.Account)
		if exists {
			if account.ID != assignment.AccountID {
				c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "unauthorized access"})
				return
			}
		}

		//Check if deadline is in future
		if !isFutureTime(assignment.Deadline) {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Assignment deadline cannot be in past"})
			return
		}

		updatedAssignment, err := services.AssignmentService.UpdateAssignment(assignment)
		if err != nil {
			fmt.Printf("Error updating the assignment with id %v, %v\n", assignmentID, err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, updatedAssignment)
	}
}

func AssignmentDeleteHandler(services services.APIServices) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check if body is empty
		all, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		} else {
			if len(all) > 0 {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "body is not allowed for DELETE request"})
				return
			}
		}

		// Get params
		assignmentID := c.Param("id")

		// Get is UUID is invalid and sent bac request
		if !IsValidUUID(assignmentID) {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID"})
			return
		}

		// Get assignment from DB
		assignment, err := services.AssignmentService.GetAssignmentByID(assignmentID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Assignment not found"})
			return
		}

		//Checking authorization,
		account, exists := c.Value("loggedInAccount").(models.Account)
		if exists {
			if account.ID != assignment.AccountID {
				fmt.Printf("%v, %v\n", account.ID, assignment.AccountID)
				fmt.Printf("Authorization failed, User %v cannot update the assignment", account.FirstName)
				c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "unauthorized access"})
				return
			}
		}

		err = services.AssignmentService.DeleteAssignment(assignment)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.Status(http.StatusNoContent)
		return
	}
}

func AssignmentPatchHandler(services services.APIServices) gin.HandlerFunc {
	return func(c *gin.Context) {
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
