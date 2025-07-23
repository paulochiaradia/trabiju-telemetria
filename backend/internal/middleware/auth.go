package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/paulochiaradia/trabiju-telemetria/internal/services"
)

// AuthMiddleware middleware de autenticação JWT
func AuthMiddleware(jwtService *services.JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Token de autorização requerido",
			})
			c.Abort()
			return
		}

		// Extrair token do header "Bearer token"
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Formato de token inválido",
			})
			c.Abort()
			return
		}

		token := tokenParts[1]

		// Validar token
		claims, err := jwtService.ValidateToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "Token inválido",
				"details": err.Error(),
			})
			c.Abort()
			return
		}

		// Adicionar informações do usuário ao contexto
		c.Set("user_id", claims.UserID)
		c.Set("user_email", claims.Email)
		c.Set("user_role_id", claims.RoleID)
		c.Set("user_empresa_id", claims.EmpresaID)

		c.Next()
	}
}

// RoleMiddleware middleware para verificar permissões de role
func RoleMiddleware(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Primeiro verifica se passou pelo AuthMiddleware
		userRoleID, exists := c.Get("user_role_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Usuário não autenticado",
			})
			c.Abort()
			return
		}

		// TODO: Implementar verificação de roles por nome
		// Por enquanto, permitir acesso baseado no role_id
		// 1 = admin, 2 = gestor, 3 = entregador, 4 = ajudante

		roleID := userRoleID.(int)

		// Verificar se o role está permitido
		allowed := false
		for _, role := range allowedRoles {
			switch role {
			case "admin":
				if roleID == 1 {
					allowed = true
				}
			case "gestor":
				if roleID == 1 || roleID == 2 {
					allowed = true
				}
			case "entregador":
				if roleID == 1 || roleID == 2 || roleID == 3 {
					allowed = true
				}
			case "ajudante":
				if roleID == 1 || roleID == 2 || roleID == 3 || roleID == 4 {
					allowed = true
				}
			}
			if allowed {
				break
			}
		}

		if !allowed {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Acesso negado",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// CompanyMiddleware middleware para verificar se usuário pertence à empresa
func CompanyMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userEmpresaID, exists := c.Get("user_empresa_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Usuário não autenticado",
			})
			c.Abort()
			return
		}

		// Verificar se a empresa do usuário está ativa
		// TODO: Implementar verificação de empresa ativa

		// Por enquanto, apenas adiciona ao contexto
		c.Set("empresa_id", userEmpresaID)
		c.Next()
	}
}

// OptionalAuthMiddleware middleware de autenticação opcional
func OptionalAuthMiddleware(jwtService *services.JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}

		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			c.Next()
			return
		}

		token := tokenParts[1]
		claims, err := jwtService.ValidateToken(token)
		if err != nil {
			c.Next()
			return
		}

		// Adicionar informações do usuário ao contexto se token válido
		c.Set("user_id", claims.UserID)
		c.Set("user_email", claims.Email)
		c.Set("user_role_id", claims.RoleID)
		c.Set("user_empresa_id", claims.EmpresaID)

		c.Next()
	}
}
