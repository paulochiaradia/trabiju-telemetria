# backend/Dockerfile

FROM golang:1.21.1

# Diretório de trabalho dentro do container
WORKDIR /app

# Copia todos os arquivos para dentro do container
COPY . .

# Inicializa o módulo Go e baixa as dependências
RUN go mod tidy

# Compila o projeto a partir do arquivo principal
RUN go build -o trabiju ./cmd/main.go

# Expõe a porta usada pelo Gin
EXPOSE 8080

# Comando para executar o binário
CMD ["./trabiju"]
