package handler

import (
	"errors"
	"net/http"
	"strconv"

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

	listing, err := h.listingService.CreateListing(c.Request().Context(), userID, req)
	if err != nil {
		if errors.Is(err, service.ErrCategoryNotfound) {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "category not found",
			})
		}
		if errors.Is(err, service.ErrInvalidCondition) {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "condition must be 'new' or 'used'",
			})
		}

		if errors.Is(err, service.ErrTooManyListingImages) {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "maximum 8 images allowed",
			})
		}

		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}
	return c.JSON(http.StatusCreated, listing)
}

func (h *ListingHandler) GetListings(c *echo.Context) error {
	filters, err := buildListingFilters(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	listings, err := h.listingService.GetListings(c.Request().Context(), filters)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to fecth listings",
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

func buildListingFilters(c *echo.Context) (dto.ListingFilterRequest, error) {
	var filters dto.ListingFilterRequest

	categoryID := c.QueryParam("category_id")
	if categoryID != "" {
		filters.CategoryID = &categoryID
	}

	province := c.QueryParam("province")
	if province != "" {
		filters.Province = &province
	}

	city := c.QueryParam("city")
	if city != "" {
		filters.City = &city
	}

	search := c.QueryParam("search")
	if search != "" {
		filters.Search = &search
	}

	minPriceValue := c.QueryParam("min_price")
	if minPriceValue != "" {
		minPrice, err := strconv.ParseFloat(minPriceValue, 64)
		if err != nil {
			return filters, errors.New("max_price must be a valid number")
		}

		filters.MaxPrice = &minPrice
	}

	maxPriceValue := c.QueryParam("max_price")
	if maxPriceValue != "" {
		maxPrice, err := strconv.ParseFloat(maxPriceValue, 64)
		if err != nil {
			return filters, errors.New("max-price must be a valid number")
		}

		filters.MaxPrice = &maxPrice
	}

	if filters.MinPrice != nil && filters.MaxPrice != nil {
		if *filters.MinPrice > *filters.MaxPrice {
			return filters, errors.New("min_price cannot be greater than max_price")
		}
	}

	return filters, nil
}
