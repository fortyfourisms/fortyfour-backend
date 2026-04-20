package handlers

import (
	"net/http"
	"strings"
	"time"

	"fortyfour-backend/internal/dto"
	"fortyfour-backend/internal/middleware"
	"fortyfour-backend/internal/services"
	"fortyfour-backend/internal/utils"
	"strconv"
)

type NotificationHandler struct {
	svc *services.NotificationService
}

func NewNotificationHandler(svc *services.NotificationService) *NotificationHandler {
	return &NotificationHandler{svc: svc}
}

// GetAll godoc
// @Summary      Get all notifications
// @Description  Mengambil semua notifikasi milik user yang sedang login
// @Tags         Notifications
// @Security     BearerAuth
// @Produce      json
// @Success      200  {object}  dto.NotificationListResponse
// @Failure      401  {object}  dto.ErrorResponse
// @Failure      500  {object}  dto.ErrorResponse
// @Router       /api/notifications [get]
func (h *NotificationHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	if userID == "" {
		utils.RespondError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	notifs, err := h.svc.GetAll(userID)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	unread := 0
	items := make([]dto.NotificationResponse, 0, len(notifs))
	for _, n := range notifs {
		if !n.Read {
			unread++
		}
		items = append(items, dto.NotificationResponse{
			ID:        n.ID,
			Type:      string(n.Type),
			Message:   n.Message,
			Read:      n.Read,
			CreatedAt: n.CreatedAt.Format(time.RFC3339),
		})
	}

	utils.RespondJSON(w, http.StatusOK, dto.NotificationListResponse{
		Notifications: items,
		UnreadCount:   unread,
	})
}

// MarkRead godoc
// @Summary      Mark notification as read
// @Description  Menandai satu notifikasi sebagai sudah dibaca
// @Tags         Notifications
// @Security     BearerAuth
// @Produce      json
// @Param        id  path  uint64  true  "Notification ID"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  dto.ErrorResponse
// @Failure      401  {object}  dto.ErrorResponse
// @Failure      500  {object}  dto.ErrorResponse
// @Router       /api/notifications/{id}/read [patch]
func (h *NotificationHandler) MarkRead(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	if userID == "" {
		utils.RespondError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	idStr := strings.TrimPrefix(r.URL.Path, "/api/notifications/")
	idStr = strings.TrimSuffix(idStr, "/read")
	if idStr == "" {
		utils.RespondError(w, http.StatusBadRequest, "notification id wajib diisi")
		return
	}

	notifID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "invalid notification id format")
		return
	}

	if err := h.svc.MarkRead(userID, notifID); err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondJSON(w, http.StatusOK, map[string]string{"message": "notifikasi ditandai sudah dibaca"})
}

// MarkAllRead godoc
// @Summary      Mark all notifications as read
// @Description  Menandai semua notifikasi user sebagai sudah dibaca
// @Tags         Notifications
// @Security     BearerAuth
// @Produce      json
// @Success      200  {object}  map[string]string
// @Failure      401  {object}  dto.ErrorResponse
// @Failure      500  {object}  dto.ErrorResponse
// @Router       /api/notifications/read-all [patch]
func (h *NotificationHandler) MarkAllRead(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	if userID == "" {
		utils.RespondError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	if err := h.svc.MarkAllRead(userID); err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondJSON(w, http.StatusOK, map[string]string{"message": "semua notifikasi ditandai sudah dibaca"})
}

// Delete godoc
// @Summary      Delete a notification
// @Description  Menghapus satu notifikasi berdasarkan ID
// @Tags         Notifications
// @Security     BearerAuth
// @Produce      json
// @Param        id  path  uint64  true  "Notification ID"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  dto.ErrorResponse
// @Failure      401  {object}  dto.ErrorResponse
// @Failure      500  {object}  dto.ErrorResponse
// @Router       /api/notifications/{id} [delete]
func (h *NotificationHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	if userID == "" {
		utils.RespondError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	idStr := strings.TrimPrefix(r.URL.Path, "/api/notifications/")
	idStr = strings.TrimSuffix(idStr, "/")
	if idStr == "" {
		utils.RespondError(w, http.StatusBadRequest, "notification id wajib diisi")
		return
	}

	notifID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "invalid notification id format")
		return
	}

	if err := h.svc.Delete(userID, notifID); err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondJSON(w, http.StatusOK, map[string]string{"message": "notifikasi berhasil dihapus"})
}

// DeleteAll godoc
// @Summary      Delete all notifications
// @Description  Menghapus semua notifikasi milik user
// @Tags         Notifications
// @Security     BearerAuth
// @Produce      json
// @Success      200  {object}  map[string]string
// @Failure      401  {object}  dto.ErrorResponse
// @Failure      500  {object}  dto.ErrorResponse
// @Router       /api/notifications [delete]
func (h *NotificationHandler) DeleteAll(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	if userID == "" {
		utils.RespondError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	if err := h.svc.DeleteAll(userID); err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondJSON(w, http.StatusOK, map[string]string{"message": "semua notifikasi berhasil dihapus"})
}

// ServeHTTP routes semua request notifikasi
func (h *NotificationHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

	switch {
	// PATCH /api/notifications/read-all
	case r.Method == http.MethodPatch && path == "/api/notifications/read-all":
		h.MarkAllRead(w, r)

	// PATCH /api/notifications/{id}/read
	case r.Method == http.MethodPatch && strings.HasSuffix(path, "/read"):
		h.MarkRead(w, r)

	// DELETE /api/notifications (semua)
	case r.Method == http.MethodDelete && (path == "/api/notifications" || path == "/api/notifications/"):
		h.DeleteAll(w, r)

	// DELETE /api/notifications/{id}
	case r.Method == http.MethodDelete:
		h.Delete(w, r)

	// GET /api/notifications
	case r.Method == http.MethodGet && (path == "/api/notifications" || path == "/api/notifications/"):
		h.GetAll(w, r)

	default:
		utils.RespondError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}
