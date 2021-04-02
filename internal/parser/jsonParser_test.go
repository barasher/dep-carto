package parser

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestCrawlJson(t *testing.T) {
	f, err := os.Open("../../testdata/samplefiles/a.json")
	assert.Nil(t, err)
	defer f.Close()
	got, err := NewJsonParser(refExtractorMock{}).Parse(f)
	assert.Nil(t, err)
	assert.ElementsMatch(t, got, []string{"a", "b", "c", "d", "e", "f", "g", "h", "i"})
}
