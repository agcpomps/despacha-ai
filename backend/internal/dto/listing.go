package dto

type CreateListingRequest struct {
	CategoryID       *string  `json:"category_id"`
	Title            string   `json:"title"`
	Description      string   `json:"description"`
	Price            float64  `json:"price"`
	Province         string   `json:"province"`
	City             *string  `json:"city"`
	AddressReference *string  `json:"address_reference"`
	WhatsAppPhone    *string  `json:"whatsapp_phone"`
	Phone            *string  `json:"phone"`
	Condition        string   `json:"condition"`
	Images           []string `json:"images"`
}

type UpdateListingRequest struct {
	CategoryID       *string  `json:"category_id"`
	Title            *string  `json:"title"`
	Description      *string  `json:"description"`
	Price            *float64 `json:"price"`
	Province         *string  `json:"province"`
	City             *string  `json:"city"`
	AddressReference *string  `json:"address_reference"`
	WhatsAppPhone    *string  `json:"whatsapp_phone"`
	Phone            *string  `json:"phone"`
	Condition        *string  `json:"condition"`
	Status           *string  `json:"status"`
}

type ListingResponse struct {
	ID               string                 `json:"id"`
	UserID           string                 `json:"user_id"`
	CategoryID       *string                `json:"category_id,omitempty"`
	Title            string                 `json:"title"`
	Description      string                 `json:"description"`
	Price            float64                `json:"price"`
	Currency         string                 `json:"currency"`
	Province         string                 `json:"province"`
	City             *string                `json:"city,omitempty"`
	AddressReference *string                `json:"address_reference,omitempty"`
	WhatsAppPhone    *string                `json:"whatsapp_phone,omitempty"`
	Phone            *string                `json:"phone,omitempty"`
	Condition        string                 `json:"condition"`
	Status           string                 `json:"status"`
	ViewsCount       int                    `json:"views_count"`
	Images           []ListingImageResponse `json:"images"`
	CreatedAt        string                 `json:"created_at"`
	UpdatedAt        string                 `json:"updated_at"`
}

type ListingImageResponse struct {
	ID       string `json:"id"`
	ImageURL string `json:"image_url"`
	Position int    `json:"position"`
}
