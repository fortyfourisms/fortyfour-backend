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

	proxy := &httputil.ReverseProxy{
		Rewrite: func(req *httputil.ProxyRequest) {
			req.SetURL(target)
			req.SetXForwarded()

			// Get user info from context (set by AuthMiddleware in Main API)
			userID, _ := req.In.Context().Value(middleware.UserIDKey).(string)
			role, _ := req.In.Context().Value(middleware.RoleKey).(string)

			// Inject headers for IKAS to trust
			req.Out.Header.Set("X-User-ID", userID)
			req.Out.Header.Set("X-User-Role", role)
			req.Out.Header.Set("X-Internal-Key", internalKey)

			log.Printf("Proxying request: %s %s -> %s (User: %s, Role: %s)", req.Out.Method, req.Out.URL.Path, targetURL, userID, role)
		},
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
