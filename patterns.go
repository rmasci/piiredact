package piiredact

import "regexp"

var Patterns = map[string]*regexp.Regexp{
	"ssn":          regexp.MustCompile(`\b\d{3}[\s\-\.]?\d{2}[\s\-\.]?\d{4}\b`),
	"email":        regexp.MustCompile(`[\w\._%+-]+@[\w\.-]+\.\w{2,}`),
	"phone":        regexp.MustCompile(`(?i)\b(?:\+?1[\s\-\.]?)?\(?\d{3}\)?[\s\-\.]?\d{3}[\s\-\.]?\d{4}\b`),
	"credit_card":  regexp.MustCompile(`\b(?:\d[\s\-]?){13,16}\b`),
	"ip":           regexp.MustCompile(`\b\d{1,3}(?:\.\d{1,3}){3}\b`),
	"dob":          regexp.MustCompile(`(?i)\b(?:jan|feb|mar|apr|may|jun|jul|aug|sep|oct|nov|dec)[a-z]*[\s\-\.]?\d{1,2},?[\s\-\.]?\d{2,4}\b`),
	"zipcode":      regexp.MustCompile(`\b\d{5}(?:[-\s]?\d{4})?\b`),
}
