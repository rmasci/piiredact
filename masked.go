package piiredact

import (
	"regexp"
)

// MaskSSN returns an SSN in the format XXX-XX-1234
func MaskSSN(ssn string) string {
	re := regexp.MustCompile(`(\d{3})[\s\-\.]?(\d{2})[\s\-\.]?(\d{4})`)
	if match := re.FindStringSubmatch(ssn); match != nil {
		return "XXX-XX-" + match[3]
	}
	return ssn
}

// MaskCreditCard returns a masked credit card like XXXX-XXXX-XXXX-5678
func MaskCreditCard(card string) string {
	re := regexp.MustCompile(`(\d{0,4})[\s\-]?(\d{0,4})[\s\-]?(\d{0,4})[\s\-]?(\d{4})`)
	if match := re.FindStringSubmatch(card); match != nil && len(match) == 5 {
		return "XXXX-XXXX-XXXX-" + match[4]
	}
	return card
}
