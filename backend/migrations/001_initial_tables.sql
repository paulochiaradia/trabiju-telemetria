-- Migration 001: Create initial tables
-- Created: 2025-07-22
-- Description: Sistema completo de usuários e telemetria

-- Usar o banco de dados existente (criado pelo docker-compose)
USE trabiju_telemetria;

-- Tabela de tipos/roles de usuários
CREATE TABLE roles (
    id INT PRIMARY KEY AUTO_INCREMENT,
    nome VARCHAR(50) NOT NULL UNIQUE,
    descricao TEXT,
    permissoes JSON,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

-- Dados iniciais dos roles
INSERT INTO roles (nome, descricao, permissoes) VALUES 
('entregador', 'Acesso limitado aos próprios dados', 
 '{"dashboard": {"own_data": true, "all_data": false}, "vehicles": {"view_assigned": true}, "reports": {"own_only": true}}'),
('ajudante', 'Acesso às entregas onde participa', 
 '{"dashboard": {"assigned_deliveries": true, "all_data": false}, "vehicles": {"view_assigned": true}, "reports": {"assigned_only": true}}'),
('gestor', 'Acesso completo ao sistema', 
 '{"dashboard": {"all_data": true}, "vehicles": {"full_access": true}, "reports": {"full_access": true}, "users": {"manage": true}}');

-- Tabela de usuários
CREATE TABLE usuarios (
    id INT PRIMARY KEY AUTO_INCREMENT,
    nome VARCHAR(100) NOT NULL,
    email VARCHAR(100) NOT NULL UNIQUE,
    senha VARCHAR(255) NOT NULL,
    telefone VARCHAR(20),
    cpf VARCHAR(14) UNIQUE,
    avatar VARCHAR(255),
    role_id INT NOT NULL,
    ativo BOOLEAN DEFAULT TRUE,
    ultimo_login TIMESTAMP NULL,
    configuracoes_dashboard JSON,
    api_token VARCHAR(255) UNIQUE,
    tentativas_login INT DEFAULT 0,
    bloqueado_ate TIMESTAMP NULL,
    senha_alterada_em TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (role_id) REFERENCES roles(id),
    INDEX idx_email (email),
    INDEX idx_api_token (api_token),
    INDEX idx_role_ativo (role_id, ativo),
    INDEX idx_cpf (cpf)
);

-- Tabela de sessões
CREATE TABLE sessoes_usuario (
    id VARCHAR(128) PRIMARY KEY,
    usuario_id INT NOT NULL,
    ip_address VARCHAR(45),
    user_agent TEXT,
    dados_sessao JSON,
    expira_em TIMESTAMP NOT NULL,
    ativo BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (usuario_id) REFERENCES usuarios(id) ON DELETE CASCADE,
    INDEX idx_usuario (usuario_id),
    INDEX idx_expiracao (expira_em),
    INDEX idx_usuario_ativo (usuario_id, ativo)
);

-- Tabela de logs de autenticação
CREATE TABLE logs_autenticacao (
    id INT PRIMARY KEY AUTO_INCREMENT,
    usuario_id INT,
    email_tentativa VARCHAR(100) NOT NULL,
    sucesso BOOLEAN NOT NULL,
    ip_address VARCHAR(45),
    user_agent TEXT,
    motivo_falha VARCHAR(100),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (usuario_id) REFERENCES usuarios(id) ON DELETE SET NULL,
    INDEX idx_usuario_data (usuario_id, created_at),
    INDEX idx_email_data (email_tentativa, created_at),
    INDEX idx_sucesso_data (sucesso, created_at)
);

-- Tabela de telemetria
CREATE TABLE telemetry_data (
    id INT AUTO_INCREMENT PRIMARY KEY,
    device_id VARCHAR(100) NOT NULL,
    usuario_id INT,
    data JSON NOT NULL,
    timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    processed BOOLEAN DEFAULT FALSE,
    FOREIGN KEY (usuario_id) REFERENCES usuarios(id) ON DELETE SET NULL,
    INDEX idx_device_timestamp (device_id, timestamp),
    INDEX idx_usuario_timestamp (usuario_id, timestamp),
    INDEX idx_processed (processed)
);
