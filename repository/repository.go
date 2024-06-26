package repository

import (
	"database/sql"
	"github.com/mashmorsik/banners-service/pkg/models"
)

type Repository interface {
	GetForUser(b *models.Banner) (*models.Banner, error)
	GetForAdmin(b *models.Banner, limit, offset int) ([]*models.Banner, error)
	CreateBanner(tx *sql.Tx, b *models.Banner) (int, error)
	CreateContent(tx *sql.Tx, b *models.Banner) error
	CreateFeatureTags(tx *sql.Tx, b *models.Banner) error
	Create(b *models.Banner) error
	MergeUpdateVersion(tx *sql.Tx, b *models.Banner, lastVersion int) (*models.Banner, error)
	UpdateBanner(tx *sql.Tx, b *models.Banner, lastVersion int) error
	UpdateFeatureTag(tx *sql.Tx, b *models.Banner) error
	UpdateBannerContent(tx *sql.Tx, b *models.Banner) error
	Update(b *models.Banner) error
	Delete(bannerID int) error
	CheckTagFeatureOverlap(b *models.Banner) (int, error)
	GetBannerActiveVersions() (map[int]int, error)
	SetVersionActive(bannerID, version int) error
	AddNewTag(banner *models.Banner) error
	AddNewFeature(banner *models.Banner) error
}
