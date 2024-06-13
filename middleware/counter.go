package middleware

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
)

var counter *prometheus.CounterVec

func init() {
	counter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "How many HTTP requests processed, partitioned by status code and HTTP method.",
		},
		[]string{"code", "method", "path"},
	)
	prometheus.MustRegister(counter)
}

func Counter() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		status := strconv.Itoa(c.Writer.Status())
		counter.WithLabelValues(status, c.Request.Method, c.Request.URL.Path).Inc()
	}
}
