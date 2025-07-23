package services

import (
	"fmt"
	"log"
	"net/smtp"
	"os"
)

type EmailService struct {
	smtpHost     string
	smtpPort     string
	smtpUser     string
	smtpPassword string
	fromEmail    string
	fromName     string
}

func NewEmailService() *EmailService {
	return &EmailService{
		smtpHost:     os.Getenv("SMTP_HOST"),
		smtpPort:     os.Getenv("SMTP_PORT"),
		smtpUser:     os.Getenv("SMTP_USER"),
		smtpPassword: os.Getenv("SMTP_PASSWORD"),
		fromEmail:    os.Getenv("FROM_EMAIL"),
		fromName:     os.Getenv("FROM_NAME"),
	}
}

// SendWelcomeEmail envia email de boas-vindas
func (s *EmailService) SendWelcomeEmail(to, name string) error {
	subject := "Bem-vindo(a) ao Sistema de Telemetria!"
	body := fmt.Sprintf(`
		<html>
		<body>
			<h2>Olá, %s!</h2>
			<p>Seja bem-vindo(a) ao nosso sistema de telemetria.</p>
			<p>Sua conta foi criada com sucesso e você já pode começar a usar o sistema.</p>
			<p>Se tiver alguma dúvida, não hesite em entrar em contato conosco.</p>
			<br>
			<p>Atenciosamente,<br>Equipe Gestão Telemetria</p>
		</body>
		</html>
	`, name)

	return s.sendEmail(to, subject, body)
}

// SendEmailConfirmation envia email de confirmação
func (s *EmailService) SendEmailConfirmation(to, name, token string) error {
	subject := "Confirme seu email - Sistema de Telemetria"
	confirmationURL := fmt.Sprintf("%s/confirm-email?token=%s", os.Getenv("FRONTEND_URL"), token)

	body := fmt.Sprintf(`
		<html>
		<body>
			<h2>Olá, %s!</h2>
			<p>Obrigado por se cadastrar em nosso sistema!</p>
			<p>Para ativar sua conta, clique no link abaixo:</p>
			<p><a href="%s" style="background-color: #4CAF50; color: white; padding: 15px 32px; text-decoration: none; display: inline-block; font-size: 16px; margin: 4px 2px; cursor: pointer; border-radius: 4px;">Confirmar Email</a></p>
			<p>Ou copie e cole este link em seu navegador: %s</p>
			<p>Este link expira em 24 horas.</p>
			<br>
			<p>Se você não se cadastrou em nosso sistema, ignore este email.</p>
			<br>
			<p>Atenciosamente,<br>Equipe Gestão Telemetria</p>
		</body>
		</html>
	`, name, confirmationURL, confirmationURL)

	return s.sendEmail(to, subject, body)
}

// SendInviteEmail envia convite por email
func (s *EmailService) SendInviteEmail(to, name, inviterName, companyName, token string) error {
	subject := fmt.Sprintf("Convite para participar da %s - Sistema de Telemetria", companyName)
	inviteURL := fmt.Sprintf("%s/accept-invite?token=%s", os.Getenv("FRONTEND_URL"), token)

	body := fmt.Sprintf(`
		<html>
		<body>
			<h2>Olá!</h2>
			<p>Você foi convidado(a) por <strong>%s</strong> para participar da empresa <strong>%s</strong> em nosso sistema de telemetria.</p>
			<p>Para aceitar o convite e criar sua conta, clique no link abaixo:</p>
			<p><a href="%s" style="background-color: #2196F3; color: white; padding: 15px 32px; text-decoration: none; display: inline-block; font-size: 16px; margin: 4px 2px; cursor: pointer; border-radius: 4px;">Aceitar Convite</a></p>
			<p>Ou copie e cole este link em seu navegador: %s</p>
			<p>Este convite expira em 7 dias.</p>
			<br>
			<p>Se você não esperava este convite, ignore este email.</p>
			<br>
			<p>Atenciosamente,<br>Equipe Gestão Telemetria</p>
		</body>
		</html>
	`, inviterName, companyName, inviteURL, inviteURL)

	return s.sendEmail(to, subject, body)
}

