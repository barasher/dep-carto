package crawler

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestExtract(t *testing.T) {
	var noPrefix []refExtractorOption
	var tcs = []struct {
		tcID       string
		tcOpts     []refExtractorOption
		tcIn       string
		expMatches string
	}{
		{"ip1", noPrefix, "1.2.3.255", "1.2.3.255"},
		{"ip3", noPrefix, "_1.2.3.4_", "1.2.3.4"},
		{"prefix", []refExtractorOption{WithSuffix(".google.com")}, "http://pre.m_-ail.google.com:80", "pre.m_-ail.google.com"},
	}

	for _, tc := range tcs {
		t.Run(tc.tcID, func(t *testing.T) {
			re, err := NewRefExtractor(tc.tcOpts...)
			assert.Nil(t, err)
			got:= re.Extract(tc.tcIn)
			assert.Equal(t, tc.expMatches, got)
		})
	}
}
