-- Migration 002: Create admin user
-- Created: 2025-07-22
-- Description: Criar usuário administrador inicial

USE trabiju_telemetria;

-- Inserir usuário administrador padrão
INSERT INTO usuarios (
    nome, 
    email, 
    senha, 
    telefone,
    cpf,
    role_id,
    ativo,
    api_token
) VALUES (
    'Administrador do Sistema',
    'admin@gestaotelemetria.com',
    '$2a$10$N9qo8uLOickgx2ZMRZoMye.mQKyVZV8Cgp1.JEE8SL8BI7.FjAy6.', -- senha: admin123
    '(11) 99999-9999',
    '000.000.000-00',
    (SELECT id FROM roles WHERE nome = 'gestor'),
    TRUE,
    UUID()
);

-- Inserir configuração inicial para o admin
UPDATE usuarios 
SET configuracoes_dashboard = '{"theme": "light", "language": "pt-BR", "notifications": true}'
WHERE email = 'admin@gestaotelemetria.com';
