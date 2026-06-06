package dto

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
	Role     string `json:"role"`
}

type AuthUserResponse struct {
	ID    uint   `json:"id"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	AccessToken  string           `json:"accessToken"`
	RefreshToken string           `json:"refreshToken"`
	User         AuthUserResponse `json:"user"`
}
type RefreshRequest struct {
	RefreshToken string `json:"refreshToken" binding:"required"`
}

type LogoutRequest struct {
	RefreshToken string `json:"refreshToken" binding:"required"`
}
