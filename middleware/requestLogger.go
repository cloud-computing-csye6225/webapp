package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"io"
	"time"
	"webapp/logger"
)

func LogWebRequests() gin.HandlerFunc {
	return func(c *gin.Context) {
		logger.Info(fmt.Sprintf("Incoming request at %s", c.Request.URL.Path))
		// Read the Body content
		var body []byte
		if c.Request.Body != nil {
			body, _ = io.ReadAll(c.Request.Body)
		}
		// Restore the io.ReadCloser to its original state
		c.Request.Body = io.NopCloser(bytes.NewBuffer(body))

		var dat map[string]interface{}
		if len(body) > 0 {
			if err := json.Unmarshal(body, &dat); err != nil {
				logger.Warn("Unable to parse incoming body to JSON", zap.Error(err))
			}
		}
		requestLog := struct {
			Endpoint string
			rawBody  string
			JsonBody map[string]interface{}
			Params   gin.Params
			Time     time.Time
		}{
			c.Request.URL.Path,
			string(body),
			dat,
			c.Params,
			time.Now().UTC(),
		}
		logger.Info("Incoming request data", zap.Any("data", requestLog))
		c.Next()
		logger.Info("Request completed", zap.Any("data", gin.H{"time": time.Now().UTC()}))
	}
}
