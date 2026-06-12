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
	ErrInvalidStatus        = errors.New("invalid status")
	ErrTooManyListingImages = errors.New("too many listing images")
)

type ListingService interface {
	CreateListing(ctx context.Context, userID string, req dto.CreateListingRequest) (*dto.ListingResponse, error)
	GetListings(ctx context.Context, filters dto.ListingFilterRequest) (*dto.PaginatedListingresponse, error)
	GetMyListings(ctx context.Context, userID string, filters dto.ListingFilterRequest) (*dto.PaginatedListingresponse, error)
	GetListingByID(ctx context.Context, id string) (*dto.ListingResponse, error)
	UpdateListing(ctx context.Context, userID string, id string, req dto.UpdateListingRequest) (*dto.ListingResponse, error)
	DeleteListing(ctx context.Context, userID string, id string) error

	// admin-only promotion actions
	FeatureListing(ctx context.Context, id string, days int) error
	UnfeatureListing(ctx context.Context, id string) error
	BumpListing(ctx context.Context, id string) error
}

type listingService struct {
	listingRepo  repository.ListingRepository
	categoryRepo repository.CategoryRepository
	userRepo     repository.UserRepository
}

func NewListingService(listingRepo repository.ListingRepository, categoryRepo repository.CategoryRepository, userRepo repository.UserRepository) ListingService {
	return &listingService{
		listingRepo:  listingRepo,
		categoryRepo: categoryRepo,
		userRepo:     userRepo,
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
		_, err := s.categoryRepo.FindByID(ctx, *req.CategoryID)
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

	seller, err := s.userRepo.FindByID(ctx, createdListing.UserID)
	if err != nil {
		return nil, err
	}

	if seller == nil {
		return nil, errors.New("seller not found")
	}

	var category *domain.Category
	if createdListing.CategoryID != nil {
		category, err = s.categoryRepo.FindByID(ctx, *createdListing.CategoryID)
		if err != nil {
			category = nil
		}
	}

	return toListingResponse(createdListing, images, seller, category), nil
}

func (s *listingService) GetListings(ctx context.Context, filters dto.ListingFilterRequest) (*dto.PaginatedListingresponse, error) {
	listings, err := s.listingRepo.FindAll(ctx, filters)
	if err != nil {
		return nil, err
	}

	total, err := s.listingRepo.Count(ctx, filters)
	if err != nil {
		return nil, err
	}

	responses := make([]dto.ListingResponse, 0, len(listings))

	for _, listing := range listings {
		response, err := s.buildListingResponse(ctx, &listing)
		if err != nil {
			return nil, err
		}

		responses = append(responses, *response)
	}

	totalPages := 0
	if filters.Limit > 0 {
		totalPages = (total + filters.Limit - 1) / filters.Limit
	}

	return &dto.PaginatedListingresponse{
		Data:       responses,
		Page:       filters.Page,
		Limit:      filters.Limit,
		Total:      total,
		TotalPages: totalPages,
	}, nil
}

func (s *listingService) GetMyListings(ctx context.Context, userID string, filters dto.ListingFilterRequest) (*dto.PaginatedListingresponse, error) {
	filters.UserID = &userID
	return s.GetListings(ctx, filters)
}

func (s *listingService) GetListingByID(ctx context.Context, id string) (*dto.ListingResponse, error) {
	listing, err := s.listingRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrListingNotFoud
		}

		return nil, err
	}

	return s.buildListingResponse(ctx, listing)
}

func (s *listingService) buildListingResponse(ctx context.Context, listing *domain.Listing) (*dto.ListingResponse, error) {
	images, err := s.listingRepo.FindImagesByListingID(ctx, listing.ID)
	if err != nil {
		return nil, err
	}

	seller, err := s.userRepo.FindByID(ctx, listing.UserID)
	if err != nil {
		return nil, err
	}
	if seller == nil {
		return nil, errors.New("seller not found")
	}

	var category *domain.Category
	if listing.CategoryID != nil {
		cat, err := s.categoryRepo.FindByID(ctx, *listing.CategoryID)
		if err == nil {
			category = cat
		}
	}

	return toListingResponse(listing, images, seller, category), nil
}

