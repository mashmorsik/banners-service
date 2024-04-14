package cache

import (
	"context"
	"sync"
	"time"

	"github.com/mashmorsik/banners-service/config"
	"github.com/mashmorsik/banners-service/pkg/models"
	"github.com/mashmorsik/logger"
)

var cache sync.Map

type BannerCache struct {
	Ctx                    context.Context
	evictionWorkerDuration time.Duration
	Config                 *config.Config
}

func NewBannerCache(ctx context.Context, evictionWorkerDuration time.Duration, conf *config.Config) BannerCache {
	bc := BannerCache{Ctx: ctx, evictionWorkerDuration: evictionWorkerDuration, Config: conf}
	bc.evictionWorker()
	return bc
}

type Item struct {
	BannerContent models.Content
	Eviction      time.Time
}

func (b *BannerCache) Set(key string, bannerContent models.Content) {
	cache.Store(key, Item{
		BannerContent: bannerContent,
		Eviction:      time.Now().Add(b.Config.Cache.BannerExpiration),
	})
}

func (b *BannerCache) Get(key string) (*models.Content, bool) {
	foundItem, ok := cache.Load(key)
	if !ok {
		return nil, false
	}

	item, ok := b.isInvalidType(foundItem)
	if !ok {
		return nil, false
	}

	return &item.BannerContent, true
}

func (b *BannerCache) evictionWorker() {
	ticker := time.NewTicker(b.evictionWorkerDuration)
	for {
		select {
		case <-b.Ctx.Done():
			return
		case <-ticker.C:
			cache.Range(func(key any, value any) bool {
				item, ok := b.isInvalidType(value)
				if !ok {
					return true
				}

				if item.Eviction.Before(time.Now()) {
					cache.Delete(key)
				}

				return true
			})
		}
	}
}

func (b *BannerCache) isInvalidType(foundItem any) (item *Item, invalid bool) {
	var cacheItem Item
	switch foundItem.(type) {
	case Item:
		cacheItem = foundItem.(Item)
	default:
		logger.Errf("invalid cache item %#+v", foundItem)
		return nil, false
	}

	return &cacheItem, true
}
