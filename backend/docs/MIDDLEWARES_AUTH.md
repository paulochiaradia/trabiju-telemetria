# Middlewares de Autenticação e Autorização

Este documento explica os middlewares implementados no sistema e como utilizá-los.

## 📋 Middlewares Disponíveis

### 1. **AuthMiddleware** (Obrigatório)
Middleware padrão que exige autenticação JWT válida.

```go
// Uso nas rotas
protected := router.Group("/api/v1")
protected.Use(middleware.AuthMiddleware(jwtService))
```

**Funcionalidades:**
- Valida token JWT no header Authorization
- Adiciona informações do usuário ao contexto (user_id, user_email, user_role_id, user_empresa_id)
- Retorna erro 401 se token inválido ou ausente

---

### 2. **OptionalAuthMiddleware** (Opcional)
Middleware para endpoints que funcionam tanto para usuários autenticados quanto anônimos.

```go
// Uso nas rotas
hybrid := router.Group("/api/v1")
hybrid.Use(middleware.OptionalAuthMiddleware(jwtService))
```

**Funcionalidades:**
- Se token válido: adiciona informações do usuário + `is_authenticated: true`
- Se token inválido/ausente: define `is_authenticated: false`
- Nunca retorna erro - sempre continua a execução

**Exemplo de uso no handler:**
```go
func PublicDashboard(c *gin.Context) {
    if middleware.IsAuthenticated(c) {
        // Usuário logado - mostrar dados personalizados
        userID, _ := middleware.GetUserID(c)
        empresaID, _ := middleware.GetUserEmpresaID(c)
        // ... buscar dados específicos da empresa
    } else {
        // Usuário anônimo - mostrar dados públicos/demo
        // ... buscar dados públicos
    }
}
```

---

### 3. **CompanyMiddleware** (Validação de Empresa)
Middleware que verifica se a empresa do usuário está ativa.

```go
// Uso nas rotas (sempre depois do AuthMiddleware)
companyProtected := router.Group("/api/v1")
companyProtected.Use(middleware.AuthMiddleware(jwtService))
companyProtected.Use(middleware.CompanyMiddleware(db))
```

**Funcionalidades:**
- Verifica se empresa existe no banco de dados
- Valida se empresa está ativa (`ativa = true`)
- Adiciona informações da empresa ao contexto (empresa_id, empresa_nome, empresa_ativa)
- Retorna erro 403 se empresa inativa ou não encontrada

**Quando usar:**
- ✅ Operações de telemetria
- ✅ Gestão de veículos
- ✅ Relatórios empresariais
- ❌ Perfil do usuário (pode consultar mesmo com empresa inativa)
- ❌ Logout (sempre deve funcionar)

---

### 4. **RoleMiddleware** (Controle de Acesso)
Middleware que verifica permissões baseadas no role do usuário.

```go
// Uso nas rotas (sempre depois do AuthMiddleware)
admin := protected.Group("/admin")
admin.Use(middleware.RoleMiddleware("admin", "gestor"))
```

**Roles disponíveis:**
- `"admin"`: Acesso total (role_id = 3)
- `"gestor"`: Gestão da empresa (role_id = 3)
- `"entregador"`: Operações básicas (role_id = 1)
- `"ajudante"`: Acesso limitado (role_id = 2)

**Hierarquia:**
- Admin/Gestor pode fazer tudo
- Entregador + Gestor podem fazer operações de entregador
- Qualquer role pode fazer operações de ajudante

---

## 🔧 Funções Helper

### IsAuthenticated
```go
func HandlerExample(c *gin.Context) {
    if middleware.IsAuthenticated(c) {
        // Usuário autenticado
    }
}
```

### GetUserID
```go
func HandlerExample(c *gin.Context) {
    userID, exists := middleware.GetUserID(c)
    if exists {
        // Usar userID
    }
}
```

### GetUserEmpresaID
```go
func HandlerExample(c *gin.Context) {
    empresaID, exists := middleware.GetUserEmpresaID(c)
    if exists {
        // Usar empresaID
    }
}
```

---

## 🚀 Exemplos Práticos

### Rota Pública Simples
```go
// Sem middleware - acesso público total
public.GET("/ping", controllers.PingHandler)
```

### Rota Híbrida (Opcional Auth)
```go
// Com OptionalAuthMiddleware
hybrid := router.Group("/api/v1")
hybrid.Use(middleware.OptionalAuthMiddleware(jwtService))
hybrid.GET("/public/dashboard", handlers.PublicDashboard)
```

### Rota Autenticada Simples
```go
// Com AuthMiddleware
protected := router.Group("/api/v1")
protected.Use(middleware.AuthMiddleware(jwtService))
protected.GET("/auth/profile", authController.GetProfile)
```

### Rota de Negócio (Empresa Ativa)
```go
// Com AuthMiddleware + CompanyMiddleware
business := router.Group("/api/v1")
business.Use(middleware.AuthMiddleware(jwtService))
business.Use(middleware.CompanyMiddleware(db))
business.GET("/vehicles", vehicleController.List)
```

### Rota Administrativa
```go
// Com AuthMiddleware + CompanyMiddleware + RoleMiddleware
admin := business.Group("/admin")
admin.Use(middleware.RoleMiddleware("admin", "gestor"))
admin.GET("/reports", reportController.GetAll)
```

---

## ⚠️ Boas Práticas

1. **Ordem dos middlewares importa:**
   ```go
   ✅ Correto:
   group.Use(middleware.AuthMiddleware(jwtService))
   group.Use(middleware.CompanyMiddleware(db))
   group.Use(middleware.RoleMiddleware("gestor"))
   
   ❌ Incorreto:
   group.Use(middleware.CompanyMiddleware(db))  // Precisa de AuthMiddleware antes
   group.Use(middleware.AuthMiddleware(jwtService))
   ```

2. **Use OptionalAuthMiddleware apenas quando necessário:**
   - Dashboards públicos com dados extras para usuários logados
   - APIs que servem tanto frontend quanto integrações externas
   - Endpoints de consulta com diferentes níveis de detalhe

3. **CompanyMiddleware para operações de negócio:**
   - Sempre use para operações que envolvem dados da empresa
   - Não use para operações pessoais do usuário (perfil, logout)

4. **Teste sempre os cenários:**
   - Usuário não autenticado
   - Token inválido/expirado
   - Empresa inativa
   - Role insuficiente
