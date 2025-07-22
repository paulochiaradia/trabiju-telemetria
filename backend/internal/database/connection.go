package database

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/paulochiaradia/trabiju-telemetria/internal/config"
)

// Connection mantém a conexão com o banco de dados
type Connection struct {
	DB *sql.DB
}

// NewConnection cria uma nova conexão com o banco de dados
func NewConnection(cfg *config.Config) (*Connection, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName,
	)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("falha ao conectar com o banco de dados: %w", err)
	}

	// Testa a conexão
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("falha ao testar conexão com o banco: %w", err)
	}

	return &Connection{DB: db}, nil
}

// Close fecha a conexão com o banco de dados
func (c *Connection) Close() error {
	if c.DB != nil {
		return c.DB.Close()
	}
	return nil
}