// SendRegistrationRequestNotification notifica gestores sobre nova solicitação
func (s *EmailService) SendRegistrationRequestNotification(userEmail, userName, companyName string) error {
	subject := fmt.Sprintf("Nova solicitação de cadastro - %s", companyName)
	adminURL := fmt.Sprintf("%s/admin/registration-requests", os.Getenv("FRONTEND_URL"))

	body := fmt.Sprintf(`
		<html>
		<body>
			<h2>Nova Solicitação de Cadastro</h2>
			<p>Uma nova solicitação de cadastro foi recebida:</p>
			<ul>
				<li><strong>Nome:</strong> %s</li>
				<li><strong>Email:</strong> %s</li>
				<li><strong>Empresa:</strong> %s</li>
			</ul>
			<p>Para revisar e aprovar/rejeitar esta solicitação, acesse o painel administrativo:</p>
			<p><a href="%s" style="background-color: #FF9800; color: white; padding: 15px 32px; text-decoration: none; display: inline-block; font-size: 16px; margin: 4px 2px; cursor: pointer; border-radius: 4px;">Revisar Solicitação</a></p>
			<br>
			<p>Atenciosamente,<br>Sistema de Telemetria</p>
		</body>
		</html>
	`, userName, userEmail, companyName, adminURL)

	// TODO: Buscar emails dos gestores da empresa
	// Por enquanto, enviando para um email administrativo padrão
	adminEmail := os.Getenv("ADMIN_EMAIL")
	if adminEmail == "" {
		adminEmail = "admin@gestaotelemetria.com"
	}

	return s.sendEmail(adminEmail, subject, body)
}

// SendPasswordResetEmail envia email de reset de senha
func (s *EmailService) SendPasswordResetEmail(to, name, token string) error {
	subject := "Redefinir Senha - Sistema de Telemetria"
	resetURL := fmt.Sprintf("%s/reset-password?token=%s", os.Getenv("FRONTEND_URL"), token)

	body := fmt.Sprintf(`
		<html>
		<body>
			<h2>Olá, %s!</h2>
			<p>Recebemos uma solicitação para redefinir sua senha.</p>
			<p>Para criar uma nova senha, clique no link abaixo:</p>
			<p><a href="%s" style="background-color: #f44336; color: white; padding: 15px 32px; text-decoration: none; display: inline-block; font-size: 16px; margin: 4px 2px; cursor: pointer; border-radius: 4px;">Redefinir Senha</a></p>
			<p>Ou copie e cole este link em seu navegador: %s</p>
			<p>Este link expira em 1 hora.</p>
			<br>
			<p>Se você não solicitou a redefinição de senha, ignore este email.</p>
			<br>
			<p>Atenciosamente,<br>Equipe Gestão Telemetria</p>
		</body>
		</html>
	`, name, resetURL, resetURL)

	return s.sendEmail(to, subject, body)
}

// sendEmail método privado para enviar emails
func (s *EmailService) sendEmail(to, subject, body string) error {
	// Se não tiver configuração SMTP, apenas loga
	if s.smtpHost == "" {
		log.Printf("EMAIL (SMTP não configurado): Para: %s, Assunto: %s", to, subject)
		return nil
	}

	from := fmt.Sprintf("%s <%s>", s.fromName, s.fromEmail)

	headers := make(map[string]string)
	headers["From"] = from
	headers["To"] = to
	headers["Subject"] = subject
	headers["MIME-Version"] = "1.0"
	headers["Content-Type"] = "text/html; charset=UTF-8"

	message := ""
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + body

	auth := smtp.PlainAuth("", s.smtpUser, s.smtpPassword, s.smtpHost)

	err := smtp.SendMail(
		s.smtpHost+":"+s.smtpPort,
		auth,
		s.fromEmail,
		[]string{to},
		[]byte(message),
	)

	if err != nil {
		log.Printf("Erro ao enviar email para %s: %v", to, err)
		return err
	}

	log.Printf("Email enviado com sucesso para: %s", to)
	return nil
}
