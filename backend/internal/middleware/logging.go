package middleware

import (
	"github.com/gin-gonic/gin"
)

// CustomLogger retorna um middleware de logging que filtra requisições de health check
func CustomLogger() gin.HandlerFunc {
	return gin.LoggerWithConfig(gin.LoggerConfig{
		SkipPaths: []string{"/ping"}, // Skip logging for ping endpoint
	})
}
