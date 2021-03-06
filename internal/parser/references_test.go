package parser

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExtract(t *testing.T) {
	var noPrefix []RefExtractorOption
	var tcs = []struct {
		tcID       string
		tcOpts     []RefExtractorOption
		tcIn       string
		expMatches string
	}{
		{"ip1", noPrefix, "1.2.3.255", "1.2.3.255"},
		{"ip3", noPrefix, "_1.2.3.4_", "1.2.3.4"},
		{"prefix", []RefExtractorOption{WithSuffix(".google.com")}, "http://pre.m_-ail.google.com:80", "pre.m_-ail.google.com"},
		{"nothing", noPrefix, "blabla", ""},
	}

	for _, tc := range tcs {
		t.Run(tc.tcID, func(t *testing.T) {
			re, err := NewRefExtractor(tc.tcOpts...)
			assert.Nil(t, err)
			got := re.Extract(tc.tcIn)
			assert.Equal(t, tc.expMatches, got)
		})
	}
}

func TestNewRefExtractor_FailOnOpts(t *testing.T) {
	opt := func(extractor *RefExtractor) error {
		return fmt.Errorf("err")
	}
	_, err := NewRefExtractor(opt)
	assert.NotNil(t, err)
}
