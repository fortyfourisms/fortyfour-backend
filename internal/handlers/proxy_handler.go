package handlers

import (
	"fortyfour-backend/internal/middleware"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

type ProxyHandler struct {
	target      *url.URL
	proxy       *httputil.ReverseProxy
	internalKey string
}

func NewProxyHandler(targetURL string, internalKey string) *ProxyHandler {
	target, err := url.Parse(targetURL)
	if err != nil {
		log.Fatalf("Failed to parse proxy target URL: %v", err)
	}

	proxy := httputil.NewSingleHostReverseProxy(target)

	// Customizing the Director to ensure correct path and headers
	originalDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		originalDirector(req)
		req.Host = target.Host

		// Get user info from context (set by AuthMiddleware in Main API)
		userID, _ := req.Context().Value(middleware.UserIDKey).(string)
		role, _ := req.Context().Value(middleware.RoleKey).(string)

		// Inject headers for IKAS to trust
		req.Header.Set("X-User-ID", userID)
		req.Header.Set("X-User-Role", role)
		req.Header.Set("X-Internal-Key", internalKey)

		log.Printf("Proxying request: %s %s -> %s (User: %s, Role: %s)", req.Method, req.URL.Path, targetURL, userID, role)
	}

	return &ProxyHandler{
		target:      target,
		proxy:       proxy,
		internalKey: internalKey,
	}
}

func (h *ProxyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.proxy.ServeHTTP(w, r)
}
