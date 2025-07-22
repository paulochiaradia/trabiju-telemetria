package models

import (
	"time"
)

// Empresa representa uma empresa/organização
type Empresa struct {
	ID            int       `json:"id" db:"id"`
	Nome          string    `json:"nome" db:"nome"`
	CNPJ          string    `json:"cnpj" db:"cnpj"`
	CodigoConvite string    `json:"codigo_convite" db:"codigo_convite"`
	Ativa         bool      `json:"ativa" db:"ativa"`
	Configuracoes string    `json:"configuracoes" db:"configuracoes"` // JSON
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
}

// Usuario representa um usuário no sistema
type Usuario struct {
	ID                     int        `json:"id" db:"id"`
	Nome                   string     `json:"nome" db:"nome"`
	Email                  string     `json:"email" db:"email"`
	Senha                  string     `json:"-" db:"senha"` // Nunca retorna senha no JSON
	Telefone               string     `json:"telefone" db:"telefone"`
	CPF                    string     `json:"cpf" db:"cpf"`
	Avatar                 string     `json:"avatar" db:"avatar"`
	RoleID                 int        `json:"role_id" db:"role_id"`
	EmpresaID              int        `json:"empresa_id" db:"empresa_id"`
	Ativo                  bool       `json:"ativo" db:"ativo"`
	UltimoLogin            *time.Time `json:"ultimo_login" db:"ultimo_login"`
	ConfiguracoesDashboard string     `json:"configuracoes_dashboard" db:"configuracoes_dashboard"` // JSON
	APIToken               string     `json:"api_token" db:"api_token"`
	TentativasLogin        int        `json:"tentativas_login" db:"tentativas_login"`
	BloqueadoAte           *time.Time `json:"bloqueado_ate" db:"bloqueado_ate"`
	SenhaAlteradaEm        time.Time  `json:"senha_alterada_em" db:"senha_alterada_em"`
	CreatedAt              time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt              time.Time  `json:"updated_at" db:"updated_at"`

	// Relacionamentos
	Role    *Role    `json:"role,omitempty"`
	Empresa *Empresa `json:"empresa,omitempty"`
}

// Role representa os papéis/permissões de usuário
type Role struct {
	ID         int       `json:"id" db:"id"`
	Nome       string    `json:"nome" db:"nome"`
	Descricao  string    `json:"descricao" db:"descricao"`
	Permissoes string    `json:"permissoes" db:"permissoes"` // JSON
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
}

// ConviteUsuario representa convites de usuários
type ConviteUsuario struct {
	ID           int       `json:"id" db:"id"`
	Email        string    `json:"email" db:"email"`
	RoleID       int       `json:"role_id" db:"role_id"`
	EmpresaID    int       `json:"empresa_id" db:"empresa_id"`
	ConvidadoPor int       `json:"convidado_por" db:"convidado_por"`
	Token        string    `json:"token" db:"token"`
	Usado        bool      `json:"usado" db:"usado"`
	ExpiraEm     time.Time `json:"expira_em" db:"expira_em"`
	DadosConvite string    `json:"dados_convite" db:"dados_convite"` // JSON
	CreatedAt    time.Time `json:"created_at" db:"created_at"`

	// Relacionamentos
	Role             *Role    `json:"role,omitempty"`
	Empresa          *Empresa `json:"empresa,omitempty"`
	ConvidadoPorUser *Usuario `json:"convidado_por_user,omitempty"`
}

// SolicitacaoCadastro representa solicitações de cadastro
type SolicitacaoCadastro struct {
	ID                   int       `json:"id" db:"id"`
	Nome                 string    `json:"nome" db:"nome"`
	Email                string    `json:"email" db:"email"`
	Telefone             string    `json:"telefone" db:"telefone"`
	CPF                  string    `json:"cpf" db:"cpf"`
	RoleSolicitado       int       `json:"role_solicitado" db:"role_solicitado"`
	EmpresaID            int       `json:"empresa_id" db:"empresa_id"`
	CodigoUsado          string    `json:"codigo_usado" db:"codigo_usado"`
	Justificativa        string    `json:"justificativa" db:"justificativa"`
	Status               string    `json:"status" db:"status"` // pendente, aprovado, rejeitado
	AprovadoPor          *int      `json:"aprovado_por" db:"aprovado_por"`
	ObservacoesAprovacao string    `json:"observacoes_aprovacao" db:"observacoes_aprovacao"`
	CreatedAt            time.Time `json:"created_at" db:"created_at"`
	UpdatedAt            time.Time `json:"updated_at" db:"updated_at"`

	// Relacionamentos
	RoleSolicitadoObj *Role    `json:"role_solicitado_obj,omitempty"`
	Empresa           *Empresa `json:"empresa,omitempty"`
	AprovadoPorUser   *Usuario `json:"aprovado_por_user,omitempty"`
}

// Estruturas de Request/Response para API

// CadastroComCodigoRequest - Cadastro usando código da empresa
type CadastroComCodigoRequest struct {
	Nome          string `json:"nome" validate:"required,min=2,max=100"`
	Email         string `json:"email" validate:"required,email"`
	Senha         string `json:"senha" validate:"required,min=6"`
	Telefone      string `json:"telefone" validate:"required"`
	CPF           string `json:"cpf" validate:"required"`
	CodigoEmpresa string `json:"codigo_empresa" validate:"required"`
	RoleDesejado  string `json:"role_desejado" validate:"required,oneof=entregador ajudante"`
	Justificativa string `json:"justificativa,omitempty"`
}

// AceitarConviteRequest - Aceitar convite por token
type AceitarConviteRequest struct {
	Token    string `json:"token" validate:"required"`
	Nome     string `json:"nome" validate:"required,min=2,max=100"`
	Senha    string `json:"senha" validate:"required,min=6"`
	Telefone string `json:"telefone" validate:"required"`
	CPF      string `json:"cpf" validate:"required"`
}

// CriarConviteRequest - Gestor criar convite
type CriarConviteRequest struct {
	Email         string `json:"email" validate:"required,email"`
	RoleNome      string `json:"role_nome" validate:"required,oneof=entregador ajudante gestor"`
	Mensagem      string `json:"mensagem,omitempty"`
	ValidoPorDias int    `json:"valido_por_dias" validate:"min=1,max=30"`
}
