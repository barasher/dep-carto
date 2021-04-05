package parser

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCrawlXml(t *testing.T) {
	f, err := os.Open("../../testdata/samplefiles/a.xml")
	assert.Nil(t, err)
	defer f.Close()
	re, err := NewRefExtractor(WithSuffix(".acme"))
	assert.Nil(t, err)
	got, err := NewXmlParser(re).Parse(f)
	assert.Nil(t, err)
	exp := []string{
		"a.acme",
		"b.acme",
		"c.acme",
	}
	assert.ElementsMatch(t, exp, got)
}
