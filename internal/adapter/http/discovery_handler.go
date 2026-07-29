package http

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"prevailing-wage-service/internal/usecase"
)

type DiscoveryHandler struct {
	discoveryUC usecase.IDiscoveryUseCase
}

func NewDiscoveryHandler(discoveryUC usecase.IDiscoveryUseCase) *DiscoveryHandler {
	return &DiscoveryHandler{discoveryUC: discoveryUC}
}

// GET /api/v1/occupations/search
func (h *DiscoveryHandler) SearchOccupations(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	limitStr := r.URL.Query().Get("limit")
	limit := 10
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	results, err := h.discoveryUC.SearchOccupations(r.Context(), q, limit)
	if err != nil {
		RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	RespondJSON(w, http.StatusOK, results)
}

// GET /api/v1/occupations/{soc_code}
func (h *DiscoveryHandler) GetOccupationDetails(w http.ResponseWriter, r *http.Request) {
	socCode := chi.URLParam(r, "soc_code")
	if socCode == "" {
		RespondError(w, http.StatusBadRequest, "URL parameter 'soc_code' is required")
		return
	}

	result, err := h.discoveryUC.GetOccupationDetails(r.Context(), socCode)
	if err != nil {
		RespondError(w, http.StatusNotFound, err.Error())
		return
	}

	RespondJSON(w, http.StatusOK, result)
}

// GET /api/v1/locations/resolve
func (h *DiscoveryHandler) ResolveLocation(w http.ResponseWriter, r *http.Request) {
	zipCode := r.URL.Query().Get("zip_code")
	if zipCode == "" {
		RespondError(w, http.StatusBadRequest, "Query parameter 'zip_code' is required (e.g. ?zip_code=94103)")
		return
	}

	result, err := h.discoveryUC.ResolveLocation(r.Context(), zipCode)
	if err != nil {
		RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	RespondJSON(w, http.StatusOK, result)
}
