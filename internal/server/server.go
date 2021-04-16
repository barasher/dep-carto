package server

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/barasher/dep-carto/internal/model"
	"github.com/barasher/dep-carto/internal/output"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog/log"
)

const (
	jsonFormat    = "json"
	dotFormat     = "dot"
	jpgFormat     = "jpg"
	defaultFormat = jsonFormat
	defaultSince  = 3600 * 365 * time.Hour
	defaultDepth  = 100
)

type Server struct {
	router *mux.Router
	model  model.Model
	port   uint
}

func NewServer(model model.Model, port uint) (*Server, error) {
	s := &Server{
		model: model,
		port:  port,
	}
	s.router = mux.NewRouter()
	registerHandler(s.router, NewAddHandler(model), "add")
	registerHandler(s.router, NewClearHandler(model), "clear")
	registerHandler(s.router, NewGetHandler(model), "get")
	registerHandler(s.router, NewGetDependenciesHandler(model), "getdependencies")
	s.router.Handle("/metrics", promhttp.Handler())
	return s, nil
}

func registerHandler(r *mux.Router, h handlerInterface, promKey string) {
	metric := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:        fmt.Sprintf("depcarto_%v_request_duration_seconds", promKey),
			Help:        fmt.Sprintf("Histogram concerning %v request duration (seconds)", promKey),
			Buckets:     []float64{.0025, .005, .01, .025, .05, .1},
			ConstLabels: prometheus.Labels{"method": h.Method(), "path": h.Path()},
		},
		[]string{},
	)
	prometheus.Unregister(metric)
	prometheus.MustRegister(metric)
	h2 := promhttp.InstrumentHandlerDuration(metric, h)
	r.HandleFunc(h.Path(), h2).Methods(h.Method())
}

type handlerInterface interface {
	Path() string
	Method() string
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}

func (s *Server) Run() error {
	addr := fmt.Sprintf("0.0.0.0:%v", s.port)
	srv := &http.Server{
		Addr:         addr,
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      s.router,
	}
	log.Info().Msgf("Server running (port %v)...", s.port)
	return srv.ListenAndServe()
}

func formatter(r *http.Request) (output.Formatter, error) {
	v := defaultFormat
	if f, found := r.URL.Query()["format"]; found {
		v = strings.ToLower(f[0])
	}

	switch v {
	case jsonFormat:
		return output.NewJSONFormatter(), nil
	case dotFormat:
		return output.NewDotFormatter(), nil
	case jpgFormat:
		return output.NewJpgFormatter(), nil
	default:
		return nil, fmt.Errorf("unsupported output format (%v)", v)
	}
}

func since(r *http.Request) (time.Duration, error) {
	if f, found := r.URL.Query()["since"]; found {
		return time.ParseDuration(f[0])
	}
	return defaultSince, nil
}

func depth(r *http.Request) (int, error) {
	if f, found := r.URL.Query()["depth"]; found {
		return strconv.Atoi(f[0])
	}
	return defaultDepth, nil
}