func toListingResponse(listing *domain.Listing, images []domain.ListingImage, seller *domain.User, category *domain.Category) *dto.ListingResponse {
	imageResponses := make([]dto.ListingImageResponse, 0, len(images))

	for _, image := range images {
		imageResponses = append(imageResponses, dto.ListingImageResponse{
			ID:       image.ID,
			ImageURL: image.ImageURL,
			Position: image.Position,
		})
	}

	var sellerResponse *dto.ListingSellerResponse
	if seller != nil {
		sellerResponse = &dto.ListingSellerResponse{
			ID:    seller.ID,
			Name:  seller.Name,
			Phone: seller.Phone,
		}

		if seller.AvatarURL != nil {
			sellerResponse.AvatarURL = *seller.AvatarURL
		}
	}
	if seller == nil {
		return nil
	}

	var categoryResponse *dto.ListingCategoryResponse
	if category != nil {
		categoryResponse = &dto.ListingCategoryResponse{
			ID:   category.ID,
			Name: category.Name,
			Slug: category.Slug,
		}
	}

	isFeatured := listing.IsFeatured &&
		(listing.FeaturedUntil == nil || listing.FeaturedUntil.After(time.Now()))

	var featuredUntil *string
	if isFeatured && listing.FeaturedUntil != nil {
		formatted := formatTime(*listing.FeaturedUntil)
		featuredUntil = &formatted
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
		IsFeatured:       isFeatured,
		FeaturedUntil:    featuredUntil,
		Seller:           sellerResponse,
		Category:         categoryResponse,
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

func isValidListingStatus(status string) bool {
	allowedStatuses := map[string]bool{
		"active": true,
		"sold":   true,
		"paused": true,
	}

	return allowedStatuses[status]
}

func (s *listingService) UpdateListing(ctx context.Context, userID string, id string, req dto.UpdateListingRequest) (*dto.ListingResponse, error) {
	listing, err := s.listingRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if listing.UserID != userID {
		return nil, errors.New("not allowed")
	}

	if req.Title != nil {
		listing.Title = *req.Title
	}

	if req.Description != nil {
		listing.Description = *req.Description
	}
	if req.Price != nil {
		listing.Price = *req.Price
	}
	if req.Province != nil {
		listing.Province = *req.Province
	}
	if req.City != nil {
		listing.City = req.City
	}
	if req.AddressReference != nil {
		listing.AddressReference = req.AddressReference
	}
	if req.WhatsAppPhone != nil {
		listing.WhatsAppPhone = req.WhatsAppPhone
	}
	if req.Phone != nil {
		listing.Phone = req.Phone
	}
	if req.Condition != nil {
		if !isValidListingCondition(*req.Condition) {
			return nil, ErrInvalidCondition
		}
		listing.Condition = *req.Condition
	}
	if req.Status != nil {
		if !isValidListingStatus(*req.Status) {
			return nil, ErrInvalidStatus
		}
		listing.Status = *req.Status
	}

	updated, err := s.listingRepo.Update(ctx, listing)
	if err != nil {
		return nil, err
	}

	return s.buildListingResponse(ctx, updated)
}

func (s *listingService) FeatureListing(ctx context.Context, id string, days int) error {
	if days < 1 || days > 90 {
		return errors.New("days must be between 1 and 90")
	}

	until := time.Now().AddDate(0, 0, days)

	err := s.listingRepo.SetFeatured(ctx, id, true, &until)
	if errors.Is(err, sql.ErrNoRows) {
		return ErrListingNotFoud
	}

	return err
}

func (s *listingService) UnfeatureListing(ctx context.Context, id string) error {
	err := s.listingRepo.SetFeatured(ctx, id, false, nil)
	if errors.Is(err, sql.ErrNoRows) {
		return ErrListingNotFoud
	}

	return err
}

func (s *listingService) BumpListing(ctx context.Context, id string) error {
	err := s.listingRepo.Bump(ctx, id)
	if errors.Is(err, sql.ErrNoRows) {
		return ErrListingNotFoud
	}

	return err
}

func (s *listingService) DeleteListing(ctx context.Context, userID string, id string) error {
	listing, err := s.listingRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	if listing.UserID != userID {
		return errors.New("not allowed")

	}

	return s.listingRepo.Delete(ctx, id, userID)
}
