package server

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/barasher/dep-carto/internal/model"
	"github.com/stretchr/testify/assert"
)

type modelMock struct {
	add struct {
		inServer model.Server
		outErr   error
	}
	getAll struct {
		inSince time.Duration
		servers []model.Server
		err     error
	}
	clear struct {
		err error
	}
	getDepending struct {
		inIdent string
		inDepth int
		inSince time.Duration
		servers []model.Server
		err     error
	}
	getDependencies struct {
		inIdent string
		inDepth int
		inSince time.Duration
		servers []model.Server
		err     error
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

func (m *modelMock) GetAll(ctx context.Context, d time.Duration) ([]model.Server, error) {
	m.getAll.inSince = d
	return m.getAll.servers, m.getAll.err
}

func (m *modelMock) MockClear(err error) *modelMock {
	m.clear.err = err
	return m
}

func (m *modelMock) Clear(ctx context.Context) error {
	return m.clear.err
}

func (m *modelMock) MockGetDepending(s []model.Server, err error) *modelMock {
	m.getDepending.servers = s
	m.getDepending.err = err
	return m
}

func (m *modelMock) GetDepending(ctx context.Context, ident string, depth int, since time.Duration) ([]model.Server, error) {
	m.getDepending.inIdent = ident
	m.getDepending.inDepth = depth
	m.getDepending.inSince = since
	return m.getDepending.servers, m.getDepending.err
}

func (m *modelMock) MockGetDependencies(s []model.Server, err error) *modelMock {
	m.getDependencies.servers = s
	m.getDependencies.err = err
	return m
}

func (m *modelMock) GetDependencies(ctx context.Context, ident string, depth int, since time.Duration) ([]model.Server, error) {
	m.getDependencies.inIdent = ident
	m.getDependencies.inDepth = depth
	m.getDependencies.inSince = since
	return m.getDependencies.servers, m.getDependencies.err
}

func TestServer(t *testing.T) {
	m := model.NewMemoryModel()
	s, err := NewServer(m, 8080)
	assert.Nil(t, err)
	h := httptest.NewServer(s.router)
	defer h.Close()

	// create
	elt := model.Server{Hostname: "host", LastUpdate: time.Now()}
	b, err := json.Marshal(elt)
	assert.Nil(t, err)
	u := fmt.Sprintf("%v/server", h.URL)
	t.Logf("url create: %v", u)
	r, err := http.Post(u, "application/json", bytes.NewReader(b))
	assert.Nil(t, err)
	assert.Equal(t, http.StatusNoContent, r.StatusCode)

	// get
	u = fmt.Sprintf("%v/servers", h.URL)
	r, err = http.Get(u)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, r.StatusCode)
	var servers []model.Server
	assert.Nil(t, json.NewDecoder(r.Body).Decode(&servers))
	defer r.Body.Close()
	assert.Len(t, servers, 1)

	// clear
	u = fmt.Sprintf("%v/servers", h.URL)
	req, err := http.NewRequest(http.MethodDelete, u, nil)
	assert.Nil(t, err)
	c := http.Client{Timeout: time.Second}
	r, err = c.Do(req)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusNoContent, r.StatusCode)

	// get (expected empty)
	u = fmt.Sprintf("%v/servers", h.URL)
	r, err = http.Get(u)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, r.StatusCode)
	assert.Nil(t, json.NewDecoder(r.Body).Decode(&servers))
	defer r.Body.Close()
	assert.Len(t, servers, 0)
}

func TestFormat(t *testing.T) {
	var tcs = []struct {
		inTC     string
		inSuffix string
		expOk    bool
		expType  string
	}{
		{"notProvided", "", true, "output.JSONFormatter"},
		{"JPG", "format=JPG", true, "output.JpgFormatter"},
		{"JSON", "format=JSON", true, "output.JSONFormatter"},
		{"DOT", "format=DOT", true, "output.DotFormatter"},
		{"blabla", "format=blabla", false, ""},
	}
	for _, tc := range tcs {
		t.Run(tc.inTC, func(t *testing.T) {
			u := fmt.Sprintf("http://www.google.com/?%v", tc.inSuffix)
			r, err := http.NewRequest(http.MethodGet, u, nil)
			assert.Nil(t, err)
			f, err := formatter(r)
			if tc.expOk {
				assert.Nil(t, err)
				assert.Equal(t, tc.expType, fmt.Sprintf("%T", f))
			} else {
				assert.NotNil(t, err)
			}
		})
	}
}

func TestSince(t *testing.T) {
	var tcs = []struct {
		inTC     string
		inSuffix string
		expOk    bool
		expOut   time.Duration
	}{
		{"notProvided", "", true, defaultSince},
		{"1m", "since=1m", true, time.Minute},
		{"blabla", "since=blabla", false, time.Minute},
	}
	for _, tc := range tcs {
		t.Run(tc.inTC, func(t *testing.T) {
			u := fmt.Sprintf("http://www.google.com/?%v", tc.inSuffix)
			r, err := http.NewRequest(http.MethodGet, u, nil)
			s, err := since(r)
			if tc.expOk {
				assert.Nil(t, err)
				assert.Equal(t, tc.expOut, s)
			} else {
				assert.NotNil(t, err)
			}
		})
	}
}

func TestDepth(t *testing.T) {
	var tcs = []struct {
		inTC     string
		inSuffix string
		expOk    bool
		expOut   int
	}{
		{"notProvided", "", true, defaultDepth},
		{"42", "depth=42", true, 42},
		{"blabla", "depth=blabla", false, 0},
	}
	for _, tc := range tcs {
		t.Run(tc.inTC, func(t *testing.T) {
			u := fmt.Sprintf("http://www.google.com/?%v", tc.inSuffix)
			r, err := http.NewRequest(http.MethodGet, u, nil)
			s, err := depth(r)
			if tc.expOk {
				assert.Nil(t, err)
				assert.Equal(t, tc.expOut, s)
			} else {
				assert.NotNil(t, err)
			}
		})
	}
}
