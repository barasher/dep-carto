package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/barasher/dep-carto/internal/model"
	"github.com/rs/zerolog/log"

	"github.com/gorilla/mux"
)

type getDependenciesHandler struct {
	model model.Model
}

func NewGetDependenciesHandler(m model.Model) getDependenciesHandler {
	return getDependenciesHandler{m}
}

func (h getDependenciesHandler) Path() string {
	return "/dependencies/{ident}"
}

func (h getDependenciesHandler) Method() string {
	return http.MethodGet
}

func (h getDependenciesHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	ident := mux.Vars(r)["ident"]

	since, err := since(r)
	if err != nil {
		err = fmt.Errorf("error while parsing duration %v: %w", since, err)
		log.Error().Msgf("%v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	depth, err := depth(r)
	if err != nil {
		err = fmt.Errorf("error while depth %v: %w", depth, err)
		log.Error().Msgf("%v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	s, err := h.model.GetDependencies(ctx, ident, depth, since)
	if err != nil {
		log.Error().Msgf("%v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	f, err := formatter(r)
	if err != nil {
		log.Error().Msgf("%v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	f.Format(s, w)
}
