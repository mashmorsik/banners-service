package banners_service_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/mashmorsik/banners-service/config"
	"github.com/mashmorsik/banners-service/infrastructure/data"
	"github.com/mashmorsik/banners-service/infrastructure/data/cache"
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

	conf.Postgres.Host = "localhost"
	path := "http://localhost"

	ctx := context.Background()

	conn := data.MustConnectPostgres(ctx, conf)
	data.MustMigrate(conn)

	dat := data.NewData(ctx, conn)

	bannerCache := cache.NewBannerCache(ctx, conf.Cache.EvictionWorkerDuration, conf)

	bannerRepo := repository.NewBannerRepo(ctx, dat)
	_ = banner.NewBanner(ctx, bannerRepo, conf, &bannerCache)

	token.NewTokenManager(conf.Auth.TokenSecret)

	time.Sleep(2 * time.Second)

	// get an admin token
	getAdminToken := fmt.Sprintf("%s%s/token?role=admin", path, conf.Server.Port)
	adminTokenResp, err := http.Get(getAdminToken)
	if err != nil {
		t.Errorf("Error sending GET request: %v", err)
		return
	}
	defer adminTokenResp.Body.Close()

	adminToken, err := io.ReadAll(adminTokenResp.Body)
	if err != nil {
		t.Errorf("Error reading response body: %v", err)
	}

	// create a new active banner
	newBanner := &models.Banner{
		TagIDs:    []int{70, 71},
		FeatureID: 70,
		IsActive:  true,
		Content: models.Content{
			Title: "New_super_banner",
			Text:  "Super_description",
			URL:   "Mega_super_URL",
		},
	}

	postDataJSON, err := json.Marshal(newBanner)
	if err != nil {
		t.Errorf("Error marshaling banner: %v", err)
	}

	reqBody := bytes.NewBuffer(postDataJSON)

	client := &http.Client{}
	req, err := http.NewRequestWithContext(ctx, "POST", path+conf.Server.Port+"/banner", reqBody)
	if err != nil {
		logger.Errf("Error sending POST request: %v", err)
		return
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+string(adminToken))

	resp, err := client.Do(req)
	if err != nil {
		logger.Errf("Error sending POST request: %v", err)
		return
	}
	defer resp.Body.Close()

	// get a user token
	getUserToken := fmt.Sprintf("%s%s/token?role=user", path, conf.Server.Port)
	userTokenResp, err := http.Get(getUserToken)
	if err != nil {
		t.Errorf("Error sending GET request: %v", err)
		return
	}
	defer userTokenResp.Body.Close()

	userToken, err := io.ReadAll(userTokenResp.Body)
	if err != nil {
		t.Errorf("Error reading response body: %v", err)
	}

	// get a banner for the user
	clientUser := &http.Client{}
	userBannerReq, err := http.NewRequestWithContext(ctx, "GET",
		path+conf.Server.Port+"/user_banner?tag_id=70&feature_id=70", nil)
	if err != nil {
		logger.Errf(err.Error(), "Error sending GET request: %v", err)
		return
	}
	userBannerReq.Header.Add("Content-Type", "application/json")
	userBannerReq.Header.Add("Authorization", "Bearer "+string(userToken))

	userBannerResp, err := clientUser.Do(userBannerReq)
	if err != nil {
		logger.Errf(err.Error(), "Error sending GET request: %v", err)
		return
	}
	defer userBannerResp.Body.Close()

	var userBanner *models.Content
	err = json.NewDecoder(userBannerResp.Body).Decode(&userBanner)
	if err != nil {
		t.Errorf("Error decoding response body: %v", err)
	}

	// compare the content
	if newBanner.Content != *userBanner {
		t.Errorf("GetUserBanner failed: expected %v, got %v", newBanner.Content, userBanner)
	}

	// delete newBanner
	defer func() {
		bannerDelete, err := bannerRepo.GetForUser(newBanner)
		if err != nil {
			return
		}

		err = bannerRepo.Delete(bannerDelete.ID)
		if err != nil {
			logger.Errf(err.Error(), "Error sending DELETE request: %v", err)
			return
		}
	}()
}
