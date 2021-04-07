package server

import (
	"context"
	"net/http"

	"github.com/barasher/dep-carto/internal/model"
	"github.com/rs/zerolog/log"
)

type clearHandler struct {
	model model.Model
}

func NewClearHandler(m model.Model) clearHandler {
	return clearHandler{m}
}

func (h clearHandler) Path() string {
	return "/servers"
}

func (h clearHandler) Method() string {
	return http.MethodDelete
}

func (h clearHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := h.model.Clear(context.Background()); err != nil {
		log.Error().Msgf("Error while clearing servers: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusNoContent)
}
