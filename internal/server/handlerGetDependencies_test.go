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

func TestGetDependenciesHandler_Nominal(t *testing.T) {
	s := []model.Server{model.Server{Hostname: "h"}}
	m := (&modelMock{}).MockGetDependencies(s, nil)
	h := NewGetDependenciesHandler(m)
	req, err := http.NewRequest(h.Method(), "/dependencies/blabla?since=3600s&depth=2&format=json", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.Handle(h.Path(), h)
	router.ServeHTTP(rr, req)

	assert.Equal(t, 3600*time.Second, m.getDependencies.inSince)
	assert.Equal(t, 2, m.getDependencies.inDepth)
	assert.Equal(t, http.StatusOK, rr.Code)
	var got []model.Server
	assert.Nil(t, json.Unmarshal(rr.Body.Bytes(), &got))
	assert.Equal(t, s, got)
}

func TestGetDependenciesHandler_StatusCode(t *testing.T) {
	sOk := []model.Server{model.Server{Hostname: "h"}}
	mOk := (&modelMock{}).MockGetDependencies(sOk, nil)
	mKo := (&modelMock{}).MockGetDependencies([]model.Server{}, fmt.Errorf("err"))

	var tcs = []struct {
		inTC          string
		inMock        *modelMock
		inUrlParams   string
		expStatusCode int
	}{
		{"nominal", mOk, "since=3600s&depth=2&format=json", http.StatusOK},
		{"processError", mKo, "since=3600s&depth=2&format=json", http.StatusInternalServerError},
		{"sinceError", mOk, "since=blabla&depth=2&format=json", http.StatusBadRequest},
		{"depthError", mOk, "since=3600s&depth=blabla&format=json", http.StatusBadRequest},
		{"formatError", mOk, "since=3600s&depth=2&format=blabla", http.StatusBadRequest},
	}
	for _, tc := range tcs {
		t.Run(tc.inTC, func(t *testing.T) {
			h := NewGetDependenciesHandler(tc.inMock)
			req, err := http.NewRequest(h.Method(), "/dependencies/blabla?"+tc.inUrlParams, nil)
			if err != nil {
				t.Fatal(err)
			}
			rr := httptest.NewRecorder()
			router := mux.NewRouter()
			router.Handle(h.Path(), h)
			router.ServeHTTP(rr, req)
			assert.Equal(t, tc.expStatusCode, rr.Code)
		})
	}

}
