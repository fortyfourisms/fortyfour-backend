package handlers

import (
	"encoding/json"
	"fortyfour-backend/internal/dto"
	"fortyfour-backend/internal/middleware"
	"fortyfour-backend/internal/services"
	"fortyfour-backend/internal/utils"
	"net/http"

	"fortyfour-backend/pkg/logger"
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
	rle, roleOk := r.Context().Value(middleware.RoleKey).(string)

	if !uidOk || !roleOk {
		return "", "", false
	}
	return uid, rle, true
}

// @Summary Add Casbin policy
// @Description Menambahkan satu permission policy (admin only)
// @Tags Casbin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param policy body dto.AddPolicyRequest true "Policy data"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} dto.ErrorResponse
// @Failure 403 {object} dto.ErrorResponse
// @Failure 409 {object} dto.ErrorResponse
// @Router /api/casbin/policies/add [post]
func (h *CasbinHandler) AddPolicy(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.RespondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// Get user from context
	userID, userRole, ok := getUserFromContext(r)
	if ok && userRole != "admin" {
		utils.RespondError(w, http.StatusForbidden, "Only admin can manage policies")
		return
	}

	var req dto.AddPolicyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request body")
		logger.Error(err, "operation failed")
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
		logger.Error(err, "operation failed")
		return
	}

	if !added {
		utils.RespondError(w, http.StatusConflict, "Policy already exists")
		return
	}

	// Send SSE notification if service available and user context exists
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

// @Summary Bulk add Casbin policies
// @Description Menambahkan banyak permission policies sekaligus (admin only)
// @Tags Casbin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param policies body dto.BulkAddPolicyRequest true "Bulk policy data"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} dto.ErrorResponse
// @Failure 403 {object} dto.ErrorResponse
// @Failure 409 {object} dto.ErrorResponse
// @Router /api/casbin/policies/bulk [post]
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
		logger.Error(err, "operation failed")
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
		logger.Error(err, "operation failed")
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
	if h.sseService != nil && ok && len(result.Added) > 0 {
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

// @Summary Remove Casbin policy
// @Description Menghapus satu permission policy (admin only)
// @Tags Casbin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param policy body dto.RemovePolicyRequest true "Policy to remove"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} dto.ErrorResponse
// @Failure 403 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Router /api/casbin/policies/remove [delete]
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
		logger.Error(err, "operation failed")
		utils.RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	removed, err := h.casbinService.RemovePolicy(req.Role, req.Resource, req.Action)
	if err != nil {
		logger.Error(err, "operation failed")
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

// @Summary Get all Casbin policies
// @Description Mengambil semua permission policies
// @Tags Casbin
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Router /api/casbin/policies [get]
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

// @Summary Get role permissions
// @Description Mengambil permissions untuk role tertentu
// @Tags Casbin
// @Produce json
// @Security BearerAuth
// @Param role query string true "Role name"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} dto.ErrorResponse
// @Router /api/casbin/permissions [get]
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
