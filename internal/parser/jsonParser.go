package parser

import (
	"encoding/json"
	"io"
)

func NewJsonParser(refExt RefExtractorInterface) *jsonParser {
	return &jsonParser{
		refExt: refExt,
	}
}

type jsonParser struct {
	refExt RefExtractorInterface
}

func (c *jsonParser) Parse(in io.Reader) ([]string, error) {
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
			if ext := c.refExt.Extract(t); ext != "" {
				res = append(res, ext)
			}
		}
	}
	return res, nil
}
