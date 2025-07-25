package services

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/paulochiaradia/trabiju-telemetria/internal/models"
	"github.com/paulochiaradia/trabiju-telemetria/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo         *repository.UserRepository
	inviteRepo       *repository.InviteRepository
	companyRepo      *repository.CompanyRepository
	registrationRepo *repository.RegistrationRequestRepository
	roleRepo         *repository.RoleRepository
	emailService     *EmailService
	jwtService       *JWTService
}

func NewAuthService(
	userRepo *repository.UserRepository,
	inviteRepo *repository.InviteRepository,
	companyRepo *repository.CompanyRepository,
	registrationRepo *repository.RegistrationRequestRepository,
	roleRepo *repository.RoleRepository,
	emailService *EmailService,
	jwtService *JWTService,
) *AuthService {
	return &AuthService{
		userRepo:         userRepo,
		inviteRepo:       inviteRepo,
		companyRepo:      companyRepo,
		registrationRepo: registrationRepo,
		roleRepo:         roleRepo,
		emailService:     emailService,
		jwtService:       jwtService,
	}
}

// RegisterWithInvite registra usuário usando token de convite
func (s *AuthService) RegisterWithInvite(req models.AceitarConviteRequest) (*models.Usuario, string, error) {
	// 1. Validar convite
	invite, err := s.inviteRepo.GetInviteByToken(req.Token)
	if err != nil {
		return nil, "", err
	}
	if invite == nil {
		return nil, "", errors.New("convite inválido ou expirado")
	}

	// 2. Verificar se email já existe
	existingUser, err := s.userRepo.GetUserByEmail(invite.Email)
	if err != nil {
		return nil, "", err
	}
	if existingUser != nil {
		return nil, "", errors.New("email já está em uso")
	}

	// 3. Verificar CPF
	cpfExists, err := s.userRepo.CPFExists(req.CPF)
	if err != nil {
		return nil, "", err
	}
	if cpfExists {
		return nil, "", errors.New("CPF já está em uso")
	}

	// 4. Criptografar senha
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Senha), bcrypt.DefaultCost)
	if err != nil {
		return nil, "", err
	}

	// 5. Criar usuário
	defaultConfig := "{}"
	apiToken := s.generateAPIToken()
	user := &models.Usuario{
		Nome:                   req.Nome,
		Email:                  invite.Email,
		Senha:                  string(hashedPassword),
		Telefone:               req.Telefone,
		CPF:                    req.CPF,
		RoleID:                 invite.RoleID,
		EmpresaID:              invite.EmpresaID,
		Ativo:                  true, // Ativo direto por convite
		ConfiguracoesDashboard: &defaultConfig,
		APIToken:               &apiToken,
		SenhaAlteradaEm:        time.Now(),
	}

	err = s.userRepo.CreateUser(user)
	if err != nil {
		return nil, "", err
	}

	// 6. Marcar convite como usado
	err = s.inviteRepo.MarkInviteAsUsed(invite.ID)
	if err != nil {
		return nil, "", err
	}

	// 7. Atualizar último login
	err = s.userRepo.UpdateLastLogin(user.ID)
	if err != nil {
		return nil, "", err
	}

	// 8. Gerar JWT
	token, err := s.jwtService.GenerateToken(user.ID, user.Email, user.RoleID, user.EmpresaID)
	if err != nil {
		return nil, "", err
	}

	// 9. Enviar email de boas-vindas
	go s.emailService.SendWelcomeEmail(user.Email, user.Nome)

	return user, token, nil
}

