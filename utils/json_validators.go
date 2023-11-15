package utils

import (
	"fmt"
	"github.com/xeipuuv/gojsonschema"
	"go.uber.org/zap"
	"webapp/logger"
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

func ValidateAssignmentInput(data string) (bool, []string, error) {
	schemaLoader := gojsonschema.NewStringLoader(assignmentSchema)
	payloadLoader := gojsonschema.NewStringLoader(data)

	result, err := gojsonschema.Validate(schemaLoader, payloadLoader)
	if err != nil {
		logger.Error("Failed validating the schema", zap.Error(err))
		return false, nil, err
	}

	if result.Valid() {
		logger.Info("The incoming payload is VALID")
		return true, nil, nil
	} else {
		errors := result.Errors()
		var errorSlice []string

		for i := 0; i < len(errors); i++ {
			errorSlice = append(errorSlice, fmt.Sprintf("%v, %v", errors[i].Field(), errors[i].Description()))
			fmt.Printf("validation error: %v - %v\n", errors[i].Field(), errors[i].Description())
		}
		logger.Warn("The incoming payload is INVALID", zap.Any("Validation errors", errorSlice))
		return false, errorSlice, nil
	}

}
