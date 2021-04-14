package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/barasher/dep-carto/internal/model"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestGetHandler_All_Nominal(t *testing.T) {
	s := []model.Server{model.Server{Hostname: "h"}}
	m := (&modelMock{}).MockGetAll(s, nil)
	h := NewGetHandler(m)
	req, err := http.NewRequest(h.Method(), h.Path(), nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.Handle(h.Path(), h)
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	var got []model.Server
	assert.Nil(t, json.Unmarshal(rr.Body.Bytes(), &got))
	assert.Equal(t, s, got)
}

func TestGetHandler_All_Error(t *testing.T) {
	m := (&modelMock{}).MockGetAll(nil, fmt.Errorf("err"))
	h := NewGetHandler(m)
	req, err := http.NewRequest(h.Method(), h.Path(), nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.Handle(h.Path(), h)
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestGetHandler_Since_Nominal(t *testing.T) {
	s := []model.Server{model.Server{Hostname: "h"}}
	m := (&modelMock{}).MockGetAll(s, nil)
	h := NewGetHandler(m)
	req, err := http.NewRequest(h.Method(), h.Path()+"?since=3600s", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.Handle(h.Path(), h)
	router.ServeHTTP(rr, req)

	assert.Equal(t, 3600*time.Second, m.getAll.inSince)
	assert.Equal(t, http.StatusOK, rr.Code)
	var got []model.Server
	assert.Nil(t, json.Unmarshal(rr.Body.Bytes(), &got))
	assert.Equal(t, s, got)
}

func TestGetHandler_Since_Error(t *testing.T) {
	m := (&modelMock{}).MockGetAll(nil, fmt.Errorf("err"))
	h := NewGetHandler(m)
	req, err := http.NewRequest(h.Method(), h.Path()+"?since=3600s", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.Handle(h.Path(), h)
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestGetHandler_Since_DurationParseError(t *testing.T) {
	m := (&modelMock{}).MockGetAll(nil, fmt.Errorf("err"))
	h := NewGetHandler(m)
	req, err := http.NewRequest(h.Method(), h.Path()+"?since=blabla", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.Handle(h.Path(), h)
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}
