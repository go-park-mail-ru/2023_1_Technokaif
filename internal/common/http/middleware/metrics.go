package middleware

import (
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	commonHTTP "github.com/go-park-mail-ru/2023_1_Technokaif/internal/common/http"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var responseTimeMetrics = promauto.NewSummaryVec(
	prometheus.SummaryOpts{
		Namespace:  "fluire",
		Subsystem:  "http",
		Name:       "response_time",
		Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01},
	},
	[]string{"method", "route", "code"},
)

func observeResponseTime(duration time.Duration, method, route, code string) {
	responseTimeMetrics.WithLabelValues(method, route, code).
		Observe(float64(duration.Microseconds()))
}

func Metrics() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			defer func() {
				// API route (URL Pattern)
				routePattern := chi.RouteContext(r.Context()).RoutePattern()

				// Status code
				code := commonHTTP.GetResponseCodeFromRequest(r)
				codeStr := ""
				if code != 0 {
					codeStr = strconv.Itoa(code)
				}

				observeResponseTime(time.Since(start), r.Method, routePattern, codeStr)
			}()

			next.ServeHTTP(w, r)
		})
	}
}
