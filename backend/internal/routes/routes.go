package routes

import (
	"database/sql"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/paulochiaradia/trabiju-telemetria/internal/controllers"
	"github.com/paulochiaradia/trabiju-telemetria/internal/middleware"
	"github.com/paulochiaradia/trabiju-telemetria/internal/repository"
	"github.com/paulochiaradia/trabiju-telemetria/internal/services"
)

// RegisterRoutes configura as rotas básicas da aplicação (compatibilidade).
func RegisterRoutes(r *gin.Engine) {
	// Rota básica de teste
	r.GET("/ping", controllers.PingHandler)

	// Controller para testes de banco de dados
	dbController := &controllers.DatabaseController{}

	// Grupo de rotas para testes de banco
	dbGroup := r.Group("/database")
	{
		dbGroup.GET("/test", dbController.TestConnection)        // Testar conexão
		dbGroup.GET("/tables", dbController.ListTables)          // Listar tabelas
		dbGroup.GET("/table/:table", dbController.DescribeTable) // Descrever tabela
	}
}

// SetupRoutes configura as rotas completas da aplicação com autenticação
func SetupRoutes(router *gin.Engine, db *sql.DB) {
	// Configurar middlewares globais
	router.Use(middleware.CustomLogger())
	router.Use(middleware.CORSMiddleware())
	//	router.Use(middleware.SecurityMiddleware())
	//	router.Use(middleware.RateLimitMiddleware()) // Rate limiting global
	router.Use(gin.Recovery())

	// Inicializar repositórios
	userRepo := repository.NewUserRepository(db)
	inviteRepo := repository.NewInviteRepository(db)
	companyRepo := repository.NewCompanyRepository(db)
	registrationRepo := repository.NewRegistrationRequestRepository(db)
	roleRepo := repository.NewRoleRepository(db)

	// Inicializar serviços
	emailService := services.NewEmailService()
	jwtService := services.NewJWTService("your-secret-key") // TODO: Mover para variável de ambiente
	authService := services.NewAuthService(
		userRepo, inviteRepo, companyRepo, registrationRepo,
		roleRepo, emailService, jwtService,
	)

	// Inicializar controllers
	authController := controllers.NewAuthController(authService, jwtService)
	dbController := &controllers.DatabaseController{}
	healthController := controllers.NewHealthController(db)

	// Rotas públicas com rate limiting padrão
	public := router.Group("/api/v1")
	{
		// Saúde da aplicação (sem rate limit adicional)
		public.GET("/ping", controllers.PingHandler)
		public.GET("/health", healthController.HealthCheck) // Health check completo

		// Confirmação de email (sem auth para permitir clique direto do email)
		public.GET("/auth/confirm-email", authController.ConfirmarEmail)
	}

	// Rotas híbridas (funcionam com ou sem autenticação)
	hybrid := router.Group("/api/v1")
	hybrid.Use(middleware.OptionalAuthMiddleware(jwtService))
	{
		// Exemplo: Dashboard público com dados extras para usuários logados
		// hybrid.GET("/public/dashboard", handlers.PublicDashboard)

		// Exemplo: Informações de empresa (dados básicos públicos, detalhes para usuários logados)
		// hybrid.GET("/public/company-info", handlers.CompanyInfo)
	}

	// Rotas de autenticação com rate limiting mais restritivo
	authGroup := router.Group("/api/v1")
	authGroup.Use(middleware.RateLimitWithConfig(20, time.Minute)) // 20 req/min para auth
	{
		// Endpoints sensíveis de autenticação
		authGroup.POST("/auth/login", authController.Login)
		authGroup.POST("/auth/register/code", authController.CadastroComCodigo)
		authGroup.POST("/auth/invite/accept", authController.AceitarConvite)
		authGroup.POST("/auth/refresh", authController.RefreshToken)
	}

	// Rotas protegidas (requer autenticação)
	protected := router.Group("/api/v1")
	protected.Use(middleware.AuthMiddleware(jwtService))
	{
		// Perfil do usuário
		protected.GET("/auth/profile", authController.GetProfile)
		protected.POST("/auth/logout", authController.Logout)
	}

	// Rotas que requerem empresa ativa (operações de negócio)
	companyProtected := router.Group("/api/v1")
	companyProtected.Use(middleware.AuthMiddleware(jwtService))
	companyProtected.Use(middleware.CompanyMiddleware(db))
	{
		// Database operations (apenas para admin/gestor)
		admin := companyProtected.Group("/admin")
		admin.Use(middleware.RoleMiddleware("admin", "gestor"))
		{
			admin.GET("/db/test", dbController.TestConnection)
			admin.GET("/db/tables", dbController.ListTables)
			admin.GET("/db/table/:table", dbController.DescribeTable)
		}

		// Futuras rotas de telemetria e veículos devem usar companyProtected
		// telemetry := companyProtected.Group("/telemetry")
		// vehicles := companyProtected.Group("/vehicles")
	}
}
