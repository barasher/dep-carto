package output

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDotToJpg_Nominal(t *testing.T) {
	bIn := bytes.Buffer{}
	bOut := bytes.Buffer{}
	_, err := bIn.WriteString("digraph G { Hello -> World }")
	assert.Nil(t, err)
	assert.Nil(t, dotToJpg(&bIn, &bOut))
	assert.True(t, len(bOut.Bytes()) > 0)
}

func TestDotToJpg_Error(t *testing.T) {
	bIn := bytes.Buffer{}
	bOut := bytes.Buffer{}
	_, err := bIn.WriteString("digraph G { ")
	assert.Nil(t, err)
	assert.NotNil(t, dotToJpg(&bIn, &bOut))
}

func TestJpgFormatter(t *testing.T) {
	rec := httptest.NewRecorder()
	NewJpgFormatter().Format(generateServers(), rec)
	contentType := rec.Header().Get("Content-Type")
	assert.Truef(t, strings.HasPrefix(contentType, "image/jpeg"), "got: %v", contentType)
	assert.Equal(t, http.StatusOK, rec.Code)
}
