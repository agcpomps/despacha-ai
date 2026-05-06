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

	type categoryNode struct {
		response dto.CategoryResponse
		children []*categoryNode
	}

	categoryMap := make(map[string]*categoryNode)
	var rootCategories []*categoryNode

	for _, category := range categories {
		categoryMap[category.ID] = &categoryNode{
			response: dto.CategoryResponse{
				ID:       category.ID,
				Name:     category.Name,
				Slug:     category.Slug,
				ParentID: category.ParentID,
			},
		}
	}

	for _, category := range categories {
		current := categoryMap[category.ID]
		if category.ParentID == nil {
			rootCategories = append(rootCategories, current)
			continue
		}

		parent, exists := categoryMap[*category.ParentID]
		if exists {
			parent.children = append(parent.children, current)
		}
	}

	var buildResponse func(node *categoryNode) dto.CategoryResponse
	buildResponse = func(node *categoryNode) dto.CategoryResponse {
		response := node.response
		for _, child := range node.children {
			response.Children = append(response.Children, buildResponse(child))
		}
		return response
	}

	responses := make([]dto.CategoryResponse, 0, len(rootCategories))
	for _, category := range rootCategories {
		responses = append(responses, buildResponse(category))
	}

	return responses, nil
}
