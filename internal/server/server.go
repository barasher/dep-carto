package server

import (
	"fmt"
	"github.com/barasher/dep-carto/internal/model"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
	"net/http"
	"time"
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
	registerHandler(s.router, NewAddHandler(model))
	registerHandler(s.router, NewClearHandler(model))
	registerHandler(s.router, NewGetHandler(model))
	return s, nil
}

func registerHandler(r *mux.Router, h handlerInterface) {
	r.Handle(h.Path(), h).Methods(h.Method())
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
