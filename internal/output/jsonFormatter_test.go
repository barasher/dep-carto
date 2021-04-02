package output

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestJsonFormatter(t *testing.T) {
	rec := httptest.NewRecorder()
	NewJSONFormatter().Format(generateServers(), rec)
	assert.Equal(t, http.StatusOK, rec.Code)
	contentType := rec.Header().Get("Content-Type")
	assert.Truef(t, strings.HasPrefix(contentType, "application/json"), "got: %v", contentType)
}
