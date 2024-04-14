package banner

import (
	"context"
	"database/sql"
	"github.com/golang/mock/gomock"
	"github.com/mashmorsik/banners-service/config"
	"github.com/mashmorsik/banners-service/infrastructure/data/cache"
	"github.com/mashmorsik/banners-service/pkg/models"
	mock_repository "github.com/mashmorsik/banners-service/testdata/mock_repo"
	"github.com/mashmorsik/logger"
	errs "github.com/pkg/errors"
	"reflect"
	"strconv"
	"testing"
	"time"
)

func TestBanner_GetForUser(t *testing.T) {
	logger.BuildLogger(nil)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	conf := config.Config{}

	ctx := context.Background()
	bannerCache := cache.NewBannerCache(ctx, time.Hour, &conf)

	banner := &models.Banner{
		TagIDs:    []int{2},
		FeatureID: 4,
		Content: models.Content{
			Title: "test_1_title",
			Text:  "test_1_text",
			URL:   "test_1_url",
		},
	}

	cacheKey := strconv.Itoa(banner.FeatureID) + strconv.Itoa(banner.TagIDs[0])
	bannerCache.Set(cacheKey, banner.Content)

	mockRepo := mock_repository.NewMockRepository(ctrl)
	mockRepo.EXPECT().GetForUser(&models.Banner{
		TagIDs:    []int{5},
		FeatureID: 1,
	}).Return(nil, errs.New("banner not found"))

	type args struct {
		req *models.Banner
	}

	tests := []struct {
		name    string
		args    args
		want    *models.Content
		wantErr bool
	}{
		{
			name: "get_banner_for_user_success",
			args: args{req: &models.Banner{
				TagIDs:    []int{2},
				FeatureID: 4,
			}},
			want: &models.Content{
				Title: "test_1_title",
				Text:  "test_1_text",
				URL:   "test_1_url",
			},
			wantErr: false,
		},
		{
			name: "get_banner_for_user_fail",
			args: args{req: &models.Banner{
				TagIDs:    []int{5},
				FeatureID: 1,
			}},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Banner{
				Ctx:    ctx,
				Repo:   mockRepo,
				Config: &conf,
				Cache:  &bannerCache,
			}
			got, err := b.GetForUser(tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetForUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetForUser() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBanner_Create(t *testing.T) {
	logger.BuildLogger(nil)

	banner := &models.Banner{
		TagIDs:    []int{10, 15},
		FeatureID: 1,
		IsActive:  true,
		Content: models.Content{
			Title: "test_create_title",
			Text:  "test_create_text",
			URL:   "test_create_url",
		},
		CreatedAt: time.Time{},
		UpdatedAt: time.Time{},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	conf := config.Config{}

	ctx := context.Background()
	bannerCache := cache.NewBannerCache(ctx, time.Hour, &conf)

	mockRepo := mock_repository.NewMockRepository(ctrl)
	mockRepo.EXPECT().CheckTagFeatureOverlap(banner).Return(0, sql.ErrNoRows)
	mockRepo.EXPECT().Create(banner).Return(nil)

	type args struct {
		req *models.Banner
	}

	tests := []struct {
		name    string
		args    args
		want    *models.Content
		wantErr bool
	}{
		{
			name:    "create_banner",
			args:    args{req: banner},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Banner{
				Ctx:    ctx,
				Repo:   mockRepo,
				Config: &conf,
				Cache:  &bannerCache,
			}
			err := b.Create(tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetForUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestBanner_Delete(t *testing.T) {
	logger.BuildLogger(nil)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	conf := config.Config{}

	ctx := context.Background()
	bannerCache := cache.NewBannerCache(ctx, time.Hour, &conf)

	mockRepo := mock_repository.NewMockRepository(ctrl)
	mockRepo.EXPECT().Delete(2).Return(nil)
	mockRepo.EXPECT().Delete(3).Return(errs.New("banner not found"))

	tests := []struct {
		name     string
		bannerID int
		wantErr  bool
	}{
		{
			name:     "delete_banner_success",
			bannerID: 2,
			wantErr:  false,
		},
		{
			name:     "delete_banner_fail",
			bannerID: 3,
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Banner{
				Ctx:    ctx,
				Repo:   mockRepo,
				Config: &conf,
				Cache:  &bannerCache,
			}
			err := b.Delete(tt.bannerID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetForUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
