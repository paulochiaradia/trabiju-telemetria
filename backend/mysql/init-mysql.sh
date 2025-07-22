#!/bin/bash
# Script para reduzir warnings do MySQL

echo "Configurando MySQL para reduzir warnings..."

# Aguardar o MySQL estar pronto
while ! mysqladmin ping -h"localhost" --silent; do
    echo "Aguardando MySQL ficar disponível..."
    sleep 1
done

# Configurar timezone tables se disponível (reduz warnings de timezone)
if [ -d "/usr/share/zoneinfo" ]; then
    echo "Configurando dados de timezone..."
    mysql_tzinfo_to_sql /usr/share/zoneinfo 2>/dev/null | mysql -u root -p"$MYSQL_ROOT_PASSWORD" mysql 2>/dev/null || echo "Dados de timezone não disponíveis ou já configurados"
fi

echo "Configuração do MySQL concluída"
