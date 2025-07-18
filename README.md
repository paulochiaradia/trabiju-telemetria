# 🚀 Trabiju Telemetria

Projeto backend em Go (Gin) com banco de dados MySQL, rodando totalmente via Docker. Este repositório está estruturado para suportar autenticação segura, arquitetura escalável e fácil integração com o frontend.

---

## 📦 Tecnologias Utilizadas

- [Go 1.21.1](https://golang.org/)
- [Gin Web Framework](https://github.com/gin-gonic/gin)
- [Docker](https://www.docker.com/)
- [Docker Compose](https://docs.docker.com/compose/)
- [MySQL 8.0](https://www.mysql.com/)
- [JWT](https://github.com/golang-jwt/jwt)
- [godotenv](https://github.com/joho/godotenv)

---

## ⚙️ Como Rodar Localmente

1. Clone o repositório:
   ```bash
   git clone https://github.com/paulochiaradia/trabiju-telemetria.git
   cd trabiju-telemetria
   ```

2. Crie um arquivo `.env` na raiz do projeto com:
   ```env
   MYSQL_ROOT_PASSWORD=root
   MYSQL_DATABASE=trabiju
   MYSQL_USER=teste
   MYSQL_PASSWORD=senha123
   ```

3. Suba o ambiente com Docker:
   ```bash
   docker compose up --build
   ```

4. Acesse o endpoint de teste:
   [http://localhost:8080/ping](http://localhost:8080/ping)

---

## 📁 Estrutura de Pastas

```bash
backend/
├── cmd/             # Ponto de entrada (main.go)
├── config/          # Configurações de ambiente e conexão
├── internal/
│   ├── controllers/ # Handlers das rotas
│   ├── routes/      # Registro das rotas
│   ├── models/      # Structs e acesso ao banco
│   ├── middleware/  # Autenticação, CORS, etc.
│   └── utils/       # Funções auxiliares
```

---

## 📡 Rotas Atuais

| Método | Rota     | Descrição      |
|--------|----------|----------------|
| GET    | `/ping`  | Teste da API   |

---

## 🛣️ Roadmap

- [x] Estrutura inicial com Docker e backend
- [ ] MySQL rodando via container
- [ ] Conexão ao banco com Go
- [ ] Cadastro de usuários
- [ ] Autenticação JWT
- [ ] CRUD completo
- [ ] Frontend em React
- [ ] Deploy em ambiente cloud

---

## 🧙 Autor

**Paulo Chiaradia**  
🔗 [github.com/paulochiaradia](https://github.com/paulochiaradia)

---

> Este projeto está em desenvolvimento contínuo. Sinta-se à vontade para contribuir, abrir issues ou dar sugestões. 🚧✨
