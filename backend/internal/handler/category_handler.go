package handler

import (
	"net/http"

	"github.com/agcpomps/despacha-ai/backend/internal/service"
	"github.com/labstack/echo/v5"
)

type CategoryHandler struct {
	categoryServise service.CategoryService
}

func NewCategoryHandler(categoryService service.CategoryService) *CategoryHandler {
	return &CategoryHandler{
		categoryServise: categoryService,
	}
}

func (h *CategoryHandler) GetCategories(c *echo.Context) error {
	categories, err := h.categoryServise.GetCategories(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to fecth categories",
		})
	}

	return c.JSON(http.StatusOK, categories)
}
