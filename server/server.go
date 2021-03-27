package server

import (
	"github.com/barasher/dep-carto/internal/model"
	"github.com/gorilla/mux"
	"net/http"
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
