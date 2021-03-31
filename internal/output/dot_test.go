package output

import (
	"bytes"
	"github.com/barasher/dep-carto/internal/model"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBuildDotGraph(t *testing.T) {
	s := generateServers()
	dg := newDotGraph(s)
	expServ := []model.Server{
		model.Server{Hostname: "s1.domain", IPs: []string{"ip1a", "ip1b"}},
		model.Server{Hostname: "s2.domain", IPs: []string{"ip2"}},
		model.Server{Hostname: "s3.domain", IPs: []string{"ip3"}},
	}
	assert.ElementsMatch(t, expServ, dg.Servers)
	assert.ElementsMatch(t, []string{"s.otherdomain"}, dg.ExternalServers)
	expDep := []DotGraphDep{
		{"s1.domain", "s2.domain"},
		{"s1.domain", "s.otherdomain"},
		{"s1.domain", "s3.domain"},
	}
	assert.ElementsMatch(t, expDep, dg.Dependencies)
}

func TestDotFormat(t *testing.T) {
	buf := bytes.Buffer{}
	ct, err := DotOutput{}.Format(generateServers(), &buf)
	assert.Nil(t, err)
	assert.Equal(t, "text/plain", ct)
	t.Logf("%v", buf.String())
	assert.Fail(t, "aaaa")
}

