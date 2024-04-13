package banner

import (
	"context"
	"database/sql"
	"github.com/mashmorsik/banners-service/config"
	"github.com/mashmorsik/banners-service/infrastructure/data/cache"
	"github.com/mashmorsik/banners-service/pkg/models"
	"github.com/mashmorsik/banners-service/repository"
	errs "github.com/pkg/errors"
	"strconv"
)

type Banner struct {
	Ctx    context.Context
	Repo   repository.Repository
	Config *config.Config
	Cache  *cache.BannerCache
}

func NewBanner(ctx context.Context, repo repository.Repository, conf *config.Config, cache *cache.BannerCache) *Banner {
	return &Banner{Ctx: ctx, Repo: repo, Config: conf, Cache: cache}
}

func (b *Banner) GetForUser(req *models.Banner) (*models.Content, error) {
	cacheKey := strconv.Itoa(req.FeatureID) + strconv.Itoa(req.TagIDs[0])
	content, ok := b.Cache.Get(cacheKey)
	if !ok {
		banner, err := b.Repo.GetForUser(req)
		if err != nil {
			return nil, errs.WithMessage(err, "banner not found")
		}

		b.Cache.Set(cacheKey, *banner)

		return banner, nil
	}

	return content, nil
}

func (b *Banner) GetForUserLatest(req *models.Banner) (*models.Content, error) {
	content, err := b.Repo.GetForUser(req)
	if err != nil {
		return nil, errs.WithMessage(err, "banner not found")
	}

	return content, nil
}

func (b *Banner) GetForAdmin(req *models.Banner, limit, offset int) ([]*models.Banner, error) {
	banners, err := b.Repo.GetForAdmin(req, limit, offset)
	if err != nil {
		return nil, errs.WithMessage(err, "banner not found")
	}

	return banners, nil
}

func (b *Banner) Create(req *models.Banner) error {
	_, err := b.Repo.CheckTagFeatureOverlap(req)
	if err != nil {
		if errs.Is(err, sql.ErrNoRows) {
			err = b.Repo.Create(req)
			if err != nil {
				return errs.WithMessagef(err, "fail to create banner with id: %d", req.ID)
			}
			return nil
		}
		return errs.WithMessagef(err, "fail to execute CheckTagOverlap request with id: %d", req.ID)
	}

	return errs.Errorf("banner with tag: %d and feature: %d already exists", req.TagIDs, req.FeatureID)
}

func (b *Banner) Update(req *models.Banner) error {
	bannerID, err := b.Repo.CheckTagFeatureOverlap(req)
	if err != nil || bannerID == req.ID {
		if errs.Is(err, sql.ErrNoRows) || bannerID == req.ID {
			err = b.Repo.Update(req)
			if err != nil {
				return errs.WithMessagef(err, "fail to update banner with id: %d", req.ID)
			}
			return nil
		}
		return errs.WithMessagef(err, "fail to execute CheckTagOverlap request with id: %d", req.ID)
	}

	return errs.Errorf("banner with tag: %d and feature: %d already exists and active", req.TagIDs, req.FeatureID)
}

func (b *Banner) Delete(bannerID int) error {
	err := b.Repo.Delete(bannerID)
	if err != nil {
		return errs.WithMessage(err, "banner not found")
	}

	return nil
}
