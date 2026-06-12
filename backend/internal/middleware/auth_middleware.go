package middleware

import (
	"errors"
	"net/http"
	"strings"

	"github.com/agcpomps/despacha-ai/backend/internal/auth"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v5"
)

type JWTCustomClaims struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
	Phone  string `json:"phone"`
	jwt.RegisteredClaims
}

func AuthMiddleware(jwtSecret string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": "missing authorization header",
				})
			}

			parts := strings.Fields(authHeader)
			if len(parts) != 2 || parts[0] != "Bearer" {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": "invalid authorization header format",
				})
			}

			tokenString := parts[1]

			token, err := jwt.ParseWithClaims(
				tokenString,
				&JWTCustomClaims{},
				func(token *jwt.Token) (interface{}, error) {
					if token.Method != jwt.SigningMethodHS256 {
						return nil, errors.New("unexpected token signing method")
					}

					return []byte(jwtSecret), nil
				},
			)

			if err != nil || !token.Valid {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": "invalid or expired token",
				})
			}

			claims, ok := token.Claims.(*JWTCustomClaims)
			if !ok {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": "invalid token claims",
				})
			}

			c.Set(auth.ContextUserKey, &auth.UserContext{
				UserID: claims.UserID,
				Role:   claims.Role,
				Phone:  claims.Phone,
			})
			c.Set("user_id", claims.UserID)
			c.Set("user_role", claims.Role)
			c.Set("user_phone", claims.Phone)

			return next(c)
		}
	}
}

func RequireRole(allowedRoles ...string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			userRoleValue := c.Get("user_role")

			if userRoleValue == nil {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": "missing user role",
				})
			}

			userRole, ok := userRoleValue.(string)
			if !ok {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": "invalid user role",
				})
			}

			for _, role := range allowedRoles {
				if userRole == role {
					return next(c)
				}
			}

			return c.JSON(http.StatusForbidden, map[string]string{
				"error": "permission denied",
			})
		}
	}
}
