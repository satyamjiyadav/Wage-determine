package http

import (
	"net/http"
	"runtime"
	"strconv"
	"time"
)

type HealthHandler struct {
	startTime time.Time
}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{startTime: time.Now()}
}

// GET /healthz
func (h *HealthHandler) Healthz(w http.ResponseWriter, r *http.Request) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	RespondJSON(w, http.StatusOK, map[string]interface{}{
		"status":          "UP",
		"service":         "dol-prevailing-wage-service",
		"uptime_seconds":  int(time.Since(h.startTime).Seconds()),
		"goroutines":      runtime.NumGoroutine(),
		"allocated_bytes": m.Alloc,
	})
}

// GET /metrics
func (h *HealthHandler) Metrics(w http.ResponseWriter, r *http.Request) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	w.Header().Set("Content-Type", "text/plain; version=0.0.4")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(
		"# HELP go_goroutines Number of goroutines currently existing.\n" +
			"# TYPE go_goroutines gauge\n" +
			"go_goroutines " + strconv.Itoa(runtime.NumGoroutine()) + "\n" +
			"# HELP go_memstats_alloc_bytes Number of bytes allocated and still in use.\n" +
			"# TYPE go_memstats_alloc_bytes gauge\n" +
			"go_memstats_alloc_bytes " + strconv.FormatUint(m.Alloc, 10) + "\n",
	))
}
