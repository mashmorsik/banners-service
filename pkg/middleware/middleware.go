package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-openapi/runtime"
	"github.com/mashmorsik/banners-service/pkg/token"
	"github.com/mashmorsik/logger"
)

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		logger.Infof("Started %s %s", r.Method, r.URL.Path)

		next.ServeHTTP(w, r)
		logger.Infof("Completed %s %s in %v", r.Method, r.URL.Path, time.Since(start))
	})
}

func AdminAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h, ok := r.Header[runtime.HeaderAuthorization]
		if !ok {
			http.Error(w, fmt.Sprintf("%s header is missing ", runtime.HeaderAuthorization), http.StatusForbidden)
			return
		}
		if len(h) == 0 {
			http.Error(w, fmt.Sprintf("%s header is empty", runtime.HeaderAuthorization), http.StatusForbidden)
			return
		}

		claims, err := token.Validate(h[0])
		if err != nil {
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}

		roles, err := token.GetRoles(claims)
		if err != nil {
			http.Error(w, fmt.Sprint("get auth roles failed"), http.StatusForbidden)
			return
		}

		if !token.IsAdmin(roles) {
			http.Error(w, fmt.Sprintf("%s role is missing", token.RoleAdmin), http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func UserAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h, ok := r.Header[runtime.HeaderAuthorization]
		if !ok {
			http.Error(w, fmt.Sprintf("%s header is missing ", runtime.HeaderAuthorization), http.StatusForbidden)
			return
		}
		if len(h) == 0 {
			http.Error(w, fmt.Sprintf("%s header is empty", runtime.HeaderAuthorization), http.StatusForbidden)
			return
		}

		claims, err := token.Validate(h[0])
		if err != nil {
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}

		roles, err := token.GetRoles(claims)
		if err != nil {
			http.Error(w, fmt.Sprint("get auth roles failed"), http.StatusForbidden)
			return
		}

		if !token.IsUser(roles) {
			http.Error(w, fmt.Sprintf("%s role is missing", token.RoleUser), http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}
