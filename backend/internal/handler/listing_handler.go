package handler

import (
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/agcpomps/despacha-ai/backend/internal/dto"
	"github.com/agcpomps/despacha-ai/backend/internal/imageutil"
	"github.com/agcpomps/despacha-ai/backend/internal/service"
	"github.com/google/uuid"
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

func buildCreateListingRequest(c *echo.Context) (dto.CreateListingRequest, error) {
	contentType := c.Request().Header.Get(echo.HeaderContentType)
	if !strings.HasPrefix(contentType, "multipart/form-data") {
		var req dto.CreateListingRequest
		if err := c.Bind(&req); err != nil {
			return req, errors.New("invalid request body")
		}
		return req, nil
	}

	price, err := strconv.ParseFloat(c.FormValue("price"), 64)
	if err != nil {
		return dto.CreateListingRequest{}, errors.New("price must be a valid number")
	}

	images, err := saveUploadedListingImages(c)
	if err != nil {
		return dto.CreateListingRequest{}, err
	}

	return dto.CreateListingRequest{
		CategoryID:       optionalStringPointer(c.FormValue("category_id")),
		Title:            strings.TrimSpace(c.FormValue("title")),
		Description:      strings.TrimSpace(c.FormValue("description")),
		Price:            price,
		Province:         strings.TrimSpace(c.FormValue("province")),
		City:             optionalStringPointer(c.FormValue("city")),
		AddressReference: optionalStringPointer(c.FormValue("address_reference")),
		WhatsAppPhone:    optionalStringPointer(c.FormValue("whatsapp_phone")),
		Phone:            optionalStringPointer(c.FormValue("phone")),
		Condition:        strings.TrimSpace(c.FormValue("condition")),
		Images:           images,
	}, nil
}

func optionalStringPointer(value string) *string {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return nil
	}
	return &trimmed
}

func saveUploadedListingImages(c *echo.Context) ([]string, error) {
	form, err := c.MultipartForm()
	if err != nil {
		return nil, errors.New("invalid multipart form")
	}

	files := form.File["images"]
	if len(files) > 8 {
		return nil, service.ErrTooManyListingImages
	}

	images := make([]string, 0, len(files))
	for _, file := range files {
		if file == nil || file.Size == 0 {
			continue
		}

		imageURL, err := saveUploadedFile(c, file)
		if err != nil {
			return nil, err
		}
		images = append(images, imageURL)
	}

	return images, nil
}

func saveUploadedFile(c *echo.Context, file *multipart.FileHeader) (string, error) {
	if !strings.HasPrefix(file.Header.Get("Content-Type"), "image/") {
		return "", errors.New("only image uploads are allowed")
	}

	source, err := file.Open()
	if err != nil {
		return "", errors.New("failed to open uploaded image")
	}
	defer source.Close()

	dir := filepath.Join("uploads", "listings")

	// auto-orient, downscale and recompress to a web-friendly JPEG
	fileName, err := imageutil.SaveCompressed(source, dir, uuid.NewString())
	if err != nil {
		c.Logger().Error(
			"image processing failed",
			"filename", file.Filename,
			"content_type", file.Header.Get("Content-Type"),
			"size", file.Size,
			"error", err,
		)
		return "", errors.New("unsupported or corrupt image file")
	}

	scheme := "http"
	if c.Request().TLS != nil {
		scheme = "https"
	}
	// behind a reverse proxy (Caddy) TLS terminates upstream
	if forwarded := c.Request().Header.Get("X-Forwarded-Proto"); forwarded != "" {
		scheme = forwarded
	}

	host := c.Request().Host
	if forwardedHost := c.Request().Header.Get("X-Forwarded-Host"); forwardedHost != "" {
		host = forwardedHost
	}

	return fmt.Sprintf("%s://%s/uploads/listings/%s", scheme, host, fileName), nil
}

