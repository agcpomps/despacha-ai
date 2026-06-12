package service

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/agcpomps/despacha-ai/backend/internal/domain"
	"github.com/agcpomps/despacha-ai/backend/internal/dto"
	"github.com/agcpomps/despacha-ai/backend/internal/repository"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("invalid phone or password")
	ErrPhoneAlreadyExists = errors.New("phone already exists")
	ErrUserSuspended      = errors.New("user account is suspended")
)

type AuthService interface {
	Register(ctx context.Context, req dto.RegisterRequest) (*dto.AuthResponse, error)
	Login(ctx context.Context, req dto.LoginRequest) (*dto.AuthResponse, error)
}

type authService struct {
	userRepo  repository.UserRepository
	jwtSecret string
}

func NewAuthService(userRepo repository.UserRepository, jwtSecret string) AuthService {
	return &authService{
		userRepo:  userRepo,
		jwtSecret: jwtSecret,
	}
}

func (s *authService) Register(ctx context.Context, req dto.RegisterRequest) (*dto.AuthResponse, error) {
	existingUser, err := s.userRepo.FindByPhone(ctx, req.Phone)
	if err == nil && existingUser != nil {
		return nil, ErrPhoneAlreadyExists
	}

	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &domain.User{
		Name:         req.Name,
		Phone:        req.Phone,
		Email:        req.Email,
		PasswordHash: string(passwordHash),
		Role:         domain.RoleUser,
		Status:       domain.UserStatusActive,
	}

	cretaedUser, err := s.userRepo.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	token, err := s.generateToken(cretaedUser)
	if err != nil {
		return nil, err
	}

	return &dto.AuthResponse{
		AccessToken: token,
		User:        toUserResponse(cretaedUser),
	}, nil
}

func (s *authService) Login(ctx context.Context, req dto.LoginRequest) (*dto.AuthResponse, error) {
	user, err := s.userRepo.FindByPhone(ctx, req.Phone)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	if user.Status == domain.UserStatusSuspended {
		return nil, ErrUserSuspended
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	token, err := s.generateToken(user)
	if err != nil {
		return nil, err
	}

	return &dto.AuthResponse{
		AccessToken: token,
		User:        toUserResponse(user),
	}, nil
}

func (s *authService) generateToken(user *domain.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"role":    string(user.Role),
		"phone":   user.Phone,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(s.jwtSecret))
}

func toUserResponse(user *domain.User) dto.UserResponse {
	return dto.UserResponse{
		ID:         user.ID,
		Name:       user.Name,
		Phone:      user.Phone,
		Email:      user.Email,
		AvatarURL:  user.AvatarURL,
		Role:       string(user.Role),
		Status:     string(user.Status),
		IsVerified: user.IsVerified,
		CreatedAt:  user.CreatedAt.Format(time.RFC3339),
	}
}
