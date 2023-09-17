package metrics

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Server struct {
	server *http.Server
}

func NewServer(listenAddr string) (*Server, error) {
	server := &http.Server{
		Addr:              listenAddr,
		ReadTimeout:       5 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      5 * time.Second,
		IdleTimeout:       5 * time.Second,
	}

	router := http.NewServeMux()
	server.Handler = router

	app := &Server{server: server}

	registry := prometheus.NewRegistry()

	if err := registry.Register(collectors.NewGoCollector()); err != nil {
		return nil, fmt.Errorf("registering metrics for collectors.NewGoCollector %w", err)
	}
	if err := registry.Register(totalRequests); err != nil {
		return nil, fmt.Errorf("registering metrics for totalRequests %w", err)
	}
	if err := registry.Register(responseStatus); err != nil {
		return nil, fmt.Errorf("registering metrics for responseStatus %w", err)
	}
	if err := registry.Register(httpDuration); err != nil {
		return nil, fmt.Errorf("registering metrics for httpDuration %w", err)
	}

	router.Handle("/metrics", promhttp.HandlerFor(registry, promhttp.HandlerOpts{
		Registry:          registry,
		EnableOpenMetrics: true,
	}))

	return app, nil
}

func (s *Server) ListenAndServe() error              { return s.server.ListenAndServe() }
func (s *Server) Shutdown(ctx context.Context) error { return s.server.Shutdown(ctx) }
func (s *Server) Close() error                       { return s.server.Close() }
