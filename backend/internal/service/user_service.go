package service

import (
	"context"
	"database/sql"
	"errors"

	"github.com/agcpomps/despacha-ai/backend/internal/domain"
	"github.com/agcpomps/despacha-ai/backend/internal/dto"
	"github.com/agcpomps/despacha-ai/backend/internal/repository"
)

var (
	ErrUserNotFound    = errors.New("user not found")
	ErrInvalidRole     = errors.New("invalid role")
	ErrCannotChangeOwn = errors.New("cannot change your own role")
)

type UserService interface {
	GetProfile(ctx context.Context, userID string) (*dto.UserResponse, error)
	ListUsers(ctx context.Context, search string, page, limit int) (*dto.PaginatedUsersResponse, error)
	UpdateUserRole(ctx context.Context, actorID, targetID, role string) (*dto.UserResponse, error)
}

type userService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{
		userRepo: userRepo,
	}
}

func (s *userService) GetProfile(ctx context.Context, userID string) (*dto.UserResponse, error) {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}

		return nil, err
	}

	response := toUserResponse(user)
	return &response, nil
}

func (s *userService) ListUsers(ctx context.Context, search string, page, limit int) (*dto.PaginatedUsersResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 50 {
		limit = 20
	}

	offset := (page - 1) * limit

	users, err := s.userRepo.List(ctx, search, limit, offset)
	if err != nil {
		return nil, err
	}

	total, err := s.userRepo.Count(ctx, search)
	if err != nil {
		return nil, err
	}

	responses := make([]dto.UserResponse, 0, len(users))
	for _, user := range users {
		responses = append(responses, toUserResponse(&user))
	}

	totalPages := 0
	if limit > 0 {
		totalPages = (total + limit - 1) / limit
	}

	return &dto.PaginatedUsersResponse{
		Data:       responses,
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: totalPages,
	}, nil
}

func (s *userService) UpdateUserRole(ctx context.Context, actorID, targetID, role string) (*dto.UserResponse, error) {
	if actorID == targetID {
		return nil, ErrCannotChangeOwn
	}

	userRole := domain.UserRole(role)
	if userRole != domain.RoleUser && userRole != domain.RoleModerator && userRole != domain.RoleAdmin {
		return nil, ErrInvalidRole
	}

	if err := s.userRepo.UpdateRole(ctx, targetID, userRole); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	user, err := s.userRepo.FindByID(ctx, targetID)
	if err != nil {
		return nil, err
	}

	response := toUserResponse(user)
	return &response, nil
}
