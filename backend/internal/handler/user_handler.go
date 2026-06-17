package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/agcpomps/despacha-ai/backend/internal/dto"
	"github.com/agcpomps/despacha-ai/backend/internal/service"
	"github.com/labstack/echo/v5"
)

type UserHandler struct {
	userService service.UserService
}

func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (h *UserHandler) GetMe(c *echo.Context) error {
	userIDValue := c.Get("user_id")
	userID, ok := userIDValue.(string)
	if !ok || userID == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "unauthorized",
		})
	}

	user, err := h.userService.GetProfile(c.Request().Context(), userID)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			return c.JSON(http.StatusNotFound, map[string]string{
				"error": "user not found",
			})
		}

		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to fetch profile",
		})
	}

	return c.JSON(http.StatusOK, user)
}

func (h *UserHandler) GetUsers(c *echo.Context) error {
	search := c.QueryParam("search")

	page := 1
	if value := c.QueryParam("page"); value != "" {
		parsed, err := strconv.Atoi(value)
		if err != nil || parsed < 1 {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "page must be a positive number",
			})
		}
		page = parsed
	}

	limit := 20
	if value := c.QueryParam("limit"); value != "" {
		parsed, err := strconv.Atoi(value)
		if err != nil || parsed < 1 || parsed > 50 {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "limit must be between 1 and 50",
			})
		}
		limit = parsed
	}

	users, err := h.userService.ListUsers(c.Request().Context(), search, page, limit)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to fetch users",
		})
	}

	return c.JSON(http.StatusOK, users)
}

type changePasswordRequest struct {
	CurrentPassword string `json:"current_password"`
	NewPassword     string `json:"new_password"`
}

func (h *UserHandler) ChangePassword(c *echo.Context) error {
	userID, ok := c.Get("user_id").(string)
	if !ok || userID == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
	}

	var req changePasswordRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}

	err := h.userService.ChangePassword(c.Request().Context(), userID, req.CurrentPassword, req.NewPassword)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrIncorrectPassword):
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "palavra-passe actual incorrecta"})
		case errors.Is(err, service.ErrWeakPassword):
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "a nova palavra-passe deve ter pelo menos 6 caracteres"})
		case errors.Is(err, service.ErrUserNotFound):
			return c.JSON(http.StatusNotFound, map[string]string{"error": "user not found"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to change password"})
	}

	return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
}

func (h *UserHandler) ResetUserPassword(c *echo.Context) error {
	targetID := c.Param("id")

	tempPassword, err := h.userService.ResetUserPassword(c.Request().Context(), targetID)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			return c.JSON(http.StatusNotFound, map[string]string{
				"error": "user not found",
			})
		}

		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to reset password",
		})
	}

	return c.JSON(http.StatusOK, map[string]string{"password": tempPassword})
}

func (h *UserHandler) UpdateUserRole(c *echo.Context) error {
	actorID, ok := c.Get("user_id").(string)
	if !ok || actorID == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "unauthorized",
		})
	}

	targetID := c.Param("id")

	var req dto.UpdateUserRoleRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid request body",
		})
	}

	user, err := h.userService.UpdateUserRole(c.Request().Context(), actorID, targetID, req.Role)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrUserNotFound):
			return c.JSON(http.StatusNotFound, map[string]string{
				"error": "user not found",
			})
		case errors.Is(err, service.ErrInvalidRole):
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "role must be 'user', 'moderator' or 'admin'",
			})
		case errors.Is(err, service.ErrCannotChangeOwn):
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "cannot change your own role",
			})
		}

		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to update role",
		})
	}

	return c.JSON(http.StatusOK, user)
}
