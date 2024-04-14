package banners_service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/mashmorsik/banners-service/config"
	"github.com/mashmorsik/banners-service/infrastructure/data"
	"github.com/mashmorsik/banners-service/infrastructure/data/cache"
	"github.com/mashmorsik/banners-service/infrastructure/server"
	"github.com/mashmorsik/banners-service/internal/banner"
	"github.com/mashmorsik/banners-service/pkg/models"
	"github.com/mashmorsik/banners-service/pkg/token"
	"github.com/mashmorsik/banners-service/repository"
	"github.com/mashmorsik/logger"
	"io"
	"net/http"
	"testing"
	"time"
)

func Test_GetUserBanner(t *testing.T) {
	logger.BuildLogger(nil)

	conf, err := config.LoadConfig()
	if err != nil {
		logger.Errf("Error loading config: %v", err)
		return
	}

	ctx := context.Background()

	conn := data.MustConnectPostgres(ctx, conf)
	data.MustMigrate(conn)
	dat := data.NewData(ctx, conn)

	bannerCache := cache.NewBannerCache(ctx, conf.Cache.EvictionWorkerDuration, conf)
	bannerRepo := repository.NewBannerRepo(ctx, dat)
	bb := banner.NewBanner(ctx, bannerRepo, conf, &bannerCache)

	token.NewTokenManager(conf.Auth.TokenSecret)

	go func() {
		httpServer := server.NewServer(conf, *bb)
		if err = httpServer.StartServer(ctx); err != nil {
			logger.Warn(err.Error())
		}
	}()

	time.Sleep(1 * time.Second)

	// get an admin token
	getAdminToken := fmt.Sprintf("http://localhost%s/token?role=admin", conf.Server.Port)
	resp, err := http.Get(getAdminToken)
	if err != nil {
		t.Errorf("Error sending GET request: %v", err)
		return
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			return
		}
	}(resp.Body)

	adminToken, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("Error reading response body: %v", err)
	}

	// create a new active banner with TagIDs 9, 10 and FeatureID 3
	newBanner := models.Banner{
		ID:        0,
		TagIDs:    []int{9, 10},
		FeatureID: 3,
		IsActive:  true,
		Content: models.Content{
			Title: "New_super_banner",
			Text:  "Super_description",
			URL:   "Mega_super_URL",
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	postDataJSON, err := json.Marshal(newBanner)
	if err != nil {
		t.Errorf("Error marshalling banner: %v", err)
	}

	reqBody := bytes.NewBuffer(postDataJSON)

	req, err := http.Post("http://localhost"+conf.Server.Port+"/banner", "application/json", reqBody)
	if err != nil {
		t.Errorf("Error sending POST request: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+string(adminToken))

	// get a user token
	getUserToken := fmt.Sprintf("http://localhost%s/token?role=admin", conf.Server.Port)
	resp, err = http.Get(getUserToken)
	if err != nil {
		t.Errorf("Error sending GET request: %v", err)
		return
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			return
		}
	}(resp.Body)

	userToken, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("Error reading response body: %v", err)
	}

	// get a banner for the user
	req, err = http.Post("http://localhost"+conf.Server.Port+"/user_banner?tag_id=10&feature_id=3", "application/json", reqBody)
	if err != nil {
		t.Errorf("Error sending POST request: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+string(userToken))

	var userBanner models.Content
	err = json.NewDecoder(resp.Body).Decode(&userBanner)
	if err != nil {
		t.Errorf("Error decoding response body: %v", err)
	}

	// compare the content
	if newBanner.Content != userBanner {
		t.Errorf("GetUserBanner failed: expected %v, got %v", newBanner.Content, userBanner)
	}
}
