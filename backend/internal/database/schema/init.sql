-- Inicialização do banco de dados trabiju_db
-- Sistema de Gestão de Telemetria

-- Criação do banco de dados
CREATE DATABASE IF NOT EXISTS gestao_telemetria;
USE gestao_telemetria;

-- ==========================================
-- 1. TABELAS DE USUÁRIOS E AUTENTICAÇÃO
-- ==========================================

-- Tabela de tipos/roles de usuários
CREATE TABLE roles (
    id INT PRIMARY KEY AUTO_INCREMENT,
    nome VARCHAR(50) NOT NULL UNIQUE,
    descricao TEXT,
    permissoes JSON, -- Armazenar permissões específicas
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

-- Inserir os tipos de usuários com permissões específicas
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
    senha VARCHAR(255) NOT NULL, -- Hash da senha
    telefone VARCHAR(20),
    cpf VARCHAR(14) UNIQUE,
    avatar VARCHAR(255), -- URL da foto do usuário
    role_id INT NOT NULL,
    ativo BOOLEAN DEFAULT TRUE,
    ultimo_login TIMESTAMP NULL,
    configuracoes_dashboard JSON, -- Preferências personalizadas
    api_token VARCHAR(255) UNIQUE, -- Token para API externa
    tentativas_login INT DEFAULT 0, -- Controle de tentativas de login
    bloqueado_ate TIMESTAMP NULL, -- Bloqueio temporário por tentativas
    senha_alterada_em TIMESTAMP DEFAULT CURRENT_TIMESTAMP, -- Controle de expiração de senha
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (role_id) REFERENCES roles(id),
    INDEX idx_email (email),
    INDEX idx_api_token (api_token),
    INDEX idx_role_ativo (role_id, ativo),
    INDEX idx_cpf (cpf)
);

-- Tabela de sessões de usuário
CREATE TABLE sessoes_usuario (
    id VARCHAR(128) PRIMARY KEY,
    usuario_id INT NOT NULL,
    ip_address VARCHAR(45),
    user_agent TEXT,
    dados_sessao JSON,
    expira_em TIMESTAMP NOT NULL,
    ativo BOOLEAN DEFAULT TRUE, -- Para invalidar sessão sem deletar
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (usuario_id) REFERENCES usuarios(id) ON DELETE CASCADE,
    INDEX idx_usuario (usuario_id),
    INDEX idx_expiracao (expira_em),
    INDEX idx_usuario_ativo (usuario_id, ativo)
);

-- Tabela de logs de autenticação (auditoria)
CREATE TABLE logs_autenticacao (
    id INT PRIMARY KEY AUTO_INCREMENT,
    usuario_id INT,
    email_tentativa VARCHAR(100) NOT NULL,
    sucesso BOOLEAN NOT NULL,
    ip_address VARCHAR(45),
    user_agent TEXT,
    motivo_falha VARCHAR(100), -- 'senha_incorreta', 'usuario_bloqueado', etc.
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (usuario_id) REFERENCES usuarios(id) ON DELETE SET NULL,
    INDEX idx_usuario_data (usuario_id, created_at),
    INDEX idx_email_data (email_tentativa, created_at),
    INDEX idx_sucesso_data (sucesso, created_at)
);

-- ==========================================
-- 2. TABELAS DE TELEMETRIA (PLACEHOLDER)
-- ==========================================

-- Exemplo de tabela para dados de telemetria
CREATE TABLE telemetry_data (
    id INT AUTO_INCREMENT PRIMARY KEY,
    device_id VARCHAR(100) NOT NULL,
    usuario_id INT, -- Relacionar com o usuário responsável
    data JSON NOT NULL,
    timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    processed BOOLEAN DEFAULT FALSE, -- Se os dados foram processados
    FOREIGN KEY (usuario_id) REFERENCES usuarios(id) ON DELETE SET NULL,
    INDEX idx_device_timestamp (device_id, timestamp),
    INDEX idx_usuario_timestamp (usuario_id, timestamp),
    INDEX idx_processed (processed)
);
