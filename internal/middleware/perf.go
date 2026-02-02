package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

// PerfMiddleware logge les requêtes lentes pour aider au diagnostic
func PerfMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()

		duration := time.Since(start)
		if duration >= 500*time.Millisecond {
			log.Printf("SLOW REQ %s %s status=%d dur=%s", c.Request.Method, c.Request.URL.Path, c.Writer.Status(), duration)
		}
	}
}
