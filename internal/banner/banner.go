package banner

import (
	"context"
	"github.com/mashmorsik/banners-service/config"
	"github.com/mashmorsik/banners-service/pkg/models"
	"github.com/mashmorsik/banners-service/repository"
	errs "github.com/pkg/errors"
)

type Banner struct {
	Ctx    context.Context
	Repo   repository.Repository
	Config *config.Config
}

func NewBanner(ctx context.Context, repo repository.Repository, conf *config.Config) *Banner {
	return &Banner{Ctx: ctx, Repo: repo, Config: conf}
}

func (b *Banner) GetForUser(req *models.Banner) (*models.Banner, error) {
	banner, err := b.Repo.GetForUser(req)
	if err != nil {
		return nil, errs.WithMessage(err, "banner not found")
	}

	return banner, nil
}

func (b *Banner) GetForUserLatest(req *models.Banner) (*models.Banner, error) {
	banner, err := b.Repo.GetForUserLatest(req)
	if err != nil {
		return nil, errs.WithMessage(err, "banner not found")
	}

	return banner, nil
}

func (b *Banner) GetForAdmin(req *models.Banner) ([]*models.Banner, error) {
	banners, err := b.Repo.GetForAdmin(req)
	if err != nil {
		return nil, errs.WithMessage(err, "banner not found")
	}

	return banners, nil
}
