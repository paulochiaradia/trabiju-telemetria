# 🚚 Sistema de Gestão de Telemetria

Sistema completo de gestão de telemetria para empresas de logística, desenvolvido em Go com arquitetura robusta, sistema de autenticação avançado e controle granular de permissões.

---

## 🎯 Características Principais

### 🔐 Sistema de Autenticação Multinível
- **3 Tipos de Usuário**: Entregador, Ajudante e Gestor
- **Sistema de Convites**: Gestores podem convidar usuários com roles específicos
- **Códigos de Empresa**: Cadastro usando códigos únicos por organização
- **Aprovação Manual**: Solicitações para roles sensíveis
- **Controle de Sessões**: Gerenciamento completo de sessões ativas
- **Logs de Auditoria**: Rastreamento completo de tentativas de login

### 🏢 Gestão Multi-Empresa
- **Códigos Únicos**: Cada empresa possui código de convite próprio
- **Configurações Personalizadas**: Políticas de cadastro por empresa
- **Isolamento de Dados**: Usuários só acessam dados da própria empresa

### 🔧 Arquitetura Robusta
- **Clean Architecture**: Separação clara de responsabilidades
- **Database First**: Migrations versionadas e controladas
- **API RESTful**: Endpoints bem definidos e documentados
- **Configuração por Ambiente**: Variáveis de ambiente para flexibilidade

---

## 📦 Tecnologias Utilizadas

