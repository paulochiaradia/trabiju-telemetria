package middleware

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

// CustomLogger retorna um middleware de logging que filtra requisições de health check
// e adiciona informações úteis de performance e debugging
func CustomLogger() gin.HandlerFunc {
	return gin.LoggerWithConfig(gin.LoggerConfig{
		// Pular logs para endpoints de health check
		SkipPaths: []string{"/ping", "/health", "/healthz"},

		// Formato customizado com mais informações
		Formatter: func(param gin.LogFormatterParams) string {
			// Não logar se for um dos paths skipados
			for _, path := range []string{"/ping", "/health", "/healthz"} {
				if param.Path == path {
					return ""
				}
			}

			// Formato: [TIME] STATUS METHOD PATH (LATENCY) from IP "USER_AGENT"
			return fmt.Sprintf("[%s] %3d %s %s (%v) from %s \"%s\"\n",
				param.TimeStamp.Format("2006/01/02 - 15:04:05"),
				param.StatusCode,
				param.Method,
				param.Path,
				param.Latency.Truncate(time.Microsecond),
				param.ClientIP,
				param.Request.UserAgent(),
			)
		},
	})
}
