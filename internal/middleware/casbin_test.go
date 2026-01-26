package middleware

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/casbin/casbin/v3"
	"github.com/casbin/casbin/v3/model"
	"github.com/stretchr/testify/assert"
)

// Helper function to create a test handler
func casbinTestHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("success"))
}

// Helper function to create a basic Casbin enforcer for testing
func createTestEnforcer() *casbin.Enforcer {
	m := model.NewModel()
	m.AddDef("r", "r", "sub, obj, act")
	m.AddDef("p", "p", "sub, obj, act")
	m.AddDef("e", "e", "some(where (p.eft == allow))")
	m.AddDef("m", "m", "r.sub == p.sub && r.obj == p.obj && r.act == p.act")

	enforcer, _ := casbin.NewEnforcer(m)
	return enforcer
}

func TestNewCasbinMiddleware(t *testing.T) {
	enforcer := createTestEnforcer()
	middleware := NewCasbinMiddleware(enforcer)

	assert.NotNil(t, middleware)
	assert.Equal(t, enforcer, middleware.enforcer)
}

func TestAuthorize_NoRoleInHeader(t *testing.T) {
	enforcer := createTestEnforcer()
	middleware := NewCasbinMiddleware(enforcer)

	req := httptest.NewRequest(http.MethodGet, "/api/posts", nil)
	w := httptest.NewRecorder()

	handler := middleware.Authorize(casbinTestHandler)
	handler(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "Role not found in token")
}

func TestAuthorize_EmptyRoleInHeader(t *testing.T) {
	enforcer := createTestEnforcer()
	middleware := NewCasbinMiddleware(enforcer)

	req := httptest.NewRequest(http.MethodGet, "/api/posts", nil)
	req.Header.Set("X-User-Role", "")
	w := httptest.NewRecorder()

	handler := middleware.Authorize(casbinTestHandler)
	handler(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "Role not found in token")
}

