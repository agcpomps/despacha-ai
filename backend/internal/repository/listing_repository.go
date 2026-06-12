package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/agcpomps/despacha-ai/backend/internal/domain"
	"github.com/agcpomps/despacha-ai/backend/internal/dto"
	"github.com/jmoiron/sqlx"
)

type ListingRepository interface {
	Create(ctx context.Context, listing *domain.Listing) (*domain.Listing, error)
	CreateImage(ctx context.Context, image *domain.ListingImage) (*domain.ListingImage, error)
	FindAll(ctx context.Context, filters dto.ListingFilterRequest) ([]domain.Listing, error)
	Count(ctx context.Context, filters dto.ListingFilterRequest) (int, error)
	FindByID(ctx context.Context, id string) (*domain.Listing, error)
	FindImagesByListingID(ctx context.Context, listingID string) ([]domain.ListingImage, error)
	Update(ctx context.Context, listing *domain.Listing) (*domain.Listing, error)
	Delete(ctx context.Context, id string, userID string) error
	SetFeatured(ctx context.Context, id string, featured bool, featuredUntil *time.Time) error
	Bump(ctx context.Context, id string) error
}

type listingRepository struct {
	db *sqlx.DB
}

func NewListingRepository(db *sqlx.DB) ListingRepository {
	return &listingRepository{
		db: db,
	}
}

func (r *listingRepository) Create(ctx context.Context, listing *domain.Listing) (*domain.Listing, error) {
	query := `
		INSERT INTO listings (
			user_id,
			category_id,
			title,
			description,
			price,
			currency,
			province,
			city,
			address_reference,
			whatsapp_phone,
			phone,
			condition,
			status
		)
		VALUES (
			:user_id,
			:category_id,
			:title,
			:description,
			:price,
			:currency,
			:province,
			:city,
			:address_reference,
			:whatsapp_phone,
			:phone,
			:condition,
			:status
		)
		RETURNING
			id,
			user_id,
			category_id,
			title,
			description,
			price,
			currency,
			province,
			city,
			address_reference,
			whatsapp_phone,
			phone,
			condition,
			status,
			views_count,
			is_featured,
			featured_until,
			bumped_at,
			created_at,
			updated_at;
	`

	rows, err := r.db.NamedQueryContext(ctx, query, listing)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if rows.Next() {
		var createdListing domain.Listing

		if err := rows.StructScan(&createdListing); err != nil {
			return nil, err
		}

		return &createdListing, nil
	}

	return nil, sql.ErrNoRows

}

func (r *listingRepository) CreateImage(ctx context.Context, image *domain.ListingImage) (*domain.ListingImage, error) {
	query := `
		INSERT INTO listing_images (
			listing_id,
			image_url,
			position
		)
		VALUES (
			:listing_id,
			:image_url,
			:position
		)
		RETURNING
			id,
			listing_id,
			image_url,
			position,
			created_at;
	`

	rows, err := r.db.NamedQueryContext(ctx, query, image)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if rows.Next() {
		var createdImage domain.ListingImage

		if err := rows.StructScan(&createdImage); err != nil {
			return nil, err
		}

		return &createdImage, nil
	}

	return nil, sql.ErrNoRows

}

func (r *listingRepository) FindAll(ctx context.Context, filters dto.ListingFilterRequest) ([]domain.Listing, error) {

	whereClause, args, argPosition := buildListingWhereClause(filters, 1)

	query := `
		SELECT
			id,
			user_id,
			category_id,
			title,
			description,
			price,
			currency,
			province,
			city,
			address_reference,
			whatsapp_phone,
			phone,
			condition,
			status,
			views_count,
			is_featured,
			featured_until,
			bumped_at,
			created_at,
			updated_at
		FROM listings
	`
	query += whereClause
	// active highlights always rank first, regardless of the chosen sort
	orderBy := "(is_featured AND (featured_until IS NULL OR featured_until > NOW())) DESC, " + listingOrderBy(filters.Sort)
	query += fmt.Sprintf(" ORDER BY %s LIMIT $%d OFFSET $%d", orderBy, argPosition, argPosition+1)

	args = append(args, filters.Limit, filters.Offset)

	var listings []domain.Listing

	if err := r.db.SelectContext(ctx, &listings, query, args...); err != nil {
		return nil, err
	}

	return listings, nil
}

func (r *listingRepository) FindByID(ctx context.Context, id string) (*domain.Listing, error) {
	query := `
		SELECT
			id,
			user_id,
			category_id,
			title,
			description,
			price,
			currency,
			province,
			city,
			address_reference,
			whatsapp_phone,
			phone,
			condition,
			status,
			views_count,
			is_featured,
			featured_until,
			bumped_at,
			created_at,
			updated_at
		FROM listings
		WHERE id = $1
		LIMIT 1;
	`

	var listing domain.Listing

	if err := r.db.GetContext(ctx, &listing, query, id); err != nil {
		return nil, err
	}

	return &listing, nil
}

