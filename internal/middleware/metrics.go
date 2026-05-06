package middleware

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	httpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "path", "status"},
	)

	httpRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request latency in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path"},
	)

	httpRequestsInFlight = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "http_requests_in_flight",
			Help: "Current number of HTTP requests in flight",
		},
	)

	imageUploadTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "image_upload_total",
			Help: "Total number of image uploads",
		},
	)

	imageDeleteTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "image_delete_total",
			Help: "Total number of image deletions",
		},
	)

	diskUsagePercent = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "disk_usage_percent",
			Help: "Current disk usage percentage",
		},
	)

	fileCount = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "file_count",
			Help: "Current number of stored files",
		},
	)
)

func init() {
	prometheus.MustRegister(httpRequestsTotal)
	prometheus.MustRegister(httpRequestDuration)
	prometheus.MustRegister(httpRequestsInFlight)
	prometheus.MustRegister(imageUploadTotal)
	prometheus.MustRegister(imageDeleteTotal)
	prometheus.MustRegister(diskUsagePercent)
	prometheus.MustRegister(fileCount)
}

func PrometheusHandler() gin.HandlerFunc {
	h := promhttp.Handler()
	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

func PrometheusMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		httpRequestsInFlight.Inc()
		start := time.Now()

		c.Next()

		httpRequestsInFlight.Dec()
		duration := time.Since(start).Seconds()
		status := strconv.Itoa(c.Writer.Status())
		path := c.FullPath()

		httpRequestsTotal.WithLabelValues(c.Request.Method, path, status).Inc()
		httpRequestDuration.WithLabelValues(c.Request.Method, path).Observe(duration)
	}
}

func RecordImageUpload() {
	imageUploadTotal.Inc()
}

func RecordImageDelete() {
	imageDeleteTotal.Inc()
}

func UpdateDiskMetrics(usagePct float64, totalFiles int) {
	diskUsagePercent.Set(usagePct)
	fileCount.Set(float64(totalFiles))
}