// RegisterWithCompanyCode registra usuário usando código da empresa
func (s *AuthService) RegisterWithCompanyCode(req models.CadastroComCodigoRequest) (*models.Usuario, string, error) {
	// 1. Validar empresa
	company, err := s.companyRepo.GetCompanyByInviteCode(req.CodigoEmpresa)
	if err != nil {
		return nil, "", err
	}
	if company == nil {
		return nil, "", errors.New("código da empresa inválido")
	}

	// 2. Verificar se email já existe
	existingUser, err := s.userRepo.GetUserByEmail(req.Email)
	if err != nil {
		return nil, "", err
	}
	if existingUser != nil {
		return nil, "", errors.New("email já está em uso")
	}

	// 3. Verificar CPF
	cpfExists, err := s.userRepo.CPFExists(req.CPF)
	if err != nil {
		return nil, "", err
	}
	if cpfExists {
		return nil, "", errors.New("CPF já está em uso")
	}

	// 4. Buscar role
	role, err := s.roleRepo.GetRoleByName(req.RoleDesejado)
	if err != nil {
		return nil, "", err
	}
	if role == nil {
		return nil, "", errors.New("papel inválido")
	}

	// 5. Verificar se role permite cadastro direto
	needsApproval := s.needsApproval(req.RoleDesejado, company)

	if needsApproval {
		return s.createRegistrationRequest(req, company, role)
	}

	// 6. Cadastro direto - criptografar senha
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Senha), bcrypt.DefaultCost)
	if err != nil {
		return nil, "", err
	}

	// 7. Criar usuário (inativo até confirmar email)
	defaultConfig := "{}"
	apiToken := s.generateAPIToken()
	user := &models.Usuario{
		Nome:                   req.Nome,
		Email:                  req.Email,
		Senha:                  string(hashedPassword),
		Telefone:               req.Telefone,
		CPF:                    req.CPF,
		RoleID:                 role.ID,
		EmpresaID:              company.ID,
		Ativo:                  false, // Inativo até confirmar email
		ConfiguracoesDashboard: &defaultConfig,
		APIToken:               &apiToken,
		SenhaAlteradaEm:        time.Now(),
	}

	err = s.userRepo.CreateUser(user)
	if err != nil {
		return nil, "", err
	}

	// 8. Gerar token de confirmação de email
	confirmationToken, err := s.generateEmailConfirmationToken(user.ID, user.Email)
	if err != nil {
		return nil, "", err
	}

	// 9. Enviar email de confirmação
	go s.emailService.SendEmailConfirmation(user.Email, user.Nome, confirmationToken)

	// 10. Gerar JWT (mesmo inativo, para poder confirmar email)
	token, err := s.jwtService.GenerateToken(user.ID, user.Email, user.RoleID, user.EmpresaID)
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}

// Login autentica usuário
func (s *AuthService) Login(email, password string) (*models.Usuario, string, error) {
	// 1. Buscar usuário
	user, err := s.userRepo.GetUserByEmail(email)
	if err != nil {
		return nil, "", err
	}
	if user == nil {
		return nil, "", errors.New("credenciais inválidas")
	}

	// 2. Verificar se está bloqueado
	if user.BloqueadoAte != nil && user.BloqueadoAte.After(time.Now()) {
		return nil, "", fmt.Errorf("usuário bloqueado até %s", user.BloqueadoAte.Format("02/01/2006 15:04"))
	}

	// 3. Verificar senha
	err = bcrypt.CompareHashAndPassword([]byte(user.Senha), []byte(password))
	if err != nil {
		// Incrementar tentativas de login
		attempts := user.TentativasLogin + 1
		s.userRepo.UpdateLoginAttempts(user.ID, attempts)

		// Bloquear após 5 tentativas
		if attempts >= 5 {
			blockUntil := time.Now().Add(30 * time.Minute)
			s.userRepo.BlockUser(user.ID, blockUntil)
			return nil, "", errors.New("muitas tentativas incorretas. Usuário bloqueado por 30 minutos")
		}

		return nil, "", errors.New("credenciais inválidas")
	}

	// 4. Verificar se está ativo
	if !user.Ativo {
		return nil, "", errors.New("usuário inativo. Confirme seu email primeiro")
	}

	// 5. Reset tentativas e atualizar último login
	s.userRepo.UpdateLoginAttempts(user.ID, 0)
	s.userRepo.UpdateLastLogin(user.ID)

	// 6. Gerar JWT
	token, err := s.jwtService.GenerateToken(user.ID, user.Email, user.RoleID, user.EmpresaID)
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}

