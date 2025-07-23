package routes

import (
	"database/sql"

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
	router.Use(middleware.SecurityMiddleware())
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

	// Rotas públicas (sem autenticação)
	public := router.Group("/api/v1")
	{
		// Saúde da aplicação
		public.GET("/ping", controllers.PingHandler)
		public.GET("/health", controllers.PingHandler) // Usando PingHandler como health check

		// Autenticação
		public.POST("/auth/login", authController.Login)
		public.POST("/auth/register/code", authController.CadastroComCodigo)
		public.POST("/auth/invite/accept", authController.AceitarConvite)
		public.POST("/auth/refresh", authController.RefreshToken)

		// Confirmação de email (sem auth para permitir clique direto do email)
		public.GET("/auth/confirm-email", authController.ConfirmarEmail)
	}

	// Rotas protegidas (requer autenticação)
	protected := router.Group("/api/v1")
	protected.Use(middleware.AuthMiddleware(jwtService))
	{
		// Perfil do usuário
		protected.GET("/auth/profile", authController.GetProfile)
		protected.POST("/auth/logout", authController.Logout)

		// Database operations (apenas para admin/gestor)
		admin := protected.Group("/admin")
		admin.Use(middleware.RoleMiddleware("admin", "gestor"))
		{
			admin.GET("/db/test", dbController.TestConnection)
			admin.GET("/db/tables", dbController.ListTables)
			admin.GET("/db/table/:table", dbController.DescribeTable)
		}
	}

	// Manter compatibilidade com rotas antigas (evitar duplicação)
	// RegisterRoutes(router) - Removido para evitar duplicação
}