func TestAuthorizeCasbin_PermissionDenied(t *testing.T) {
	enforcer := createTestEnforcer()
	// Add policy: only admin can POST to /api/posts
	enforcer.AddPolicy("admin", "/api/posts", "POST")

	middleware := NewCasbinMiddleware(enforcer)

	// Capture log output
	var logBuf bytes.Buffer
	log.SetOutput(&logBuf)
	defer log.SetOutput(nil)

	req := httptest.NewRequest(http.MethodPost, "/api/posts", nil)
	req.Header.Set("X-User-Role", "user") // user role trying to POST
	w := httptest.NewRecorder()

	handler := middleware.Authorize(casbinTestHandler)
	handler(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
	assert.Contains(t, w.Body.String(), "You don't have permission to access this resource")
	assert.Contains(t, logBuf.String(), "Access denied")
	assert.Contains(t, logBuf.String(), "role=user")
	assert.Contains(t, logBuf.String(), "resource=/api/posts")
	assert.Contains(t, logBuf.String(), "action=POST")
}

func TestAuthorizeCasbin_PermissionGranted(t *testing.T) {
	enforcer := createTestEnforcer()
	// Add policy: admin can POST to /api/posts
	enforcer.AddPolicy("admin", "/api/posts", "POST")

	middleware := NewCasbinMiddleware(enforcer)

	// Capture log output
	var logBuf bytes.Buffer
	log.SetOutput(&logBuf)
	defer log.SetOutput(nil)

	req := httptest.NewRequest(http.MethodPost, "/api/posts", nil)
	req.Header.Set("X-User-Role", "admin")
	w := httptest.NewRecorder()

	handler := middleware.Authorize(casbinTestHandler)
	handler(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "success", w.Body.String())
	assert.Contains(t, logBuf.String(), "Access granted")
	assert.Contains(t, logBuf.String(), "role=admin")
	assert.Contains(t, logBuf.String(), "resource=/api/posts")
	assert.Contains(t, logBuf.String(), "action=POST")
}

// func TestAuthorize_MultipleRolesAndResources(t *testing.T) {
// 	enforcer := createTestEnforcer()

// 	// Add policies
// 	enforcer.AddPolicy("user", "/api/posts", "GET")
// 	enforcer.AddPolicy("user", "/api/posts", "POST")
// 	enforcer.AddPolicy("admin", "/api/posts", "GET")
// 	enforcer.AddPolicy("admin", "/api/posts", "POST")
// 	enforcer.AddPolicy("admin", "/api/posts", "DELETE")
// 	enforcer.AddPolicy("admin", "/api/users", "GET")
// 	enforcer.AddPolicy("admin", "/api/users", "POST")
// 	enforcer.AddPolicy("admin", "/api/users", "DELETE")

// 	middleware := NewCasbinMiddleware(enforcer)

// 	tests := []struct {
// 		name           string
// 		role           string
// 		resource       string
// 		method         string
// 		expectedStatus int
// 		shouldSucceed  bool
// 	}{
// 		{
// 			name:           "user can GET posts",
// 			role:           "user",
// 			resource:       "/api/posts",
// 			method:         http.MethodGet,
// 			expectedStatus: http.StatusOK,
// 			shouldSucceed:  true,
// 		},
// 		{
// 			name:           "user can POST posts",
// 			role:           "user",
// 			resource:       "/api/posts",
// 			method:         http.MethodPost,
// 			expectedStatus: http.StatusOK,
// 			shouldSucceed:  true,
// 		},
// 		{
// 			name:           "user cannot DELETE posts",
// 			role:           "user",
// 			resource:       "/api/posts",
// 			method:         http.MethodDelete,
// 			expectedStatus: http.StatusForbidden,
// 			shouldSucceed:  false,
// 		},
// 		{
// 			name:           "user cannot access users endpoint",
// 			role:           "user",
// 			resource:       "/api/users",
// 			method:         http.MethodGet,
// 			expectedStatus: http.StatusForbidden,
// 			shouldSucceed:  false,
// 		},
// 		{
// 			name:           "admin can DELETE posts",
// 			role:           "admin",
// 			resource:       "/api/posts",
// 			method:         http.MethodDelete,
// 			expectedStatus: http.StatusOK,
// 			shouldSucceed:  true,
// 		},
// 		{
// 			name:           "admin can GET users",
// 			role:           "admin",
// 			resource:       "/api/users",
// 			method:         http.MethodGet,
// 			expectedStatus: http.StatusOK,
// 			shouldSucceed:  true,
// 		},
// 		{
// 			name:           "admin can DELETE users",
// 			role:           "admin",
// 			resource:       "/api/users",
// 			method:         http.MethodDelete,
// 			expectedStatus: http.StatusOK,
// 			shouldSucceed:  true,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			req := httptest.NewRequest(tt.method, tt.resource, nil)
// 			req.Header.Set("X-User-Role", tt.role)
// 			w := httptest.NewRecorder()

// 			handler := middleware.Authorize(casbinTestHandler)
// 			handler(w, req)

// 			assert.Equal(t, tt.expectedStatus, w.Code)
// 			if tt.shouldSucceed {
// 				assert.Equal(t, "success", w.Body.String())
// 			} else {
// 				assert.Contains(t, w.Body.String(), "You don't have permission to access this resource")
// 			}
// 		})
// 	}
// }

// func TestAuthorize_DifferentHTTPMethods(t *testing.T) {
// 	enforcer := createTestEnforcer()

// 	// Add selective permissions
// 	enforcer.AddPolicy("viewer", "/api/posts", "GET")
// 	enforcer.AddPolicy("editor", "/api/posts", "GET")
// 	enforcer.AddPolicy("editor", "/api/posts", "POST")
// 	enforcer.AddPolicy("editor", "/api/posts", "PUT")

// 	middleware := NewCasbinMiddleware(enforcer)

// 	tests := []struct {
// 		name           string
// 		role           string
// 		method         string
// 		expectedStatus int
// 	}{
// 		{
// 			name:           "viewer GET allowed",
// 			role:           "viewer",
// 			method:         http.MethodGet,
// 			expectedStatus: http.StatusOK,
// 		},
// 		{
// 			name:           "viewer POST denied",
// 			role:           "viewer",
// 			method:         http.MethodPost,
// 			expectedStatus: http.StatusForbidden,
// 		},
// 		{
// 			name:           "viewer PUT denied",
// 			role:           "viewer",
// 			method:         http.MethodPut,
// 			expectedStatus: http.StatusForbidden,
// 		},
// 		{
// 			name:           "viewer DELETE denied",
// 			role:           "viewer",
// 			method:         http.MethodDelete,
// 			expectedStatus: http.StatusForbidden,
// 		},
// 		{
// 			name:           "editor GET allowed",
// 			role:           "editor",
// 			method:         http.MethodGet,
// 			expectedStatus: http.StatusOK,
// 		},
// 		{
// 			name:           "editor POST allowed",
// 			role:           "editor",
// 			method:         http.MethodPost,
// 			expectedStatus: http.StatusOK,
// 		},
// 		{
// 			name:           "editor PUT allowed",
// 			role:           "editor",
// 			method:         http.MethodPut,
// 			expectedStatus: http.StatusOK,
// 		},
// 		{
// 			name:           "editor DELETE denied",
// 			role:           "editor",
// 			method:         http.MethodDelete,
// 			expectedStatus: http.StatusForbidden,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			req := httptest.NewRequest(tt.method, "/api/posts", nil)
// 			req.Header.Set("X-User-Role", tt.role)
// 			w := httptest.NewRecorder()

// 			handler := middleware.Authorize(casbinTestHandler)
// 			handler(w, req)

// 			assert.Equal(t, tt.expectedStatus, w.Code)
// 		})
// 	}
// }

// func TestAuthorize_DifferentResourcePaths(t *testing.T) {
// 	enforcer := createTestEnforcer()

// 	// Add policies for different resources
// 	enforcer.AddPolicy("user", "/api/posts", "GET")
// 	enforcer.AddPolicy("user", "/api/comments", "GET")
// 	enforcer.AddPolicy("user", "/api/profile", "GET")
// 	enforcer.AddPolicy("admin", "/api/settings", "GET")

// 	middleware := NewCasbinMiddleware(enforcer)

// 	tests := []struct {
// 		name           string
// 		role           string
// 		path           string
// 		expectedStatus int
// 	}{
// 		{
// 			name:           "user can access posts",
// 			role:           "user",
// 			path:           "/api/posts",
// 			expectedStatus: http.StatusOK,
// 		},
// 		{
// 			name:           "user can access comments",
// 			role:           "user",
// 			path:           "/api/comments",
// 			expectedStatus: http.StatusOK,
// 		},
// 		{
// 			name:           "user can access profile",
// 			role:           "user",
// 			path:           "/api/profile",
// 			expectedStatus: http.StatusOK,
// 		},
// 		{
// 			name:           "user cannot access settings",
// 			role:           "user",
// 			path:           "/api/settings",
// 			expectedStatus: http.StatusForbidden,
// 		},
// 		{
// 			name:           "admin can access settings",
// 			role:           "admin",
// 			path:           "/api/settings",
// 			expectedStatus: http.StatusOK,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			req := httptest.NewRequest(http.MethodGet, tt.path, nil)
// 			req.Header.Set("X-User-Role", tt.role)
// 			w := httptest.NewRecorder()

// 			handler := middleware.Authorize(casbinTestHandler)
// 			handler(w, req)

// 			assert.Equal(t, tt.expectedStatus, w.Code)
// 		})
// 	}
// }

func TestAuthorize_LoggingBehavior(t *testing.T) {
	enforcer := createTestEnforcer()
	enforcer.AddPolicy("admin", "/api/posts", "POST")

	middleware := NewCasbinMiddleware(enforcer)

	tests := []struct {
		name              string
		role              string
		shouldGrantAccess bool
		expectedLogMsg    string
	}{
		{
			name:              "logs access granted",
			role:              "admin",
			shouldGrantAccess: true,
			expectedLogMsg:    "Access granted",
		},
		{
			name:              "logs access denied",
			role:              "user",
			shouldGrantAccess: false,
			expectedLogMsg:    "Access denied",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var logBuf bytes.Buffer
			log.SetOutput(&logBuf)
			defer log.SetOutput(nil)

			req := httptest.NewRequest(http.MethodPost, "/api/posts", nil)
			req.Header.Set("X-User-Role", tt.role)
			w := httptest.NewRecorder()

			handler := middleware.Authorize(casbinTestHandler)
			handler(w, req)

			logOutput := logBuf.String()
			assert.Contains(t, logOutput, tt.expectedLogMsg)
			assert.Contains(t, logOutput, "role="+tt.role)
			assert.Contains(t, logOutput, "resource=/api/posts")
			assert.Contains(t, logOutput, "action=POST")
		})
	}
}
