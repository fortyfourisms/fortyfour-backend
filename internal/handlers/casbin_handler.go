package handlers

import (
	"encoding/json"
	"fortyfour-backend/internal/dto"
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
	userID := r.Context().Value("user_id").(string)
	userRole := r.Context().Value("role").(string)

	if userRole != "admin" {
		utils.RespondError(w, http.StatusForbidden, "Only admin can manage policies")
		return
	}

	var req dto.AddPolicyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	added, err := h.casbinService.AddPolicy(req.Role, req.Resource, req.Action)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if !added {
		utils.RespondJSON(w, http.StatusOK, map[string]interface{}{
			"message": "Policy already exists",
			"added":   false,
		})
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
		"added": true,
	})
}

// BulkAddPolicies adds multiple policies at once
func (h *CasbinHandler) BulkAddPolicies(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.RespondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	userID := r.Context().Value("user_id").(string)
	userRole := r.Context().Value("role").(string)

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

	added, err := h.casbinService.AddPolicies(policies)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.sseService.NotifyCreate("policy.bulk", map[string]any{
		"role":  req.Role,
		"count": len(req.Policies),
	}, userID)

	utils.RespondJSON(w, http.StatusCreated, map[string]interface{}{
		"message": "Bulk policies added successfully",
		"role":    req.Role,
		"count":   len(req.Policies),
		"added":   added,
	})
}

// RemovePolicy removes a permission
func (h *CasbinHandler) RemovePolicy(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		utils.RespondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	userID := r.Context().Value("user_id").(string)
	userRole := r.Context().Value("role").(string)

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
