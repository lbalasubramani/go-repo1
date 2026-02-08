package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
	"sync"

	"github.com/lbalasubramani/go-repo1/internal/models"
	"github.com/lbalasubramani/go-repo1/pkg/logger"
)

type ItemHandler struct {
	logger *logger.Logger
	store  *ItemStore
}

type ItemStore struct {
	mu    sync.RWMutex
	items map[string]*models.Item
}

func NewItemHandler(log *logger.Logger) *ItemHandler {
	return &ItemHandler{
		logger: log,
		store: &ItemStore{
			items: make(map[string]*models.Item),
		},
	}
}

func (h *ItemHandler) HandleItems(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.listItems(w, r)
	case http.MethodPost:
		h.createItem(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *ItemHandler) HandleItemByID(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	id := strings.TrimPrefix(path, "/items/")
	if id == "" {
		http.Error(w, "Item ID required", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.getItem(w, r, path)
	case http.MethodPut:
		h.updateItem(w, r, path)
	case http.MethodDelete:
		h.deleteItem(w, r, path)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *ItemHandler) listItems(w http.ResponseWriter, r *http.Request) {
	h.store.mu.RLock()
	defer h.store.mu.RUnlock()

	items := make([]*models.Item, 0, len(h.store.items))
	for _, item := range h.store.items {
		items = append(items, item)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(items)
}

func (h *ItemHandler) createItem(w http.ResponseWriter, r *http.Request) {
	var item models.Item
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if item.ID == "" {
		http.Error(w, "Item ID required", http.StatusBadRequest)
		return
	}

	h.store.mu.Lock()
	defer h.store.mu.Unlock()

	if _, exists := h.store.items[item.ID]; exists {
		http.Error(w, "Item already exists", http.StatusConflict)
		return
	}

	h.store.items[item.ID] = &item
	h.logger.Info("Created item: " + item.ID)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(item)
}

func (h *ItemHandler) getItem(w http.ResponseWriter, r *http.Request, id string) {
	h.store.mu.RLock()
	defer h.store.mu.RUnlock()

	item, exists := h.store.items[id]
	if !exists {
		http.Error(w, "Item not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(item)
}

func (h *ItemHandler) updateItem(w http.ResponseWriter, r *http.Request, id string) {
	var item models.Item
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	h.store.mu.Lock()
	defer h.store.mu.Unlock()

	if _, exists := h.store.items[id]; !exists {
		http.Error(w, "Item not found", http.StatusNotFound)
		return
	}

	item.ID = id
	h.store.items[id] = &item
	h.logger.Info("Updated item: " + id)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(item)
}

func (h *ItemHandler) deleteItem(w http.ResponseWriter, r *http.Request, id string) {
	h.store.mu.Lock()
	defer h.store.mu.Unlock()

	if _, exists := h.store.items[id]; !exists {
		http.Error(w, "Item not found", http.StatusNotFound)
		return
	}

	delete(h.store.items, id)
	h.logger.Info("Deleted item: " + id)

	w.WriteHeader(http.StatusNoContent)
}
