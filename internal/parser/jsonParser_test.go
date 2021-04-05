package parser

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCrawlJson(t *testing.T) {
	f, err := os.Open("../../testdata/samplefiles/a.json")
	assert.Nil(t, err)
	defer f.Close()
	re, err := NewRefExtractor(WithSuffix(".acme"))
	assert.Nil(t, err)
	got, err := NewJsonParser(re).Parse(f)
	assert.Nil(t, err)
	exp := []string{
		"c.acme",
		"e.acme",
		"f.acme",
		"i.acme",
	}
	assert.ElementsMatch(t, exp, got)
}