### Backend
- **[Go 1.21.1](https://golang.org/)** - Linguagem principal
- **[Gin Web Framework](https://github.com/gin-gonic/gin)** - Framework HTTP
- **[MySQL 8.0](https://www.mysql.com/)** - Banco de dados
- **[JWT](https://github.com/golang-jwt/jwt)** - Autenticação (planejado)
- **[godotenv](https://github.com/joho/godotenv)** - Gerenciamento de configurações

### DevOps & Infraestrutura
- **[Docker](https://www.docker.com/)** - Containerização
- **[Docker Compose](https://docs.docker.com/compose/)** - Orquestração
- **[Air](https://github.com/cosmtrek/air)** - Hot reload em desenvolvimento

---

## 🚀 Instalação e Execução

### Pré-requisitos
- Docker Desktop instalado e rodando
- Git para clonar o repositório

### 1. Clone o Repositório
```bash
git clone https://github.com/paulochiaradia/trabiju-telemetria.git
cd trabiju-telemetria
```

### 2. Configuração Automática
O projeto já vem com todas as configurações necessárias:
- ✅ Arquivo `.env` configurado
- ✅ Migrations prontas para execução
- ✅ Estrutura de banco definida

### 3. Execução Completa com Docker
```bash
# Iniciar todos os serviços
docker compose -f docker-compose.dev.yml up --build

# Executar apenas o banco (para desenvolvimento local)
docker compose -f docker-compose.dev.yml up mysql -d
```

### 4. Executar Migrations
```powershell
# Migration 1: Estrutura inicial
Get-Content "backend/migrations/001_initial_tables.sql" | docker exec -i mysql_db mysql -uroot -proot123

# Migration 2: Usuário admin
Get-Content "backend/migrations/002_create_admin_user.sql" | docker exec -i mysql_db mysql -uroot -proot123

# Migration 3: Sistema de convites
Get-Content "backend/migrations/003_invite_system.sql" | docker exec -i mysql_db mysql -uroot -proot123
```

### 5. Desenvolvimento Local (Opcional)
```bash
cd backend
go run cmd/main.go
```

---

## 📡 API Endpoints

### 🔍 Testes e Monitoramento
| Método | Endpoint | Descrição |
|--------|----------|-----------|
| `GET` | `/ping` | Teste básico da API |
| `GET` | `/database/test` | Teste de conexão com banco |
| `GET` | `/database/tables` | Lista todas as tabelas |
| `GET` | `/database/table/{nome}` | Estrutura de uma tabela |

### 👥 Autenticação e Usuários (Em Desenvolvimento)
| Método | Endpoint | Descrição |
|--------|----------|-----------|
| `POST` | `/auth/cadastro-com-codigo` | Cadastro usando código da empresa |
| `POST` | `/auth/aceitar-convite` | Aceitar convite por token |
| `POST` | `/auth/criar-convite` | Criar convite (apenas gestores) |
| `GET` | `/auth/verificar-codigo/{codigo}` | Validar código da empresa |
| `GET` | `/auth/solicitacoes-pendentes` | Listar solicitações (gestores) |
| `POST` | `/auth/aprovar-solicitacao/{id}` | Aprovar/rejeitar solicitação |

### 📊 Exemplos de Uso

#### Testar Conexão
```bash
curl http://localhost:8080/database/test
```

#### Listar Tabelas
```bash
curl http://localhost:8080/database/tables
```

#### Ver Estrutura de Usuários
```bash
curl http://localhost:8080/database/table/usuarios
```

---

## 🗄️ Estrutura do Banco de Dados

### 📋 Tabelas Principais

#### `usuarios` - Usuários do Sistema
- Informações pessoais (nome, email, telefone, CPF)
- Controle de acesso (role, empresa, status)
- Segurança (tentativas de login, bloqueios)
- Auditoria (criação, última atualização)

#### `roles` - Papéis e Permissões
- **Entregador**: Acesso limitado aos próprios dados
- **Ajudante**: Acesso às entregas onde participa
- **Gestor**: Acesso completo ao sistema

#### `empresas` - Organizações
- Dados da empresa (nome, CNPJ)
- Código único para convites
- Configurações personalizadas

#### `convites_usuario` - Sistema de Convites
- Convites por email com tokens únicos
- Controle de expiração e uso
- Relacionamento com role específico

#### `solicitacoes_cadastro` - Aprovações
- Solicitações de cadastro pendentes
- Justificativas e observações
- Controle de aprovação por gestores

### 🔧 Dados Pré-configurados

#### Empresa Padrão
- **Nome**: Gestão Telemetria LTDA
- **Código**: `TELEME2025`
- **CNPJ**: 12.345.678/0001-99

#### Usuário Administrador
- **Email**: admin@gestaotelemetria.com
- **Role**: Gestor
- **Senha**: admin123 (hash: $2a$10$N9q...)

---

## 📁 Estrutura do Projeto

```
trabiju-telemetria/
├── 📄 .env                          # Configurações de ambiente
├── 🐳 docker-compose.dev.yml        # Configuração Docker
├── 📋 README.md                     # Este arquivo
└── backend/
    ├── 📁 cmd/
    │   └── main.go                  # Ponto de entrada
    ├── 📁 internal/
    │   ├── 📁 config/               # Configurações
    │   │   └── config.go
    │   ├── 📁 controllers/          # Controladores HTTP
    │   │   ├── auth.go              # Autenticação
    │   │   ├── database.go          # Testes de banco
    │   │   └── ping.go              # Teste básico
    │   ├── 📁 database/             # Conexão e schemas
    │   │   ├── connection.go        # Conexão MySQL
    │   │   └── schema/
    │   │       └── init.sql         # Schema inicial
    │   ├── 📁 models/               # Estruturas de dados
    │   │   └── user.go              # Models de usuário
    │   ├── 📁 routes/               # Definição de rotas
    │   │   └── routes.go
    │   └── 📁 [middleware, repository, services, utils]/
    ├── 📁 migrations/               # Migrations SQL
    │   ├── 001_initial_tables.sql   # Estrutura base
    │   ├── 002_create_admin_user.sql # Usuário admin
    │   ├── 003_invite_system.sql    # Sistema convites
    │   └── README.md                # Guia de migrations
    ├── 📁 docs/                     # Documentação
    │   └── SISTEMA_CADASTRO.md      # Fluxos de cadastro
    ├── 🔧 .air.toml                 # Configuração hot reload
    ├── 🐳 Dockerfile.dev            # Container desenvolvimento
    ├── 📦 go.mod                    # Dependências Go
    └── 📦 go.sum                    # Lock das dependências
```

---

## 🛣️ Roadmap de Desenvolvimento

### ✅ Concluído
- [x] Estrutura inicial com Docker
- [x] Conexão MySQL funcional
- [x] Sistema de migrations versionadas
- [x] Arquitetura clean com separação de responsabilidades
- [x] Estrutura completa de usuários e roles
- [x] Sistema de empresas e códigos de convite
- [x] Endpoints de teste e monitoramento
- [x] Logs de auditoria e controle de sessões

### 🚧 Em Desenvolvimento
- [ ] Implementação completa da autenticação JWT
- [ ] Endpoints de cadastro e login
- [ ] Middleware de autorização por role
- [ ] Sistema de convites funcional
- [ ] Validações avançadas de CPF/CNPJ

### 📋 Próximas Implementações
- [ ] API de telemetria para dispositivos
- [ ] Dashboard web para gestores
- [ ] Relatórios e analytics
- [ ] Sistema de notificações
- [ ] Integração com mapas
- [ ] Frontend React/Next.js

### 🚀 Futuro
- [ ] Deploy em cloud (AWS/GCP)
- [ ] Monitoramento com Prometheus/Grafana
- [ ] CI/CD com GitHub Actions
- [ ] Testes automatizados
- [ ] API Gateway
- [ ] Microserviços

---

## 🧪 Testes

### Verificação da Instalação
```bash
# 1. Testar API básica
curl http://localhost:8080/ping
# Resposta esperada: {"message":"pong"}

# 2. Testar conexão com banco
curl http://localhost:8080/database/test
# Resposta: status "success" com versão do MySQL

# 3. Verificar tabelas criadas
curl http://localhost:8080/database/tables
# Resposta: 8 tabelas (usuarios, roles, empresas, etc.)
```

### Verificação no Banco
```sql
-- Conectar no MySQL
docker exec -it mysql_db mysql -uroot -proot123

-- Verificar dados
USE gestao_telemetria;
SELECT * FROM roles;
SELECT * FROM empresas;
SELECT nome, email, role_id FROM usuarios;
```

---

## 🔒 Segurança

### Implementado
- ✅ Senhas hasheadas (bcrypt)
- ✅ Controle de tentativas de login
- ✅ Bloqueio temporário por tentativas
- ✅ Logs de auditoria completos
- ✅ Validação de entrada de dados
- ✅ Isolamento por empresa

### Planejado
- 🔄 Autenticação JWT
- 🔄 Rate limiting
- 🔄 CORS configurado
- 🔄 Validação HTTPS
- 🔄 Sanitização de inputs

---

## 📚 Documentação Adicional

- **[Sistema de Cadastro](backend/docs/SISTEMA_CADASTRO.md)** - Fluxos detalhados de cadastro e autenticação
- **[Guia de Migrations](backend/migrations/README.md)** - Como gerenciar o banco de dados

---

## 🤝 Contribuição

1. Faça um fork do projeto
2. Crie uma branch para sua feature (`git checkout -b feature/AmazingFeature`)
3. Commit suas mudanças (`git commit -m 'Add some AmazingFeature'`)
4. Push para a branch (`git push origin feature/AmazingFeature`)
5. Abra um Pull Request

---

## 🧙 Autor

**Paulo Chiaradia**  
🔗 [github.com/paulochiaradia](https://github.com/paulochiaradia)

---

## 📄 Licença

Este projeto está sob a licença MIT. Veja o arquivo [LICENSE](LICENSE) para detalhes.

---

> 🚧 **Status**: Projeto em desenvolvimento ativo. Sistema de autenticação e banco de dados funcionais. Próxima etapa: implementação completa da API de usuários.

---

**Última atualização**: 22 de Julho de 2025
