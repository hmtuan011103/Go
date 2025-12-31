package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"sync"

	"github.com/gostructure/app/internal/model"
	"github.com/gostructure/app/pkg/response"
)

// In-memory user storage (replace with database in production)
var (
	users     = make(map[int64]*model.User)
	userMutex sync.RWMutex
	userIDSeq int64 = 1
)

// ListUsers returns all users
func (h *Handler) ListUsers(w http.ResponseWriter, r *http.Request) {
	userMutex.RLock()
	defer userMutex.RUnlock()

	userList := make([]*model.User, 0, len(users))
	for _, user := range users {
		userList = append(userList, user)
	}

	response.JSON(w, http.StatusOK, map[string]interface{}{
		"users": userList,
		"total": len(userList),
	})
}

// GetUser returns a specific user by ID
func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	userMutex.RLock()
	user, exists := users[id]
	userMutex.RUnlock()

	if !exists {
		response.Error(w, http.StatusNotFound, "User not found")
		return
	}

	response.JSON(w, http.StatusOK, user)
}

// CreateUser creates a new user
func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req model.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Name == "" || req.Email == "" {
		response.Error(w, http.StatusBadRequest, "Name and email are required")
		return
	}

	userMutex.Lock()
	user := &model.User{
		ID:    userIDSeq,
		Name:  req.Name,
		Email: req.Email,
	}
	users[userIDSeq] = user
	userIDSeq++
	userMutex.Unlock()

	response.JSON(w, http.StatusCreated, user)
}

// UpdateUser updates an existing user
func (h *Handler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	var req model.UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	userMutex.Lock()
	defer userMutex.Unlock()

	user, exists := users[id]
	if !exists {
		response.Error(w, http.StatusNotFound, "User not found")
		return
	}

	if req.Name != "" {
		user.Name = req.Name
	}
	if req.Email != "" {
		user.Email = req.Email
	}

	response.JSON(w, http.StatusOK, user)
}

// DeleteUser deletes a user
func (h *Handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	userMutex.Lock()
	defer userMutex.Unlock()

	if _, exists := users[id]; !exists {
		response.Error(w, http.StatusNotFound, "User not found")
		return
	}

	delete(users, id)
	response.JSON(w, http.StatusOK, map[string]string{
		"message": "User deleted successfully",
	})
}
