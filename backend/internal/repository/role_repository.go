package repository

import (
	"database/sql"

	"github.com/paulochiaradia/trabiju-telemetria/internal/models"
)

type RoleRepository struct {
	db *sql.DB
}

func NewRoleRepository(db *sql.DB) *RoleRepository {
	return &RoleRepository{db: db}
}

// GetRoleByName busca role por nome
func (r *RoleRepository) GetRoleByName(name string) (*models.Role, error) {
	role := &models.Role{}
	query := `
		SELECT id, nome, descricao, permissoes, created_at, updated_at
		FROM roles 
		WHERE nome = ?
	`

	row := r.db.QueryRow(query, name)

	err := row.Scan(
		&role.ID, &role.Nome, &role.Descricao, &role.Permissoes,
		&role.CreatedAt, &role.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return role, nil
}

// GetRoleByID busca role por ID
func (r *RoleRepository) GetRoleByID(id int) (*models.Role, error) {
	role := &models.Role{}
	query := `
		SELECT id, nome, descricao, permissoes, created_at, updated_at
		FROM roles 
		WHERE id = ?
	`

	row := r.db.QueryRow(query, id)

	err := row.Scan(
		&role.ID, &role.Nome, &role.Descricao, &role.Permissoes,
		&role.CreatedAt, &role.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return role, nil
}

// GetAllRoles busca todas as roles
func (r *RoleRepository) GetAllRoles() ([]models.Role, error) {
	var roles []models.Role
	query := `
		SELECT id, nome, descricao, permissoes, created_at, updated_at
		FROM roles 
		ORDER BY nome
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var role models.Role
		err := rows.Scan(
			&role.ID, &role.Nome, &role.Descricao, &role.Permissoes,
			&role.CreatedAt, &role.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		roles = append(roles, role)
	}

	return roles, nil
}
