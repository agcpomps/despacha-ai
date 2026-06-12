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
	ID         string  `json:"id"`
	UserID     string  `json:"user_id"`
	CategoryID *string `json:"category_id,omitempty"`

	Title       string  `json:"title"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Currency    string  `json:"currency"`

	Province         string  `json:"province"`
	City             *string `json:"city,omitempty"`
	AddressReference *string `json:"address_reference,omitempty"`

	WhatsAppPhone *string `json:"whatsapp_phone,omitempty"`
	Phone         *string `json:"phone,omitempty"`

	Condition  string `json:"condition"`
	Status     string `json:"status"`
	ViewsCount int    `json:"views_count"`

	IsFeatured    bool    `json:"is_featured"`
	FeaturedUntil *string `json:"featured_until,omitempty"`

	Seller   *ListingSellerResponse   `json:"seller,omitempty"`
	Category *ListingCategoryResponse `json:"category,omitempty"`
	Images   []ListingImageResponse   `json:"images"`

	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type ListingImageResponse struct {
	ID       string `json:"id"`
	ImageURL string `json:"image_url"`
	Position int    `json:"position"`
}

type ListingFilterRequest struct {
	CategoryID *string  `json:"category_id"`
	Province   *string  `json:"province"`
	City       *string  `json:"city"`
	MinPrice   *float64 `json:"min_price"`
	MaxPrice   *float64 `json:"max_price"`
	Search     *string  `json:"search"`
	Sort       string   `json:"sort"`

	Featured *bool `json:"featured"`

	// UserID and Status are set internally (e.g. /me/listings), not from query params.
	UserID *string `json:"-"`
	Status *string `json:"-"`

	Page   int `json:"page"`
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}

type PaginatedListingresponse struct {
	Data       []ListingResponse `json:"data"`
	Page       int               `json:"page"`
	Limit      int               `json:"limit"`
	Total      int               `json:"total"`
	TotalPages int               `json:"total_pages"`
}

type ListingSellerResponse struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Phone     string `json:"phone"`
	AvatarURL string `json:"avatar_url,omitempty"`
}

type ListingCategoryResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}
