package repository

import (
	"database/sql"
	"github.com/mashmorsik/banners-service/pkg/models"
)

type Repository interface {
	GetForUser(bannerID int) (*models.Content, error)
	GetForAdmin(b *models.Banner, limit, offset int) ([]*models.Banner, error)
	CreateBanner(tx *sql.Tx, b *models.Banner) (int, error)
	CreateVersion(tx *sql.Tx, b *models.Banner) error
	Create(b *models.Banner) error
	UpdateBanner(tx *sql.Tx, b *models.Banner) error
	UpdateBannerVersion(tx *sql.Tx, b *models.Banner) error
	SetOldVersionInactive(tx *sql.Tx, b *models.Banner) error
	Update(b *models.Banner) error
	Delete(bannerID int) error
	CheckTagFeatureOverlap(b *models.Banner) (int, error)
}
