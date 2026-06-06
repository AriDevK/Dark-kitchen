package services

import (
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/aridevk/dark-kitchen/packages/go/common/auth"
	"github.com/aridevk/dark-kitchen/services/auth-service/internal/dto"
	"github.com/aridevk/dark-kitchen/services/auth-service/internal/models"
	"github.com/aridevk/dark-kitchen/services/auth-service/internal/repositories"
)

var ErrEmailAlreadyExists = errors.New("email already exists")
var ErrInvalidCredentials = errors.New("invalid credentials")
var ErrUserNotFound = errors.New("user not found")
var ErrInvalidRefreshToken = errors.New("invalid refresh token")

type AuthService struct {
	userRepo         *repositories.UserRepository
	jwtSecret        string
	jwtTTL           int
	refreshTokenRepo *repositories.RefreshTokenRepository
	refreshTTLDays   int
}

func NewAuthService(userRepo *repositories.UserRepository, refreshTokenRepo *repositories.RefreshTokenRepository, jwtSecret string, jwtTTL int, refreshTTLDays int) *AuthService {
	return &AuthService{
		userRepo:         userRepo,
		jwtSecret:        jwtSecret,
		jwtTTL:           jwtTTL,
		refreshTokenRepo: refreshTokenRepo,
		refreshTTLDays:   refreshTTLDays,
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

	refreshToken, err := s.createRefreshToken(user.ID)
	if err != nil {
		return nil, err
	}

	return &dto.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User: dto.AuthUserResponse{
			ID:    user.ID,
			Email: user.Email,
			Role:  user.Role,
		},
	}, nil
}

func (s *AuthService) Me(userID uint) (*dto.AuthUserResponse, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, ErrUserNotFound
	}

	return &dto.AuthUserResponse{
		ID:    user.ID,
		Email: user.Email,
		Role:  user.Role,
	}, nil
}

func (s *AuthService) createRefreshToken(userID uint) (string, error) {
	rawToken, err := auth.GenerateRefreshToken()
	if err != nil {
		return "", err
	}

	tokenHash := auth.HashRefreshToken(rawToken)

	refreshToken := &models.RefreshToken{
		UserID:    userID,
		TokenHash: tokenHash,
		ExpiresAt: time.Now().AddDate(0, 0, s.refreshTTLDays),
	}

	if err := s.refreshTokenRepo.Create(refreshToken); err != nil {
		return "", err
	}

	return rawToken, nil
}

func (s *AuthService) Refresh(req dto.RefreshRequest) (*dto.LoginResponse, error) {
	tokenHash := auth.HashRefreshToken(req.RefreshToken)

	storedToken, err := s.refreshTokenRepo.FindValidByHash(tokenHash)
	if err != nil {
		return nil, ErrInvalidRefreshToken
	}

	user, err := s.userRepo.FindByID(storedToken.UserID)
	if err != nil {
		return nil, ErrInvalidRefreshToken
	}

	if err := s.refreshTokenRepo.RevokeByID(storedToken.ID); err != nil {
		return nil, err
	}

	newAccessToken, err := auth.GenerateAccessToken(
		user.ID,
		user.Email,
		user.Role,
		s.jwtSecret,
		s.jwtTTL,
	)
	if err != nil {
		return nil, err
	}

	newRefreshToken, err := s.createRefreshToken(user.ID)
	if err != nil {
		return nil, err
	}

	return &dto.LoginResponse{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
		User: dto.AuthUserResponse{
			ID:    user.ID,
			Email: user.Email,
			Role:  user.Role,
		},
	}, nil
}

func (s *AuthService) Logout(req dto.LogoutRequest) error {
	tokenHash := auth.HashRefreshToken(req.RefreshToken)

	storedToken, err := s.refreshTokenRepo.FindValidByHash(tokenHash)
	if err != nil {
		return ErrInvalidRefreshToken
	}

	return s.refreshTokenRepo.RevokeByID(storedToken.ID)
}