func (h *ListingHandler) CreateListing(c *echo.Context) error {
	req, err := buildCreateListingRequest(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
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

func (h *ListingHandler) GetMyListings(c *echo.Context) error {
	userIDValue := c.Get("user_id")
	userID, ok := userIDValue.(string)
	if !ok || userID == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "unauthorized",
		})
	}

	filters, err := buildListingFilters(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	status := c.QueryParam("status")
	if status != "" {
		filters.Status = &status
	}

	listings, err := h.listingService.GetMyListings(c.Request().Context(), userID, filters)
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

func buildListingFilters(c *echo.Context) (dto.ListingFilterRequest, error) {
	var filters dto.ListingFilterRequest

	filters.Page = 1
	filters.Limit = 12

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
			return filters, errors.New("min_price must be a valid number")
		}

		filters.MinPrice = &minPrice
	}

	maxPriceValue := c.QueryParam("max_price")
	if maxPriceValue != "" {
		maxPrice, err := strconv.ParseFloat(maxPriceValue, 64)
		if err != nil {
			return filters, errors.New("max-price must be a valid number")
		}

		filters.MaxPrice = &maxPrice
	}

	featuredValue := c.QueryParam("featured")
	if featuredValue != "" {
		featured, err := strconv.ParseBool(featuredValue)
		if err != nil {
			return filters, errors.New("featured must be true or false")
		}
		filters.Featured = &featured
	}

	sort := c.QueryParam("sort")
	switch sort {
	case "", "newest", "oldest", "price_asc", "price_desc":
		filters.Sort = sort
	default:
		return filters, errors.New("sort must be one of: newest, oldest, price_asc, price_desc")
	}

	if filters.MinPrice != nil && filters.MaxPrice != nil {
		if *filters.MinPrice > *filters.MaxPrice {
			return filters, errors.New("min_price cannot be greater than max_price")
		}
	}

	pageValue := c.QueryParam("page")
	if pageValue != "" {
		page, err := strconv.Atoi(pageValue)
		if err != nil {
			return filters, errors.New("page must be valid number")
		}

		if page < 1 {
			return filters, errors.New("page must be greater than zero")
		}

		filters.Page = page
	}

	limitValue := c.QueryParam("limit")
	if limitValue != "" {
		limit, err := strconv.Atoi(limitValue)
		if err != nil {
			return filters, errors.New("limit must be a valid number")
		}

		if limit < 1 {
			return filters, errors.New("limit must be greater than zero")
		}

		if limit > 50 {
			return filters, errors.New("limit cannot be greater than 50")
		}

		filters.Limit = limit
	}

	filters.Offset = (filters.Page - 1) * filters.Limit

	return filters, nil
}

func (h *ListingHandler) UpdateListing(c *echo.Context) error {
	id := c.Param("id")
	var req dto.UpdateListingRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(400, map[string]string{
			"error": "inavalid request",
		})
	}

	userID := c.Get("user_id").(string)

	listing, err := h.listingService.UpdateListing(c.Request().Context(), userID, id, req)
	if err != nil {
		if errors.Is(err, service.ErrInvalidStatus) {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "status must be 'active', 'sold' or 'paused'",
			})
		}
		if errors.Is(err, service.ErrInvalidCondition) {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "condition must be 'new' or 'used'",
			})
		}
		return c.JSON(400, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, listing)
}

type featureListingRequest struct {
	Days int `json:"days"`
}

func (h *ListingHandler) FeatureListing(c *echo.Context) error {
	id := c.Param("id")

	req := featureListingRequest{Days: 7}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid request body",
		})
	}

	if err := h.listingService.FeatureListing(c.Request().Context(), id, req.Days); err != nil {
		if errors.Is(err, service.ErrListingNotFoud) {
			return c.JSON(http.StatusNotFound, map[string]string{
				"error": "listing not found",
			})
		}

		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"status": "featured"})
}

func (h *ListingHandler) UnfeatureListing(c *echo.Context) error {
	id := c.Param("id")

	if err := h.listingService.UnfeatureListing(c.Request().Context(), id); err != nil {
		if errors.Is(err, service.ErrListingNotFoud) {
			return c.JSON(http.StatusNotFound, map[string]string{
				"error": "listing not found",
			})
		}

		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to unfeature listing"})
	}

	return c.JSON(http.StatusOK, map[string]string{"status": "unfeatured"})
}

func (h *ListingHandler) BumpListing(c *echo.Context) error {
	id := c.Param("id")

	if err := h.listingService.BumpListing(c.Request().Context(), id); err != nil {
		if errors.Is(err, service.ErrListingNotFoud) {
			return c.JSON(http.StatusNotFound, map[string]string{
				"error": "listing not found",
			})
		}

		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to bump listing"})
	}

	return c.JSON(http.StatusOK, map[string]string{"status": "bumped"})
}

func (h *ListingHandler) DeleteListing(c *echo.Context) error {
	id := c.Param("id")

	userID := c.Get("user_id").(string)

	err := h.listingService.DeleteListing(c.Request().Context(), userID, id)
	if err != nil {
		return c.JSON(400, map[string]string{"error": err.Error()})
	}

	return c.JSON(204, nil)
}
