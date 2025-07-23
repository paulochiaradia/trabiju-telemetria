package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/paulochiaradia/trabiju-telemetria/internal/config"
	"github.com/paulochiaradia/trabiju-telemetria/internal/database"
	"github.com/paulochiaradia/trabiju-telemetria/internal/middleware"
	"github.com/paulochiaradia/trabiju-telemetria/internal/routes"
)

// main inicializa o roteador Gin e registra as rotas da aplicação.
// Em seguida, inicia o servidor na porta 8080.
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

	log.Println("✅ Conexão com banco de dados estabelecida com sucesso!")
	log.Printf("🚀 Servidor iniciando na porta %s", cfg.ServerPort)

	// Configurar Gin com middleware customizado (sem logs de health check)
	r := gin.New()
	r.Use(middleware.CustomLogger())
	r.Use(gin.Recovery())

	// Configurar proxies confiáveis para segurança
	if len(cfg.TrustedProxies) > 0 {
		err = r.SetTrustedProxies(cfg.TrustedProxies)
		if err != nil {
			log.Fatal("Erro ao configurar proxies confiáveis:", err)
		}
	}

	// Registrar rotas
	routes.RegisterRoutes(r)

	// Iniciar servidor
	r.Run(":" + cfg.ServerPort)
}
