package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/paulochiaradia/trabiju-telemetria/internal/controllers"
)

// RegisterRoutes configura as rotas da aplicação.
func RegisterRoutes(r *gin.Engine) {
	// Rota básica de teste
	r.GET("/ping", controllers.PingHandler)

	// Controller para testes de banco de dados
	dbController := &controllers.DatabaseController{}
	
	// Grupo de rotas para testes de banco
	dbGroup := r.Group("/database")
	{
		dbGroup.GET("/test", dbController.TestConnection)     // Testar conexão
		dbGroup.GET("/tables", dbController.ListTables)      // Listar tabelas
		dbGroup.GET("/table/:table", dbController.DescribeTable) // Descrever tabela
	}
}
