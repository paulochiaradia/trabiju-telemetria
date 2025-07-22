-- Migration 003: Sistema de convites e códigos de empresa
-- Created: 2025-07-22
-- Description: Tabelas para controlar cadastro de usuários com roles específicos

USE trabiju_telemetria;

-- Tabela de empresas/organizações
CREATE TABLE empresas (
    id INT PRIMARY KEY AUTO_INCREMENT,
    nome VARCHAR(100) NOT NULL,
    cnpj VARCHAR(18) UNIQUE,
    codigo_convite VARCHAR(10) UNIQUE NOT NULL, -- Código único para cadastros
    ativa BOOLEAN DEFAULT TRUE,
    configuracoes JSON, -- Configurações específicas da empresa
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_codigo_convite (codigo_convite),
    INDEX idx_cnpj (cnpj)
);

-- Inserir empresa padrão
INSERT INTO empresas (nome, cnpj, codigo_convite, configuracoes) VALUES 
('Gestão Telemetria LTDA', '12.345.678/0001-99', 'TELEME2025', 
 '{"allow_self_register": true, "default_role": "entregador", "require_approval": false}');

-- Atualizar tabela de usuários para incluir empresa
ALTER TABLE usuarios ADD COLUMN empresa_id INT AFTER role_id;
ALTER TABLE usuarios ADD FOREIGN KEY (empresa_id) REFERENCES empresas(id);
ALTER TABLE usuarios ADD INDEX idx_empresa (empresa_id);

-- Tabela de convites de usuários
CREATE TABLE convites_usuario (
    id INT PRIMARY KEY AUTO_INCREMENT,
    email VARCHAR(100) NOT NULL,
    role_id INT NOT NULL,
    empresa_id INT NOT NULL,
    convidado_por INT NOT NULL, -- ID do usuário que fez o convite
    token VARCHAR(255) UNIQUE NOT NULL, -- Token único para aceitar convite
    usado BOOLEAN DEFAULT FALSE,
    expira_em TIMESTAMP NOT NULL,
    dados_convite JSON, -- Informações adicionais do convite
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (role_id) REFERENCES roles(id),
    FOREIGN KEY (empresa_id) REFERENCES empresas(id),
    FOREIGN KEY (convidado_por) REFERENCES usuarios(id),
    INDEX idx_email_token (email, token),
    INDEX idx_empresa_role (empresa_id, role_id),
    INDEX idx_expiracao (expira_em, usado)
);

-- Tabela de solicitações de cadastro (para aprovação manual)
CREATE TABLE solicitacoes_cadastro (
    id INT PRIMARY KEY AUTO_INCREMENT,
    nome VARCHAR(100) NOT NULL,
    email VARCHAR(100) NOT NULL,
    telefone VARCHAR(20),
    cpf VARCHAR(14),
    role_solicitado INT NOT NULL,
    empresa_id INT NOT NULL,
    codigo_usado VARCHAR(10), -- Código da empresa usado
    justificativa TEXT, -- Por que quer este role
    status ENUM('pendente', 'aprovado', 'rejeitado') DEFAULT 'pendente',
    aprovado_por INT NULL, -- ID do gestor que aprovou
    observacoes_aprovacao TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (role_solicitado) REFERENCES roles(id),
    FOREIGN KEY (empresa_id) REFERENCES empresas(id),
    FOREIGN KEY (aprovado_por) REFERENCES usuarios(id),
    INDEX idx_status_empresa (status, empresa_id),
    INDEX idx_email (email)
);

-- Atualizar usuário admin para ter empresa
UPDATE usuarios 
SET empresa_id = (SELECT id FROM empresas WHERE codigo_convite = 'TELEME2025')
WHERE email = 'admin@gestaotelemetria.com';
