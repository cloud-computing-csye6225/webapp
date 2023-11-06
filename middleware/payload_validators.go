package middleware

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/xeipuuv/gojsonschema"
	"go.uber.org/zap"
	"io"
	"net/http"
	"webapp/logger"
	"webapp/services"
)

const assignmentSchema = `
	{
	  "type": "object",
	  "properties": {
		"name": {
		  "type": "string"
		},
		"points": {
		  "type": "integer",
		  "minimum": 1,
		  "maximum": 100
		},
		"num_of_attempts": {
		  "type": "integer",
		  "minimum": 1,
		  "maximum": 100
		},
		"deadline": {
		  "type": "string",
		  "format": "date-time"
		},
		"assignment_created": {
		  "type": "string",
		  "format": "date-time"
		},
		"assignment_updated": {
		  "type": "string",
		  "format": "date-time"
		}
	  },
	  "required": ["name", "points", "num_of_attempts", "deadline"],
	  "additionalProperties": false
	}
`

func ValidateAssignmentsPayload(services services.APIServices) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Read the Body content
		var body []byte
		if c.Request.Body != nil {
			body, _ = io.ReadAll(c.Request.Body)
		}
		// Restore the io.ReadCloser to its original state
		c.Request.Body = io.NopCloser(bytes.NewBuffer(body))

		schemaLoader := gojsonschema.NewStringLoader(assignmentSchema)
		payloadLoader := gojsonschema.NewStringLoader(string(body))

		result, err := gojsonschema.Validate(schemaLoader, payloadLoader)
		if err != nil {
			logger.Error("Failed validating the schema", zap.Error(err))
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if result.Valid() {
			logger.Info("The incoming payload is VALID")
			c.Next()
			return
		} else {
			errors := result.Errors()
			var errorSlice []string

			for i := 0; i < len(errors); i++ {
				errorSlice = append(errorSlice, fmt.Sprintf("%v, %v", errors[i].Field(), errors[i].Description()))
				fmt.Printf("validation error: %v - %v\n", errors[i].Field(), errors[i].Description())
			}
			logger.Warn("The incoming payload is INVALID", zap.Any("Validation errors", errorSlice))
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errors": errorSlice})
			return
		}
	}
}
