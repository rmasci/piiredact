package piiredact

import (
	"regexp"
)

// builtinPatterns defines the standard PII detection patterns.
//
// Each pattern includes a name, regex pattern, and optional validation function.
// These are the core patterns used by the redaction engine unless disabled
// in the configuration.
var builtinPatterns = []PatternDef{
	// Social Security Number (SSN)
	// Matches formats like 123-45-6789 or 123456789
	{
		Name:     "SSN",
		Regex:    regexp.MustCompile(`\b(?:\d{3}-\d{2}-\d{4}|\d{9})\b`),
		Validate: validateSSN,
	},

	// Credit Card Number (CC)
	// Matches major card formats with appropriate prefixes
	{
		Name:     "CC",
		Regex:    regexp.MustCompile(`\b(?:\d{4}[- ]?\d{4}[- ]?\d{4}[- ]?\d{4}|\d{16})\b`),
		Validate: validateLuhn,
	},

	// Phone Number (PHONE)
	// Matches various US formats
	{
		Name:     "PHONE",
		Regex:    regexp.MustCompile(`\b(?:\+?1[- ]?)?(?:\([0-9]{3}\)[- ]?|[0-9]{3}[- ]?)[0-9]{3}[- ]?[0-9]{4}\b`),
		Validate: nil,
	},

	// Bank Routing Number (ABA)
	// Matches 9-digit ABA routing numbers
	{
		Name:     "ABA",
		Regex:    regexp.MustCompile(`\b[0-9]{9}\b`),
		Validate: validateABA,
	},

	// Driver's License (DL)
	// Matches common formats across multiple states
	{
		Name:     "DL",
		Regex:    regexp.MustCompile(`\b(?:[A-Z][0-9]{7}|[A-Z][0-9]{8}|[A-Z]{2}[0-9]{6}|[0-9]{9})\b`),
		Validate: nil,
	},

	// Email Address (EMAIL)
	// Matches standard email address format
	{
		Name:     "EMAIL",
		Regex:    regexp.MustCompile(`\b[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}\b`),
		Validate: nil,
	},

	// IP Address (IP)
	// Matches IPv4 addresses
	{
		Name:     "IP",
		Regex:    regexp.MustCompile(`\b(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\b`),
		Validate: nil,
	},

	// Passport Number (PASSPORT)
	// Matches common US passport format
	{
		Name:     "PASSPORT",
		Regex:    regexp.MustCompile(`\b[A-Z][0-9]{8}\b`),
		Validate: nil,
	},

	// Date of Birth (DOB)
	// Matches common date formats
	{
		Name:     "DOB",
		Regex:    regexp.MustCompile(`\b(?:0[1-9]|1[0-2])[/.-](?:0[1-9]|[12][0-9]|3[01])[/.-](?:19|20)\d{2}\b`),
		Validate: nil,
	},
}
