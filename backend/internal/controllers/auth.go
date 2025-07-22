package controllers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/paulochiaradia/trabiju-telemetria/internal/models"
)

// AuthController gerencia autenticação e cadastro
type AuthController struct {
	// Aqui você injetaria os services/repositories necessários
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

	// TODO: Implementar validação do código da empresa
	// TODO: Verificar se email já existe
	// TODO: Validar CPF

	// Fluxo baseado no papel desejado
	switch req.RoleDesejado {
	case "entregador", "ajudante":
		// Estes papéis podem se cadastrar diretamente (dependendo da config da empresa)
		// ou criar uma solicitação para aprovação
		ac.processarCadastroBasico(c, req)
	case "gestor":
		// Gestores sempre precisam de aprovação
		ac.criarSolicitacaoAprovacao(c, req)
	default:
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Papel inválido. Use: entregador, ajudante ou gestor",
		})
		return
	}
}

func (ac *AuthController) processarCadastroBasico(c *gin.Context, req models.CadastroComCodigoRequest) {
	// TODO: Implementar lógica de cadastro direto ou solicitação

	// Exemplo de resposta para cadastro direto
	c.JSON(http.StatusCreated, gin.H{
		"message": "Usuário cadastrado com sucesso",
		"user_id": 123, // ID do usuário criado
		"status":  "ativo",
	})
}

func (ac *AuthController) criarSolicitacaoAprovacao(c *gin.Context, req models.CadastroComCodigoRequest) {
	// TODO: Criar registro na tabela solicitacoes_cadastro

	c.JSON(http.StatusAccepted, gin.H{
		"message":        "Solicitação enviada para aprovação",
		"solicitacao_id": 456,
		"status":         "pendente",
		"info":           "Um gestor analisará sua solicitação em até 24 horas",
	})
}

// AceitarConvite permite aceitar convite por token
func (ac *AuthController) AceitarConvite(c *gin.Context) {
	var req models.AceitarConviteRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Dados inválidos",
			"details": err.Error(),
		})
		return
	}

	// TODO: Validar token do convite
	// TODO: Verificar se não expirou
	// TODO: Criar usuário com papel definido no convite

	c.JSON(http.StatusCreated, gin.H{
		"message": "Convite aceito com sucesso",
		"user_id": 789,
		"role":    "entregador",
	})
}

// CriarConvite (apenas para gestores)
func (ac *AuthController) CriarConvite(c *gin.Context) {
	var req models.CriarConviteRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Dados inválidos",
			"details": err.Error(),
		})
		return
	}

	// TODO: Verificar se usuário logado é gestor
	// TODO: Gerar token único
	// TODO: Salvar convite no banco
	// TODO: Enviar email com link de convite

	expiraEm := time.Now().AddDate(0, 0, req.ValidoPorDias)

	c.JSON(http.StatusCreated, gin.H{
		"message":      "Convite criado com sucesso",
		"convite_id":   321,
		"email":        req.Email,
		"expira_em":    expiraEm,
		"link_convite": "https://app.gestaotelemetria.com/aceitar-convite?token=abc123",
	})
}

// VerificarCodigoEmpresa valida se código da empresa existe
func (ac *AuthController) VerificarCodigoEmpresa(c *gin.Context) {
	//codigo := c.Param("codigo")

	// TODO: Buscar empresa por código
	// TODO: Verificar se está ativa

	c.JSON(http.StatusOK, gin.H{
		"valido": true,
		"empresa": gin.H{
			"nome":                  "Gestão Telemetria LTDA",
			"permite_auto_cadastro": true,
			"papeis_disponiveis":    []string{"entregador", "ajudante"},
		},
	})
}

// ListarSolicitacoesPendentes (apenas para gestores)
func (ac *AuthController) ListarSolicitacoesPendentes(c *gin.Context) {
	// TODO: Verificar se usuário é gestor
	// TODO: Buscar solicitações pendentes da empresa do gestor

	c.JSON(http.StatusOK, gin.H{
		"solicitacoes": []gin.H{
			{
				"id":               1,
				"nome":             "João Silva",
				"email":            "joao@email.com",
				"papel_solicitado": "gestor",
				"justificativa":    "Preciso gerenciar a equipe de entregas",
				"created_at":       "2025-07-22T10:00:00Z",
			},
		},
	})
}

// AprovarSolicitacao (apenas para gestores)
func (ac *AuthController) AprovarSolicitacao(c *gin.Context) {
	solicitacaoID := c.Param("id")

	var req struct {
		Aprovado    bool   `json:"aprovado"`
		Observacoes string `json:"observacoes,omitempty"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos"})
		return
	}

	// TODO: Verificar se usuário é gestor
	// TODO: Atualizar status da solicitação
	// TODO: Se aprovado, criar usuário
	// TODO: Enviar email de notificação

	status := "rejeitado"
	if req.Aprovado {
		status = "aprovado"
	}

	c.JSON(http.StatusOK, gin.H{
		"message":        "Solicitação processada",
		"solicitacao_id": solicitacaoID,
		"status":         status,
	})
}
