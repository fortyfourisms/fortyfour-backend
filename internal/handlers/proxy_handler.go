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
			perusahaanID, _ := req.In.Context().Value(middleware.IDPerusahaanKey).(string)

			// Inject headers for IKAS to trust
			req.Out.Header.Set("X-User-ID", userID)
			req.Out.Header.Set("X-User-Role", role)
			req.Out.Header.Set("X-Perusahaan-ID", perusahaanID)
			req.Out.Header.Set("X-Internal-Key", internalKey)

			log.Printf("Proxying request: %s %s -> %s (User: %s, Role: %s, Perusahaan: %s)", req.Out.Method, req.Out.URL.Path, targetURL, userID, role, perusahaanID)
		},
	}

	return &ProxyHandler{
		target:      target,
		proxy:       proxy,
		internalKey: internalKey,
	}
}

// ProxyMaturityRoot godoc
//
//	@Summary		Proxy request ke layanan IKAS
//	@Description	Meneruskan request ke layanan IKAS melalui gateway internal untuk endpoint root maturity.
//	@Tags			IKAS Proxy
//	@Produce		json
//	@Security		BearerAuth
//	@Success		200	{object}	map[string]interface{}
//	@Failure		401	{object}	map[string]interface{}
//	@Failure		403	{object}	map[string]interface{}
//	@Failure		500	{object}	map[string]interface{}
//	@Router			/api/maturity [get]
//	@Router			/api/maturity [post]
//	@Router			/api/maturity [put]
//	@Router			/api/maturity [delete]
//
// ProxyMaturityPath godoc
//
//	@Summary		Proxy request ke layanan IKAS subpath
//	@Description	Meneruskan request ke layanan IKAS untuk seluruh subpath di bawah endpoint maturity.
//	@Tags			IKAS Proxy
//	@Produce		json
//	@Security		BearerAuth
//	@Param			path	path		string	true	"Subpath endpoint IKAS"
//	@Success		200		{object}	map[string]interface{}
//	@Failure		401		{object}	map[string]interface{}
//	@Failure		403		{object}	map[string]interface{}
//	@Failure		500		{object}	map[string]interface{}
//	@Router			/api/maturity/{path} [get]
//	@Router			/api/maturity/{path} [post]
//	@Router			/api/maturity/{path} [put]
//	@Router			/api/maturity/{path} [delete]
func (h *ProxyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.proxy.ServeHTTP(w, r)
}
