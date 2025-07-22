# Sistema de Cadastro e Identificação de Roles

## 🎯 **Como o usuário identifica seu papel (entregador, gestor, ajudante)?**

### **3 Fluxos Principais:**

---

## 🔗 **1. CONVITE POR GESTOR (Mais Seguro)**

### Como funciona:
1. **Gestor** acessa o sistema
2. **Gestor** convida usuário por email
3. **Gestor** define qual será o role (entregador/ajudante/gestor)
4. **Sistema** envia email com link único
5. **Usuário** clica no link e completa cadastro
6. **Usuário** já fica com role definido

### Endpoint:
```bash
POST /api/auth/criar-convite
{
  "email": "joao@email.com",
  "role_nome": "entregador",
  "mensagem": "Bem-vindo à equipe!",
  "valido_por_dias": 7
}
```

### Vantagens:
- ✅ Máxima segurança
- ✅ Controle total do gestor
- ✅ Sem cadastros indevidos

---

## 🏢 **2. CÓDIGO DA EMPRESA (Prático)**

### Como funciona:
1. **Empresa** possui código único (ex: TELEME2025)
2. **Usuário** acessa tela de cadastro
3. **Usuário** informa código da empresa
4. **Sistema** valida código e mostra roles disponíveis
5. **Usuário** escolhe role (entregador/ajudante)
6. **Sistema** pode criar diretamente ou solicitar aprovação

### Endpoint:
```bash
POST /api/auth/cadastro-com-codigo
{
  "nome": "João Silva",
  "email": "joao@email.com",
  "senha": "senha123",
  "codigo_empresa": "TELEME2025",
  "role_desejado": "entregador",
  "justificativa": "Trabalho nas entregas da região sul"
}
```

### Configuração por empresa:
```json
{
  "allow_self_register": true,
  "default_role": "entregador",
  "require_approval": false,
  "available_roles": ["entregador", "ajudante"]
}
```

---

## 📋 **3. SOLICITAÇÃO COM APROVAÇÃO**

### Como funciona:
1. **Usuário** se cadastra com código da empresa
2. **Usuário** solicita role de "gestor" ou empresa exige aprovação
3. **Sistema** cria solicitação pendente
4. **Gestor** recebe notificação
5. **Gestor** aprova ou rejeita
6. **Usuário** recebe notificação do resultado

### Endpoints:
```bash
# Listar solicitações (gestor)
GET /api/auth/solicitacoes-pendentes

# Aprovar/Rejeitar (gestor)
POST /api/auth/aprovar-solicitacao/123
{
  "aprovado": true,
  "observacoes": "Aprovado para coordenar entregas"
}
```

---

## 🚀 **Fluxo Recomendado para seu Sistema:**

### **Para Entregadores e Ajudantes:**
```
Usuário → Código da Empresa → Escolhe Role → Cadastro Direto
```

### **Para Gestores:**
```
Usuário → Código da Empresa → Solicita Gestor → Aprovação Manual
```

### **Para Convites:**
```
Gestor → Criar Convite → Email → Usuário Aceita → Cadastro Automático
```

---

## 🔧 **Implementação no Frontend:**

### **Tela de Cadastro:**
```javascript
// 1. Verificar código da empresa
const verificarCodigo = async (codigo) => {
  const response = await fetch(`/api/auth/verificar-codigo/${codigo}`);
  const data = await response.json();
  
  if (data.valido) {
    // Mostrar roles disponíveis
    setRolesDisponiveis(data.roles_disponiveis);
  }
};

// 2. Cadastro com role
const cadastrar = async (dados) => {
  const response = await fetch('/api/auth/cadastro-com-codigo', {
    method: 'POST',
    body: JSON.stringify(dados)
  });
  
  const result = await response.json();
  
  if (result.status === 'ativo') {
    // Cadastro aprovado imediatamente
    redirectToLogin();
  } else if (result.status === 'pendente') {
    // Aguardando aprovação
    showPendingMessage();
  }
};
```

---

## 📊 **Tabelas do Sistema:**

### **empresas** - Controle de organizações
### **convites_usuario** - Convites por email
### **solicitacoes_cadastro** - Solicitações para aprovação
### **usuarios** - Usuários com role definido

---

## 🛡️ **Segurança:**

### **Códigos de Empresa:**
- Únicos por empresa
- Podem ser regenerados
- Controlam quais roles são permitidos

### **Tokens de Convite:**
- Únicos e com expiração
- Uso único (marcados como "usado")
- Ligados a email específico

### **Aprovações:**
- Apenas gestores podem aprovar
- Log completo de quem aprovou
- Notificações por email

---

## ✅ **Vantagens desta Abordagem:**

1. **Flexível** - Múltiplas formas de cadastro
2. **Segura** - Controle total sobre quem acessa
3. **Escalável** - Funciona para múltiplas empresas
4. **Auditável** - Log completo de cadastros
5. **UX Amigável** - Processo claro para usuários

Esta estrutura garante que apenas pessoas autorizadas tenham acesso ao sistema e com o role correto! 🎯
