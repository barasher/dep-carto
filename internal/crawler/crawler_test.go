package crawler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/barasher/dep-carto/internal/model"
	"github.com/barasher/dep-carto/internal/parser"
	"github.com/stretchr/testify/assert"
)

func buildCrawler(t *testing.T, u string) Crawler {
	re, err := parser.NewRefExtractor()
	assert.Nil(t, err)
	c, err := NewCrawler(u, *re)
	assert.Nil(t, err)
	return c
}

func TestBuildServer(t *testing.T) {
	c := buildCrawler(t, "url")
	deps := []string{"a", "b"}
	s := c.buildServer(deps)
	assert.ElementsMatch(t, deps, s.Dependencies)
	assert.NotZero(t, s.Hostname)
	assert.True(t, len(s.IPs) > 0)
	assert.NotZero(t, s.LastUpdate)
}

func TestPushServer_Nominal(t *testing.T) {
	s := model.Server{
		Hostname:     "h",
		IPs:          []string{"i1", "i2"},
		Dependencies: []string{"d1", "d2"},
		LastUpdate:   time.Now(),
		Key:          "k",
	}
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		assert.Equal(t, "/server", req.URL.String())
		defer req.Body.Close()
		var s2 model.Server
		assert.Nil(t, json.NewDecoder(req.Body).Decode(&s2))
		assert.Equal(t, s.Hostname, s2.Hostname)
		assert.ElementsMatch(t, s.Dependencies, s2.Dependencies)
		assert.ElementsMatch(t, s.IPs, s2.IPs)
		assert.Equal(t, s.Key, s2.Key)
		assert.Equal(t, s.LastUpdate.Unix(), s2.LastUpdate.Unix())
		rw.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()
	c := buildCrawler(t, server.URL)
	assert.Nil(t, c.pushServer(s))
}

func TestPushServer_WrongStatus(t *testing.T) {
	s := model.Server{Hostname: "h"}
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(http.StatusBadRequest)
	}))
	defer server.Close()
	c := buildCrawler(t, server.URL)
	assert.NotNil(t, c.pushServer(s))
}

func TestPushServer_BadUrl(t *testing.T) {
	s := model.Server{Hostname: "h"}
	c := buildCrawler(t, "blabla")
	assert.NotNil(t, c.pushServer(s))
}

func TestCrawl_Nominal(t *testing.T) {
	re, err := parser.NewRefExtractor(parser.WithSuffix(".acme"))
	assert.Nil(t, err)
	postCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		postCount++
		expDeps := []string{
			"c.acme",
			"e.acme",
			"f.acme",
			"i.acme",
			"a.acme",
			"b.acme",
		}
		defer req.Body.Close()
		var s2 model.Server
		assert.Nil(t, json.NewDecoder(req.Body).Decode(&s2))
		assert.ElementsMatch(t, expDeps, s2.Dependencies)
		rw.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()
	c, err := NewCrawler(server.URL, *re)
	assert.Nil(t, err)
	err = c.Crawl([]string{"../../testdata/samplefiles"})
	assert.Equal(t, 1, postCount)
	assert.Nil(t, err)
}

func TestCrawl_NonExistingFolder(t *testing.T) {
	re, err := parser.NewRefExtractor(parser.WithSuffix(".acme"))
	assert.Nil(t, err)
	postCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		postCount++
		rw.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()
	c, err := NewCrawler(server.URL, *re)
	assert.Nil(t, err)
	err = c.Crawl([]string{"../../testdata/nonExistingFolder"})
	assert.NotNil(t, err)
	assert.Equal(t, 0, postCount)
}

func TestParseFile_NonExistingFile(t *testing.T) {
	re, err := parser.NewRefExtractor(parser.WithSuffix(".acme"))
	assert.Nil(t, err)
	c, err := NewCrawler("", *re)
	assert.Nil(t, err)
	_, err = c.parseFile("nonExistingFile.json")
	assert.NotNil(t, err)
}
