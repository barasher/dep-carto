package crawler

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestCrawlJson(t *testing.T) {
	f, err := os.Open("../../testdata/samplefiles/a.json")
	assert.Nil(t, err)
	defer f.Close()
	got, err := NewJsonCrawler(refExtractorMock{}).Crawl(f)
	assert.Nil(t, err)
	assert.ElementsMatch(t, got, []string{"a", "b", "c", "d", "e", "f", "g", "h", "i"})
}
