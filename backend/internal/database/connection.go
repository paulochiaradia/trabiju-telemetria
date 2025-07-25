package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/paulochiaradia/trabiju-telemetria/internal/config"
)

// Connection mantém a conexão com o banco de dados
type Connection struct {
	DB *sql.DB
}

// NewConnection cria uma nova conexão com o banco de dados com retry automático
func NewConnection(cfg *config.Config) (*Connection, error) {
	return NewConnectionWithRetry(cfg, 30, 2*time.Second)
}

// NewConnectionWithRetry cria uma nova conexão com retry configurável
func NewConnectionWithRetry(cfg *config.Config, maxRetries int, retryInterval time.Duration) (*Connection, error) {
	// Usa o método do config para obter a DSN
	dsn := cfg.GetDatabaseDSN()

	log.Printf("🔄 Tentando conectar ao banco de dados...")
	log.Printf("📍 Host: %s:%s", cfg.DBHost, cfg.DBPort)

	var db *sql.DB
	var err error

	for attempt := 1; attempt <= maxRetries; attempt++ {
		log.Printf("🔄 Tentativa %d/%d de conexão com o banco...", attempt, maxRetries)

		db, err = sql.Open("mysql", dsn)
		if err != nil {
			log.Printf("❌ Falha ao abrir conexão (tentativa %d): %v", attempt, err)
			time.Sleep(retryInterval)
			continue
		}

		// Configurações de pool de conexões
		configureConnectionPool(db)

		// Testa a conexão
		if err := db.Ping(); err != nil {
			log.Printf("❌ Falha no ping do banco (tentativa %d): %v", attempt, err)
			db.Close()

			if attempt < maxRetries {
				log.Printf("⏳ Aguardando %v antes da próxima tentativa...", retryInterval)
				time.Sleep(retryInterval)
			}
			continue
		}

		// Conexão bem-sucedida
		log.Printf("✅ Conexão com banco estabelecida com sucesso na tentativa %d!", attempt)
		return &Connection{DB: db}, nil
	}

	return nil, fmt.Errorf("falha ao conectar com o banco após %d tentativas: %w", maxRetries, err)
}

// configureConnectionPool configura o pool de conexões do banco
func configureConnectionPool(db *sql.DB) {
	// Máximo de conexões abertas
	db.SetMaxOpenConns(25)

	// Máximo de conexões inativas
	db.SetMaxIdleConns(25)

	// Tempo máximo que uma conexão pode ficar aberta
	db.SetConnMaxLifetime(5 * time.Minute)

	// Tempo máximo que uma conexão pode ficar inativa
	db.SetConnMaxIdleTime(5 * time.Minute)
}

// Close fecha a conexão com o banco de dados
func (c *Connection) Close() error {
	if c.DB != nil {
		return c.DB.Close()
	}
	return nil
}
