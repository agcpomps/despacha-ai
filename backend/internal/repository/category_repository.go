package repository

import (
	"context"

	"github.com/agcpomps/despacha-ai/backend/internal/domain"
	"github.com/jmoiron/sqlx"
)

type CategoryRepository interface {
	FindAll(ctx context.Context) ([]domain.Category, error)
	FindByID(ctx context.Context, id string) (*domain.Category, error)
	FindBySlug(ctx context.Context, slug string) (*domain.Category, error)
}

type categoryRepository struct {
	db *sqlx.DB
}

func NewCategoryRepository(db *sqlx.DB) CategoryRepository {
	return &categoryRepository{
		db: db,
	}
}

func (r *categoryRepository) FindAll(ctx context.Context) ([]domain.Category, error) {
	query := `
		SELECT
			id,
			name,
			slug,
			parent_id,
			created_at
		FROM categories
		ORDER BY
			parent_id NULLS FIRST,
			name ASC;
	`
	var categories []domain.Category

	if err := r.db.SelectContext(ctx, &categories, query); err != nil {
		return nil, err
	}

	return categories, nil
}

func (r *categoryRepository) FindByID(ctx context.Context, id string) (*domain.Category, error) {
	query := `
		SELECT
			id,
			name,
			slug,
			parent_id,
			created_at
		FROM categories
		WHERE id = $1
		LIMIT 1;
	`

	var category domain.Category

	if err := r.db.GetContext(ctx, &category, query, id); err != nil {
		return nil, err
	}

	return &category, nil
}

func (r *categoryRepository) FindBySlug(ctx context.Context, slug string) (*domain.Category, error) {
	query := `
		SELECT
			id,
			name,
			slug,
			parent_id,
			created_at
		FROM categories
		WHERE slug = $1
		LIMIT 1;
	`
	var category domain.Category

	if err := r.db.GetContext(ctx, &category, query, slug); err != nil {
		return nil, err
	}

	return &category, nil
}
