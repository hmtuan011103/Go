package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"sync"

	"github.com/gostructure/app/internal/model"
	"github.com/gostructure/app/pkg/response"
)

// In-memory item storage (replace with database in production)
var (
	items     = make(map[int64]*model.Item)
	itemMutex sync.RWMutex
	itemIDSeq int64 = 1
)

// ListItems returns all items
func (h *Handler) ListItems(w http.ResponseWriter, r *http.Request) {
	itemMutex.RLock()
	defer itemMutex.RUnlock()

	itemList := make([]*model.Item, 0, len(items))
	for _, item := range items {
		itemList = append(itemList, item)
	}

	response.JSON(w, http.StatusOK, map[string]interface{}{
		"items": itemList,
		"total": len(itemList),
	})
}

// GetItem returns a specific item by ID
func (h *Handler) GetItem(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid item ID")
		return
	}

	itemMutex.RLock()
	item, exists := items[id]
	itemMutex.RUnlock()

	if !exists {
		response.Error(w, http.StatusNotFound, "Item not found")
		return
	}

	response.JSON(w, http.StatusOK, item)
}

// CreateItem creates a new item
func (h *Handler) CreateItem(w http.ResponseWriter, r *http.Request) {
	var req model.CreateItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Name == "" {
		response.Error(w, http.StatusBadRequest, "Name is required")
		return
	}

	itemMutex.Lock()
	item := &model.Item{
		ID:          itemIDSeq,
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Quantity:    req.Quantity,
	}
	items[itemIDSeq] = item
	itemIDSeq++
	itemMutex.Unlock()

	response.JSON(w, http.StatusCreated, item)
}

// UpdateItem updates an existing item
func (h *Handler) UpdateItem(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid item ID")
		return
	}

	var req model.UpdateItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	itemMutex.Lock()
	defer itemMutex.Unlock()

	item, exists := items[id]
	if !exists {
		response.Error(w, http.StatusNotFound, "Item not found")
		return
	}

	if req.Name != "" {
		item.Name = req.Name
	}
	if req.Description != "" {
		item.Description = req.Description
	}
	if req.Price != nil {
		item.Price = *req.Price
	}
	if req.Quantity != nil {
		item.Quantity = *req.Quantity
	}

	response.JSON(w, http.StatusOK, item)
}

// DeleteItem deletes an item
func (h *Handler) DeleteItem(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid item ID")
		return
	}

	itemMutex.Lock()
	defer itemMutex.Unlock()

	if _, exists := items[id]; !exists {
		response.Error(w, http.StatusNotFound, "Item not found")
		return
	}

	delete(items, id)
	response.JSON(w, http.StatusOK, map[string]string{
		"message": "Item deleted successfully",
	})
}
