package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/paulochiaradia/trabiju-telemetria/internal/models"
	"github.com/paulochiaradia/trabiju-telemetria/internal/services"
	"github.com/paulochiaradia/trabiju-telemetria/internal/validation"
)

// AuthController gerencia autenticação e cadastro
type AuthController struct {
	authService *services.AuthService
	jwtService  *services.JWTService
}

func NewAuthController(authService *services.AuthService, jwtService *services.JWTService) *AuthController {
	return &AuthController{
		authService: authService,
		jwtService:  jwtService,
	}
}

// CadastroComCodigo permite cadastro usando código da empresa
func (ac *AuthController) CadastroComCodigo(c *gin.Context) {
	var req models.CadastroComCodigoRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Dados inválidos",
			"details": err.Error(),
		})
		return
	}

	// Validar os dados usando o validador customizado
	if err := validation.ValidateStruct(req); err != nil {
		validationErrors := validation.GetValidationErrors(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Dados inválidos",
			"details": validationErrors,
		})
		return
	}

	user, token, err := ac.authService.RegisterWithCompanyCode(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Se user é nil, significa que foi criada uma solicitação de aprovação
	if user == nil {
		c.JSON(http.StatusAccepted, gin.H{
			"message": "Solicitação de cadastro enviada para aprovação",
			"status":  "pending_approval",
		})
		return
	}

	// Usuário criado com sucesso
	response := gin.H{
		"message": "Cadastro realizado com sucesso",
		"user": gin.H{
			"id":         user.ID,
			"nome":       user.Nome,
			"email":      user.Email,
			"role_id":    user.RoleID,
			"empresa_id": user.EmpresaID,
			"ativo":      user.Ativo,
		},
		"token": token,
	}

	if !user.Ativo {
		response["message"] = "Cadastro realizado. Verifique seu email para ativar a conta"
		response["status"] = "email_confirmation_required"
	}

	c.JSON(http.StatusCreated, response)
}

// AceitarConvite permite aceitar convite via token
func (ac *AuthController) AceitarConvite(c *gin.Context) {
	var req models.AceitarConviteRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Dados inválidos",
			"details": err.Error(),
		})
		return
	}

	// Validar os dados usando o validador customizado
	if err := validation.ValidateStruct(req); err != nil {
		validationErrors := validation.GetValidationErrors(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Dados inválidos",
			"details": validationErrors,
		})
		return
	}

	user, token, err := ac.authService.RegisterWithInvite(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Convite aceito com sucesso! Você já está logado",
		"user": gin.H{
			"id":         user.ID,
			"nome":       user.Nome,
			"email":      user.Email,
			"role_id":    user.RoleID,
			"empresa_id": user.EmpresaID,
			"ativo":      user.Ativo,
		},
		"token": token,
	})
}

// Login autentica usuário
func (ac *AuthController) Login(c *gin.Context) {
	var req struct {
		Email string `json:"email" validate:"required,email"`
		Senha string `json:"senha" validate:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Dados inválidos",
			"details": err.Error(),
		})
		return
	}

	user, token, err := ac.authService.Login(req.Email, req.Senha)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login realizado com sucesso",
		"user": gin.H{
			"id":           user.ID,
			"nome":         user.Nome,
			"email":        user.Email,
			"role_id":      user.RoleID,
			"empresa_id":   user.EmpresaID,
			"ativo":        user.Ativo,
			"ultimo_login": user.UltimoLogin,
		},
		"token": token,
	})
}

// RefreshToken renova o token de acesso
func (ac *AuthController) RefreshToken(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token" validate:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Refresh token requerido",
			"details": err.Error(),
		})
		return
	}

	// Buscar informações do usuário através do refresh token
	claims, err := ac.jwtService.ValidateToken(req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Refresh token inválido",
		})
		return
	}

	// TODO: Buscar role_id e empresa_id atualizados do banco
	// Por enquanto, usando valores do token antigo
	newToken, err := ac.jwtService.RefreshToken(req.RefreshToken, claims.RoleID, claims.EmpresaID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Erro ao renovar token",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Token renovado com sucesso",
		"token":   newToken,
	})
}

// ConfirmarEmail confirma email do usuário
func (ac *AuthController) ConfirmarEmail(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Token de confirmação requerido",
		})
		return
	}

	// Buscar usuário pelo contexto (deve estar autenticado)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Usuário não autenticado",
		})
		return
	}

	err := ac.authService.ConfirmEmail(userID.(int), token)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Erro ao confirmar email",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Email confirmado com sucesso",
	})
}

// Logout realiza logout do usuário
func (ac *AuthController) Logout(c *gin.Context) {
	// Com JWT stateless, o logout é feito no frontend removendo o token
	// Aqui podemos apenas registrar o logout ou invalidar tokens se implementarmos blacklist

	c.JSON(http.StatusOK, gin.H{
		"message": "Logout realizado com sucesso",
	})
}

// GetProfile retorna perfil do usuário autenticado
func (ac *AuthController) GetProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Usuário não autenticado",
		})
		return
	}

	// TODO: Buscar dados completos do usuário do banco
	// Por enquanto, retornar dados do token
	c.JSON(http.StatusOK, gin.H{
		"user": gin.H{
			"id":         userID,
			"email":      c.GetString("user_email"),
			"role_id":    c.GetInt("user_role_id"),
			"empresa_id": c.GetInt("user_empresa_id"),
		},
	})
}
