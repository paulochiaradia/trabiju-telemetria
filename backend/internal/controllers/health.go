package controllers

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"runtime"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/paulochiaradia/trabiju-telemetria/internal/config"
)

// HealthController gerencia o endpoint de health check
type HealthController struct {
	startTime time.Time
	db        *sql.DB
}

// NewHealthController cria uma nova instância do HealthController
func NewHealthController(db *sql.DB) *HealthController {
	return &HealthController{
		startTime: time.Now(),
		db:        db,
	}
}

// HealthStatus representa o status de saúde de um componente
type HealthStatus struct {
	Status  string                 `json:"status"`
	Details map[string]interface{} `json:"details,omitempty"`
	Error   string                 `json:"error,omitempty"`
}

// HealthResponse representa a resposta completa do health check
type HealthResponse struct {
	Status      string                  `json:"status"`
	Timestamp   time.Time               `json:"timestamp"`
	Uptime      string                  `json:"uptime"`
	Version     string                  `json:"version"`
	Environment string                  `json:"environment"`
	Components  map[string]HealthStatus `json:"components"`
	System      SystemInfo              `json:"system"`
}

// SystemInfo contém informações do sistema
type SystemInfo struct {
	GoVersion    string   `json:"go_version"`
	NumCPU       int      `json:"num_cpu"`
	NumGoroutine int      `json:"num_goroutine"`
	MemoryStats  MemStats `json:"memory_stats"`
}

// MemStats contém estatísticas de memória
type MemStats struct {
	Alloc        uint64 `json:"alloc_mb"`
	TotalAlloc   uint64 `json:"total_alloc_mb"`
	Sys          uint64 `json:"sys_mb"`
	NumGC        uint32 `json:"num_gc"`
	HeapAlloc    uint64 `json:"heap_alloc_mb"`
	HeapSys      uint64 `json:"heap_sys_mb"`
	HeapInuse    uint64 `json:"heap_inuse_mb"`
	HeapReleased uint64 `json:"heap_released_mb"`
}

// HealthCheck executa verificações completas de saúde da aplicação
func (hc *HealthController) HealthCheck(c *gin.Context) {
	timestamp := time.Now()
	uptime := timestamp.Sub(hc.startTime)

	// Carregar configurações
	cfg, err := config.LoadConfig()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":    "error",
			"timestamp": timestamp,
			"error":     "Erro ao carregar configurações: " + err.Error(),
		})
		return
	}

	// Verificar componentes
	components := make(map[string]HealthStatus)
	overallStatus := "healthy"

	// 1. Verificar banco de dados
	dbStatus := hc.checkDatabase(cfg)
	components["database"] = dbStatus
	if dbStatus.Status != "healthy" {
		overallStatus = "unhealthy"
	}

	// 2. Verificar sistema de arquivos (logs, temp, etc.)
	fsStatus := hc.checkFileSystem()
	components["filesystem"] = fsStatus
	if fsStatus.Status != "healthy" {
		overallStatus = "degraded"
	}

	// 3. Verificar memória
	memStatus := hc.checkMemory()
	components["memory"] = memStatus
	if memStatus.Status != "healthy" {
		if overallStatus == "healthy" {
			overallStatus = "degraded"
		}
	}

	// Obter informações do sistema
	systemInfo := hc.getSystemInfo()

	// Construir resposta
	response := HealthResponse{
		Status:      overallStatus,
		Timestamp:   timestamp,
		Uptime:      hc.formatUptime(uptime),
		Version:     "1.0.0", // Você pode pegar isso de uma variável de build
		Environment: cfg.Environment,
		Components:  components,
		System:      systemInfo,
	}

	// Definir status HTTP baseado na saúde geral
	statusCode := http.StatusOK
	if overallStatus == "degraded" {
		statusCode = http.StatusOK // 200 para degraded (ainda funcional)
	} else if overallStatus == "unhealthy" {
		statusCode = http.StatusServiceUnavailable // 503 para unhealthy
	}

	c.JSON(statusCode, response)
}

