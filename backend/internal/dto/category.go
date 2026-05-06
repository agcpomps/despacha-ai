package dto

type CategoryResponse struct {
	ID       string             `json:"id"`
	Name     string             `json:"name"`
	Slug     string             `json:"slug"`
	ParentID *string            `json:"parent_id,omitempty"`
	Children []CategoryResponse `json:"children,omitempty"`
}
