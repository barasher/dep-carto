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

func TestGetHandler_Nominal(t *testing.T) {
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

func TestGetHandler_StatusCode(t *testing.T) {
	sOk := []model.Server{model.Server{Hostname: "h"}}
	mOk := (&modelMock{}).MockGetAll(sOk, nil)
	mKo := (&modelMock{}).MockGetAll([]model.Server{}, fmt.Errorf("err"))

	var tcs = []struct {
		inTC          string
		inMock        *modelMock
		inUrlParams   string
		expStatusCode int
	}{
		{"nominal", mOk, "since=3600s", http.StatusOK},
		{"noParam", mOk, "", http.StatusOK},
		{"processError", mKo, "since=3600s", http.StatusInternalServerError},
		{"sinceError", mOk, "since=blabla", http.StatusBadRequest},
	}
	for _, tc := range tcs {
		t.Run(tc.inTC, func(t *testing.T) {
			h := NewGetHandler(tc.inMock)
			req, err := http.NewRequest(h.Method(), "/servers?"+tc.inUrlParams, nil)
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
