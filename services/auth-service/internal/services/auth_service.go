package services

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/aridevk/dark-kitchen/packages/go/common/auth"
	"github.com/aridevk/dark-kitchen/services/auth-service/internal/dto"
	"github.com/aridevk/dark-kitchen/services/auth-service/internal/models"
	"github.com/aridevk/dark-kitchen/services/auth-service/internal/repositories"
)

var ErrEmailAlreadyExists = errors.New("email already exists")
var ErrInvalidCredentials = errors.New("invalid credentials")

type AuthService struct {
	userRepo  *repositories.UserRepository
	jwtSecret string
	jwtTTL    int
}

func NewAuthService(userRepo *repositories.UserRepository, jwtSecret string, jwtTTL int) *AuthService {
	return &AuthService{
		userRepo:  userRepo,
		jwtSecret: jwtSecret,
		jwtTTL:    jwtTTL,
	}
}

func (s *AuthService) Register(req dto.RegisterRequest) (*dto.AuthUserResponse, error) {
	existingUser, err := s.userRepo.FindByEmail(req.Email)
	if err == nil && existingUser != nil {
		return nil, ErrEmailAlreadyExists
	}

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	role := req.Role
	if role == "" {
		role = "customer"
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Email:    req.Email,
		Password: string(hashedPassword),
		Role:     role,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	return &dto.AuthUserResponse{
		ID:    user.ID,
		Email: user.Email,
		Role:  user.Role,
	}, nil
}

func (s *AuthService) Login(req dto.LoginRequest) (*dto.LoginResponse, error) {
	user, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	accessToken, err := auth.GenerateAccessToken(
		user.ID,
		user.Email,
		user.Role,
		s.jwtSecret,
		s.jwtTTL,
	)
	if err != nil {
		return nil, err
	}

	return &dto.LoginResponse{
		AccessToken: accessToken,
		User: dto.AuthUserResponse{
			ID:    user.ID,
			Email: user.Email,
			Role:  user.Role,
		},
	}, nil
}
