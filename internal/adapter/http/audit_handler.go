package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"prevailing-wage-service/internal/usecase"
)

type AuditHandler struct {
	auditUC usecase.IDeterminationAuditUseCase
}

func NewAuditHandler(auditUC usecase.IDeterminationAuditUseCase) *AuditHandler {
	return &AuditHandler{auditUC: auditUC}
}

// GET /api/v1/determinations/{determination_number}
func (h *AuditHandler) GetDetermination(w http.ResponseWriter, r *http.Request) {
	number := chi.URLParam(r, "determination_number")
	if number == "" {
		RespondError(w, http.StatusBadRequest, "URL parameter 'determination_number' is required")
		return
	}

	result, err := h.auditUC.GetDeterminationByNumber(r.Context(), number)
	if err != nil {
		RespondError(w, http.StatusNotFound, err.Error())
		return
	}

	RespondJSON(w, http.StatusOK, result)
}
