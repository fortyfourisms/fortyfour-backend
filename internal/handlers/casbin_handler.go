package handlers

import (
	"encoding/json"
	"fortyfour-backend/internal/dto"
	"fortyfour-backend/internal/middleware"
	"fortyfour-backend/internal/services"
	"fortyfour-backend/internal/utils"
	"net/http"
)

type CasbinHandler struct {
	casbinService *services.CasbinService
	sseService    *services.SSEService
}

func NewCasbinHandler(casbinService *services.CasbinService, sseService *services.SSEService) *CasbinHandler {
	return &CasbinHandler{
		casbinService: casbinService,
		sseService:    sseService,
	}
}

// AddPolicy adds a single permission
func (h *CasbinHandler) AddPolicy(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.RespondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// Only admin can manage policies
	userID := r.Context().Value(middleware.UserIDKey).(string)
	userRole := r.Context().Value(middleware.Role).(string)

	if userRole != "admin" {
		utils.RespondError(w, http.StatusForbidden, "Only admin can manage policies")
		return
	}

	var req dto.AddPolicyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate request
	if req.Role == "" || req.Resource == "" || req.Action == "" {
		utils.RespondError(w, http.StatusBadRequest, "Role, resource, and action are required")
		return
	}

	added, err := h.casbinService.AddPolicy(req.Role, req.Resource, req.Action)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if !added {
		utils.RespondError(w, http.StatusConflict, "Policy already exists")
		return
	}

	h.sseService.NotifyCreate("policy", map[string]any{
		"role":     req.Role,
		"resource": req.Resource,
		"action":   req.Action,
	}, userID)

	utils.RespondJSON(w, http.StatusCreated, map[string]interface{}{
		"message": "Policy added successfully",
		"policy": map[string]string{
			"role":     req.Role,
			"resource": req.Resource,
			"action":   req.Action,
		},
	})
}

// BulkAddPolicies adds multiple policies at once
func (h *CasbinHandler) BulkAddPolicies(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.RespondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	userID := r.Context().Value(middleware.UserIDKey).(string)
	userRole := r.Context().Value(middleware.Role).(string)

	if userRole != "admin" {
		utils.RespondError(w, http.StatusForbidden, "Only admin can manage policies")
		return
	}

	var req dto.BulkAddPolicyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Convert to casbin format [][]string
	var policies [][]string
	for _, p := range req.Policies {
		policies = append(policies, []string{req.Role, p.Resource, p.Action})
	}

	result, err := h.casbinService.BulkAddPolicies(policies)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Jika semua data sudah ada, return conflict error
	if len(result.Added) == 0 && len(result.Existing) > 0 {
		utils.RespondJSON(w, http.StatusConflict, map[string]interface{}{
			"message":  "All policies already exist",
			"role":     req.Role,
			"existing": result.Existing,
			"total":    len(req.Policies),
			"added":    0,
		})
		return
	}

	// Jika ada yang berhasil ditambahkan, kirim notifikasi
	if len(result.Added) > 0 {
		h.sseService.NotifyCreate("policy.bulk", map[string]any{
			"role":  req.Role,
			"count": len(result.Added),
		}, userID)
	}

	statusCode := http.StatusCreated
	message := "Bulk policies added successfully"

	// Jika ada yang sudah exist, ubah status dan message
	if len(result.Existing) > 0 {
		statusCode = http.StatusPartialContent
		message = "Bulk policies partially added"
	}

	utils.RespondJSON(w, statusCode, map[string]interface{}{
		"message":  message,
		"role":     req.Role,
		"added":    result.Added,
		"existing": result.Existing,
		"summary": map[string]int{
			"total":    len(req.Policies),
			"added":    len(result.Added),
			"existing": len(result.Existing),
		},
	})
}

// RemovePolicy removes a permission
func (h *CasbinHandler) RemovePolicy(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		utils.RespondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	userID := r.Context().Value(middleware.UserIDKey).(string)
	userRole := r.Context().Value(middleware.Role).(string)

	if userRole != "admin" {
		utils.RespondError(w, http.StatusForbidden, "Only admin can manage policies")
		return
	}

	var req dto.RemovePolicyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	removed, err := h.casbinService.RemovePolicy(req.Role, req.Resource, req.Action)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if !removed {
		utils.RespondError(w, http.StatusNotFound, "Policy not found")
		return
	}

	h.sseService.NotifyDelete("policy", map[string]any{
		"role":     req.Role,
		"resource": req.Resource,
		"action":   req.Action,
	}, userID)

	utils.RespondJSON(w, http.StatusOK, map[string]interface{}{
		"message": "Policy removed successfully",
		"removed": true,
	})
}

// GetAllPolicies returns all policies
func (h *CasbinHandler) GetAllPolicies(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.RespondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	policies := h.casbinService.GetAllPolicies()

	utils.RespondJSON(w, http.StatusOK, map[string]interface{}{
		"policies": policies,
		"count":    len(policies),
	})
}

// GetRolePermissions returns permissions for a specific role
func (h *CasbinHandler) GetRolePermissions(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.RespondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	role := r.URL.Query().Get("role")
	if role == "" {
		utils.RespondError(w, http.StatusBadRequest, "Role parameter is required")
		return
	}

	permissions := h.casbinService.GetRolePermissions(role)

	utils.RespondJSON(w, http.StatusOK, map[string]interface{}{
		"role":        role,
		"permissions": permissions,
		"count":       len(permissions),
	})
}
