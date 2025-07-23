package repository

import (
	"database/sql"
	"time"

	"github.com/paulochiaradia/trabiju-telemetria/internal/models"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

// CreateUser cria um novo usuário
func (r *UserRepository) CreateUser(user *models.Usuario) error {
	query := `
		INSERT INTO usuarios (
			nome, email, senha, telefone, cpf, avatar, role_id, empresa_id, 
			ativo, configuracoes_dashboard, api_token, senha_alterada_em,
			created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NOW(), NOW())
	`

	result, err := r.db.Exec(query,
		user.Nome, user.Email, user.Senha, user.Telefone, user.CPF,
		user.Avatar, user.RoleID, user.EmpresaID, user.Ativo,
		user.ConfiguracoesDashboard, user.APIToken, time.Now(),
	)

	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	user.ID = int(id)
	return nil
}

// GetUserByEmail busca usuário por email
func (r *UserRepository) GetUserByEmail(email string) (*models.Usuario, error) {
	user := &models.Usuario{}
	query := `
		SELECT u.id, u.nome, u.email, u.senha, u.telefone, u.cpf, u.avatar,
			   u.role_id, u.empresa_id, u.ativo, u.ultimo_login, 
			   u.configuracoes_dashboard, u.api_token, u.tentativas_login,
			   u.bloqueado_ate, u.senha_alterada_em, u.created_at, u.updated_at,
			   r.nome as role_nome, r.descricao as role_descricao,
			   e.nome as empresa_nome, e.codigo_convite
		FROM usuarios u
		LEFT JOIN roles r ON u.role_id = r.id
		LEFT JOIN empresas e ON u.empresa_id = e.id
		WHERE u.email = ?
	`

	row := r.db.QueryRow(query, email)

	var ultimoLogin, bloqueadoAte sql.NullTime
	var roleNome, roleDescricao, empresaNome, codigoConvite sql.NullString

	err := row.Scan(
		&user.ID, &user.Nome, &user.Email, &user.Senha, &user.Telefone,
		&user.CPF, &user.Avatar, &user.RoleID, &user.EmpresaID, &user.Ativo,
		&ultimoLogin, &user.ConfiguracoesDashboard, &user.APIToken,
		&user.TentativasLogin, &bloqueadoAte, &user.SenhaAlteradaEm,
		&user.CreatedAt, &user.UpdatedAt, &roleNome, &roleDescricao,
		&empresaNome, &codigoConvite,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	// Converter NullTime para *time.Time
	if ultimoLogin.Valid {
		user.UltimoLogin = &ultimoLogin.Time
	}
	if bloqueadoAte.Valid {
		user.BloqueadoAte = &bloqueadoAte.Time
	}

	// Adicionar relacionamentos se existirem
	if roleNome.Valid {
		user.Role = &models.Role{
			ID:        user.RoleID,
			Nome:      roleNome.String,
			Descricao: roleDescricao.String,
		}
	}

	if empresaNome.Valid {
		user.Empresa = &models.Empresa{
			ID:            user.EmpresaID,
			Nome:          empresaNome.String,
			CodigoConvite: codigoConvite.String,
		}
	}

	return user, nil
}

// GetUserByID busca usuário por ID
func (r *UserRepository) GetUserByID(id int) (*models.Usuario, error) {
	user := &models.Usuario{}
	query := `
		SELECT u.id, u.nome, u.email, u.telefone, u.cpf, u.avatar,
			   u.role_id, u.empresa_id, u.ativo, u.ultimo_login, 
			   u.configuracoes_dashboard, u.api_token, u.tentativas_login,
			   u.bloqueado_ate, u.senha_alterada_em, u.created_at, u.updated_at
		FROM usuarios u
		WHERE u.id = ?
	`

	row := r.db.QueryRow(query, id)

	var ultimoLogin, bloqueadoAte sql.NullTime

	err := row.Scan(
		&user.ID, &user.Nome, &user.Email, &user.Telefone,
		&user.CPF, &user.Avatar, &user.RoleID, &user.EmpresaID, &user.Ativo,
		&ultimoLogin, &user.ConfiguracoesDashboard, &user.APIToken,
		&user.TentativasLogin, &bloqueadoAte, &user.SenhaAlteradaEm,
		&user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	if ultimoLogin.Valid {
		user.UltimoLogin = &ultimoLogin.Time
	}
	if bloqueadoAte.Valid {
		user.BloqueadoAte = &bloqueadoAte.Time
	}

	return user, nil
}

// EmailExists verifica se email já existe
func (r *UserRepository) EmailExists(email string) (bool, error) {
	var count int
	query := "SELECT COUNT(*) FROM usuarios WHERE email = ?"
	err := r.db.QueryRow(query, email).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// CPFExists verifica se CPF já existe
func (r *UserRepository) CPFExists(cpf string) (bool, error) {
	var count int
	query := "SELECT COUNT(*) FROM usuarios WHERE cpf = ?"
	err := r.db.QueryRow(query, cpf).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// UpdateLastLogin atualiza último login
func (r *UserRepository) UpdateLastLogin(userID int) error {
	query := "UPDATE usuarios SET ultimo_login = NOW() WHERE id = ?"
	_, err := r.db.Exec(query, userID)
	return err
}

// UpdateLoginAttempts atualiza tentativas de login
func (r *UserRepository) UpdateLoginAttempts(userID int, attempts int) error {
	query := "UPDATE usuarios SET tentativas_login = ? WHERE id = ?"
	_, err := r.db.Exec(query, attempts, userID)
	return err
}

// BlockUser bloqueia usuário por tempo determinado
func (r *UserRepository) BlockUser(userID int, until time.Time) error {
	query := "UPDATE usuarios SET bloqueado_ate = ? WHERE id = ?"
	_, err := r.db.Exec(query, until, userID)
	return err
}

// ActivateUser ativa usuário (confirma email)
func (r *UserRepository) ActivateUser(userID int) error {
	query := "UPDATE usuarios SET ativo = TRUE WHERE id = ?"
	_, err := r.db.Exec(query, userID)
	return err
}
