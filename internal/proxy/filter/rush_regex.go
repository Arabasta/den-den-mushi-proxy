package filter

import "regexp"

func init() {
	changeFilter.ouGroupRegexFiltersMap["default"] = []regexp.Regexp{
		*regexp.MustCompile(`(?i)^\s*su\s*$`),         // su
		*regexp.MustCompile(`(?i)^\s*su\s*-$`),        // su -
		*regexp.MustCompile(`(?i)^\s*su\s+-\s*root$`), // su - root
		*regexp.MustCompile(`(?i)^\s*su\s+root$`),     // su root
	}
}
