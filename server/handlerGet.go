package server

import (
	"encoding/json"
	"github.com/barasher/dep-carto/internal/model"
	"github.com/rs/zerolog/log"
	"net/http"
	"time"
)

type getHandler struct {
	model model.Model
}

func NewGetHandler(m model.Model) getHandler {
	return getHandler{m}
}

func (h getHandler) Path() string {
	return "/server"
}

func (h getHandler) Method() string {
	return http.MethodGet
}

func (h getHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	since, found := r.URL.Query()["since"]
	if found {
		h.getSince(w, since[0])
	} else {
		h.getAll(w)
	}
}

func (h getHandler) getAll(w http.ResponseWriter) {
	s, err := h.model.GetAll()
	if err != nil {
		log.Error().Msgf("Error while getting all servers: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	json.NewEncoder(w).Encode(s)
	w.WriteHeader(http.StatusOK)
}

func (h getHandler) getSince(w http.ResponseWriter, dur string) {
	since, err := time.ParseDuration(dur)
	if err != nil {
		log.Error().Msgf("Error while parsing duration (%v): %v", since, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	s, err := h.model.GetSince(since)
	if err != nil {
		log.Error().Msgf("Error while getting servers since %v: %v", dur, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	json.NewEncoder(w).Encode(s)
	w.WriteHeader(http.StatusOK)
}
