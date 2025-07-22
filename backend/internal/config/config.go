package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerPort     string
	DBHost         string
	DBPort         string
	DBUser         string
	DBPassword     string
	DBName         string
	JWTSecret      string
	Environment    string
	TrustedProxies []string
}

// GetDatabaseDSN constrói e retorna a string de conexão do banco
func (c *Config) GetDatabaseDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		c.DBUser,
		c.DBPassword,
		c.DBHost,
		c.DBPort,
		c.DBName,
	)
}

// GetDatabaseDSNWithOptions constrói DSN com opções customizadas
func (c *Config) GetDatabaseDSNWithOptions(options map[string]string) string {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		c.DBUser, c.DBPassword, c.DBHost, c.DBPort, c.DBName)

	for key, value := range options {
		dsn += fmt.Sprintf("&%s=%s", key, value)
	}

	return dsn
}

func LoadConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

	return &Config{
		ServerPort:     getEnv("SERVER_PORT", "8080"),
		DBHost:         getEnv("DB_HOST", "localhost"),
		DBPort:         getEnv("DB_PORT", "3306"),
		DBUser:         getEnv("DB_USER", "user"),
		DBPassword:     getEnv("DB_PASSWORD", "password"),
		DBName:         getEnv("DB_NAME", "dbname"),
		JWTSecret:      getEnv("JWT_SECRET", ""),
		Environment:    getEnv("ENVIRONMENT", "development"),
		TrustedProxies: getTrustedProxies(),
	}, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getTrustedProxies retorna a lista de proxies confiáveis baseada na configuração
func getTrustedProxies() []string {
	// Verifica se há configuração customizada via variável de ambiente
	if proxies := os.Getenv("TRUSTED_PROXIES"); proxies != "" {
		return strings.Split(proxies, ",")
	}

	// Configuração padrão baseada no ambiente
	environment := getEnv("ENVIRONMENT", "development")

	switch environment {
	case "production":
		// Em produção, deve ser configurado via TRUSTED_PROXIES
		// Por segurança, não confiamos em nenhum proxy por padrão
		return []string{}
	case "development", "dev":
		// Em desenvolvimento, permitimos redes locais e Docker
		return []string{
			"127.0.0.1",      // localhost IPv4
			"::1",            // localhost IPv6
			"172.16.0.0/12",  // Docker default bridge networks
			"192.168.0.0/16", // Private networks (VirtualBox, etc)
			"10.0.0.0/8",     // Private networks (Docker, etc)
		}
	default:
		// Para outros ambientes, usar configuração conservadora
		return []string{"127.0.0.1", "::1"}
	}
}
