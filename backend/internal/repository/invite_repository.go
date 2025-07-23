package repository

import (
	"database/sql"

	"github.com/paulochiaradia/trabiju-telemetria/internal/models"
)

type InviteRepository struct {
	db *sql.DB
}

func NewInviteRepository(db *sql.DB) *InviteRepository {
	return &InviteRepository{db: db}
}

// CreateInvite cria um novo convite
func (r *InviteRepository) CreateInvite(invite *models.ConviteUsuario) error {
	query := `
		INSERT INTO convites_usuario (
			email, role_id, empresa_id, convidado_por, token, 
			expira_em, dados_convite, created_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, NOW())
	`

	result, err := r.db.Exec(query,
		invite.Email, invite.RoleID, invite.EmpresaID, invite.ConvidadoPor,
		invite.Token, invite.ExpiraEm, invite.DadosConvite,
	)

	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	invite.ID = int(id)
	return nil
}

// GetInviteByToken busca convite por token
func (r *InviteRepository) GetInviteByToken(token string) (*models.ConviteUsuario, error) {
	invite := &models.ConviteUsuario{}
	query := `
		SELECT c.id, c.email, c.role_id, c.empresa_id, c.convidado_por,
			   c.token, c.usado, c.expira_em, c.dados_convite, c.created_at,
			   r.nome as role_nome, e.nome as empresa_nome,
			   u.nome as convidado_por_nome
		FROM convites_usuario c
		LEFT JOIN roles r ON c.role_id = r.id
		LEFT JOIN empresas e ON c.empresa_id = e.id
		LEFT JOIN usuarios u ON c.convidado_por = u.id
		WHERE c.token = ? AND c.usado = FALSE AND c.expira_em > NOW()
	`

	row := r.db.QueryRow(query, token)

	var roleNome, empresaNome, convidadoPorNome sql.NullString

	err := row.Scan(
		&invite.ID, &invite.Email, &invite.RoleID, &invite.EmpresaID,
		&invite.ConvidadoPor, &invite.Token, &invite.Usado, &invite.ExpiraEm,
		&invite.DadosConvite, &invite.CreatedAt, &roleNome, &empresaNome,
		&convidadoPorNome,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	// Adicionar relacionamentos
	if roleNome.Valid {
		invite.Role = &models.Role{
			ID:   invite.RoleID,
			Nome: roleNome.String,
		}
	}

	if empresaNome.Valid {
		invite.Empresa = &models.Empresa{
			ID:   invite.EmpresaID,
			Nome: empresaNome.String,
		}
	}

	if convidadoPorNome.Valid {
		invite.ConvidadoPorUser = &models.Usuario{
			ID:   invite.ConvidadoPor,
			Nome: convidadoPorNome.String,
		}
	}

	return invite, nil
}

// MarkInviteAsUsed marca convite como usado
func (r *InviteRepository) MarkInviteAsUsed(inviteID int) error {
	query := "UPDATE convites_usuario SET usado = TRUE WHERE id = ?"
	_, err := r.db.Exec(query, inviteID)
	return err
}

// GetPendingInvitesByEmail verifica convites pendentes por email
func (r *InviteRepository) GetPendingInvitesByEmail(email string) ([]models.ConviteUsuario, error) {
	var invites []models.ConviteUsuario
	query := `
		SELECT id, email, role_id, empresa_id, convidado_por, token, 
			   usado, expira_em, dados_convite, created_at
		FROM convites_usuario 
		WHERE email = ? AND usado = FALSE AND expira_em > NOW()
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query, email)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var invite models.ConviteUsuario
		err := rows.Scan(
			&invite.ID, &invite.Email, &invite.RoleID, &invite.EmpresaID,
			&invite.ConvidadoPor, &invite.Token, &invite.Usado, &invite.ExpiraEm,
			&invite.DadosConvite, &invite.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		invites = append(invites, invite)
	}

	return invites, nil
}

// CleanExpiredInvites remove convites expirados
func (r *InviteRepository) CleanExpiredInvites() error {
	query := "DELETE FROM convites_usuario WHERE expira_em < NOW()"
	_, err := r.db.Exec(query)
	return err
}
