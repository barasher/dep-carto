package model

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func getDay(t *testing.T, d string) time.Time {
	parsed, err := time.Parse(time.RFC3339, "2020-03-"+d+"T01:01:01Z")
	assert.Nil(t, err)
	return parsed
}

func testModelWorkflow(t *testing.T, m Model) {
	// create first server (s1 - no key)
	s1 := Server{
		Hostname:     "h1",
		IPs:          []string{"h1ip1", "h1ip2"},
		Dependencies: []string{"d1","d2"},
		LastUpdate:   getDay(t, "01"),
	}
	assert.Nil(t, m.AddServer(s1))
	servers, err := m.GetAllServers()
	assert.Nil(t, err)
	assert.ElementsMatch(t, []Server{s1}, servers)

	// override server s1 - no key
	s1 = Server{
		Hostname:     "h1",
		IPs:          []string{"h1ip1", "h1ip2"},
		Dependencies: []string{"d1"},
		LastUpdate:   getDay(t, "02"),
	}
	assert.Nil(t, m.AddServer(s1))
	servers, err = m.GetAllServers()
	assert.Nil(t, err)
	assert.ElementsMatch(t, []Server{s1}, servers)

	// create new server (s2 - no key)
	s2 := Server{
		Hostname:     "h2",
		IPs:          []string{"h2ip1"},
		Dependencies: []string{"d3"},
		LastUpdate:   getDay(t, "03"),
	}
	assert.Nil(t, m.AddServer(s2))
	servers, err = m.GetAllServers()
	assert.Nil(t, err)
	assert.ElementsMatch(t, []Server{s1, s2}, servers)

	// create new server (s1 - key)
	s1b := Server{
		Hostname:     "h1",
		Key: "k",
		IPs:          []string{"h1ip1", "h1ip2"},
		Dependencies: []string{"d1","d2"},
		LastUpdate:   getDay(t, "04"),
	}
	assert.Nil(t, m.AddServer(s1b))
	servers, err = m.GetAllServers()
	assert.Nil(t, err)
	assert.ElementsMatch(t, []Server{s1, s2, s1b}, servers)
}
