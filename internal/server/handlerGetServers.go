package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/barasher/dep-carto/internal/model"
	"github.com/rs/zerolog/log"
)

type getHandler struct {
	model model.Model
}

func NewGetHandler(m model.Model) getHandler {
	return getHandler{m}
}

func (h getHandler) Path() string {
	return "/servers"
}

func (h getHandler) Method() string {
	return http.MethodGet
}

func (h getHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	since, err := since(r)
	if err != nil {
		err = fmt.Errorf("error while parsing duration %v: %w", since, err)
		log.Error().Msgf("%v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	s, err := h.model.GetAll(ctx, since)
	if err != nil {
		err = fmt.Errorf("error while getting servers since %v: %w", since, err)
		log.Error().Msgf("%v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	f, err := formatter(r)
	if err != nil {
		log.Error().Msgf("%v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	f.Format(s, w)
}
