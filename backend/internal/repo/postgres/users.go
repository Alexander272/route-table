package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/Alexander272/route-table/internal/models"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type UserRepo struct {
	db *sqlx.DB
}

func NewUserRepo(db *sqlx.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) Get(ctx context.Context, u models.SignIn) (user models.UserWithRole, err error) {
	query := fmt.Sprintf("SELECT id, password, role_id FROM %s WHERE login=$1", UsersTable)

	if err = r.db.Get(&user, query, u.Login); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return user, err
		}
		return user, fmt.Errorf("failed to execute query. error: %w", err)
	}

	return user, nil
}

func (r *UserRepo) GetAll(ctx context.Context) (users []models.User, err error) {
	query := fmt.Sprintf("SELECT %s.id, login, role FROM %s INNER JOIN %s ON role_id=%s.id",
		UsersTable, UsersTable, RolesTable, RolesTable)

	if err := r.db.Select(&users, query); err != nil {
		return nil, fmt.Errorf("failed to execute query. error: %w", err)
	}
	return users, nil
}

func (r *UserRepo) Create(ctx context.Context, user models.UserDTO) (id uuid.UUID, err error) {
	query := fmt.Sprintf("INSERT INTO %s (id, login, password, role_id) VALUES ($1, $2, $3, $4)", UsersTable)

	id = uuid.New()

	_, err = r.db.Exec(query, id, user.Login, user.Password, user.RoleId)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to execute query. error: %w", err)
	}

	return id, nil
}

func (r *UserRepo) Update(ctx context.Context, user models.UserDTO) error {
	query := fmt.Sprintf("UPDATE %s SET login=$1, password=$2, role_id=$3 WHERE id=$4", UsersTable)

	_, err := r.db.Exec(query, user.Login, user.Password, user.RoleId, user.Id)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}

	return nil
}

func (r *UserRepo) Delete(ctx context.Context, user models.UserDTO) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id=$1", UsersTable)

	if _, err := r.db.Exec(query, user.Id); err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}
