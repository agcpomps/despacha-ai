package handler

import (
	"errors"
	"net/http"

	"github.com/agcpomps/despacha-ai/backend/internal/dto"
	"github.com/agcpomps/despacha-ai/backend/internal/service"
	"github.com/labstack/echo/v5"
)

type ListingHandler struct {
	listingService service.ListingService
}

func NewListingHandler(listingService service.ListingService) *ListingHandler {
	return &ListingHandler{
		listingService: listingService,
	}
}

func (h *ListingHandler) CreateListing(c *echo.Context) error {
	var req dto.CreateListingRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid request body",
		})
	}

	userIDValue := c.Get("user_id")
	if userIDValue == nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "unauthorized",
		})
	}
	userID, ok := userIDValue.(string)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "invalid user context",
		})
	}

	listing, err := h.listingService.CreatedListing(c.Request().Context(), userID, req)
	if err != nil {
		if errors.Is(err, service.ErrInvalidListing) {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "invalid listing data",
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to create listing",
		})
	}

	return c.JSON(http.StatusCreated, listing)
}

func (h *ListingHandler) GetListings(c *echo.Context) error {
	listings, err := h.listingService.GetListings(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to fetch listings",
		})
	}

	return c.JSON(http.StatusOK, listings)
}

func (h *ListingHandler) GetListingByID(c *echo.Context) error {
	id := c.Param("id")
	listing, err := h.listingService.GetListingByID(c.Request().Context(), id)
	if err != nil {
		if errors.Is(err, service.ErrListingNotFoud) {
			return c.JSON(http.StatusNotFound, map[string]string{
				"error": "listing not found",
			})
		}

		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to fetch listing",
		})
	}

	return c.JSON(http.StatusOK, listing)
}
