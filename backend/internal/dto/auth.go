package dto

type RegisterRequest struct {
	Name     string  `json:"name"`
	Phone    string  `json:"phone"`
	Email    *string `json:"email"`
	Password string  `json:"password"`
}

type LoginRequest struct {
	Phone    string `json:"phone"`
	Password string `json:"password"`
}

type AuthResponse struct {
	AccessToken string       `json:"access_token"`
	User        UserResponse `json:"user"`
}

type UserResponse struct {
	ID         string  `json:"id"`
	Name       string  `json:"name"`
	Phone      string  `json:"phone"`
	Email      *string `json:"email,omitempty"`
	AvatarURL  *string `json:"avatar_url,omitempty"`
	Role       string  `json:"role"`
	Status     string  `json:"status"`
	IsVerified bool    `json:"is_verified"`
	CreatedAt  string  `json:"created_at,omitempty"`
}

type PaginatedUsersResponse struct {
	Data       []UserResponse `json:"data"`
	Page       int            `json:"page"`
	Limit      int            `json:"limit"`
	Total      int            `json:"total"`
	TotalPages int            `json:"total_pages"`
}

type UpdateUserRoleRequest struct {
	Role string `json:"role"`
}
