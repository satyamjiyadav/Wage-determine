package http

import (
	"encoding/json"
	"net/http"

	"prevailing-wage-service/internal/domain/model"
	"prevailing-wage-service/internal/usecase"
)

type WageHandler struct {
	lookupUC usecase.IWageLookupUseCase
	calcUC   usecase.ILevelCalculatorUseCase
}

func NewWageHandler(lookupUC usecase.IWageLookupUseCase, calcUC usecase.ILevelCalculatorUseCase) *WageHandler {
	return &WageHandler{
		lookupUC: lookupUC,
		calcUC:   calcUC,
	}
}

// GET /api/v1/wages/lookup
func (h *WageHandler) LookupWage(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	socCode := query.Get("soc_code")
	zipCode := query.Get("zip_code")
	fipsCode := query.Get("fips_code")
	areaCode := query.Get("area_code")
	programYear := query.Get("program_year")

	if socCode == "" {
		RespondError(w, http.StatusBadRequest, "Query parameter 'soc_code' is required (e.g. ?soc_code=15-1252.00)")
		return
	}

	result, err := h.lookupUC.LookupWage(r.Context(), socCode, zipCode, fipsCode, areaCode, programYear)
	if err != nil {
		RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	RespondJSON(w, http.StatusOK, result)
}

// POST /api/v1/wages/determine-level
func (h *WageHandler) DetermineWageLevel(w http.ResponseWriter, r *http.Request) {
	var payload model.JobRequirementPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		RespondError(w, http.StatusBadRequest, "Invalid JSON payload format: "+err.Error())
		return
	}

	if payload.SOCCode == "" {
		RespondError(w, http.StatusBadRequest, "Field 'soc_code' is required in request body")
		return
	}

	result, err := h.calcUC.DetermineWageLevel(r.Context(), payload)
	if err != nil {
		RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	RespondJSON(w, http.StatusOK, result)
}

// POST /api/v1/wages/batch-lookup
func (h *WageHandler) BatchLookupWages(w http.ResponseWriter, r *http.Request) {
	var requests []model.WageLookupRequest
	if err := json.NewDecoder(r.Body).Decode(&requests); err != nil {
		RespondError(w, http.StatusBadRequest, "Invalid JSON array payload: "+err.Error())
		return
	}

	if len(requests) == 0 {
		RespondError(w, http.StatusBadRequest, "Batch array cannot be empty")
		return
	}

	if len(requests) > 100 {
		RespondError(w, http.StatusBadRequest, "Batch query limit exceeded (maximum 100 items per request)")
		return
	}

	results, err := h.lookupUC.BatchLookupWages(r.Context(), requests)
	if err != nil {
		RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	RespondJSON(w, http.StatusOK, results)
}
