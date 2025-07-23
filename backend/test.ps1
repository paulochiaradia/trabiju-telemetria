# Script PowerShell para testes do projeto Trabiju Telemetria
# Equivalente ao Makefile para Windows

param(
    [string]$Command = "help"
)

function Show-Help {
    Write-Host "🧪 Trabiju Telemetria - Scripts de Teste" -ForegroundColor Green
    Write-Host ""
    Write-Host "Comandos disponíveis:" -ForegroundColor Yellow
    Write-Host "  test              - Executar todos os testes"
    Write-Host "  test-verbose      - Executar testes com saída detalhada"
    Write-Host "  test-coverage     - Executar testes com relatório de cobertura"
    Write-Host "  test-validation   - Executar apenas testes de validação"
    Write-Host "  build             - Compilar aplicação"
    Write-Host "  run               - Executar aplicação"
    Write-Host "  clean             - Limpar arquivos gerados"
    Write-Host "  fmt               - Formatar código"
    Write-Host "  vet               - Executar go vet"
    Write-Host ""
    Write-Host "Exemplo: .\test.ps1 test-validation" -ForegroundColor Cyan
}

function Run-Tests {
    Write-Host "🧪 Executando todos os testes..." -ForegroundColor Cyan
    go test ./...
}

function Run-TestsVerbose {
    Write-Host "🧪 Executando testes com saída detalhada..." -ForegroundColor Cyan
    go test -v ./...
}

function Run-TestsCoverage {
    Write-Host "🧪 Executando testes com cobertura..." -ForegroundColor Cyan
    go test -v -cover ./...
    go test -coverprofile=coverage.out ./...
    go tool cover -html=coverage.out -o coverage.html
    Write-Host "📊 Relatório de cobertura gerado: coverage.html" -ForegroundColor Green
}

function Run-TestsValidation {
    Write-Host "🧪 Executando testes de validação..." -ForegroundColor Cyan
    go test -v ./internal/validation/
}

function Build-App {
    Write-Host "🔨 Compilando aplicação..." -ForegroundColor Cyan
    if (-not (Test-Path "build")) {
        New-Item -ItemType Directory -Path "build"
    }
    go build -o build/trabiju.exe ./cmd/main.go
    Write-Host "✅ Compilação concluída: build/trabiju.exe" -ForegroundColor Green
}

function Run-App {
    Write-Host "🚀 Executando aplicação..." -ForegroundColor Cyan
    go run ./cmd/main.go
}

function Clean-Files {
    Write-Host "🧹 Limpando arquivos gerados..." -ForegroundColor Cyan
    if (Test-Path "build") {
        Remove-Item -Recurse -Force "build"
    }
    if (Test-Path "coverage.out") {
        Remove-Item "coverage.out"
    }
    if (Test-Path "coverage.html") {
        Remove-Item "coverage.html"
    }
    Write-Host "✅ Limpeza concluída" -ForegroundColor Green
}

function Format-Code {
    Write-Host "📝 Formatando código..." -ForegroundColor Cyan
    go fmt ./...
    Write-Host "✅ Código formatado" -ForegroundColor Green
}

function Run-Vet {
    Write-Host "🔍 Executando go vet..." -ForegroundColor Cyan
    go vet ./...
    Write-Host "✅ Go vet concluído" -ForegroundColor Green
}

# Execução baseada no comando
switch ($Command) {
    "test" { Run-Tests }
    "test-verbose" { Run-TestsVerbose }
    "test-coverage" { Run-TestsCoverage }
    "test-validation" { Run-TestsValidation }
    "build" { Build-App }
    "run" { Run-App }
    "clean" { Clean-Files }
    "fmt" { Format-Code }
    "vet" { Run-Vet }
    "help" { Show-Help }
    default {
        Write-Host "❌ Comando desconhecido: $Command" -ForegroundColor Red
        Show-Help
    }
}
