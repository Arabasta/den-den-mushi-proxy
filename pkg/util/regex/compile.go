package regex

import (
	"errors"
	"regexp"
)

func CompilePattern(pattern string) (*regexp.Regexp, error) {
	if pattern == "" {
		return nil, errors.New("empty regex pattern")
	}
	return regexp.Compile(pattern)
}

func DecompilePattern(r *regexp.Regexp) string {
	if r == nil {
		return ""
	}
	return r.String()
}
