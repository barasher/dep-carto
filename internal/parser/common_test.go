package parser

import "strings"

type refExtractorMock struct{}

func (refExtractorMock) Extract(s string) string {
	if len(strings.TrimSpace(s)) > 0 {
		return s
	}
	return ""
}