func (r *listingRepository) FindImagesByListingID(ctx context.Context, listingID string) ([]domain.ListingImage, error) {
	query := `
		SELECT
			id,
			listing_id,
			image_url,
			position,
			created_at
		FROM listing_images
		WHERE listing_id = $1
		ORDER BY position ASC;
	`
	var images []domain.ListingImage

	if err := r.db.SelectContext(ctx, &images, query, listingID); err != nil {
		return nil, err
	}

	return images, nil
}

func (r *listingRepository) Count(ctx context.Context, filters dto.ListingFilterRequest) (int, error) {
	whereClause, args, _ := buildListingWhereClause(filters, 1)

	query := `
	  SELECT COUNT(*)
	  FROM listings
	`
	query += whereClause

	var total int

	if err := r.db.GetContext(ctx, &total, query, args...); err != nil {
		return 0, err
	}

	return total, nil
}

func buildListingWhereClause(filters dto.ListingFilterRequest, startPosition int) (string, []any, int) {
	args := []any{}
	argPosition := startPosition

	var query string
	switch {
	case filters.Status != nil && *filters.Status != "":
		query = fmt.Sprintf(" WHERE status = $%d", argPosition)
		args = append(args, *filters.Status)
		argPosition++
	case filters.UserID != nil:
		// owner views include everything except deleted listings
		query = " WHERE status <> 'deleted'"
	default:
		query = " WHERE status = 'active'"
	}

	if filters.UserID != nil && *filters.UserID != "" {
		query += fmt.Sprintf(" AND user_id = $%d", argPosition)
		args = append(args, *filters.UserID)
		argPosition++
	}

	if filters.CategoryID != nil && *filters.CategoryID != "" {
		query += fmt.Sprintf(" AND category_id = $%d", argPosition)
		args = append(args, *filters.CategoryID)
		argPosition++
	}

	if filters.Province != nil && *filters.Province != "" {
		query += fmt.Sprintf(" AND province ILIKE $%d", argPosition)
		args = append(args, *filters.Province)
		argPosition++
	}

	if filters.City != nil && *filters.City != "" {
		query += fmt.Sprintf(" AND city ILIKE $%d", argPosition)
		args = append(args, *filters.City)
		argPosition++
	}

	if filters.MinPrice != nil {
		query += fmt.Sprintf(" AND price >= $%d", argPosition)
		args = append(args, *filters.MinPrice)
		argPosition++
	}

	if filters.MaxPrice != nil {
		query += fmt.Sprintf(" AND price <= $%d", argPosition)
		args = append(args, *filters.MaxPrice)
		argPosition++
	}

	if filters.Featured != nil && *filters.Featured {
		query += " AND is_featured = TRUE AND (featured_until IS NULL OR featured_until > NOW())"
	}

	if filters.Search != nil && strings.TrimSpace(*filters.Search) != "" {
		searchTerm := "%" + strings.TrimSpace(*filters.Search) + "%"

		query += fmt.Sprintf(
			" AND (title ILIKE $%d OR description ILIKE $%d)",
			argPosition,
			argPosition,
		)

		args = append(args, searchTerm)
		argPosition++
	}

	return query, args, argPosition
}

func listingOrderBy(sort string) string {
	switch sort {
	case "price_asc":
		return "price ASC, created_at DESC"
	case "price_desc":
		return "price DESC, created_at DESC"
	case "oldest":
		return "created_at ASC"
	default: // newest — bumps count as fresh publications
		return "COALESCE(bumped_at, created_at) DESC"
	}
}

func (r *listingRepository) Update(ctx context.Context, listing *domain.Listing) (*domain.Listing, error) {
	query := `
		UPDATE listings
		SET
			title = :title,
			description = :description,
			price = :price,
			province = :province,
			city = :city,
			address_reference = :address_reference,
			whatsapp_phone = :whatsapp_phone,
			phone = :phone,
			condition = :condition,
			status = :status,
			updated_at = NOW()
		WHERE id = :id AND user_id = :user_id
		RETURNING *
	`
	rows, err := r.db.NamedQueryContext(ctx, query, listing)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	if rows.Next() {
		var updated domain.Listing
		if err := rows.StructScan(&updated); err != nil {
			return nil, err
		}

		return &updated, nil
	}

	return nil, sql.ErrNoRows
}

func (r *listingRepository) SetFeatured(ctx context.Context, id string, featured bool, featuredUntil *time.Time) error {
	query := `
		UPDATE listings
		SET is_featured = $2, featured_until = $3, updated_at = NOW()
		WHERE id = $1 AND status <> 'deleted'
	`

	result, err := r.db.ExecContext(ctx, query, id, featured, featuredUntil)
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

func (r *listingRepository) Bump(ctx context.Context, id string) error {
	query := `
		UPDATE listings
		SET bumped_at = NOW(), updated_at = NOW()
		WHERE id = $1 AND status = 'active'
	`

	result, err := r.db.ExecContext(ctx, query, id)
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

func (r *listingRepository) Delete(ctx context.Context, id string, userID string) error {
	query := `
		UPDATE listings
		SET status = 'deleted'
		WHERE id = $1 AND user_id = $2
	`

	_, err := r.db.ExecContext(ctx, query, id, userID)
	return err

}
