package crawler

import (
	"io"
)

type Crawler interface {
	Crawl(in io.Reader) ([]string, error)
}

