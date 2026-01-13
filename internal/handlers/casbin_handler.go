package handlers

import (
	"encoding/json"
	"net/http"

	"fortyfour-backend/internal/dto"
	"fortyfour-backend/internal/middleware"
	"fortyfour-backend/internal/services"
	"fortyfour-backend/internal/utils"
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

// ===== helper (UNTUK TEST) =====

func getUserFromContext(r *http.Request) (userID string, role string, ok bool) {
	uid, uidOk := r.Context().Value(middleware.UserIDKey).(string)
	rle, roleOk := r.Context().Value(middleware.Role).(string)

	if !uidOk || !roleOk {
		return "", "", false
	}
	return uid, rle, true
}

// ===== HANDLERS =====

// AddPolicy adds a single permission
func (h *CasbinHandler) AddPolicy(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.RespondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	userID, userRole, ok := getUserFromContext(r)
	if ok && userRole != "admin" {
		utils.RespondError(w, http.StatusForbidden, "Only admin can manage policies")
		return
	}

	var req dto.AddPolicyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

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

	if h.sseService != nil && ok {
		h.sseService.NotifyCreate("policy", map[string]any{
			"role":     req.Role,
			"resource": req.Resource,
			"action":   req.Action,
		}, userID)
	}

	utils.RespondJSON(w, http.StatusCreated, map[string]interface{}{
		"message": "Policy added successfully",
		"policy": map[string]string{
			"role":     req.Role,
			"resource": req.Resource,
			"action":   req.Action,
		},
	})
}

// BulkAddPolicies adds multiple policies
func (h *CasbinHandler) BulkAddPolicies(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.RespondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	userID, userRole, ok := getUserFromContext(r)
	if ok && userRole != "admin" {
		utils.RespondError(w, http.StatusForbidden, "Only admin can manage policies")
		return
	}

	var req dto.BulkAddPolicyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	var policies [][]string
	for _, p := range req.Policies {
		policies = append(policies, []string{req.Role, p.Resource, p.Action})
	}

	result, err := h.casbinService.BulkAddPolicies(policies)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

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

	if h.sseService != nil && ok && len(result.Added) > 0 {
		h.sseService.NotifyCreate("policy.bulk", map[string]any{
			"role":  req.Role,
			"count": len(result.Added),
		}, userID)
	}

	status := http.StatusCreated
	message := "Bulk policies added successfully"
	if len(result.Existing) > 0 {
		status = http.StatusPartialContent
		message = "Bulk policies partially added"
	}

	utils.RespondJSON(w, status, map[string]interface{}{
		"message":  message,
		"role":     req.Role,
		"added":    result.Added,
		"existing": result.Existing,
	})
}

// RemovePolicy removes a permission
func (h *CasbinHandler) RemovePolicy(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		utils.RespondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	userID, userRole, ok := getUserFromContext(r)
	if ok && userRole != "admin" {
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

	if h.sseService != nil && ok {
		h.sseService.NotifyDelete("policy", map[string]any{
			"role":     req.Role,
			"resource": req.Resource,
			"action":   req.Action,
		}, userID)
	}

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

// GetRolePermissions returns permissions for a role
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

	perms := h.casbinService.GetRolePermissions(role)

	utils.RespondJSON(w, http.StatusOK, map[string]interface{}{
		"role":        role,
		"permissions": perms,
		"count":       len(perms),
	})
}
