# Configurações MySQL - Ambiente de Desenvolvimento

## Warnings Resolvidos

Este documento explica as configurações aplicadas ao MySQL para reduzir warnings no ambiente de desenvolvimento.

### Warnings Originais e Suas Soluções

1. **`--skip-host-cache` is deprecated**
   - **Warning**: `The syntax '--skip-host-cache' is deprecated and will be removed in a future release. Please use SET GLOBAL host_cache_size=0 instead.`
   - **Solução**: Substituído por `--host-cache-size=0`
   - **Benefício**: Remove o cache de resolução de hosts, adequado para desenvolvimento

2. **CA certificate ca.pem is self signed**
   - **Warning**: `CA certificate ca.pem is self signed.`
   - **Solução**: Desabilitado SSL com `--tls-version=''`
   - **Benefício**: Remove warnings de SSL em ambiente de desenvolvimento

3. **Insecure configuration for --pid-file**
   - **Warning**: `Insecure configuration for --pid-file: Location '/var/run/mysqld' in the path is accessible to all OS users.`
   - **Solução**: Alterado para `--pid-file=/var/lib/mysql/mysqld.pid`
   - **Benefício**: PID file em localização mais segura

4. **X Plugin SSL warnings**
   - **Warning**: `Plugin mysqlx reported: 'Failed at SSL configuration'`
   - **Solução**: Desabilitado X Plugin com `--mysqlx=OFF`
   - **Benefício**: Remove warnings relacionados ao X Plugin que não é usado em desenvolvimento

5. **Timezone warnings**
   - **Warning**: Vários warnings sobre arquivos de timezone não encontrados
   - **Solução**: Configurado timezone padrão com `--default-time-zone=+00:00`
   - **Benefício**: Define timezone consistente

6. **README.md warning**
   - **Warning**: `/usr/local/bin/docker-entrypoint.sh: ignoring /docker-entrypoint-initdb.d/README.md`
   - **Solução**: Renomeado arquivo para `README.md.bak`
   - **Benefício**: Remove warning sobre arquivos não executáveis na pasta de inicialização

### Configurações Adicionais

- **Character Set**: `utf8mb4` para suporte completo a Unicode
- **Collation**: `utf8mb4_unicode_ci` para ordenação correta
- **Log Verbosity**: Nível 2 para reduzir logs desnecessários
- **Name Resolution**: Desabilitado para melhor performance em desenvolvimento

### Como Aplicar

As configurações são aplicadas automaticamente ao executar:

```bash
docker-compose -f docker-compose.dev.yml up -d
```

### Nota Importante

Estas configurações são otimizadas para **ambiente de desenvolvimento**. Para produção, considere:
- Habilitar SSL/TLS adequadamente
- Configurar autenticação mais robusta
- Ajustar configurações de performance conforme necessário
- Habilitar logs de auditoria se necessário

### Arquivos Modificados

- `docker-compose.dev.yml`: Configurações principais do MySQL
- `backend/migrations/README.md`: Renomeado para evitar warning
- `backend/mysql/`: Pasta criada para configurações personalizadas (opcional)
