package service

import (
	"testing"
	"time"

	"github.com/agcpomps/despacha-ai/backend/internal/domain"
)

func TestToListingResponseAllowsSellerWithoutAvatar(t *testing.T) {
	now := time.Date(2026, time.May, 21, 13, 0, 0, 0, time.UTC)

	listing := &domain.Listing{
		ID:          "listing-1",
		UserID:      "user-1",
		Title:       "Phone",
		Description: "Working phone",
		Price:       100,
		Currency:    "AOA",
		Province:    "Luanda",
		Condition:   "used",
		Status:      "active",
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	seller := &domain.User{
		ID:    "user-1",
		Name:  "Ana",
		Phone: "+244900000000",
	}

	response := toListingResponse(listing, nil, seller, nil)

	if response == nil {
		t.Fatal("expected listing response, got nil")
	}
	if response.Seller == nil {
		t.Fatal("expected seller response, got nil")
	}
	if response.Seller.AvatarURL != "" {
		t.Fatalf("expected empty avatar URL, got %q", response.Seller.AvatarURL)
	}
}
