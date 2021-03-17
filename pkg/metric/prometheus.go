package metric

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/push"
	"github.com/sgraham785/gocleanarch-example/pkg/config"
)

// service implements Service interface
type service struct {
	config               *config.Specification
	pHistogram           *prometheus.HistogramVec
	httpRequestHistogram *prometheus.HistogramVec
}

//NewPrometheusService create a new prometheus service
func NewPrometheusService(c *config.Specification) (*service, error) {
	cli := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "pushgateway",
		Name:      "cmd_duration_seconds",
		Help:      "CLI application execution in seconds",
		Buckets:   prometheus.DefBuckets,
	}, []string{"name"})
	http := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "http",
		Name:      "request_duration_seconds",
		Help:      "The latency of the HTTP requests.",
		Buckets:   prometheus.DefBuckets,
	}, []string{"handler", "method", "code"})

	s := &service{
		config:               c,
		pHistogram:           cli,
		httpRequestHistogram: http,
	}
	err := prometheus.Register(s.pHistogram)
	if err != nil && err.Error() != "duplicate metrics collector registration attempted" {
		return nil, err
	}
	err = prometheus.Register(s.httpRequestHistogram)
	if err != nil && err.Error() != "duplicate metrics collector registration attempted" {
		return nil, err
	}
	return s, nil
}

//SaveCLI send metrics to server
func (s *service) SaveCLI(c *CLI) error {
	s.pHistogram.WithLabelValues(c.Name).Observe(c.Duration)
	return push.New(s.config.PrometheusPushgateway, "cmd_job").Collector(s.pHistogram).Push()
}

//SaveHTTP send metrics to server
func (s *service) SaveHTTP(h *HTTP) {
	s.httpRequestHistogram.WithLabelValues(h.Handler, h.Method, h.StatusCode).Observe(h.Duration)
}
