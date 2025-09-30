package traceid

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/oklog/ulid/v2"
)

func TraceID() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := c.GetHeader("X-Trace-ID")
		if traceID == "" {
			traceID = fmt.Sprintf("REQ-%s", ulid.Make().String())
			c.Request.Header.Set("X-Trace-ID", traceID)
		}
		c.Writer.Header().Set("X-Trace-ID", traceID)
		c.Next()
	}
}
