package handler

import (
	"errors"
	"net/http"

	"github.com/agcpomps/despacha-ai/backend/internal/dto"
	"github.com/agcpomps/despacha-ai/backend/internal/service"
	"github.com/labstack/echo/v5"
)

type AuthHandler struct {
	authService service.AuthService
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

func (h *AuthHandler) Register(c *echo.Context) error {
	var req dto.RegisterRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid request body",
		})
	}

	res, err := h.authService.Register(c.Request().Context(), req)
	if err != nil {
		if err != nil {
			if errors.Is(err, service.ErrPhoneAlreadyExists) {
				return c.JSON(http.StatusConflict, map[string]string{
					"error": "phone already exists",
				})
			}
		}

		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to register user",
		})
	}

	return c.JSON(http.StatusCreated, res)
}

func (h *AuthHandler) Login(c *echo.Context) error {
	var req dto.LoginRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid request body",
		})
	}

	res, err := h.authService.Login(c.Request().Context(), req)
	if err != nil {
		if errors.Is(err, service.ErrInvalidCredentials) {
			return c.JSON(http.StatusUnauthorized, map[string]string{
				"error": "invalid phone or password",
			})
		}
		if errors.Is(err, service.ErrUserSuspended) {
			return c.JSON(http.StatusForbidden, map[string]string{
				"error": "user account is suspended",
			})
		}

		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to login",
		})
	}

	return c.JSON(http.StatusOK, res)
}
