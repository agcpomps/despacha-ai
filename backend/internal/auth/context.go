package auth

import "github.com/labstack/echo/v5"

const ContextUserKey = "user"

type UserContext struct {
	UserID string
	Role   string
	Phone  string
}

func GetUser(c *echo.Context) (*UserContext, bool) {
	user, ok := c.Get(ContextUserKey).(*UserContext)
	return user, ok
}
