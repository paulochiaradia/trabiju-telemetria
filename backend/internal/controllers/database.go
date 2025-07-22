package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/paulochiaradia/trabiju-telemetria/internal/config"
	"github.com/paulochiaradia/trabiju-telemetria/internal/database"
)

// DatabaseController gerencia endpoints relacionados ao banco de dados
type DatabaseController struct{}

// TestConnection testa a conexão com o banco de dados
func (dc *DatabaseController) TestConnection(c *gin.Context) {
	// Carregar configurações
	cfg, err := config.LoadConfig()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Erro ao carregar configurações",
			"details": err.Error(),
		})
		return
	}

	// Tentar conectar ao banco
	conn, err := database.NewConnection(cfg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Erro ao conectar com o banco de dados",
			"details": err.Error(),
		})
		return
	}
	defer conn.Close()

	// Testar conexão executando uma query simples
	var version string
	err = conn.DB.QueryRow("SELECT VERSION()").Scan(&version)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Erro ao executar query de teste",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":           "success",
		"message":          "Conexão com banco de dados funcionando!",
		"database_version": version,
		"config": gin.H{
			"host":     cfg.DBHost,
			"port":     cfg.DBPort,
			"database": cfg.DBName,
			"user":     cfg.DBUser,
		},
	})
}

// ListTables lista todas as tabelas do banco de dados
func (dc *DatabaseController) ListTables(c *gin.Context) {
	// Carregar configurações
	cfg, err := config.LoadConfig()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Erro ao carregar configurações",
			"details": err.Error(),
		})
		return
	}

	// Conectar ao banco
	conn, err := database.NewConnection(cfg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Erro ao conectar com o banco de dados",
			"details": err.Error(),
		})
		return
	}
	defer conn.Close()

	// Buscar todas as tabelas
	query := "SELECT TABLE_NAME FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_SCHEMA = ?"
	rows, err := conn.DB.Query(query, cfg.DBName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Erro ao buscar tabelas",
			"details": err.Error(),
		})
		return
	}
	defer rows.Close()

	var tables []string
	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Erro ao ler resultado",
				"details": err.Error(),
			})
			return
		}
		tables = append(tables, tableName)
	}

	c.JSON(http.StatusOK, gin.H{
		"status":       "success",
		"database":     cfg.DBName,
		"total_tables": len(tables),
		"tables":       tables,
	})
}

// DescribeTable mostra a estrutura de uma tabela específica
func (dc *DatabaseController) DescribeTable(c *gin.Context) {
	tableName := c.Param("table")

	// Carregar configurações
	cfg, err := config.LoadConfig()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Erro ao carregar configurações",
			"details": err.Error(),
		})
		return
	}

	// Conectar ao banco
	conn, err := database.NewConnection(cfg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Erro ao conectar com o banco de dados",
			"details": err.Error(),
		})
		return
	}
	defer conn.Close()

	// Descrever estrutura da tabela
	query := "DESCRIBE " + tableName
	rows, err := conn.DB.Query(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Erro ao descrever tabela",
			"details": err.Error(),
		})
		return
	}
	defer rows.Close()

	type Column struct {
		Field   string  `json:"field"`
		Type    string  `json:"type"`
		Null    string  `json:"null"`
		Key     string  `json:"key"`
		Default *string `json:"default"`
		Extra   string  `json:"extra"`
	}

	var columns []Column
	for rows.Next() {
		var col Column
		if err := rows.Scan(&col.Field, &col.Type, &col.Null, &col.Key, &col.Default, &col.Extra); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Erro ao ler estrutura da tabela",
				"details": err.Error(),
			})
			return
		}
		columns = append(columns, col)
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"table":   tableName,
		"columns": columns,
	})
}
