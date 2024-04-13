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
	"github.com/mashmorsik/banners-service/pkg/token"
	"github.com/mashmorsik/logger"
	"github.com/pkg/errors"
	"github.com/rs/cors"
	"golang.org/x/sync/errgroup"
	"net/http"
	"slices"
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
	router.HandleFunc("/banner/{id}", s.UpdateBanner).Methods(http.MethodPatch)
	router.HandleFunc("/banner/{id}", s.DeleteBanner).Methods(http.MethodDelete)
	router.HandleFunc("/token", s.MakeToken).Methods(http.MethodGet)

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
		TagIDs:    append([]int{}, tagID),
		FeatureID: featureID,
		IsActive:  true,
		Latest:    useLatest,
	}

	var respBanner *models.Content
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

func (s *HTTPServer) GetAdminBanner(w http.ResponseWriter, r *http.Request) {
	_ = r.Header.Get("Authorization")
	// TODO token validation

	var tagID int
	tagIDStr := r.URL.Query().Get("tag_id")
	if tagIDStr != "" {
		var err error
		tagID, err = strconv.Atoi(tagIDStr)
		if err != nil {
			http.Error(w, "invalid tagID", http.StatusBadRequest)
			return
		}
	}

	var featureID int
	featureIDStr := r.URL.Query().Get("feature_id")
	if tagIDStr != "" {
		var err error
		featureID, err = strconv.Atoi(featureIDStr)
		if err != nil {
			http.Error(w, "invalid featureID", http.StatusBadRequest)
			return
		}
	}

	var limit int
	limitStr := r.URL.Query().Get("limit")
	if limitStr != "" {
		var err error
		limit, err = strconv.Atoi(limitStr)
		if err != nil {
			http.Error(w, "invalid limit", http.StatusBadRequest)
			return
		}
	} else {
		limit = 0
	}

	var offset int
	offsetStr := r.URL.Query().Get("offset")
	if limitStr != "" {
		var err error
		offset, err = strconv.Atoi(offsetStr)
		if err != nil {
			http.Error(w, "invalid limit", http.StatusBadRequest)
			return
		}
	} else {
		offset = 0
	}

	reqBanner := &models.Banner{
		TagIDs:    append([]int{}, tagID),
		FeatureID: featureID,
		IsActive:  true,
	}

	var respBanner []*models.Banner
	respBanner, err := s.Banners.GetForAdmin(reqBanner, limit, offset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

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

func (s *HTTPServer) CreateBanner(w http.ResponseWriter, r *http.Request) {
	_ = r.Header.Get("Authorization")
	// TODO token validation

	var b *models.Banner
	err := json.NewDecoder(r.Body).Decode(&b)
	if err != nil {
		http.Error(w, "Failed to parse request body", http.StatusBadRequest)
		return
	}

	err = s.Banners.Create(b)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	w.WriteHeader(http.StatusCreated)
}

func (s *HTTPServer) UpdateBanner(w http.ResponseWriter, r *http.Request) {
	_ = r.Header.Get("Authorization")
	// TODO token validation

	var b *models.Banner
	err := json.NewDecoder(r.Body).Decode(&b)
	if err != nil {
		http.Error(w, "Failed to parse request body", http.StatusBadRequest)
		return
	}

	params := mux.Vars(r)
	bannerIDStr := params["id"]

	b.ID, err = strconv.Atoi(bannerIDStr)
	if err != nil {
		http.Error(w, "invalid bannerID", http.StatusBadRequest)
	}

	err = s.Banners.Update(b)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}

func (s *HTTPServer) DeleteBanner(w http.ResponseWriter, r *http.Request) {
	_ = r.Header.Get("Authorization")
	// TODO token validation

	params := mux.Vars(r)
	bannerIDStr := params["id"]

	bannerID, err := strconv.Atoi(bannerIDStr)
	if err != nil {
		http.Error(w, "invalid bannerID", http.StatusBadRequest)
	}
	// TODO tagID validation

	err = s.Banners.Delete(bannerID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	w.WriteHeader(http.StatusNoContent)
}

func (s *HTTPServer) MakeToken(w http.ResponseWriter, r *http.Request) {
	role := r.URL.Query().Get("role")
	if role == "" || !slices.Contains([]token.Role{token.RoleAdmin, token.RoleUser}, token.Role(role)) {
		http.Error(w, "invalid role", http.StatusBadRequest)
		return
	}

	sign, err := token.Create(token.Role(role))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	s.writeResponse(w, []byte(sign))
}

func (s *HTTPServer) writeResponse(w http.ResponseWriter, response []byte) {
	_, err := w.Write(response)
	if err != nil {
		logger.Errf("failed to write response: %v", err)
		http.Error(w, "failed to write response", http.StatusInternalServerError)
		return
	}
}
