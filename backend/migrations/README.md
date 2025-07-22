# Guia de Migrations

## Como usar as migrations deste projeto

### Ordem de execução:
1. `001_initial_tables.sql` - Estrutura base do banco
2. `002_create_admin_user.sql` - Usuário administrador inicial

### Aplicar manualmente:
```bash
# Conectar no MySQL
mysql -u root -p

# Executar as migrations em ordem
source backend/migrations/001_initial_tables.sql
source backend/migrations/002_create_admin_user.sql
```

### Verificar se foi aplicada:
```sql
SELECT * FROM usuarios WHERE email = 'admin@gestaotelemetria.com';
```

### Futuras migrations:
- `003_add_vehicles_table.sql`
- `004_add_delivery_routes.sql`
- `005_add_telemetry_processing.sql`

### Convenções:
- Sempre numerar sequencialmente (001, 002, 003...)
- Um propósito por migration
- Incluir comentários explicativos
- Testar antes de aplicar em produção

### Rollback:
Para desfazer mudanças, criar migrations específicas:
- `006_rollback_vehicles_table.sql`
