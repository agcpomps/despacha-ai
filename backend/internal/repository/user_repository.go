package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/agcpomps/despacha-ai/backend/internal/domain"
	"github.com/jmoiron/sqlx"
)

type UserRepository interface {
	Create(ctx context.Context, user *domain.User) (*domain.User, error)
	FindByPhone(ctx context.Context, phone string) (*domain.User, error)
	FindByID(ctx context.Context, id string) (*domain.User, error)
	List(ctx context.Context, search string, limit, offset int) ([]domain.User, error)
	Count(ctx context.Context, search string) (int, error)
	UpdateRole(ctx context.Context, id string, role domain.UserRole) error
}

type userRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) UserRepository {
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) List(ctx context.Context, search string, limit, offset int) ([]domain.User, error) {
	query := `
		SELECT
			id,
			name,
			phone,
			email,
			password_hash,
			avatar_url,
			role,
			status,
			is_verified,
			created_at,
			updated_at
		FROM users
	`

	args := []any{}
	if search != "" {
		query += ` WHERE name ILIKE $1 OR phone ILIKE $1`
		args = append(args, "%"+search+"%")
	}

	query += fmt.Sprintf(" ORDER BY created_at DESC LIMIT $%d OFFSET $%d", len(args)+1, len(args)+2)
	args = append(args, limit, offset)

	var users []domain.User
	if err := r.db.SelectContext(ctx, &users, query, args...); err != nil {
		return nil, err
	}

	return users, nil
}

func (r *userRepository) Count(ctx context.Context, search string) (int, error) {
	query := `SELECT COUNT(*) FROM users`

	args := []any{}
	if search != "" {
		query += ` WHERE name ILIKE $1 OR phone ILIKE $1`
		args = append(args, "%"+search+"%")
	}

	var total int
	if err := r.db.GetContext(ctx, &total, query, args...); err != nil {
		return 0, err
	}

	return total, nil
}

func (r *userRepository) UpdateRole(ctx context.Context, id string, role domain.UserRole) error {
	query := `
		UPDATE users
		SET role = $2, updated_at = NOW()
		WHERE id = $1
	`

	result, err := r.db.ExecContext(ctx, query, id, role)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *userRepository) Create(ctx context.Context, user *domain.User) (*domain.User, error) {
	query := `
		INSERT INTO users (
			name,
			phone,
			email,
			password_hash,
			role,
			status
		)
		VALUES (
			:name,
			:phone,
			:email,
			:password_hash,
			:role,
			:status
		)
		RETURNING
			id,
			name,
			phone,
			email,
			password_hash,
			avatar_url,
			role,
			status,
			is_verified,
			created_at,
			updated_at;
	`

	rows, err := r.db.NamedQueryContext(ctx, query, user)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if rows.Next() {
		var createdUser domain.User
		if err := rows.StructScan(&createdUser); err != nil {
			return nil, err
		}

		return &createdUser, nil
	}

	return nil, sql.ErrNoRows
}

func (r *userRepository) FindByPhone(ctx context.Context, phone string) (*domain.User, error) {
	query := `
		SELECT
			id,
			name,
			phone,
			email,
			password_hash,
			avatar_url,
			role,
			status,
			is_verified,
			created_at,
			updated_at
		FROM users
		WHERE phone = $1
		LIMIT 1;
    `
	var user domain.User

	if err := r.db.GetContext(ctx, &user, query, phone); err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) FindByID(ctx context.Context, id string) (*domain.User, error) {
	query := `
		SELECT
			id,
			name,
			phone,
			email,
			password_hash,
			avatar_url,
			role,
			status,
			is_verified,
			created_at,
			updated_at
		FROM users
		WHERE id = $1
		LIMIT 1;
	`

	var user domain.User

	if err := r.db.GetContext(ctx, &user, query, id); err != nil {
		return nil, err
	}

	return &user, nil
}
