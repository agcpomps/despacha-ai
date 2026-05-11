package repository

import (
	"context"
	"database/sql"

	"github.com/agcpomps/despacha-ai/backend/internal/domain"
	"github.com/jmoiron/sqlx"
)

type ListingRepository interface {
	Create(ctx context.Context, listing *domain.Listing) (*domain.Listing, error)
	CreateImage(ctx context.Context, image *domain.ListingImage) (*domain.ListingImage, error)
	FindAll(ctx context.Context) ([]domain.Listing, error)
	FindByID(ctx context.Context, id string) (*domain.Listing, error)
	FindImagesByListingID(ctx context.Context, listingID string) ([]domain.ListingImage, error)
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

func (r *listingRepository) FindAll(ctx context.Context) ([]domain.Listing, error) {
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
			created_at,
			updated_at
		FROM listings
		WHERE status = 'active'
		ORDER BY created_at DESC;
	`
	var listings []domain.Listing

	if err := r.db.SelectContext(ctx, &listings, query); err != nil {
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
