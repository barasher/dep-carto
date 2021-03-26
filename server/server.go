package server

import (
	"github.com/barasher/dep-carto/internal/model"
	"github.com/gorilla/mux"
)

type Server struct {
	router *mux.Router
	model  model.Model
	port   uint
}

func NewServer(model model.Model, port uint) (*Server, error) {
	s := &Server{}
	s.router = mux.NewRouter()
	s.model = model
	return s, nil
}
