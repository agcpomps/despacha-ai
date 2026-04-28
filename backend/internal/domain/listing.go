package domain

import "time"

type Listing struct {
	ID         string  `db:"id" json:"id"`
	UserID     string  `db:"user_id" json:"user_id"`
	CategoryID *string `db:"category_id" json:"category_id,omitempty"`

	Title       string  `db:"title" json:"title"`
	Description string  `db:"description" json:"description"`
	Price       float64 `db:"price" json:"price"`
	Currency    string  `db:"currency" json:"currency"`

	Province         string  `db:"province" json:"province"`
	City             *string `db:"city" json:"city,omitempty"`
	AddressReference *string `db:"address_reference" json:"address_reference,omitempty"`

	WhatsAppPhone *string `db:"whatsapp_phone" json:"whatsapp_phone,omitempty"`
	Phone         *string `db:"phone" json:"phone,omitempty"`

	Condition  string `db:"condition" json:"condition"`
	Status     string `db:"status" json:"status"`
	ViewsCount int    `db:"views_count" json:"views_count"`

	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}