// checkDatabase verifica a saúde do banco de dados
func (hc *HealthController) checkDatabase(cfg *config.Config) HealthStatus {
	status := HealthStatus{
		Status:  "healthy",
		Details: make(map[string]interface{}),
	}

	// Verificar se a conexão está disponível
	if hc.db == nil {
		status.Status = "unhealthy"
		status.Error = "Conexão com banco não inicializada"
		return status
	}

	// Testar ping
	start := time.Now()
	err := hc.db.Ping()
	latency := time.Since(start)

	if err != nil {
		status.Status = "unhealthy"
		status.Error = "Falha ao conectar: " + err.Error()
		return status
	}

	// Obter versão do MySQL
	var version string
	err = hc.db.QueryRow("SELECT VERSION()").Scan(&version)
	if err != nil {
		status.Status = "degraded"
		status.Error = "Falha ao obter versão: " + err.Error()
	}

	// Obter estatísticas de conexão
	stats := hc.db.Stats()

	status.Details["latency_ms"] = latency.Milliseconds()
	status.Details["version"] = version
	status.Details["host"] = cfg.DBHost + ":" + cfg.DBPort
	status.Details["database"] = cfg.DBName
	status.Details["connections"] = map[string]interface{}{
		"open":             stats.OpenConnections,
		"in_use":           stats.InUse,
		"idle":             stats.Idle,
		"max_open":         stats.MaxOpenConnections,
		"wait_count":       stats.WaitCount,
		"wait_duration_ms": stats.WaitDuration.Milliseconds(),
	}

	// Verificar se a latência está alta
	if latency > 100*time.Millisecond {
		status.Status = "degraded"
		status.Details["warning"] = "Latência alta detectada"
	}

	return status
}

// checkFileSystem verifica a saúde do sistema de arquivos
func (hc *HealthController) checkFileSystem() HealthStatus {
	status := HealthStatus{
		Status:  "healthy",
		Details: make(map[string]interface{}),
	}

	// Verificar se consegue escrever arquivos temporários
	start := time.Now()
	testFile := "/tmp/health_check_" + time.Now().Format("20060102150405")

	// No Windows, usar diretório temp diferente
	if runtime.GOOS == "windows" {
		testFile = "C:\\Windows\\Temp\\health_check_" + time.Now().Format("20060102150405")
	}

	err := writeTestFile(testFile)
	if err != nil {
		status.Status = "degraded"
		status.Error = "Falha ao escrever arquivo de teste: " + err.Error()
	} else {
		// Limpar arquivo de teste
		removeTestFile(testFile)
	}

	ioLatency := time.Since(start)
	status.Details["io_latency_ms"] = ioLatency.Milliseconds()
	status.Details["working_directory"] = getCurrentWorkingDir()

	return status
}

// checkMemory verifica o uso de memória
func (hc *HealthController) checkMemory() HealthStatus {
	status := HealthStatus{
		Status:  "healthy",
		Details: make(map[string]interface{}),
	}

	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	allocMB := bToMb(m.Alloc)
	sysMB := bToMb(m.Sys)

	status.Details["alloc_mb"] = allocMB
	status.Details["sys_mb"] = sysMB
	status.Details["heap_alloc_mb"] = bToMb(m.HeapAlloc)
	status.Details["num_gc"] = m.NumGC

	// Verificar se o uso de memória está alto (mais de 1GB)
	if allocMB > 1024 {
		status.Status = "degraded"
		status.Details["warning"] = "Alto uso de memória detectado"
	}

	return status
}

// getSystemInfo obtém informações do sistema
func (hc *HealthController) getSystemInfo() SystemInfo {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	return SystemInfo{
		GoVersion:    runtime.Version(),
		NumCPU:       runtime.NumCPU(),
		NumGoroutine: runtime.NumGoroutine(),
		MemoryStats: MemStats{
			Alloc:        bToMb(m.Alloc),
			TotalAlloc:   bToMb(m.TotalAlloc),
			Sys:          bToMb(m.Sys),
			NumGC:        m.NumGC,
			HeapAlloc:    bToMb(m.HeapAlloc),
			HeapSys:      bToMb(m.HeapSys),
			HeapInuse:    bToMb(m.HeapInuse),
			HeapReleased: bToMb(m.HeapReleased),
		},
	}
}

// formatUptime formata o tempo de atividade
func (hc *HealthController) formatUptime(uptime time.Duration) string {
	days := int(uptime.Hours()) / 24
	hours := int(uptime.Hours()) % 24
	minutes := int(uptime.Minutes()) % 60
	seconds := int(uptime.Seconds()) % 60

	if days > 0 {
		return fmt.Sprintf("%dd %dh %dm %ds", days, hours, minutes, seconds)
	} else if hours > 0 {
		return fmt.Sprintf("%dh %dm %ds", hours, minutes, seconds)
	} else if minutes > 0 {
		return fmt.Sprintf("%dm %ds", minutes, seconds)
	}
	return fmt.Sprintf("%ds", seconds)
}

// Funções auxiliares
func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}

func writeTestFile(filename string) error {
	// Implementação simplificada - em produção você pode querer algo mais robusto
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString("health check test")
	return err
}

func removeTestFile(filename string) {
	os.Remove(filename)
}

func getCurrentWorkingDir() string {
	wd, err := os.Getwd()
	if err != nil {
		return "unknown"
	}
	return wd
}
