# Configuração de Portas - Trabiju Telemetria

## 📡 Mapeamento de Portas

| Ambiente | Porta Externa | Porta Interna | Descrição |
|----------|---------------|---------------|-----------|
| **Desenvolvimento (Docker)** | `8081` | `8080` | docker-compose.dev.yml |
| **Produção (Docker)** | `8082` | `8080` | docker-compose.yml |
| **Execução Direta** | `8080` | `8080` | go run cmd/main.go |

## 🔧 URLs de Acesso

### Desenvolvimento
```bash
# API Base
http://localhost:8081

# Endpoints Principais
http://localhost:8081/ping
http://localhost:8081/api/v1/auth/login
http://localhost:8081/api/v1/database/test
```

### Produção
```bash
# API Base
http://localhost:8082

# Endpoints Principais
http://localhost:8082/ping
http://localhost:8082/api/v1/auth/login
http://localhost:8082/api/v1/database/test
```

### Execução Direta (sem Docker)
```bash
# API Base
http://localhost:8080

# Endpoints Principais
http://localhost:8080/ping
http://localhost:8080/api/v1/auth/login
http://localhost:8080/api/v1/database/test
```

## ⚙️ Configuração de Variáveis

### Docker Development (.env)
```env
SERVER_PORT=8080  # Porta interna do container
# Mapeamento: 8081:8080 (definido no docker-compose.dev.yml)
```

### Docker Production (.env)
```env
SERVER_PORT=8080  # Porta interna do container
# Mapeamento: 8082:8080 (definido no docker-compose.yml)
```

### Execução Direta (.env)
```env
SERVER_PORT=8080  # Porta direta na máquina host
```

## 🚀 Como Usar

### Iniciar Desenvolvimento
```bash
docker-compose -f docker-compose.dev.yml up
# Acesse: http://localhost:8081
```

### Iniciar Produção
```bash
docker-compose up
# Acesse: http://localhost:8082
```

### Executar Direto
```bash
cd backend
go run cmd/main.go
# Acesse: http://localhost:8080
```

## 🔍 CORS Configurado

O middleware CORS já está configurado para todas as portas:

```go
allowedOrigins := []string{
    "http://localhost:3000", // Frontend React/Next.js
    "http://localhost:3001", // Frontend alternativo
    "http://localhost:8082", // API própria (Docker Prod)
    "http://localhost:8080", // API própria (direto)
    "http://localhost:8081", // API desenvolvimento (Docker Dev)
}
```

## 📝 Scripts de Teste

### test_dev.ps1 (Development)
```powershell
$baseUrl = "http://localhost:8081"
Invoke-RestMethod -Uri "$baseUrl/ping" -Method GET
```

### test_prod.ps1 (Production)
```powershell
$baseUrl = "http://localhost:8082"
Invoke-RestMethod -Uri "$baseUrl/ping" -Method GET
```

### test_direct.ps1 (Direct)
```powershell
$baseUrl = "http://localhost:8080"
Invoke-RestMethod -Uri "$baseUrl/ping" -Method GET
```
