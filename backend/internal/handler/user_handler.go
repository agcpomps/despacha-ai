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
