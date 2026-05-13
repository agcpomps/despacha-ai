package service

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/agcpomps/despacha-ai/backend/internal/domain"
	"github.com/agcpomps/despacha-ai/backend/internal/dto"
	"github.com/agcpomps/despacha-ai/backend/internal/repository"
)

var (
	ErrListingNotFoud       = errors.New("listing not found")
	ErrInvalidListing       = errors.New("invalid listing data")
	ErrCategoryNotfound     = errors.New("category not found")
	ErrInvalidCondition     = errors.New("invalid condition")
	ErrTooManyListingImages = errors.New("too many listing images")
)

type ListingService interface {
	CreateListing(ctx context.Context, userID string, req dto.CreateListingRequest) (*dto.ListingResponse, error)
	GetListings(ctx context.Context, filters dto.ListingFilterRequest) ([]dto.ListingResponse, error)
	GetListingByID(ctx context.Context, id string) (*dto.ListingResponse, error)
}

type listingService struct {
	listingRepo  repository.ListingRepository
	categoryRepo repository.CategoryRepository
}

func NewListingService(listingRepo repository.ListingRepository, categoryRepo repository.CategoryRepository) ListingService {
	return &listingService{
		listingRepo:  listingRepo,
		categoryRepo: categoryRepo,
	}
}

func (s *listingService) CreateListing(ctx context.Context, userID string, req dto.CreateListingRequest) (*dto.ListingResponse, error) {
	if req.Title == "" {
		return nil, errors.New("title is required")
	}

	if len(req.Title) < 3 {
		return nil, errors.New("title must have at least 3 characters")
	}

	if req.Description == "" {
		return nil, errors.New("description is required")
	}

	if len(req.Description) < 10 {
		return nil, errors.New("description must have at least 10 characters")
	}

	if req.Price <= 0 {
		return nil, errors.New("price must be greater than zero")
	}

	if req.Province == "" {
		return nil, errors.New("province is required")
	}

	if len(req.Images) > 8 {
		return nil, ErrTooManyListingImages
	}

	if req.CategoryID != nil {
		_, err := s.categoryRepo.FindById(ctx, *req.CategoryID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, ErrCategoryNotfound
			}

			return nil, err
		}
	}

	condition := req.Condition
	if condition == "" {
		condition = "used"
	}

	if !isValidListingCondition(condition) {
		return nil, ErrInvalidCondition
	}

	listing := &domain.Listing{
		UserID:           userID,
		CategoryID:       req.CategoryID,
		Title:            req.Title,
		Description:      req.Description,
		Price:            req.Price,
		Currency:         "AOA",
		Province:         req.Province,
		City:             req.City,
		AddressReference: req.AddressReference,
		WhatsAppPhone:    req.WhatsAppPhone,
		Phone:            req.Phone,
		Condition:        condition,
		Status:           "active",
	}

	createdListing, err := s.listingRepo.Create(ctx, listing)
	if err != nil {
		return nil, err
	}

	images := []domain.ListingImage{}

	for index, imageURL := range req.Images {
		if imageURL == "" {
			continue
		}

		image := &domain.ListingImage{
			ListingID: createdListing.ID,
			ImageURL:  imageURL,
			Position:  index,
		}

		createdImage, err := s.listingRepo.CreateImage(ctx, image)
		if err != nil {
			return nil, err
		}

		images = append(images, *createdImage)
	}

	return toListingResponse(createdListing, images), nil
}

func (s *listingService) GetListings(ctx context.Context, filters dto.ListingFilterRequest) ([]dto.ListingResponse, error) {
	listings, err := s.listingRepo.FindAll(ctx, filters)
	if err != nil {
		return nil, err
	}

	responses := make([]dto.ListingResponse, 0, len(listings))

	for _, listing := range listings {
		images, err := s.listingRepo.FindImagesByListingID(ctx, listing.ID)
		if err != nil {
			return nil, err
		}

		responses = append(responses, *toListingResponse(&listing, images))
	}

	return responses, nil
}

func (s *listingService) GetListingByID(ctx context.Context, id string) (*dto.ListingResponse, error) {
	listing, err := s.listingRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrListingNotFoud
		}

		return nil, err
	}

	images, err := s.listingRepo.FindImagesByListingID(ctx, listing.ID)
	if err != nil {
		return nil, err
	}

	return toListingResponse(listing, images), nil
}

func toListingResponse(listing *domain.Listing, images []domain.ListingImage) *dto.ListingResponse {
	imageResponses := make([]dto.ListingImageResponse, 0, len(images))

	for _, image := range images {
		imageResponses = append(imageResponses, dto.ListingImageResponse{
			ID:       image.ID,
			ImageURL: image.ImageURL,
			Position: image.Position,
		})
	}

	return &dto.ListingResponse{
		ID:               listing.ID,
		UserID:           listing.UserID,
		CategoryID:       listing.CategoryID,
		Title:            listing.Title,
		Description:      listing.Description,
		Price:            listing.Price,
		Currency:         listing.Currency,
		Province:         listing.Province,
		City:             listing.City,
		AddressReference: listing.AddressReference,
		WhatsAppPhone:    listing.WhatsAppPhone,
		Phone:            listing.Phone,
		Condition:        listing.Condition,
		Status:           listing.Status,
		ViewsCount:       listing.ViewsCount,
		Images:           imageResponses,
		CreatedAt:        formatTime(listing.CreatedAt),
		UpdatedAt:        formatTime(listing.UpdatedAt),
	}
}

func formatTime(t time.Time) string {
	return t.Format(time.RFC3339)
}

func isValidListingCondition(condition string) bool {
	allowedConditions := map[string]bool{
		"new":  true,
		"used": true,
	}

	return allowedConditions[condition]
}
