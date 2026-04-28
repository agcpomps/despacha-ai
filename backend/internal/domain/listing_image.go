package domain

import "time"

type ListingImage struct {
	ID        string    `db:"id" json:"id"`
	ListingID string    `db:"listing_id" json:"listing_id"`
	ImageURL  string    `db:"image_url" json:"image_url"`
	Position  int       `db:"position" json:"position"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}
