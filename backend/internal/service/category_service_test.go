package service

import (
	"context"
	"testing"

	"github.com/agcpomps/despacha-ai/backend/internal/domain"
)

type fakeCategoryRepository struct {
	categories []domain.Category
}

func (r fakeCategoryRepository) FindAll(ctx context.Context) ([]domain.Category, error) {
	return r.categories, nil
}

func (r fakeCategoryRepository) FindById(ctx context.Context, id string) (*domain.Category, error) {
	return nil, nil
}

func (r fakeCategoryRepository) FindBySlug(ctx context.Context, slug string) (*domain.Category, error) {
	return nil, nil
}

func TestCategoryServiceGetCategoriesBuildsTree(t *testing.T) {
	parentID := "root-category"
	childID := "child-category"

	service := NewCategoryService(fakeCategoryRepository{
		categories: []domain.Category{
			{ID: parentID, Name: "Root", Slug: "root"},
			{ID: childID, Name: "Child", Slug: "child", ParentID: &parentID},
		},
	})

	categories, err := service.GetCategories(context.Background())
	if err != nil {
		t.Fatalf("GetCategories returned error: %v", err)
	}

	if len(categories) != 1 {
		t.Fatalf("expected 1 root category, got %d", len(categories))
	}

	if categories[0].ID != parentID {
		t.Fatalf("expected root category %q, got %q", parentID, categories[0].ID)
	}

	if len(categories[0].Children) != 1 {
		t.Fatalf("expected 1 child category, got %d", len(categories[0].Children))
	}

	if categories[0].Children[0].ID != childID {
		t.Fatalf("expected child category %q, got %q", childID, categories[0].Children[0].ID)
	}
}
