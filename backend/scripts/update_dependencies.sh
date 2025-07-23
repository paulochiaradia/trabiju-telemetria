#!/bin/bash

# Script para atualizar dependências do Go
echo "🔧 Atualizando dependências do projeto..."

# Ir para o diretório do backend
cd backend

# Adicionar dependências necessárias
echo "📦 Adicionando dependências JWT..."
go get github.com/golang-jwt/jwt/v5

echo "📦 Adicionando dependências de criptografia..."
go get golang.org/x/crypto/bcrypt

echo "📦 Adicionando dependências CORS..."
go get github.com/gin-contrib/cors

echo "📦 Verificando dependências existentes..."
go mod tidy

echo "✅ Dependências atualizadas com sucesso!"
echo ""
echo "📋 Dependências adicionadas:"
echo "  - github.com/golang-jwt/jwt/v5 (JWT)"
echo "  - golang.org/x/crypto/bcrypt (Hash de senhas)"
echo "  - github.com/gin-contrib/cors (CORS)"
echo ""
echo "🚀 Projeto pronto para compilar!"
