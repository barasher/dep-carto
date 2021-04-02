package output

import (
	"encoding/json"
	"fmt"
	"github.com/barasher/dep-carto/internal/model"
	"github.com/rs/zerolog/log"
	"net/http"
)

type JSONFormatter struct{}

func NewJSONFormatter() JSONFormatter {
	return JSONFormatter{}
}

func (JSONFormatter) Format(s []model.Server, w http.ResponseWriter) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(s)
	if err != nil {
		err := fmt.Errorf("error while encoding JSON: %w", err)
		log.Error().Msgf("%v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