// ConfirmEmail confirma email do usuário usando token JWT
func (s *AuthService) ConfirmEmail(token string) (*models.Usuario, error) {
	// Validar e extrair informações do token de confirmação
	userID, email, err := s.validateEmailConfirmationToken(token)
	if err != nil {
		return nil, errors.New("token de confirmação inválido ou expirado")
	}

	// Buscar usuário para verificar se ainda existe e se precisa de confirmação
	user, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("usuário não encontrado")
	}

	// Verificar se o email do token confere com o do usuário
	if user.Email != email {
		return nil, errors.New("token inválido para este usuário")
	}

	// Verificar se o usuário já está ativo
	if user.Ativo {
		return user, nil // Já confirmado, retorna sucesso
	}

	// Ativar usuário
	err = s.userRepo.ActivateUser(userID)
	if err != nil {
		return nil, err
	}

	// Buscar usuário atualizado
	user, err = s.userRepo.GetUserByID(userID)
	if err != nil {
		return nil, err
	}

	// Enviar email de boas-vindas
	go s.emailService.SendWelcomeEmail(user.Email, user.Nome)

	return user, nil
}

// Helper functions
func (s *AuthService) needsApproval(role string, company *models.Empresa) bool {
	// Gestores sempre precisam de aprovação
	if role == "gestor" {
		return true
	}

	// TODO: Verificar configurações da empresa
	// Por enquanto, apenas gestores precisam de aprovação
	return false
}

func (s *AuthService) createRegistrationRequest(req models.CadastroComCodigoRequest, company *models.Empresa, role *models.Role) (*models.Usuario, string, error) {
	// Criar solicitação de aprovação
	registrationReq := &models.SolicitacaoCadastro{
		Nome:           req.Nome,
		Email:          req.Email,
		Telefone:       req.Telefone,
		CPF:            req.CPF,
		RoleSolicitado: role.ID,
		EmpresaID:      company.ID,
		CodigoUsado:    req.CodigoEmpresa,
		Justificativa:  req.Justificativa,
	}

	err := s.registrationRepo.CreateRegistrationRequest(registrationReq)
	if err != nil {
		return nil, "", err
	}

	// Enviar email de notificação para gestores
	go s.emailService.SendRegistrationRequestNotification(req.Email, req.Nome, company.Nome)

	// Retornar usuário vazio (ainda não criado) e token vazio
	return nil, "", nil
}

func (s *AuthService) generateAPIToken() string {
	bytes := make([]byte, 32)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

func (s *AuthService) generateConfirmationToken() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

// generateEmailConfirmationToken gera um JWT para confirmação de email
func (s *AuthService) generateEmailConfirmationToken(userID int, email string) (string, error) {
	// Criar claims customizadas para confirmação de email
	claims := map[string]interface{}{
		"user_id": userID,
		"email":   email,
		"type":    "email_confirmation",
		"exp":     time.Now().Add(24 * time.Hour).Unix(), // Expira em 24 horas
		"iat":     time.Now().Unix(),
	}

	// Gerar token usando o JWT service
	return s.jwtService.GenerateCustomToken(claims)
}

// validateEmailConfirmationToken valida token de confirmação de email
func (s *AuthService) validateEmailConfirmationToken(tokenString string) (int, string, error) {
	// Validar token usando o JWT service
	claims, err := s.jwtService.ValidateCustomToken(tokenString)
	if err != nil {
		return 0, "", err
	}

	// Verificar se é um token de confirmação de email
	tokenType, ok := claims["type"].(string)
	if !ok || tokenType != "email_confirmation" {
		return 0, "", errors.New("tipo de token inválido")
	}

	// Extrair user_id
	userIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		return 0, "", errors.New("user_id inválido no token")
	}
	userID := int(userIDFloat)

	// Extrair email
	email, ok := claims["email"].(string)
	if !ok {
		return 0, "", errors.New("email inválido no token")
	}

	return userID, email, nil
}
