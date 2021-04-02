package output

import (
	"github.com/barasher/dep-carto/internal/model"
	"net/http"
)

type Formatter interface {
	Format(s []model.Server, w http.ResponseWriter)
}
