package parser

import (
	"fmt"
	"regexp"
	"strings"
)

const ipRegexp = "([0-9]{1,3}\\.){3}[0-9]{1,3}"
const escChar = "\\"
const suffixRegexpPrefix = `[a-z\.\-\_]*`

var userSuffixOverride = []string{"-", "_", "."}

type RefExtractor struct {
	regexps []*regexp.Regexp
}

type RefExtractorOption func(extractor *RefExtractor) error

func NewRefExtractor(opts ...RefExtractorOption) (*RefExtractor, error) {
	refExt := &RefExtractor{}
	ipRe, err := regexp.Compile(ipRegexp)
	if err != nil {
		return nil, fmt.Errorf("error while compiling ip regexp: %w", err)
	}
	refExt.regexps = []*regexp.Regexp{ipRe}

	for _, opt := range opts {
		if err = opt(refExt); err != nil {
			return nil, fmt.Errorf("error while initializing RefExtractor: %w", err)
		}
	}
	return refExt, nil
}

func WithSuffix(s string) RefExtractorOption {
	return func(refExt *RefExtractor) error {
		suffixRe := s
		for _, v := range userSuffixOverride {
			suffixRe = strings.ReplaceAll(suffixRe, v, escChar+v)
		}
		suffixRe = suffixRegexpPrefix + suffixRe
		re, err := regexp.Compile(suffixRe)
		if err != nil {
			return fmt.Errorf("error while compiling regexp for suffix %v (%v): %w", s, suffixRe, err)
		}
		refExt.regexps = append(refExt.regexps, re)
		return nil
	}
}

// TODO serialized multivalue attributes
func (refExt *RefExtractor) Extract(s string) string {
	for _, re := range refExt.regexps {
		matched := re.FindString(s)
		if matched != "" {
			return matched
		}
	}
	return ""
}

type RefExtractorInterface interface {
	Extract(s string) string
}
