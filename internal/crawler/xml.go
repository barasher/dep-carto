package crawler

import (
	"encoding/xml"
	"io"
)

func NewXmlCrawler(refExt RefExtractorInterface) *xmlCrawler {
	return &xmlCrawler{
		refExt: refExt,
	}
}

type xmlCrawler struct {
	refExt RefExtractorInterface
}

func (c *xmlCrawler) Crawl(in io.Reader) ([]string, error) {
	var res []string
	d := xml.NewDecoder(in)
	for {
		t, err := d.Token()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		switch t := t.(type) {
		case xml.StartElement:
			for _, curAtt := range t.Attr {
				if ext := c.refExt.Extract(curAtt.Value); ext != "" {
					res = append(res, ext)
				}
			}
		case xml.CharData:
			if ext := c.refExt.Extract(string(t)); ext != "" {
				res = append(res, ext)
			}
		}
	}
	return res, nil
}
