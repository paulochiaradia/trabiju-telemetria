package repository

import (
	"database/sql"

	"github.com/paulochiaradia/trabiju-telemetria/internal/models"
)

type CompanyRepository struct {
	db *sql.DB
}

func NewCompanyRepository(db *sql.DB) *CompanyRepository {
	return &CompanyRepository{db: db}
}

// GetCompanyByInviteCode busca empresa por código de convite
func (r *CompanyRepository) GetCompanyByInviteCode(code string) (*models.Empresa, error) {
	empresa := &models.Empresa{}
	query := `
		SELECT id, nome, cnpj, codigo_convite, ativa, configuracoes,
			   created_at, updated_at
		FROM empresas 
		WHERE codigo_convite = ? AND ativa = TRUE
	`

	row := r.db.QueryRow(query, code)

	err := row.Scan(
		&empresa.ID, &empresa.Nome, &empresa.CNPJ, &empresa.CodigoConvite,
		&empresa.Ativa, &empresa.Configuracoes, &empresa.CreatedAt,
		&empresa.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return empresa, nil
}

// GetCompanyByID busca empresa por ID
func (r *CompanyRepository) GetCompanyByID(id int) (*models.Empresa, error) {
	empresa := &models.Empresa{}
	query := `
		SELECT id, nome, cnpj, codigo_convite, ativa, configuracoes,
			   created_at, updated_at
		FROM empresas 
		WHERE id = ?
	`

	row := r.db.QueryRow(query, id)

	err := row.Scan(
		&empresa.ID, &empresa.Nome, &empresa.CNPJ, &empresa.CodigoConvite,
		&empresa.Ativa, &empresa.Configuracoes, &empresa.CreatedAt,
		&empresa.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return empresa, nil
}
