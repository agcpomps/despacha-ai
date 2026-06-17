package service

import (
	"context"
	"crypto/rand"
	"database/sql"
	"errors"

	"github.com/agcpomps/despacha-ai/backend/internal/domain"
	"github.com/agcpomps/despacha-ai/backend/internal/dto"
	"github.com/agcpomps/despacha-ai/backend/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrInvalidRole       = errors.New("invalid role")
	ErrCannotChangeOwn   = errors.New("cannot change your own role")
	ErrIncorrectPassword = errors.New("incorrect current password")
	ErrWeakPassword      = errors.New("password too weak")
)

type UserService interface {
	GetProfile(ctx context.Context, userID string) (*dto.UserResponse, error)
	ListUsers(ctx context.Context, search string, page, limit int) (*dto.PaginatedUsersResponse, error)
	UpdateUserRole(ctx context.Context, actorID, targetID, role string) (*dto.UserResponse, error)
	ResetUserPassword(ctx context.Context, targetID string) (string, error)
	ChangePassword(ctx context.Context, userID, currentPassword, newPassword string) error
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

// ResetUserPassword gera uma palavra-passe temporária para o utilizador, grava o
// respectivo hash e devolve o valor em texto simples para o admin a comunicar
// (ex: via WhatsApp). O utilizador deve alterá-la depois de entrar.
func (s *userService) ResetUserPassword(ctx context.Context, targetID string) (string, error) {
	tempPassword, err := generateTempPassword()
	if err != nil {
		return "", err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(tempPassword), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	if err := s.userRepo.UpdatePassword(ctx, targetID, string(hash)); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", ErrUserNotFound
		}
		return "", err
	}

	return tempPassword, nil
}

// ChangePassword permite ao próprio utilizador trocar a palavra-passe, validando
// primeiro a password actual.
func (s *userService) ChangePassword(ctx context.Context, userID, currentPassword, newPassword string) error {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrUserNotFound
		}
		return err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(currentPassword)); err != nil {
		return ErrIncorrectPassword
	}

	if len(newPassword) < 6 {
		return ErrWeakPassword
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	return s.userRepo.UpdatePassword(ctx, userID, string(hash))
}

// generateTempPassword cria uma password de 8 caracteres legíveis (sem 0/O/1/I/l
// para evitar confusão ao ditar/escrever).
func generateTempPassword() (string, error) {
	const alphabet = "ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnpqrstuvwxyz23456789"
	const length = 8

	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	result := make([]byte, length)
	for i, b := range bytes {
		result[i] = alphabet[int(b)%len(alphabet)]
	}

	return string(result), nil
}
