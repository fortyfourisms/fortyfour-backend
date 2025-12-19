package handlers

import (
	"encoding/json"
	"fortyfour-backend/internal/dto"
	"fortyfour-backend/internal/services"
	"fortyfour-backend/internal/utils"
	"log"
	"net/http"
	"strconv"
)

type PostHandler struct {
	postService *services.PostService
	sseService  *services.SSEService
}

func NewPostHandler(postService *services.PostService, sseService *services.SSEService) *PostHandler {
	return &PostHandler{
		postService: postService,
		sseService:  sseService,
	}
}

func (h *PostHandler) CreatePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.RespondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	userID := r.Context().Value("user_id").(string)
	if userID == "" {
		utils.RespondError(w, http.StatusUnauthorized, "Missing user ID")
		return
	}

	var req dto.CreatePostRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	post, err := h.postService.CreatePost(req.Title, req.Content, userID)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	// SSE Notif Create
	h.sseService.NotifyCreate(
		"post", // resource
		post,   // payload
		userID, // actor
	)

	utils.RespondJSON(w, http.StatusCreated, post)
}

func (h *PostHandler) GetPosts(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.RespondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	posts, err := h.postService.GetAllPosts()
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondJSON(w, http.StatusOK, posts)
}

func (h *PostHandler) GetPost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.RespondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	idString := r.URL.Query().Get("id")
	if idString == "" {
		utils.RespondError(w, http.StatusBadRequest, "Post ID is required")
		return
	}
	id, err := strconv.Atoi(idString)
	if err != nil {
		log.Fatal("Error during string to int conversion:", err)
	}

	post, err := h.postService.GetPostByID(id)
	if err != nil {
		utils.RespondError(w, http.StatusNotFound, err.Error())
		return
	}

	utils.RespondJSON(w, http.StatusOK, post)
}

func (h *PostHandler) UpdatePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		utils.RespondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	userID := r.Context().Value("user_id").(string)
	if userID == "" {
		utils.RespondError(w, http.StatusUnauthorized, "Missing user ID")
		return
	}

	idString := r.URL.Query().Get("id")
	if idString == "" {
		utils.RespondError(w, http.StatusBadRequest, "Post ID is required")
		return
	}

	id, err := strconv.Atoi(idString)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid post ID")
		return
	}

	var req dto.UpdatePostRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	post, err := h.postService.UpdatePost(id, req.Title, req.Content, userID)
	if err != nil {
		if err.Error() == "unauthorized" {
			utils.RespondError(w, http.StatusForbidden, "You can only update your own posts")
			return
		}
		utils.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	// SSE Notif Update
	h.sseService.NotifyUpdate(
		"post", // resource
		post,   // payload
		userID, // actor
	)

	utils.RespondJSON(w, http.StatusOK, post)
}

func (h *PostHandler) DeletePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		utils.RespondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	userID := r.Context().Value("user_id").(string)
	if userID == "" {
		utils.RespondError(w, http.StatusUnauthorized, "Missing user ID")
		return
	}

	idString := r.URL.Query().Get("id")
	if idString == "" {
		utils.RespondError(w, http.StatusBadRequest, "Post ID is required")
		return
	}

	id, err := strconv.Atoi(idString)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid post ID")
		return
	}

	if err := h.postService.DeletePost(id, userID); err != nil {
		if err.Error() == "unauthorized" {
			utils.RespondError(w, http.StatusForbidden, "You can only delete your own posts")
			return
		}
		utils.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	// SSE Notif Delete
	h.sseService.NotifyDelete(
		"post", // resource
		id,     // payload
		userID, // actor
	)

	utils.RespondJSON(w, http.StatusOK, map[string]string{"message": "Post deleted successfully"})
}
