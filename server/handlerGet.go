package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/barasher/dep-carto/internal/model"
	"github.com/barasher/dep-carto/internal/output"
	"github.com/rs/zerolog/log"
)

const (
	jsonFormat    = "json"
	dotFormat     = "dot"
	jpgFormat     = "jpg"
	defaultFormat = jsonFormat
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

func (getHandler) format(r *http.Request) string {
	if f, found := r.URL.Query()["format"]; found {
		return f[0]
	}
	return defaultFormat
}

func (h getHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	var s []model.Server
	var err error
	since, found := r.URL.Query()["since"]
	if found { // since
		since, err := time.ParseDuration(since[0])
		if err != nil {
			err = fmt.Errorf("error while parsing duration %v: %w", since, err)
			log.Error().Msgf("%v", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if s, err = h.model.GetSince(ctx, since); err != nil {
			err = fmt.Errorf("error while getting servers since %v: %w", since, err)
			log.Error().Msgf("%v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else { // all
		if s, err = h.model.GetAll(ctx); err != nil {
			err = fmt.Errorf("error while getting all servers: %w", err)
			log.Error().Msgf("%v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	var f output.Formatter
	switch fStr := h.format(r); fStr {
	case jsonFormat:
		f = output.NewJSONFormatter()
	case dotFormat:
		f = output.NewDotFormatter()
	case jpgFormat:
		f = output.NewJpgFormatter()
	default:
		err = fmt.Errorf("unsupported output format (%v): %w", fStr, err)
		log.Error().Msgf("%v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	f.Format(s, w)

}

func (h getHandler) getAll(ctx context.Context, w http.ResponseWriter) {
	s, err := h.model.GetAll(ctx)
	if err != nil {
		log.Error().Msgf("Error while getting all servers: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(s)
}

func (h getHandler) getSince(ctx context.Context, w http.ResponseWriter, dur string) {
	since, err := time.ParseDuration(dur)
	if err != nil {
		log.Error().Msgf("Error while parsing duration (%v): %v", since, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	s, err := h.model.GetSince(ctx, since)
	if err != nil {
		log.Error().Msgf("Error while getting servers since %v: %v", dur, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	// format

	json.NewEncoder(w).Encode(s)
	w.WriteHeader(http.StatusOK)
}
