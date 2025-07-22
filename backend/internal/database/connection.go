package database

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/paulochiaradia/trabiju-telemetria/internal/config"
)

// Connection mantém a conexão com o banco de dados
type Connection struct {
	DB *sql.DB
}

// NewConnection cria uma nova conexão com o banco de dados
func NewConnection(cfg *config.Config) (*Connection, error) {
	// Usa o método do config para obter a DSN
	dsn := cfg.GetDatabaseDSN()

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("falha ao conectar com o banco de dados: %w", err)
	}

	// Configurações de pool de conexões
	configureConnectionPool(db)

	// Testa a conexão
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("falha ao testar conexão com o banco: %w", err)
	}

	return &Connection{DB: db}, nil
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
