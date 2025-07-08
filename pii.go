package piiredact

import (
	"strings"
)

// RedactPII replaces detected PII with [REDACTED:<TYPE>] using default options
func RedactPII(text string) string {
	return RedactWithOptions(text, DefaultOptions())
}

// RedactWithOptions redacts text using the provided RedactOptions
func RedactWithOptions(text string, opts *RedactOptions) string {
	for label, re := range Patterns {
		text = re.ReplaceAllStringFunc(text, func(match string) string {
			return opts.ReplaceFunc(strings.ToUpper(label), match)
		})
	}
	return text
}

// MatchPII returns all matched PII by type (for logging/stats)
func MatchPII(text string) map[string][]string {
	matches := make(map[string][]string)
	for label, re := range Patterns {
		all := re.FindAllString(text, -1)
		if len(all) > 0 {
			matches[label] = all
		}
	}
	return matches
}
