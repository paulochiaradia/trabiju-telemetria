package middleware

import (
	"database/sql"
	"errors"
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

		roleID := userRoleID.(int)

		// Verificar se o role está permitido
		allowed := false
		for _, role := range allowedRoles {
			switch role {
			case "admin":
				// Admin seria role_id = 3 (gestor) - máximo nível
				if roleID == 3 {
					allowed = true
				}
			case "gestor":
				// Gestor também é role_id = 3
				if roleID == 3 {
					allowed = true
				}
			case "entregador":
				// Entregador é role_id = 1, mas gestor também pode
				if roleID == 1 || roleID == 3 {
					allowed = true
				}
			case "ajudante":
				// Qualquer role pode fazer o que ajudante faz
				if roleID == 1 || roleID == 2 || roleID == 3 {
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

// CompanyMiddleware middleware para verificar se usuário pertence à empresa ativa
func CompanyMiddleware(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userEmpresaID, exists := c.Get("user_empresa_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Usuário não autenticado",
			})
			c.Abort()
			return
		}

		empresaID := userEmpresaID.(int)

		// Verificar se a empresa existe e está ativa
		var empresa struct {
			ID    int    `json:"id"`
			Ativa bool   `json:"ativa"`
			Nome  string `json:"nome"`
		}

		query := "SELECT id, ativa, nome FROM empresas WHERE id = ?"
		err := db.QueryRow(query, empresaID).Scan(&empresa.ID, &empresa.Ativa, &empresa.Nome)

		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				c.JSON(http.StatusForbidden, gin.H{
					"error": "Empresa não encontrada",
				})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "Erro ao verificar empresa",
				})
			}
			c.Abort()
			return
		}

		// Verificar se a empresa está ativa
		if !empresa.Ativa {
			c.JSON(http.StatusForbidden, gin.H{
				"error":   "Empresa inativa",
				"message": "Sua empresa foi desativada. Entre em contato com o suporte.",
			})
			c.Abort()
			return
		}

		// Adicionar informações da empresa ao contexto
		c.Set("empresa_id", empresa.ID)
		c.Set("empresa_nome", empresa.Nome)
		c.Set("empresa_ativa", empresa.Ativa)

		c.Next()
	}
}

// OptionalAuthMiddleware middleware de autenticação opcional
// Útil para endpoints que funcionam tanto para usuários autenticados quanto anônimos
func OptionalAuthMiddleware(jwtService *services.JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			// Definir que é usuário anônimo
			c.Set("is_authenticated", false)
			c.Next()
			return
		}

		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			// Token mal formado - continuar como anônimo
			c.Set("is_authenticated", false)
			c.Next()
			return
		}

		token := tokenParts[1]
		claims, err := jwtService.ValidateToken(token)
		if err != nil {
			// Token inválido - continuar como anônimo
			c.Set("is_authenticated", false)
			c.Next()
			return
		}

		// Token válido - adicionar informações do usuário
		c.Set("is_authenticated", true)
		c.Set("user_id", claims.UserID)
		c.Set("user_email", claims.Email)
		c.Set("user_role_id", claims.RoleID)
		c.Set("user_empresa_id", claims.EmpresaID)

		c.Next()
	}
}

// Helper function para verificar se usuário está autenticado no contexto
func IsAuthenticated(c *gin.Context) bool {
	isAuth, exists := c.Get("is_authenticated")
	if !exists {
		// Se não existe, verifica se tem user_id (AuthMiddleware obrigatório)
		_, hasUserID := c.Get("user_id")
		return hasUserID
	}
	return isAuth.(bool)
}

// Helper function para obter ID do usuário autenticado
func GetUserID(c *gin.Context) (int, bool) {
	userID, exists := c.Get("user_id")
	if !exists {
		return 0, false
	}
	return userID.(int), true
}

// Helper function para obter ID da empresa do usuário autenticado
func GetUserEmpresaID(c *gin.Context) (int, bool) {
	empresaID, exists := c.Get("user_empresa_id")
	if !exists {
		return 0, false
	}
	return empresaID.(int), true
}
