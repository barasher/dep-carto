package server

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/barasher/dep-carto/internal/model"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type modelMock struct {
	add struct {
		inServer model.Server
		outErr   error
	}
	getAll struct {
		servers []model.Server
		err     error
	}
	getSince struct {
		inDuration time.Duration
		outServers []model.Server
		outErr     error
	}
	clear struct {
		err error
	}
}

func (m *modelMock) MockAdd(err error) *modelMock {
	m.add.outErr = err
	return m
}

func (m *modelMock) Add(ctx context.Context, s model.Server) error {
	m.add.inServer = s
	return m.add.outErr
}

func (m *modelMock) MockGetAll(s []model.Server, err error) *modelMock {
	m.getAll.servers = s
	m.getAll.err = err
	return m
}

func (m *modelMock) GetAll(ctx context.Context) ([]model.Server, error) {
	return m.getAll.servers, m.getAll.err
}

func (m *modelMock) MockGetSince(s []model.Server, err error) *modelMock {
	m.getSince.outServers = s
	m.getSince.outErr = err
	return m
}

func (m *modelMock) GetSince(ctx context.Context, d time.Duration) ([]model.Server, error) {
	m.getSince.inDuration = d
	return m.getSince.outServers, m.getSince.outErr
}

func (m *modelMock) MockClear(err error) *modelMock {
	m.clear.err = err
	return m
}

func (m *modelMock) Clear(ctx context.Context) error {
	return m.clear.err
}

func TestServer(t *testing.T) {
	m := model.NewMemoryModel()
	s, err := NewServer(m, 8080)
	assert.Nil(t, err)
	h := httptest.NewServer(s.router)
	defer h.Close()

	// create
	elt := model.Server{Hostname: "host"}
	b, err := json.Marshal(elt)
	assert.Nil(t, err)
	u := fmt.Sprintf("%v/server", h.URL)
	t.Logf("url create: %v", u)
	r, err := http.Post(u, "application/json", bytes.NewReader(b))
	assert.Nil(t, err)
	assert.Equal(t, http.StatusNoContent, r.StatusCode)

	// get
	u = fmt.Sprintf("%v/server", h.URL)
	r, err = http.Get(u)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, r.StatusCode)
	var servers []model.Server
	assert.Nil(t, json.NewDecoder(r.Body).Decode(&servers))
	defer r.Body.Close()
	assert.Contains(t, servers, elt)

	// clear
	u = fmt.Sprintf("%v/server", h.URL)
	req, err := http.NewRequest(http.MethodDelete, u, nil)
	assert.Nil(t, err)
	c := http.Client{Timeout: time.Second}
	r, err = c.Do(req)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusNoContent, r.StatusCode)

	// get (expected empty)
	u = fmt.Sprintf("%v/server", h.URL)
	r, err = http.Get(u)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, r.StatusCode)
	assert.Nil(t, json.NewDecoder(r.Body).Decode(&servers))
	defer r.Body.Close()
	assert.Len(t, servers, 0)
}
