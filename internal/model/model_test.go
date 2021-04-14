package model

import (
	"context"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestJSON(t *testing.T) {
	raw := `
		{
		  "hostname": "h",
		  "key": "k",
		  "ips": ["i1", "i2"],
		  "dependencies": [
			{ "resource": "d1", "label": "l1" },
			{ "resource": "d2", "label": "l2" }
		  ],
		  "lastUpdate": "2001-02-03T04:05:06Z"
		}
		`
	var s Server
	assert.Nil(t, json.NewDecoder(strings.NewReader(raw)).Decode(&s))
	assert.Equal(t, "h", s.Hostname)
	assert.Equal(t, "k", s.Key)
	assert.ElementsMatch(t, []string{"i1", "i2"}, s.IPs)
	expDeps := []Dependency{
		{Resource: "d1", Label: "l1"},
		{Resource: "d2", Label: "l2"},
	}
	assert.ElementsMatch(t, expDeps, s.Dependencies)
	expDate, err := time.Parse(time.RFC3339, "2001-02-03T04:05:06Z")
	assert.Nil(t, err)
	assert.Equal(t, expDate, s.LastUpdate)
}

func TestMarshalingUnmarshaling(t *testing.T) {
	s1 := Server{
		Hostname: "hh",
		Key:      "k",
		IPs:      []string{"i1", "i2"},
		Dependencies: []Dependency{
			{Resource: "d1", Label: "l1"},
			{Resource: "d2", Label: "l2"},
		},
		LastUpdate: getDay(t, "01"),
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
	testGetDependencies(t, m)
}

func testCreate(t *testing.T, m Model) {
	old := 24 * 365 * 10 * time.Hour
	// create first server (s1 - no key)
	s1 := Server{
		Hostname: "h1",
		IPs:      []string{"h1ip1", "h1ip2"},
		Dependencies: []Dependency{
			{Resource: "d1", Label: "l1"},
			{Resource: "d2", Label: "l2"},
		},
		LastUpdate: getDay(t, "01"),
	}
	assert.Nil(t, m.Add(context.Background(), s1))
	servers, err := m.GetAll(context.Background(), old)
	assert.Nil(t, err)
	assert.ElementsMatch(t, []Server{s1}, servers)

	// override server s1 - no key
	s1 = Server{
		Hostname: "h1",
		IPs:      []string{"h1ip1", "h1ip2"},
		Dependencies: []Dependency{
			{Resource: "d1", Label: "l1"},
		},
		LastUpdate: getDay(t, "02"),
	}
	assert.Nil(t, m.Add(context.Background(), s1))
	servers, err = m.GetAll(context.Background(), old)
	assert.Nil(t, err)
	assert.ElementsMatch(t, []Server{s1}, servers)

	// create new server (s2 - no key)
	s2 := Server{
		Hostname: "h2",
		IPs:      []string{"h2ip1"},
		Dependencies: []Dependency{
			{Resource: "d3", Label: "l3"},
		},
		LastUpdate: getDay(t, "03"),
	}
	assert.Nil(t, m.Add(context.Background(), s2))
	servers, err = m.GetAll(context.Background(), old)
	assert.Nil(t, err)
	assert.ElementsMatch(t, []Server{s1, s2}, servers)

	// create new server (s1 - key)
	s1b := Server{
		Hostname: "h1",
		Key:      "k",
		IPs:      []string{"h1ip1", "h1ip2"},
		Dependencies: []Dependency{
			{Resource: "d1", Label: "l1"},
			{Resource: "d2", Label: "l2"},
		},
		LastUpdate: getDay(t, "04"),
	}
	assert.Nil(t, m.Add(context.Background(), s1b))
	servers, err = m.GetAll(context.Background(), old)
	assert.Nil(t, err)
	assert.ElementsMatch(t, []Server{s1, s2, s1b}, servers)
}

func testDelete(t *testing.T, m Model) {
	old := 24 * 365 * 10 * time.Hour
	// create
	s1 := Server{
		Hostname: "h1",
		Key:      "bla",
		IPs:      []string{"h1ip1", "h1ip2"},
		Dependencies: []Dependency{
			{Resource: "d1", Label: "l1"},
			{Resource: "d2", Label: "l2"},
		},
		LastUpdate: getDay(t, "01"),
	}
	assert.Nil(t, m.Add(context.Background(), s1))
	servers, err := m.GetAll(context.Background(), old)
	assert.Nil(t, err)
	assert.NotEmpty(t, servers)

	// delete
	assert.Nil(t, m.Clear(context.Background()))

	// check
	servers, err = m.GetAll(context.Background(), old)
	assert.Nil(t, err)
	assert.Empty(t, servers)
}

func testSince(t *testing.T, m Model) {
	old := 24 * 365 * 10 * time.Hour
	assert.Nil(t, m.Clear(context.Background()))

	// create (-1d)
	s := Server{
		Hostname:   "h1",
		LastUpdate: time.Now().Add(-24 * time.Hour),
	}
	assert.Nil(t, m.Add(context.Background(), s))
	servers, err := m.GetAll(context.Background(), old)
	assert.Nil(t, err)
	assert.Len(t, servers, 1)

	// check -2d
	d := 48 * time.Hour
	servers, err = m.GetAll(context.Background(), d)
	assert.Nil(t, err)
	assert.Len(t, servers, 1)

	// check -1h
	d = time.Hour
	servers, err = m.GetAll(context.Background(), d)
	assert.Nil(t, err)
	assert.Len(t, servers, 0)
}

func testGetDependencies(t *testing.T, m Model) {
	ctx := context.Background()
	s1 := Server{
		Hostname:     "h1",
		IPs:          []string{"ip1"},
		Dependencies: []Dependency{{Resource: "h2"}},
		LastUpdate:   time.Now().Add(-24 * time.Hour),
	}
	assert.Nil(t, m.Add(ctx, s1))
	s2 := Server{
		Hostname:     "h2",
		IPs:          []string{"ip2"},
		Dependencies: []Dependency{{Resource: "ip3"}, {Resource: "h1"}},
		LastUpdate:   time.Now().Add(-48 * time.Hour),
	}
	assert.Nil(t, m.Add(ctx, s2))
	s3 := Server{
		Hostname:     "h3",
		IPs:          []string{"ip3"},
		Dependencies: []Dependency{{Resource: "ip4"}},
		LastUpdate:   time.Now().Add(-72 * time.Hour),
	}
	assert.Nil(t, m.Add(ctx, s3))
	old := 24 * 365 * time.Hour

	var tcs = []struct {
		inTC       string
		inIdent    string
		inDepth    int
		inSince    time.Duration
		expServers []Server
	}{
		{"1.1", "h1", 10, old, []Server{s1, s2, s3}},
		{"1.2", "ip1", 10, old, []Server{s1, s2, s3}},
		{"1.3", "h2", 10, old, []Server{s1, s2, s3}},
		{"1.4", "ip2", 10, old, []Server{s1, s2, s3}},
		{"1.5", "h3", 10, old, []Server{s3}},
		{"1.6", "ip3", 10, old, []Server{s3}},
		{"1.7", "h4", 10, old, []Server{}},
		{"1.8", "ip4", 10, old, []Server{}},

		{"2.1", "h1", 1, old, []Server{s1, s2}},
		{"2.2", "h1", 2, old, []Server{s1, s2, s3}},
		{"2.3", "h1", 10, old, []Server{s1, s2, s3}},
		{"2.4", "h1", 0, old, []Server{s1}},

		{"3.1", "h1", 10, 50 * time.Hour, []Server{s1, s2}},
		{"3.1", "h1", 10, 300 * time.Hour, []Server{s1, s2, s3}},
	}
	for _, tc := range tcs {
		t.Run(tc.inTC, func(t *testing.T) {
			got, err := m.GetDependencies(ctx, tc.inIdent, tc.inDepth, tc.inSince)
			assert.Nil(t, err)
			assert.ElementsMatch(t, tc.expServers, got)
		})
	}

}

func TestIsIdentifiedBy(t *testing.T) {
	s := Server{
		Hostname: "h",
		IPs:      []string{"i1", "i2"},
	}
	assert.True(t, s.IsIdentifiedBy("h"))
	assert.True(t, s.IsIdentifiedBy("i1"))
	assert.True(t, s.IsIdentifiedBy("i2"))
	assert.False(t, s.IsIdentifiedBy("i3"))
}
