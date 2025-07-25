package middleware

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// CORSMiddleware configura CORS para a aplicação
func CORSMiddleware() gin.HandlerFunc {
	// URLs permitidas para desenvolvimento
	allowedOrigins := []string{
		"http://localhost:3000", // Frontend React/Next.js
		"http://localhost:3001", // Frontend alternativo
		"http://localhost:8082", // API própria (Docker)
		"http://localhost:8080", // API própria (direto)
		"http://localhost:8081", // API desenvolvimento (Docker dev)
	}

	return cors.New(cors.Config{
		AllowOrigins:     allowedOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	})
}

// SecurityMiddleware adiciona headers de segurança
func SecurityMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
		c.Next()
	}
}

// Estrutura para armazenar informações de rate limiting por IP
type rateLimitInfo struct {
	tokens     float64   // Tokens disponíveis
	lastRefill time.Time // Última vez que tokens foram reabastecidos
}

// Rate limiter global em memória
var (
	rateLimitMap   = make(map[string]*rateLimitInfo)
	rateLimitMutex = sync.RWMutex{}
)

// RateLimitMiddleware implementa rate limiting usando algoritmo token bucket
// Configuração padrão: 100 requisições por minuto por IP
func RateLimitMiddleware() gin.HandlerFunc {
	return RateLimitWithConfig(100, time.Minute) // 100 req/min por IP
}

// RateLimitWithConfig permite configurar rate limiting customizado
func RateLimitWithConfig(maxRequests int, duration time.Duration) gin.HandlerFunc {
	refillRate := float64(maxRequests) / duration.Seconds() // tokens por segundo

	return func(c *gin.Context) {
		clientIP := c.ClientIP()

		// Pular rate limiting para health checks
		if c.Request.URL.Path == "/ping" || c.Request.URL.Path == "/health" || c.Request.URL.Path == "/healthz" {
			c.Next()
			return
		}

		rateLimitMutex.Lock()
		defer rateLimitMutex.Unlock()

		now := time.Now()

		// Inicializar ou buscar informações do IP
		info, exists := rateLimitMap[clientIP]
		if !exists {
			info = &rateLimitInfo{
				tokens:     float64(maxRequests),
				lastRefill: now,
			}
			rateLimitMap[clientIP] = info
		}

		// Calcular tokens a adicionar baseado no tempo passado
		timePassed := now.Sub(info.lastRefill).Seconds()
		tokensToAdd := timePassed * refillRate

		// Atualizar tokens (máximo é o limite configurado)
		info.tokens = min(float64(maxRequests), info.tokens+tokensToAdd)
		info.lastRefill = now

		// Verificar se há tokens suficientes
		if info.tokens < 1 {
			// Rate limit excedido
			c.Header("X-RateLimit-Limit", fmt.Sprintf("%d", maxRequests))
			c.Header("X-RateLimit-Remaining", "0")
			c.Header("X-RateLimit-Reset", fmt.Sprintf("%d", now.Add(duration).Unix()))

			c.JSON(http.StatusTooManyRequests, gin.H{
				"error":       "Rate limit exceeded",
				"message":     "Too many requests. Please try again later.",
				"retry_after": duration.Seconds(),
			})
			c.Abort()
			return
		}

		// Consumir um token
		info.tokens--

		// Adicionar headers informativos
		c.Header("X-RateLimit-Limit", fmt.Sprintf("%d", maxRequests))
		c.Header("X-RateLimit-Remaining", fmt.Sprintf("%d", int(info.tokens)))

		c.Next()
	}
}

// Função auxiliar min para Go < 1.21
func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

// CleanupRateLimit remove entradas antigas do rate limiter (opcional, para limpeza de memória)
func CleanupRateLimit() {
	rateLimitMutex.Lock()
	defer rateLimitMutex.Unlock()

	now := time.Now()
	for ip, info := range rateLimitMap {
		// Remove IPs inativos por mais de 1 hora
		if now.Sub(info.lastRefill) > time.Hour {
			delete(rateLimitMap, ip)
		}
	}
}
