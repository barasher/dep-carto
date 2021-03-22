package crawler

import (
	"encoding/xml"
	"io"
)

func NewXmlCrawler() *xmlCrawler {
	return &xmlCrawler{}
}

type xmlCrawler struct {
}

func (*xmlCrawler) Crawl(in io.Reader) ([]string, error) {
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
				if is, r := GetReferences(curAtt.Value); is {
					res = append(res, r)
				}
			}
		case xml.CharData:
			if is, r := GetReferences(string(t)); is {
				res = append(res, r)
			}
		}
	}
	return res, nil
}
