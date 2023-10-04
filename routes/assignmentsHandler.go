package routes

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"webapp/models"
	"webapp/services"
)

func TestHandler(services services.APIServices) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.String(http.StatusOK, "Endpoint working")
	}
}

func AssignmentsPostHandler(services services.APIServices) gin.HandlerFunc {
	return func(c *gin.Context) {
		var assignment models.Assignment

		if err := c.ShouldBindJSON(&assignment); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		account, exists := c.Value("loggedInAccount").(models.Account)
		if exists {
			fmt.Printf("Logged in account is %v", account)
			assignment.AccountID = account.ID
			err := services.AssignmentService.AddAssignment(assignment)
			if err != nil {
				return
			}
		}
	}
}

func AssignmentGetHandler(services services.APIServices) gin.HandlerFunc {
	return func(c *gin.Context) {
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
		assignmentID := c.Param("id")
		assignment, err := services.AssignmentService.GetAssignmentByID(assignmentID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, assignment)
	}
}

func AssignmentPutHandler(services services.APIServices) gin.HandlerFunc {
	return func(c *gin.Context) {
		assignmentID := c.Param("id")

		assignment, err := services.AssignmentService.GetAssignmentByID(assignmentID)
		if err != nil {
			fmt.Printf("Error getting the assignment with id %v, %v\n", assignmentID, err)
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		if err := c.ShouldBindJSON(&assignment); err != nil {
			fmt.Printf("Error binding the assignment, %v\n", err)
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err = services.AssignmentService.UpdateAssignment(assignment)
		if err != nil {
			fmt.Printf("Error updating the assignment with id %v, %v\n", assignmentID, err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, assignment)
	}
}

func AssignmentDeleteHandler(services services.APIServices) gin.HandlerFunc {
	return func(c *gin.Context) {
		assignmentID := c.Param("id")
		assignment, err := services.AssignmentService.GetAssignmentByID(assignmentID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		err = services.AssignmentService.DeleteAssignment(assignment)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Assignment deleted successfully"})
	}
}
