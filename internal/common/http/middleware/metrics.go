package middleware

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var responseMetrics = promauto.NewSummaryVec(
	prometheus.SummaryOpts{
		Namespace:  "fluire",
		Subsystem:  "http",
		Name:       "response_time",
		Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.01},
	},
	[]string{"method", "path"},
)

func observeResponseTime(duration time.Duration, method, path string) {
	responseMetrics.WithLabelValues(method, path).
		Observe(float64(duration.Milliseconds()))
}

func Metrics() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			defer func() {
				routPattern := chi.RouteContext(r.Context()).RoutePattern
				observeResponseTime(time.Since(start), r.Method, routPattern())
			}()

			next.ServeHTTP(w, r)
		})
	}
}
