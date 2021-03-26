package server

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestClearHandler_Nominal(t *testing.T) {
	m := &modelMock{}
	h := NewClearHandler(m)
	req, err := http.NewRequest(h.Method(), h.Path(), nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.Handle(h.Path(), h)
	router.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusNoContent, rr.Code)
}

func TestClearHandler_Error(t *testing.T) {
	m := (&modelMock{}).MockClear(fmt.Errorf("err"))
	h := NewClearHandler(m)
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
