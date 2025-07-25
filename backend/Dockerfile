# backend/Dockerfile - Produção otimizada com multi-stage build

# Stage 1: Build
FROM golang:1.21.1-alpine AS builder

# Instalar dependências necessárias
RUN apk add --no-cache git ca-certificates tzdata

# Diretório de trabalho
WORKDIR /build

# Copiar arquivos de dependências primeiro (cache layer)
COPY go.mod go.sum ./

# Baixar dependências
RUN go mod download

# Copiar código fonte
COPY . .

# Compilar aplicação (binário estático)
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o trabiju ./cmd/main.go

# Stage 2: Runtime
FROM alpine:latest

# Instalar certificados, timezone e curl
RUN apk --no-cache add ca-certificates tzdata curl

# Criar usuário não-root
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

# Diretório de trabalho
WORKDIR /app

# Copiar binário do stage de build
COPY --from=builder /build/trabiju .

# Mudar ownership e permissões
RUN chown appuser:appgroup trabiju && chmod +x trabiju

# Usar usuário não-root
USER appuser

# Expor porta
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=40s --retries=3 \
  CMD curl -f http://localhost:8080/api/v1/health || exit 1

# Comando para executar
CMD ["./trabiju"]
