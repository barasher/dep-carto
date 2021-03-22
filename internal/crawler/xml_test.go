package crawler

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestCrawlXml(t *testing.T) {
	f, err := os.Open("../../testdata/samplefiles/a.xml")
	assert.Nil(t, err)
	defer f.Close()
	got, err := NewXmlCrawler().Crawl(f)
	assert.Nil(t, err)
	assert.ElementsMatch(t, got, []string{"a", "b", "c"})
}
