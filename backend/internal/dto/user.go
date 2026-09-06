package dto

// PublicUserResponse é o perfil público de um vendedor (página da loja).
// Sem telefone nem email — esses só aparecem nos anúncios, por opção do vendedor.
type PublicUserResponse struct {
	ID         string  `json:"id"`
	Name       string  `json:"name"`
	AvatarURL  *string `json:"avatar_url,omitempty"`
	IsVerified bool    `json:"is_verified"`
	CreatedAt  string  `json:"created_at"`
}
