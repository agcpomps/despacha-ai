package repository

import (
	"context"
	"database/sql"

	"github.com/agcpomps/despacha-ai/backend/internal/domain"
	"github.com/jmoiron/sqlx"
)

type UserRepository interface {
	Create(ctx context.Context, user *domain.User) (*domain.User, error)
	FindByPhone(ctx context.Context, phone string) (*domain.User, error)
	FindByID(ctx context.Context, id string) (*domain.User, error)
}

type userRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) UserRepository {
	return &userRepository{
		db: db,
	}
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
