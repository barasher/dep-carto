package model

import (
	"context"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
	"time"
)

func TestJSON(t *testing.T) {
	raw := `
		{
		  "hostname": "h",
		  "key": "k",
		  "ips": ["i1", "i2"],
		  "dependencies": ["d1", "d2"],
		  "lastUpdate": "2001-02-03T04:05:06Z"
		}
		`
	var s Server
	assert.Nil(t, json.NewDecoder(strings.NewReader(raw)).Decode(&s))
	assert.Equal(t, "h", s.Hostname)
	assert.Equal(t, "k", s.Key)
	assert.ElementsMatch(t, []string{"i1", "i2"}, s.IPs)
	assert.ElementsMatch(t, []string{"d1", "d2"}, s.Dependencies)
	expDate, err := time.Parse(time.RFC3339, "2001-02-03T04:05:06Z")
	assert.Nil(t, err)
	assert.Equal(t, expDate, s.LastUpdate)
}

func TestMarshalingUnmarshaling(t *testing.T) {
	s1 := Server{
		Hostname:     "hh",
		Key:          "k",
		IPs:          []string{"i1", "i2"},
		Dependencies: []string{"d1", "d2"},
		LastUpdate:   getDay(t, "01"),
	}
	b, err := json.Marshal(s1)
	assert.Nil(t, err)
	var s2 Server
	assert.Nil(t, json.Unmarshal(b, &s2))
	assert.Equal(t, s1, s2)
}

func getDay(t *testing.T, d string) time.Time {
	parsed, err := time.Parse(time.RFC3339, "2020-03-"+d+"T01:01:01Z")
	assert.Nil(t, err)
	return parsed
}

func testModelWorkflow(t *testing.T, m Model) {
	testCreate(t, m)
	testDelete(t, m)
	testSince(t, m)
}

func testCreate(t *testing.T, m Model) {
	// create first server (s1 - no key)
	s1 := Server{
		Hostname:     "h1",
		IPs:          []string{"h1ip1", "h1ip2"},
		Dependencies: []string{"d1", "d2"},
		LastUpdate:   getDay(t, "01"),
	}
	assert.Nil(t, m.Add(context.Background(), s1))
	servers, err := m.GetAll(context.Background())
	assert.Nil(t, err)
	assert.ElementsMatch(t, []Server{s1}, servers)

	// override server s1 - no key
	s1 = Server{
		Hostname:     "h1",
		IPs:          []string{"h1ip1", "h1ip2"},
		Dependencies: []string{"d1"},
		LastUpdate:   getDay(t, "02"),
	}
	assert.Nil(t, m.Add(context.Background(), s1))
	servers, err = m.GetAll(context.Background())
	assert.Nil(t, err)
	assert.ElementsMatch(t, []Server{s1}, servers)

	// create new server (s2 - no key)
	s2 := Server{
		Hostname:     "h2",
		IPs:          []string{"h2ip1"},
		Dependencies: []string{"d3"},
		LastUpdate:   getDay(t, "03"),
	}
	assert.Nil(t, m.Add(context.Background(), s2))
	servers, err = m.GetAll(context.Background())
	assert.Nil(t, err)
	assert.ElementsMatch(t, []Server{s1, s2}, servers)

	// create new server (s1 - key)
	s1b := Server{
		Hostname:     "h1",
		Key:          "k",
		IPs:          []string{"h1ip1", "h1ip2"},
		Dependencies: []string{"d1", "d2"},
		LastUpdate:   getDay(t, "04"),
	}
	assert.Nil(t, m.Add(context.Background(), s1b))
	servers, err = m.GetAll(context.Background())
	assert.Nil(t, err)
	assert.ElementsMatch(t, []Server{s1, s2, s1b}, servers)
}

func testDelete(t *testing.T, m Model) {
	// create
	s1 := Server{
		Hostname:     "h1",
		Key:          "bla",
		IPs:          []string{"h1ip1", "h1ip2"},
		Dependencies: []string{"d1", "d2"},
		LastUpdate:   getDay(t, "01"),
	}
	assert.Nil(t, m.Add(context.Background(), s1))
	servers, err := m.GetAll(context.Background())
	assert.Nil(t, err)
	assert.NotEmpty(t, servers)

	// delete
	assert.Nil(t, m.Clear(context.Background()))

	// check
	servers, err = m.GetAll(context.Background())
	assert.Nil(t, err)
	assert.Empty(t, servers)
}

func testSince(t *testing.T, m Model) {
	assert.Nil(t, m.Clear(context.Background()))

	// create (-1d)
	s := Server{
		Hostname:   "h1",
		LastUpdate: time.Now().Add(-24 * time.Hour),
	}
	assert.Nil(t, m.Add(context.Background(), s))
	servers, err := m.GetAll(context.Background())
	assert.Nil(t, err)
	assert.Len(t, servers, 1)

	// check -2d
	servers, err = m.GetSince(context.Background(), 48*time.Hour)
	assert.Nil(t, err)
	assert.Len(t, servers, 1)

	// check -1h
	servers, err = m.GetSince(context.Background(), time.Hour)
	assert.Nil(t, err)
	assert.Len(t, servers, 0)
}
