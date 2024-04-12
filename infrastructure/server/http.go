package server

import (
	"context"
	"encoding/json"
	"github.com/go-openapi/runtime/middleware"
	"github.com/gorilla/mux"
	"github.com/mashmorsik/banners-service/config"
	"github.com/mashmorsik/banners-service/internal/banner"
	mw "github.com/mashmorsik/banners-service/pkg/middleware"
	"github.com/mashmorsik/banners-service/pkg/models"
	"github.com/mashmorsik/logger"
	"github.com/pkg/errors"
	"github.com/rs/cors"
	"golang.org/x/sync/errgroup"
	"net/http"
	"strconv"
)

type HTTPServer struct {
	Config  *config.Config
	Banners banner.Banner
}

func NewServer(conf *config.Config, banners banner.Banner) *HTTPServer {
	return &HTTPServer{Config: conf, Banners: banners}
}

func (s *HTTPServer) StartServer(ctx context.Context) error {

	router := mux.NewRouter()

	router.Handle("/swagger.yaml", http.FileServer(http.Dir("./")))
	sh := middleware.SwaggerUI(middleware.SwaggerUIOpts{
		Path:    "/swagger",
		SpecURL: "swagger.yaml",
	}, nil)
	router.Handle("/swagger", sh)

	router.HandleFunc("/user_banner", s.GetUserBanner).Methods(http.MethodGet)
	router.HandleFunc("/banner", s.GetAdminBanner).Methods(http.MethodGet)
	router.HandleFunc("/banner", s.CreateBanner).Methods(http.MethodPost)
	router.HandleFunc("/banner", s.UpdateBanner).Methods(http.MethodPatch)
	router.HandleFunc("/banner", s.DeleteBanner).Methods(http.MethodDelete)

	logger.Infof("HTTPServer is listening on port: %s\n", s.Config.Server.Port)

	router.Use(mw.LoggingMiddleware)

	httpServer := &http.Server{
		Addr:    s.Config.Server.Port,
		Handler: cors.AllowAll().Handler(router),
	}

	g, gCtx := errgroup.WithContext(ctx)
	g.Go(func() error {
		return httpServer.ListenAndServe()
	})
	g.Go(func() error {
		<-gCtx.Done()
		return httpServer.Shutdown(ctx)
	})

	if err := g.Wait(); err != nil {
		return errors.WithMessagef(err, "exit reason: %s \n", err)
	}

	return nil
}

func (s *HTTPServer) GetUserBanner(w http.ResponseWriter, r *http.Request) {
	_ = r.Header.Get("Authorization")
	// TODO token validation

	tagID, err := strconv.Atoi(r.URL.Query().Get("tag_id"))
	if err != nil {
		http.Error(w, "invalid tagID", http.StatusBadRequest)
	}
	// TODO tagID validation

	featureID, err := strconv.Atoi(r.URL.Query().Get("feature_id"))
	if err != nil {
		http.Error(w, "invalid featureID", http.StatusBadRequest)
	}
	// TODO featureID validation

	useLatest := false
	lastRevision := r.URL.Query().Get("use_last_revision")
	if lastRevision != "" {
		useLatest = true
	}

	reqBanner := &models.Banner{
		TagID:     tagID,
		FeatureID: featureID,
		IsActive:  true,
		Latest:    useLatest,
	}

	var respBanner *models.Banner
	if useLatest {
		respBanner, err = s.Banners.GetForUserLatest(reqBanner)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	} else {
		respBanner, err = s.Banners.GetForUser(reqBanner)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	}

	s.SendRespBanner(w, respBanner)
}

func (s *HTTPServer) GetAdminBanner(w http.ResponseWriter, r *http.Request) {
	_ = r.Header.Get("Authorization")
	// TODO token validation

	tagID, err := strconv.Atoi(r.URL.Query().Get("tag_id"))
	if err != nil {
		http.Error(w, "invalid tagID", http.StatusBadRequest)
	}
	// TODO tagID validation

	featureID, err := strconv.Atoi(r.URL.Query().Get("feature_id"))
	if err != nil {
		http.Error(w, "invalid featureID", http.StatusBadRequest)
	}
	// TODO featureID validation

	useLatest := false
	lastRevision := r.URL.Query().Get("use_last_revision")
	if lastRevision != "" {
		useLatest = true
	}

	reqBanner := &models.Banner{
		TagID:     tagID,
		FeatureID: featureID,
		IsActive:  true,
		Latest:    useLatest,
	}

	var respBanner *models.Banner
	if useLatest {
		respBanner, err = s.Banners.GetForUserLatest(reqBanner)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	} else {
		respBanner, err = s.Banners.GetForUser(reqBanner)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	}

	s.SendRespBanner(w, respBanner)
}

func (s *HTTPServer) CreateBanner(w http.ResponseWriter, r *http.Request) {

}

func (s *HTTPServer) UpdateBanner(w http.ResponseWriter, r *http.Request) {

}

func (s *HTTPServer) DeleteBanner(w http.ResponseWriter, r *http.Request) {

}

func (s *HTTPServer) SendRespBanner(w http.ResponseWriter, respBanner *models.Banner) {
	jsonData, err := json.Marshal(respBanner)
	if err != nil {
		logger.Errf("failed to marshal JSON: %v", err)
		http.Error(w, "Failed to marshal JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	_, err = w.Write(jsonData)
	if err != nil {
		logger.Errf("failed to write response: %v", err)
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
		return
	}
}
