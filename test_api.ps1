# Script de Teste Completo da API de Autenticação (DESENVOLVIMENTO)
# PowerShell script para testar todos os fluxos da API no ambiente de desenvolvimento

$baseUrl = "http://localhost:8081/api/v1"
$headers = @{
    "Content-Type" = "application/json"
}

Write-Host "🚀 Iniciando testes da API de Autenticação (DESENVOLVIMENTO - Porta 8081)..." -ForegroundColor Green
Write-Host "=" * 60

# Função para fazer requisições HTTP
function Invoke-ApiRequest {
    param(
        [string]$Method,
        [string]$Url,
        [hashtable]$Headers = @{"Content-Type" = "application/json"},
        [string]$Body = $null
    )
    
    try {
        if ($Body) {
            $response = Invoke-RestMethod -Uri $Url -Method $Method -Headers $Headers -Body $Body
        } else {
            $response = Invoke-RestMethod -Uri $Url -Method $Method -Headers $Headers
        }
        return @{ Success = $true; Data = $response }
    } catch {
        return @{ Success = $false; Error = $_.Exception.Message; Response = $_.Exception.Response }
    }
}

# 1. TESTE DE HEALTH CHECK
Write-Host "📍 1. Testando Health Check..." -ForegroundColor Cyan
$healthResult = Invoke-ApiRequest -Method "GET" -Url "$baseUrl/ping"
if ($healthResult.Success) {
    Write-Host "✅ Health Check: $($healthResult.Data.message)" -ForegroundColor Green
} else {
    Write-Host "❌ Health Check falhou: $($healthResult.Error)" -ForegroundColor Red
}

# 2. TESTE DE ESTRUTURA DO BANCO
Write-Host "`n📍 2. Testando estrutura do banco..." -ForegroundColor Cyan
$dbResult = Invoke-ApiRequest -Method "GET" -Url "$baseUrl/../database/test"
if ($dbResult.Success) {
    Write-Host "✅ Banco de dados conectado" -ForegroundColor Green
} else {
    Write-Host "❌ Falha na conexão com banco: $($dbResult.Error)" -ForegroundColor Red
}

# 3. TESTE DE CADASTRO COM CÓDIGO DA EMPRESA
Write-Host "`n📍 3. Testando Cadastro com Código da Empresa..." -ForegroundColor Cyan

$cadastroData = @{
    nome = "João Silva"
    email = "joao.silva@teste.com"
    cpf = "12345678901"
    senha = "MinhaSenh@123"
    codigo_empresa = "EMP001"
} | ConvertTo-Json

$cadastroResult = Invoke-ApiRequest -Method "POST" -Url "$baseUrl/auth/register/code" -Body $cadastroData
if ($cadastroResult.Success) {
    Write-Host "✅ Cadastro com código: Usuário criado com sucesso" -ForegroundColor Green
    Write-Host "   User ID: $($cadastroResult.Data.user_id)" -ForegroundColor Yellow
} else {
    Write-Host "❌ Cadastro com código falhou: $($cadastroResult.Error)" -ForegroundColor Red
}

# 4. TESTE DE LOGIN
Write-Host "`n📍 4. Testando Login..." -ForegroundColor Cyan

$loginData = @{
    email = "joao.silva@teste.com"
    senha = "MinhaSenh@123"
} | ConvertTo-Json

$loginResult = Invoke-ApiRequest -Method "POST" -Url "$baseUrl/auth/login" -Body $loginData
$accessToken = $null
$refreshToken = $null

if ($loginResult.Success) {
    Write-Host "✅ Login bem-sucedido" -ForegroundColor Green
    $accessToken = $loginResult.Data.access_token
    $refreshToken = $loginResult.Data.refresh_token
    Write-Host "   Access Token: $($accessToken.Substring(0, 20))..." -ForegroundColor Yellow
    Write-Host "   Refresh Token: $($refreshToken.Substring(0, 20))..." -ForegroundColor Yellow
} else {
    Write-Host "❌ Login falhou: $($loginResult.Error)" -ForegroundColor Red
}

# 5. TESTE DE PERFIL (ROTA PROTEGIDA)
if ($accessToken) {
    Write-Host "`n📍 5. Testando acesso ao perfil (rota protegida)..." -ForegroundColor Cyan
    
    $authHeaders = @{
        "Content-Type" = "application/json"
        "Authorization" = "Bearer $accessToken"
    }
    
    $profileResult = Invoke-ApiRequest -Method "GET" -Url "$baseUrl/auth/profile" -Headers $authHeaders
    if ($profileResult.Success) {
        Write-Host "✅ Perfil acessado com sucesso" -ForegroundColor Green
        Write-Host "   Nome: $($profileResult.Data.nome)" -ForegroundColor Yellow
        Write-Host "   Email: $($profileResult.Data.email)" -ForegroundColor Yellow
    } else {
        Write-Host "❌ Acesso ao perfil falhou: $($profileResult.Error)" -ForegroundColor Red
    }
}

# 6. TESTE DE REFRESH TOKEN
if ($refreshToken) {
    Write-Host "`n📍 6. Testando Refresh Token..." -ForegroundColor Cyan
    
    $refreshData = @{
        refresh_token = $refreshToken
    } | ConvertTo-Json
    
    $refreshResult = Invoke-ApiRequest -Method "POST" -Url "$baseUrl/auth/refresh" -Body $refreshData
    if ($refreshResult.Success) {
        Write-Host "✅ Token renovado com sucesso" -ForegroundColor Green
        $newAccessToken = $refreshResult.Data.access_token
        Write-Host "   Novo Access Token: $($newAccessToken.Substring(0, 20))..." -ForegroundColor Yellow
    } else {
        Write-Host "❌ Renovação de token falhou: $($refreshResult.Error)" -ForegroundColor Red
    }
}

