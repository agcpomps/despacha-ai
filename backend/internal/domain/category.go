package domain

import "time"

type Category struct {
	ID        string    `db:"id" json:"id"`
	Name      string    `db:"name" json:"name"`
	Slug      string    `db:"slug" json:"slug"`
	ParentID  *string   `db:"parent_id" json:"parent_id,omitempty"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}
