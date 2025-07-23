# Dockerfile customizado para MySQL com tzdata
FROM mysql:8.0

# Instalar tzdata usando o gerenciador de pacotes correto
# MySQL 8.0 usa Oracle Linux que tem microdnf
RUN microdnf update -y && \
    microdnf install -y tzdata && \
    microdnf clean all

# Configurar timezone padrão para UTC
ENV TZ=UTC

# Expor porta padrão do MySQL
EXPOSE 3306

# Usar o entrypoint padrão do MySQL
ENTRYPOINT ["docker-entrypoint.sh"]
CMD ["mysqld"]
