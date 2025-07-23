# Makefile para o projeto Trabiju Telemetria

# Variáveis
BINARY_NAME=trabiju
MAIN_PATH=./cmd/main.go
BUILD_DIR=./build

# Comandos Go
.PHONY: build test test-verbose test-coverage clean run dev docker-build docker-run

# Build da aplicação
build:
	@echo "🔨 Building application..."
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)
	@echo "✅ Build completed: $(BUILD_DIR)/$(BINARY_NAME)"

# Executar aplicação
run:
	@echo "🚀 Running application..."
	@go run $(MAIN_PATH)

# Desenvolvimento com hot reload (requer air: go install github.com/cosmtrek/air@latest)
dev:
	@echo "🔥 Starting development server with hot reload..."
	@air

# Executar todos os testes
test:
	@echo "🧪 Running tests..."
	@go test ./...

# Executar testes com verbose
test-verbose:
	@echo "🧪 Running tests with verbose output..."
	@go test -v ./...

# Executar testes com coverage
test-coverage:
	@echo "🧪 Running tests with coverage..."
	@go test -v -cover ./...
	@go test -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "📊 Coverage report generated: coverage.html"

# Executar testes específicos
test-api:
	@echo "🧪 Running API tests..."
	@go test -v ./cmd/...

test-services:
	@echo "🧪 Running service tests..."
	@go test -v ./internal/services/...

test-validation:
	@echo "🧪 Running validation tests..."
	@go test -v ./internal/validation/...

# Limpar arquivos gerados
clean:
	@echo "🧹 Cleaning build artifacts..."
	@rm -rf $(BUILD_DIR)
	@rm -f coverage.out coverage.html
	@echo "✅ Clean completed"

# Atualizar dependências
deps:
	@echo "📦 Updating dependencies..."
	@go mod tidy
	@go mod download

# Verificar formato do código
fmt:
	@echo "📝 Formatting code..."
	@go fmt ./...

# Verificar código com vet
vet:
	@echo "🔍 Running go vet..."
	@go vet ./...

# Executar linter (requer golangci-lint)
lint:
	@echo "🔍 Running linter..."
	@golangci-lint run

# Build Docker
docker-build:
	@echo "🐳 Building Docker image..."
	@docker build -t trabiju-backend .
	@echo "✅ Docker image built: trabiju-backend"

# Executar Docker
docker-run:
	@echo "🐳 Running Docker container..."
	@docker run -p 8081:8081 trabiju-backend

# Executar Docker Compose (desenvolvimento)
docker-dev:
	@echo "🐳 Starting development environment..."
	@docker-compose -f docker-compose.dev.yml up --build

# Parar Docker Compose
docker-stop:
	@echo "🛑 Stopping development environment..."
	@docker-compose -f docker-compose.dev.yml down

# Pipeline completa de CI
ci: fmt vet test
	@echo "✅ CI pipeline completed successfully"

# Pipeline de desenvolvimento
dev-check: fmt vet test-verbose
	@echo "✅ Development checks completed"

# Ajuda
help:
	@echo "Trabiju Telemetria - Available commands:"
	@echo ""
	@echo "🔨 Build & Run:"
	@echo "  build         - Build the application"
	@echo "  run           - Run the application"
	@echo "  dev           - Start development server with hot reload"
	@echo ""
	@echo "🧪 Testing:"
	@echo "  test          - Run all tests"
	@echo "  test-verbose  - Run tests with verbose output"
	@echo "  test-coverage - Run tests with coverage report"
	@echo "  test-api      - Run API tests only"
	@echo "  test-services - Run service tests only"
	@echo "  test-validation - Run validation tests only"
	@echo ""
	@echo "🔍 Code Quality:"
	@echo "  fmt           - Format code"
	@echo "  vet           - Run go vet"
	@echo "  lint          - Run linter"
	@echo ""
	@echo "🐳 Docker:"
	@echo "  docker-build  - Build Docker image"
	@echo "  docker-run    - Run Docker container"
	@echo "  docker-dev    - Start development environment"
	@echo "  docker-stop   - Stop development environment"
	@echo ""
	@echo "🛠️  Utilities:"
	@echo "  deps          - Update dependencies"
	@echo "  clean         - Clean build artifacts"
	@echo "  ci            - Run CI pipeline"
	@echo "  dev-check     - Run development checks"
	@echo "  help          - Show this help"
