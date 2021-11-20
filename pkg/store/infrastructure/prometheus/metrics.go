package prometheus

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type EndpointLabelCollector interface {
	EndpointLabelForURI(uri string) string
}

type MetricsHandler interface {
	AddMetricsHandler(router *mux.Router, endpoint string)
	AddCommonMetricsMiddleware(router *mux.Router)
}

func NewMetricsHandler(endpointLabelCollector EndpointLabelCollector) (MetricsHandler, error) {
	handler := &metricsHandler{
		endpointLabelCollector: endpointLabelCollector,
	}
	if err := handler.initCommonMetrics(); err != nil {
		return handler, err
	}
	return handler, nil
}

type metricsHandler struct {
	endpointLabelCollector EndpointLabelCollector
	latencyHistogram       *prometheus.HistogramVec
	requestCounter         *prometheus.CounterVec
}

func (m *metricsHandler) AddMetricsHandler(router *mux.Router, endpoint string) {
	router.Handle(endpoint, promhttp.Handler())
}

func (m *metricsHandler) AddCommonMetricsMiddleware(router *mux.Router) {
	router.Use(m.getMiddlewareFunc())
}

func (m *metricsHandler) initCommonMetrics() error {
	labelNames := []string{"endpoint", "method", "status"}

	m.latencyHistogram = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "app_request_latency_seconds",
		Help:    "Application Request Latency",
		Buckets: prometheus.DefBuckets,
	}, labelNames)
	err := prometheus.Register(m.latencyHistogram)
	if err != nil {
		return err
	}

	m.requestCounter = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "app_request_count",
		Help: "Application Request Count",
	}, labelNames)
	err = prometheus.Register(m.requestCounter)
	if err != nil {
		return err
	}

	return nil
}

func (m *metricsHandler) getMiddlewareFunc() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			start := time.Now()
			code := http.StatusOK

			defer func() {
				httpDuration := time.Since(start)
				endpoint := m.endpointLabelCollector.EndpointLabelForURI(req.RequestURI)
				labels := []string{endpoint, req.Method, fmt.Sprintf("%d", code)}
				m.latencyHistogram.WithLabelValues(labels...).Observe(httpDuration.Seconds())
				m.requestCounter.WithLabelValues(labels...).Inc()
			}()

			rw := &responseWriter{w, http.StatusOK}
			next.ServeHTTP(rw, req)
			code = rw.statusCode
		})
	}
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}
