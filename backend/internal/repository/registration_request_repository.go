package repository

import (
	"database/sql"

	"github.com/paulochiaradia/trabiju-telemetria/internal/models"
)

type RegistrationRequestRepository struct {
	db *sql.DB
}

func NewRegistrationRequestRepository(db *sql.DB) *RegistrationRequestRepository {
	return &RegistrationRequestRepository{db: db}
}

// CreateRegistrationRequest cria uma solicitação de cadastro
func (r *RegistrationRequestRepository) CreateRegistrationRequest(req *models.SolicitacaoCadastro) error {
	query := `
		INSERT INTO solicitacoes_cadastro (
			nome, email, telefone, cpf, role_solicitado, empresa_id,
			codigo_usado, justificativa, status, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, 'pendente', NOW(), NOW())
	`

	result, err := r.db.Exec(query,
		req.Nome, req.Email, req.Telefone, req.CPF, req.RoleSolicitado,
		req.EmpresaID, req.CodigoUsado, req.Justificativa,
	)

	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	req.ID = int(id)
	req.Status = "pendente"
	return nil
}

// GetPendingRequestsByCompany busca solicitações pendentes por empresa
func (r *RegistrationRequestRepository) GetPendingRequestsByCompany(empresaID int) ([]models.SolicitacaoCadastro, error) {
	var requests []models.SolicitacaoCadastro
	query := `
		SELECT s.id, s.nome, s.email, s.telefone, s.cpf, s.role_solicitado,
			   s.empresa_id, s.codigo_usado, s.justificativa, s.status,
			   s.aprovado_por, s.observacoes_aprovacao, s.created_at, s.updated_at,
			   r.nome as role_nome
		FROM solicitacoes_cadastro s
		LEFT JOIN roles r ON s.role_solicitado = r.id
		WHERE s.empresa_id = ? AND s.status = 'pendente'
		ORDER BY s.created_at ASC
	`

	rows, err := r.db.Query(query, empresaID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var req models.SolicitacaoCadastro
		var aprovadoPor sql.NullInt64
		var roleName sql.NullString

		err := rows.Scan(
			&req.ID, &req.Nome, &req.Email, &req.Telefone, &req.CPF,
			&req.RoleSolicitado, &req.EmpresaID, &req.CodigoUsado,
			&req.Justificativa, &req.Status, &aprovadoPor,
			&req.ObservacoesAprovacao, &req.CreatedAt, &req.UpdatedAt,
			&roleName,
		)

		if err != nil {
			return nil, err
		}

		if aprovadoPor.Valid {
			aprovadoPorInt := int(aprovadoPor.Int64)
			req.AprovadoPor = &aprovadoPorInt
		}

		if roleName.Valid {
			req.RoleSolicitadoObj = &models.Role{
				ID:   req.RoleSolicitado,
				Nome: roleName.String,
			}
		}

		requests = append(requests, req)
	}

	return requests, nil
}

// ApproveRequest aprova uma solicitação
func (r *RegistrationRequestRepository) ApproveRequest(requestID, approverID int, observacoes string) error {
	query := `
		UPDATE solicitacoes_cadastro 
		SET status = 'aprovado', aprovado_por = ?, observacoes_aprovacao = ?, updated_at = NOW()
		WHERE id = ? AND status = 'pendente'
	`

	result, err := r.db.Exec(query, approverID, observacoes, requestID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

// RejectRequest rejeita uma solicitação
func (r *RegistrationRequestRepository) RejectRequest(requestID, approverID int, observacoes string) error {
	query := `
		UPDATE solicitacoes_cadastro 
		SET status = 'rejeitado', aprovado_por = ?, observacoes_aprovacao = ?, updated_at = NOW()
		WHERE id = ? AND status = 'pendente'
	`

	result, err := r.db.Exec(query, approverID, observacoes, requestID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

// GetRequestByID busca solicitação por ID
func (r *RegistrationRequestRepository) GetRequestByID(id int) (*models.SolicitacaoCadastro, error) {
	req := &models.SolicitacaoCadastro{}
	query := `
		SELECT s.id, s.nome, s.email, s.telefone, s.cpf, s.role_solicitado,
			   s.empresa_id, s.codigo_usado, s.justificativa, s.status,
			   s.aprovado_por, s.observacoes_aprovacao, s.created_at, s.updated_at
		FROM solicitacoes_cadastro s
		WHERE s.id = ?
	`

	row := r.db.QueryRow(query, id)

	var aprovadoPor sql.NullInt64

	err := row.Scan(
		&req.ID, &req.Nome, &req.Email, &req.Telefone, &req.CPF,
		&req.RoleSolicitado, &req.EmpresaID, &req.CodigoUsado,
		&req.Justificativa, &req.Status, &aprovadoPor,
		&req.ObservacoesAprovacao, &req.CreatedAt, &req.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	if aprovadoPor.Valid {
		aprovadoPorInt := int(aprovadoPor.Int64)
		req.AprovadoPor = &aprovadoPorInt
	}

	return req, nil
}
