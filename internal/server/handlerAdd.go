package server

import (
	"context"
	"encoding/json"
	"github.com/barasher/dep-carto/internal/model"
	"github.com/rs/zerolog/log"
	"net/http"
)

type addHandler struct {
	model model.Model
}

func NewAddHandler(m model.Model) addHandler {
	return addHandler{m}
}

func (h addHandler) Path() string {
	return "/server"
}

func (h addHandler) Method() string {
	return http.MethodPost
}

func (h addHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var s model.Server
	err := json.NewDecoder(r.Body).Decode(&s)
	if err != nil {
		log.Error().Msgf("Error while unmarshalling server: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := h.model.Add(context.Background(), s); err != nil {
		log.Error().Msgf("Error while adding server: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusNoContent)
}
