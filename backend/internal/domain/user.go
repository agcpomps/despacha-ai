package domain

import "time"

type UserRole string
type UserStatus string

const (
	RoleUser      UserRole = "user"
	RoleModerator UserRole = "moderator"
	RoleAdmin     UserRole = "admin"

	UserStatusActive    UserStatus = "active"
	UserStatusSuspended UserStatus = "suspended"
	UserStatusDeleted   UserStatus = "deleted"
)

type User struct {
	ID           string     `db:"id" json:"id"`
	Name         string     `db:"name" json:"name"`
	Phone        string     `db:"phone" json:"phone"`
	Email        *string    `db:"email" json:"email,omitempty"`
	PasswordHash string     `db:"password_hash" json:"-"`
	AvatarURL    *string    `db:"avatar_url" json:"avatar_url,omitempty"`
	Role         UserRole   `db:"role" json:"role"`
	Status       UserStatus `db:"status" json:"status"`
	IsVerified   bool       `db:"is_verified" json:"is_verified"`
	CreatedAt    time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time  `db:"updated_at" json:"updated_at"`
}
