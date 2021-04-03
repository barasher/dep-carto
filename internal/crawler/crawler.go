package crawler

import (
	"fmt"
	"github.com/barasher/dep-carto/internal/model"
	"github.com/barasher/dep-carto/internal/parser"
	"os"
	"path/filepath"
)

const (
	xmlExtension="xml"
	jsonExtension="json"
)

type Crawler struct {
	re parser.RefExtractor
}

func NewCrawler(re parser.RefExtractor) Crawler {
	return Crawler{re: re}
}

func (c Crawler) Crawl(inputs []string) error {
	for _, curInput := range inputs {
		err := filepath.Walk(curInput, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() {
				s, err := c.crawlFile(path)
				if err != nil {
					return fmt.Errorf("error while crawling %v: %w", path, err)
				}
				if s != nil {
					c.pushServer(*s)
				}
			}
			return nil
		})
		if err != nil {
			return fmt.Errorf("error while browsing %v: %w", curInput, err)
		}
	}
	return nil
}

func (c Crawler) crawlFile(f string) (*model.Server, error) {
	/*var p parser.Parser
	switch ext := filepath.Ext(f) ; ext{
	case xmlExtension:
		p = parser.NewXmlParser(&c.re)
	case jsonExtension:
		p = parser.NewJsonParser(&c.re)
	default:
		log.Info().Str("file", f).Msgf("Unsupported extension (%v)", ext)
		return nil, nil
	}

	r, err := os.Open(f)
	if err != nil {
		return nil, fmt.Errorf("error while opening %v: %w", f, err)
	}
	defer r.Close()
	deps, err := p.Parse(r)
	if err != nil {
		return nil, fmt.Errorf("error while opening %v: %w", f, err)
	}*/

	return nil, nil
}

func (Crawler) pushServer(s model.Server) error {
	return nil
}
