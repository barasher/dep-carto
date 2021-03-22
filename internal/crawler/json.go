package crawler

import (
	"encoding/json"
	"io"
)

func NewJsonCrawler() *jsonCrawler {
	return &jsonCrawler{}
}

type jsonCrawler struct {
}

func (*jsonCrawler) Crawl(in io.Reader) ([]string, error) {
	var res []string
	d := json.NewDecoder(in)
	for {
		t, err := d.Token()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		switch t := t.(type) {
		case string:
			if is, r := GetReferences(t); is {
				res = append(res, r)
			}
		}
	}
	return res, nil
}
