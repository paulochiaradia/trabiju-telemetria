package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/paulochiaradia/trabiju-telemetria/internal/config"
	"github.com/paulochiaradia/trabiju-telemetria/internal/database"
	"github.com/paulochiaradia/trabiju-telemetria/internal/routes"
)

// main inicializa o roteador Gin e registra as rotas da aplicação.
// A porta é configurável via variável de ambiente SERVER_PORT (padrão: 8080).
// Docker: Dev=8081, Prod=8082 (mapeamento externo)
func main() {
	// Carregar configurações do .env
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Erro ao carregar configurações:", err)
	}

	// Conectar ao banco de dados
	conn, err := database.NewConnection(cfg)
	if err != nil {
		log.Fatal("Erro ao conectar com o banco de dados:", err)
	}
	defer conn.Close()

	log.Printf("🚀 Servidor iniciando na porta %s", cfg.ServerPort)

	// Configurar Gin sem middleware padrão (será configurado no SetupRoutes)
	r := gin.New()
	r.Use(gin.Recovery()) // Apenas Recovery aqui

	// Configurar proxies confiáveis para segurança
	if len(cfg.TrustedProxies) > 0 {
		err = r.SetTrustedProxies(cfg.TrustedProxies)
		if err != nil {
			log.Fatal("Erro ao configurar proxies confiáveis:", err)
		}
	}

	// Registrar rotas com sistema de autenticação completo
	routes.SetupRoutes(r, conn.DB)

	// Iniciar servidor
	r.Run(":" + cfg.ServerPort)
}
