#!/bin/bash

# Script para inicializar o ambiente de desenvolvimento
# Garante que o binário seja compilado na primeira execução

echo "🚀 Iniciando ambiente de desenvolvimento..."

# Cria o diretório tmp se não existir
mkdir -p tmp

# Verifica se o binário existe, se não existir, compila
if [ ! -f "./tmp/trabiju" ]; then
    echo "📦 Compilando aplicação pela primeira vez..."
    go build -o ./tmp/trabiju ./cmd
    echo "✅ Compilação inicial concluída!"
fi

# Inicia o Air para hot reload
echo "🔥 Iniciando Air para hot reload..."
air -c .air.toml
