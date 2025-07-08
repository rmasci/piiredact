package piiredact

import "testing"

func TestRedactPII(t *testing.T) {
	input := "SSN 123 45 6789, email joe@example.com, card 4111 8888 1234 5678"
	expected := "SSN [REDACTED:SSN], email [REDACTED:EMAIL], card [REDACTED:CREDIT_CARD]"
	got := RedactPII(input)
	if got != expected {
		t.Errorf("Expected %s but got %s", expected, got)
	}
}
