package service

import (
	"context"

	"github.com/agcpomps/despacha-ai/backend/internal/dto"
	"github.com/agcpomps/despacha-ai/backend/internal/repository"
)

type CategoryService interface {
	GetCategories(ctx context.Context) ([]dto.CategoryResponse, error)
}

type categoryService struct {
	categoryRepo repository.CategoryRepository
}

func NewCategoryService(categoryRepo repository.CategoryRepository) CategoryService {
	return &categoryService{
		categoryRepo: categoryRepo,
	}
}

func (s *categoryService) GetCategories(ctx context.Context) ([]dto.CategoryResponse, error) {
	categories, err := s.categoryRepo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	categoryMap := make(map[string]*dto.CategoryResponse)

	for _, category := range categories {
		categoryMap[category.ID] = &dto.CategoryResponse{
			ID:       category.ID,
			Name:     category.Name,
			Slug:     category.Slug,
			ParentID: category.ParentID,
			Children: []dto.CategoryResponse{},
		}
	}

	for _, category := range categories {
		current := categoryMap[category.ID]
		if current == nil {
			continue
		}

		if category.ParentID != nil {
			parent := categoryMap[*category.ParentID]
			if parent != nil {
				parent.Children = append(parent.Children, *current)
			}
		}

	}

	rootCategories := []dto.CategoryResponse{}

	for _, category := range categories {
		if category.ParentID == nil {
			current := categoryMap[category.ID]

			if current != nil {
				rootCategories = append(rootCategories, *current)
			}
		}
	}

	return rootCategories, nil
}
