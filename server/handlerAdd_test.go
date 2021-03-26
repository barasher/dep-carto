package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/barasher/dep-carto/internal/model"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestAddHandler_Nominal(t *testing.T) {
	d, err := time.Parse(time.RFC3339, "2001-02-03T04:05:06Z")
	s := model.Server{
		Hostname:     "hh",
		Key:          "k",
		IPs:          []string{"i1", "i2"},
		Dependencies: []string{"d1", "d2"},
		LastUpdate:   d,
	}
	b, err := json.Marshal(s)
	assert.Nil(t, err)

	m := &modelMock{}
	h := NewAddHandler(m)
	req, err := http.NewRequest(h.Method(), h.Path(), bytes.NewReader(b))
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.Handle(h.Path(), h)
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNoContent, rr.Code)
	assert.Equal(t, s, m.add.inServer)
}

func TestAddHandler_UnmarshalError(t *testing.T) {
	m := &modelMock{}
	h := NewAddHandler(m)
	req, err := http.NewRequest(h.Method(), h.Path(), bytes.NewReader([]byte("{")))
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.Handle(h.Path(), h)
	router.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestAddHandler_Error(t *testing.T) {
	m := (&modelMock{}).MockAdd(fmt.Errorf("err"))
	h := NewAddHandler(m)
	req, err := http.NewRequest(h.Method(), h.Path(), bytes.NewReader([]byte("{}")))
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.Handle(h.Path(), h)
	router.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}
