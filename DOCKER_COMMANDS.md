# Comandos Docker para Trabiju Telemetria

## 🚀 SUBIR A APLICAÇÃO

# Subir todos os serviços (primeira vez ou após mudanças)
docker-compose up --build

# Subir em background (detached mode)
docker-compose up --build -d

# Apenas subir (sem rebuild)
docker-compose up -d

## 📋 GERENCIAMENTO

# Ver logs em tempo real
docker-compose logs -f

# Ver logs apenas do backend
docker-compose logs -f backend

# Ver logs apenas do MySQL
docker-compose logs -f mysql

# Ver status dos containers
docker-compose ps

# Parar todos os serviços
docker-compose down

# Parar e remover volumes (CUIDADO: apaga dados do banco)
docker-compose down -v

## 🔧 COMANDOS DE DESENVOLVIMENTO

# Rebuildar apenas o backend
docker-compose build backend

# Reiniciar apenas o backend
docker-compose restart backend

# Executar comando dentro do container do backend
docker-compose exec backend sh

# Executar comando dentro do container do MySQL
docker-compose exec mysql mysql -u root -p

## 🧹 LIMPEZA

# Parar e remover containers
docker-compose down

# Remover imagens não utilizadas
docker image prune

# Remover volumes órfãos
docker volume prune

# Limpeza completa (CUIDADO!)
docker system prune -a

## 📊 MONITORAMENTO

# Ver uso de recursos
docker stats

# Inspecionar container
docker inspect trabiju-backend

# Ver logs específicos
docker logs go_backend_prod -f

## 🔗 URLs DA APLICAÇÃO

# Backend API
http://localhost:8082

# Health Check
http://localhost:8082/api/v1/ping

# MySQL (se precisar conectar externamente)
mysql://localhost:3306
