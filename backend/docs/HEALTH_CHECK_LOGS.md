# Soluções para Requisições de Health Check

As requisições repetidas de `/ping` que aparecem nos logs são causadas pelo **Docker Health Check** configurado no `Dockerfile`. Estas não são "requisições fantasma", mas sim verificações automáticas do Docker para monitorar a saúde do container.

## Causa Raiz
```dockerfile
# No Dockerfile de produção:
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8080/ping || exit 1
```

## Soluções Implementadas

### ✅ Solução 1: Filtrar Logs de Health Check (IMPLEMENTADA)
- **Arquivo**: `internal/middleware/logging.go`
- **Descrição**: Middleware customizado que não loga requisições para `/ping`
- **Vantagem**: Mantém o health check funcionando, mas remove o "ruído" dos logs
- **Uso**: Aplicado automaticamente no `main.go`

```go
func CustomLogger() gin.HandlerFunc {
	return gin.LoggerWithConfig(gin.LoggerConfig{
		SkipPaths: []string{"/ping"}, // Skip logging for ping endpoint
	})
}
```

## Outras Soluções Possíveis

### Solução 2: Aumentar Intervalo do Health Check
Modificar o Dockerfile para verificar com menos frequência:
```dockerfile
HEALTHCHECK --interval=2m --timeout=3s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8080/ping || exit 1
```

### Solução 3: Desabilitar Health Check
Comentar ou remover o HEALTHCHECK do Dockerfile:
```dockerfile
# HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
#   CMD wget --no-verbose --tries=1 --spider http://localhost:8080/ping || exit 1
```

### Solução 4: Endpoint Dedicado para Health Check
Criar um endpoint específico (`/health`) que não seja logado:
```go
r.GET("/health", func(c *gin.Context) {
    c.JSON(200, gin.H{"status": "ok"})
})
```

### Solução 5: Configurar Logs por Ambiente
Aplicar filtro apenas em produção:
```go
if cfg.Environment == "production" {
    r.Use(middleware.CustomLogger())
} else {
    r.Use(gin.Logger())
}
```

## Recomendação
A **Solução 1** (já implementada) é a melhor opção porque:
- ✅ Mantém o health check ativo (importante para orquestração)
- ✅ Remove ruído desnecessário dos logs
- ✅ Preserva logs de requisições reais
- ✅ Não afeta o comportamento da aplicação

## Como Testar
1. Rebuilde a aplicação: `docker-compose up --build`
2. Observe que as requisições `/ping` não aparecem mais nos logs
3. Teste outros endpoints para confirmar que ainda são logados
