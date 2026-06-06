package repositories

import (
	"time"

	"gorm.io/gorm"

	"github.com/aridevk/dark-kitchen/services/auth-service/internal/models"
)

type RefreshTokenRepository struct {
	db *gorm.DB
}

func NewRefreshTokenRepository(db *gorm.DB) *RefreshTokenRepository {
	return &RefreshTokenRepository{db: db}
}

func (r *RefreshTokenRepository) Create(token *models.RefreshToken) error {
	return r.db.Create(token).Error
}

func (r *RefreshTokenRepository) FindValidByHash(tokenHash string) (*models.RefreshToken, error) {
	var refreshToken models.RefreshToken

	err := r.db.
		Where("token_hash = ?", tokenHash).
		Where("revoked_at IS NULL").
		Where("expires_at > ?", time.Now()).
		First(&refreshToken).Error

	if err != nil {
		return nil, err
	}

	return &refreshToken, nil
}

func (r *RefreshTokenRepository) RevokeByID(id uint) error {
	now := time.Now()

	return r.db.Model(&models.RefreshToken{}).
		Where("id = ?", id).
		Update("revoked_at", now).
		Error
}
