package http

import (
	"encoding/json"
	nethttp "net/http"

	"github.com/padiazg/broken-hexagon/internal/adapters/primary/http/dto"
	"github.com/padiazg/broken-hexagon/internal/core/services/products"
)

// Handler exposes product endpoints over HTTP.
type Handler struct {
	svc *products.Service
}

// NewHandler wires the HTTP handler.
func NewHandler(svc *products.Service) *Handler {
	return &Handler{svc: svc}
}

// GetProduct serves GET /products?id=...
func (h *Handler) GetProduct(w nethttp.ResponseWriter, r *nethttp.Request) {
	p, err := h.svc.Get(r.Context(), r.URL.Query().Get("id"))
	if err != nil {
		nethttp.Error(w, err.Error(), nethttp.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(dto.ProductResponse{ID: p.ID, Name: p.Name, Price: p.Price})
}