# 7. TESTE DE CADASTRO COM EMAIL INVÁLIDO
Write-Host "`n📍 7. Testando validação de email..." -ForegroundColor Cyan

$emailInvalidoData = @{
    nome = "Maria Santos"
    email = "email-invalido"
    cpf = "98765432100"
    senha = "OutraSenh@123"
    codigo_empresa = "EMP001"
} | ConvertTo-Json

$emailInvalidoResult = Invoke-ApiRequest -Method "POST" -Url "$baseUrl/auth/register/code" -Body $emailInvalidoData
if (!$emailInvalidoResult.Success) {
    Write-Host "✅ Validação de email funcionando: Rejeitou email inválido" -ForegroundColor Green
} else {
    Write-Host "❌ Validação de email falhou: Aceitou email inválido" -ForegroundColor Red
}

# 8. TESTE DE SENHA FRACA
Write-Host "`n📍 8. Testando validação de senha..." -ForegroundColor Cyan

$senhaFracaData = @{
    nome = "Pedro Oliveira"
    email = "pedro@teste.com"
    cpf = "11122233344"
    senha = "123"
    codigo_empresa = "EMP001"
} | ConvertTo-Json

$senhaFracaResult = Invoke-ApiRequest -Method "POST" -Url "$baseUrl/auth/register/code" -Body $senhaFracaData
if (!$senhaFracaResult.Success) {
    Write-Host "✅ Validação de senha funcionando: Rejeitou senha fraca" -ForegroundColor Green
} else {
    Write-Host "❌ Validação de senha falhou: Aceitou senha fraca" -ForegroundColor Red
}

# 9. TESTE DE LOGIN COM CREDENCIAIS INVÁLIDAS
Write-Host "`n📍 9. Testando login com credenciais inválidas..." -ForegroundColor Cyan

$loginInvalidoData = @{
    email = "usuario@inexistente.com"
    senha = "senhaerrada"
} | ConvertTo-Json

$loginInvalidoResult = Invoke-ApiRequest -Method "POST" -Url "$baseUrl/auth/login" -Body $loginInvalidoData
if (!$loginInvalidoResult.Success) {
    Write-Host "✅ Validação de login funcionando: Rejeitou credenciais inválidas" -ForegroundColor Green
} else {
    Write-Host "❌ Validação de login falhou: Aceitou credenciais inválidas" -ForegroundColor Red
}

# 10. TESTE DE ACESSO SEM TOKEN
Write-Host "`n📍 10. Testando acesso sem autenticação..." -ForegroundColor Cyan

$semTokenResult = Invoke-ApiRequest -Method "GET" -Url "$baseUrl/auth/profile"
if (!$semTokenResult.Success) {
    Write-Host "✅ Proteção de rota funcionando: Bloqueou acesso sem token" -ForegroundColor Green
} else {
    Write-Host "❌ Proteção de rota falhou: Permitiu acesso sem token" -ForegroundColor Red
}

# RESUMO DOS TESTES
Write-Host "`n" + "=" * 60
Write-Host "📊 RESUMO DOS TESTES" -ForegroundColor Magenta
Write-Host "=" * 60

$tests = @(
    @{ Name = "Health Check"; Status = $healthResult.Success }
    @{ Name = "Conexão DB"; Status = $dbResult.Success }
    @{ Name = "Cadastro com Código"; Status = $cadastroResult.Success }
    @{ Name = "Login"; Status = $loginResult.Success }
    @{ Name = "Acesso ao Perfil"; Status = ($accessToken -and $profileResult.Success) }
    @{ Name = "Refresh Token"; Status = ($refreshToken -and $refreshResult.Success) }
    @{ Name = "Validação Email"; Status = (!$emailInvalidoResult.Success) }
    @{ Name = "Validação Senha"; Status = (!$senhaFracaResult.Success) }
    @{ Name = "Validação Login"; Status = (!$loginInvalidoResult.Success) }
    @{ Name = "Proteção de Rota"; Status = (!$semTokenResult.Success) }
)

$passed = 0
$failed = 0

foreach ($test in $tests) {
    if ($test.Status) {
        Write-Host "✅ $($test.Name)" -ForegroundColor Green
        $passed++
    } else {
        Write-Host "❌ $($test.Name)" -ForegroundColor Red
        $failed++
    }
}

Write-Host "`n📈 RESULTADO FINAL:" -ForegroundColor White
Write-Host "   ✅ Passou: $passed testes" -ForegroundColor Green
Write-Host "   ❌ Falhou: $failed testes" -ForegroundColor Red
Write-Host "   📊 Taxa de sucesso: $([math]::Round(($passed / ($passed + $failed)) * 100, 2))%" -ForegroundColor Yellow

if ($failed -eq 0) {
    Write-Host "`n🎉 TODOS OS TESTES PASSARAM! API está funcionando corretamente." -ForegroundColor Green
} else {
    Write-Host "`n⚠️  Alguns testes falharam. Verifique os erros acima." -ForegroundColor Yellow
}

Write-Host "`n🔗 Para testes manuais, use os seguintes endpoints:"
Write-Host "   - Health: GET $baseUrl/ping"
Write-Host "   - Cadastro: POST $baseUrl/auth/register/code"
Write-Host "   - Login: POST $baseUrl/auth/login"
Write-Host "   - Perfil: GET $baseUrl/auth/profile (requer token)"
Write-Host "   - Refresh: POST $baseUrl/auth/refresh"
