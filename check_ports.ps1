# Exemplo de uso das diferentes portas da API
# Este script demonstra como acessar a API em diferentes ambientes

Write-Host "🚀 Demonstração das Portas da API Trabiju Telemetria" -ForegroundColor Green
Write-Host "=" * 60

# Função para testar conectividade
function Test-APIPort {
    param(
        [string]$Port,
        [string]$Description,
        [string]$Color
    )
    
    $url = "http://localhost:$Port/ping"
    
    Write-Host "`n🔍 Testando $Description (Porta $Port)..." -ForegroundColor $Color
    Write-Host "URL: $url" -ForegroundColor Gray
    
    try {
        $response = Invoke-RestMethod -Uri $url -Method GET -TimeoutSec 5
        Write-Host "✅ Conectado com sucesso!" -ForegroundColor Green
        Write-Host "   Resposta: $($response.message)" -ForegroundColor White
        return $true
    } catch {
        Write-Host "❌ Não foi possível conectar" -ForegroundColor Red
        Write-Host "   Erro: $($_.Exception.Message)" -ForegroundColor Yellow
        return $false
    }
}

# Testar todas as portas
$devOnline = Test-APIPort -Port "8081" -Description "Desenvolvimento (Docker)" -Color "Cyan"
$prodOnline = Test-APIPort -Port "8082" -Description "Produção (Docker)" -Color "Magenta"
$directOnline = Test-APIPort -Port "8080" -Description "Execução Direta" -Color "Blue"

Write-Host "`n" + "=" * 60
Write-Host "📊 Resumo dos Ambientes:" -ForegroundColor Yellow

# Resumo
Write-Host "`n🔧 DESENVOLVIMENTO (Docker Dev)" -ForegroundColor Cyan
Write-Host "   Porta: 8081" -ForegroundColor White
Write-Host "   Status: $(if($devOnline){'🟢 Online'}else{'🔴 Offline'})" -ForegroundColor White
Write-Host "   Comando: docker-compose -f docker-compose.dev.yml up" -ForegroundColor Gray
Write-Host "   Hot Reload: ✅ Ativado" -ForegroundColor Green

Write-Host "`n🚀 PRODUÇÃO (Docker Prod)" -ForegroundColor Magenta
Write-Host "   Porta: 8082" -ForegroundColor White
Write-Host "   Status: $(if($prodOnline){'🟢 Online'}else{'🔴 Offline'})" -ForegroundColor White
Write-Host "   Comando: docker-compose up" -ForegroundColor Gray
Write-Host "   Otimizado: ✅ Build otimizado" -ForegroundColor Green

Write-Host "`n⚙️ EXECUÇÃO DIRETA" -ForegroundColor Blue
Write-Host "   Porta: 8080" -ForegroundColor White
Write-Host "   Status: $(if($directOnline){'🟢 Online'}else{'🔴 Offline'})" -ForegroundColor White
Write-Host "   Comando: go run cmd/main.go" -ForegroundColor Gray
Write-Host "   Debug: ✅ Ideal para debug" -ForegroundColor Green

# Sugestões baseadas no status
Write-Host "`n💡 Sugestões:" -ForegroundColor Yellow

if (-not $devOnline -and -not $prodOnline -and -not $directOnline) {
    Write-Host "   • Nenhum ambiente está rodando" -ForegroundColor Red
    Write-Host "   • Para começar: docker-compose -f docker-compose.dev.yml up" -ForegroundColor Cyan
}
elseif ($devOnline) {
    Write-Host "   • Desenvolvimento está ativo - ideal para testes!" -ForegroundColor Green
    Write-Host "   • Teste com: .\test_api.ps1" -ForegroundColor Cyan
}
elseif ($prodOnline) {
    Write-Host "   • Produção está ativa - ideal para validação final!" -ForegroundColor Green
    Write-Host "   • Teste com: .\test_api_prod.ps1" -ForegroundColor Cyan
}

Write-Host "`n🔗 URLs Importantes:" -ForegroundColor Yellow
if ($devOnline) {
    Write-Host "   Dev API: http://localhost:8081/api/v1" -ForegroundColor Cyan
}
if ($prodOnline) {
    Write-Host "   Prod API: http://localhost:8082/api/v1" -ForegroundColor Magenta
}
if ($directOnline) {
    Write-Host "   Direct API: http://localhost:8080/api/v1" -ForegroundColor Blue
}

Write-Host "`n✨ Para mais informações, consulte PORTS_CONFIG.md" -ForegroundColor Green
