package postgres

import (
	"context"
	"fmt"

	"github.com/Alexander272/route-table/internal/models"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type RoleRepo struct {
	db *sqlx.DB
}

func NewRoleRepo(db *sqlx.DB) *RoleRepo {
	return &RoleRepo{db: db}
}

func (r *RoleRepo) Get(ctx context.Context, roleId uuid.UUID) (role models.Role, err error) {
	query := fmt.Sprintf("SELECT id, title, role, operations FROM %s WHERE id=$1", RolesTable)

	if err := r.db.Get(&role, query, roleId); err != nil {
		return role, fmt.Errorf("failed to execute query. error: %w", err)
	}
	return role, nil
}

func (r *RoleRepo) GetAll(ctx context.Context) (roles []models.Role, err error) {
	query := fmt.Sprintf("SELECT id, title, role, operations FROM %s", RolesTable)

	if err := r.db.Select(&roles, query); err != nil {
		return nil, fmt.Errorf("failed to execute query. error: %w", err)
	}
	return roles, nil
}

func (r *RoleRepo) Create(ctx context.Context, role models.RoleDTO) (id uuid.UUID, err error) {
	query := fmt.Sprintf("INSERT INTO %s (id, title, role, operations) VALUES ($1, $2, $3, $4)", RolesTable)

	id = uuid.New()

	_, err = r.db.Exec(query, id, role.Title, role.Role, pq.Array(role.Operations))
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to execute query. error: %w", err)
	}

	return id, nil
}

func (r *RoleRepo) Update(ctx context.Context, role models.RoleDTO) error {
	query := fmt.Sprintf("UPDATE %s SET title=$1, role=$2, operations=$3 WHERE id=$4", RolesTable)

	_, err := r.db.Exec(query, role.Title, role.Role, role.Operations, role.Id)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}

	return nil
}

func (r *RoleRepo) Delete(ctx context.Context, role models.RoleDTO) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id=$1", RolesTable)

	if _, err := r.db.Exec(query, role.Id); err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}
