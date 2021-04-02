package parser

import (
	"io"
)

type Parser interface {
	Parse(in io.Reader) ([]string, error)
}

