package handlers

import (
	"errors"

	"github.com/gin-gonic/gin"

	"github.com/aridevk/dark-kitchen/packages/go/common/middleware"
	commonResponse "github.com/aridevk/dark-kitchen/packages/go/common/response"
	"github.com/aridevk/dark-kitchen/services/auth-service/internal/dto"
	"github.com/aridevk/dark-kitchen/services/auth-service/internal/services"
)

type AuthHandler struct {
	authService *services.AuthService
}

func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		commonResponse.BadRequest(c, "INVALID_REQUEST", err.Error())
		return
	}

	user, err := h.authService.Register(req)
	if err != nil {
		if errors.Is(err, services.ErrEmailAlreadyExists) {
			commonResponse.BadRequest(c, "EMAIL_ALREADY_EXISTS", "Email already exists")
			return
		}

		commonResponse.InternalServerError(c)
		return
	}

	commonResponse.Created(c, user)
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		commonResponse.BadRequest(c, "INVALID_REQUEST", err.Error())
		return
	}

	result, err := h.authService.Login(req)
	if err != nil {
		if errors.Is(err, services.ErrInvalidCredentials) {
			commonResponse.Unauthorized(c, "INVALID_CREDENTIALS", "Invalid email or password")
			return
		}

		commonResponse.InternalServerError(c)
		return
	}

	commonResponse.OK(c, result)
}

func (h *AuthHandler) Me(c *gin.Context) {
	userID := middleware.GetUserID(c)

	if userID == 0 {
		commonResponse.Unauthorized(c, "UNAUTHORIZED", "Unauthorized")
		return
	}

	user, err := h.authService.Me(userID)
	if err != nil {
		if errors.Is(err, services.ErrUserNotFound) {
			commonResponse.Unauthorized(c, "USER_NOT_FOUND", "User not found")
			return
		}

		commonResponse.InternalServerError(c)
		return
	}

	commonResponse.OK(c, user)
}
