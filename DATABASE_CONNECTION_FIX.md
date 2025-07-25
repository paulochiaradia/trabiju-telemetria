# 🔧 Solução para Problemas de Conexão com Banco de Dados

## 📋 Problema Identificado

**Erro comum:**
```
go_backend_prod  | Erro ao conectar com o banco de dados: falha ao testar conexão com o banco: dial tcp 172.18.0.2:3306: connect: connection refused
```

## ✅ Soluções Implementadas

### 1. **Retry Automático na Conexão**
- ✅ Implementado retry inteligente com 30 tentativas
- ✅ Intervalo de 2 segundos entre tentativas
- ✅ Logs detalhados do processo de conexão
- ✅ Tolerância a inicializações lentas do MySQL

### 2. **Health Check Melhorado**
- ✅ `start_period: 40s` - tempo inicial para MySQL inicializar
- ✅ `interval: 10s` - verificações a cada 10 segundos
- ✅ `timeout: 10s` - timeout de 10 segundos por verificação
- ✅ `retries: 10` - até 10 tentativas antes de marcar como unhealthy

### 3. **Logs Informativos**
```
🔄 Tentativa 1/30 de conexão com o banco...
📍 Host: mysql:3306
❌ Falha no ping do banco (tentativa 1): dial tcp 172.18.0.2:3306: connect: connection refused
⏳ Aguardando 2s antes da próxima tentativa...
✅ Conexão com banco estabelecida com sucesso na tentativa 8!
```

## 🚀 Como Usar

### Para Desenvolvimento:
```bash
cd c:\Users\paulo\trabiju-telemetria
docker-compose -f docker-compose.dev.yml up
```

### Para Produção:
```bash
cd c:\Users\paulo\trabiju-telemetria
docker-compose up
```

### Para Diagnóstico:
```powershell
# Verificar status do banco
.\diagnose_db.ps1 -Environment dev   # Para desenvolvimento
.\diagnose_db.ps1 -Environment prod  # Para produção

# Monitorar logs em tempo real
docker logs -f mysql_db_dev    # Desenvolvimento
docker logs -f mysql_db        # Produção
```

## 🔍 Processo de Inicialização

### O que acontece agora:

1. **MySQL inicia** - Container é criado
2. **Inicialização interna** - MySQL configura banco de dados
3. **Migrations executam** - Scripts SQL são executados
4. **Health check ativa** - Verificações de saúde começam
5. **Backend conecta** - Com retry automático até conseguir

### Timeline típica:
- `0s` - Containers iniciam
- `0-40s` - MySQL inicializa (período de grace)
- `40s+` - Health checks começam
- `50-70s` - MySQL fica "healthy"
- `50-80s` - Backend conecta com sucesso

## ⚠️ Troubleshooting

### Se ainda houver problemas:

1. **Verificar logs do MySQL:**
   ```bash
   docker logs mysql_db_dev
   ```

2. **Verificar health status:**
   ```bash
   docker inspect mysql_db_dev --format='{{.State.Health.Status}}'
   ```

3. **Reiniciar apenas o MySQL:**
   ```bash
   docker-compose -f docker-compose.dev.yml restart mysql
   ```

4. **Reset completo:**
   ```bash
   docker-compose -f docker-compose.dev.yml down
   docker volume rm trabiju-telemetria_db_data_dev
   docker-compose -f docker-compose.dev.yml up
   ```

### Se quiser ajustar o retry:

Edite `backend/internal/database/connection.go`:
```go
// Mais conservador (mais tentativas, intervalo maior)
return NewConnectionWithRetry(cfg, 60, 5*time.Second)

// Mais agressivo (menos tentativas, intervalo menor)
return NewConnectionWithRetry(cfg, 15, 1*time.Second)
```

## 📊 Monitoramento

### Scripts úteis:

```powershell
# Status geral
.\check_ports.ps1

# Diagnóstico específico do banco
.\diagnose_db.ps1

# Teste da API após conectar
.\test_launcher.ps1
```

### Sinais de que está funcionando:

✅ **Logs do backend:**
```
✅ Conexão com banco estabelecida com sucesso na tentativa X!
🚀 Servidor iniciando na porta 8080
```

✅ **Health check:**
```bash
$ docker inspect mysql_db_dev --format='{{.State.Health.Status}}'
healthy
```

✅ **API respondendo:**
```bash
$ curl http://localhost:8081/ping
{"message":"pong"}
```

## 🎯 Resultado

Agora o sistema é **robusto** e **tolerante** a:
- ✅ Inicializações lentas do MySQL
- ✅ Problemas temporários de rede
- ✅ Reinicializações do banco
- ✅ Ambientes com recursos limitados

**Não mais erros de "connection refused"!** 🚀
