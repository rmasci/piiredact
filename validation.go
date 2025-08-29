package piiredact

import (
	"regexp"
	"strconv"
	"strings"
)

// validateSSN checks if a potential SSN follows valid format rules.
//
// It applies various validation rules to minimize false positives:
// - Rejects all-same-digit patterns (e.g., 111-11-1111)
// - Validates against SSA issuance rules (no 000, 666, 900+ area numbers)
// - Checks for valid group and serial numbers
func validateSSN(ssn string) bool {
	// Remove hyphens for validation
	cleaned := strings.ReplaceAll(ssn, "-", "")

	// Check for obviously invalid patterns (all same digit)
	if cleaned == "000000000" || cleaned == "111111111" ||
		cleaned == "222222222" || cleaned == "333333333" ||
		cleaned == "444444444" || cleaned == "555555555" ||
		cleaned == "666666666" || cleaned == "777777777" ||
		cleaned == "888888888" || cleaned == "999999999" {
		return false
	}

	// First 3 digits can't be 000, 666, or 900-999
	first3, err := strconv.Atoi(cleaned[:3])
	if err != nil || first3 == 0 || first3 == 666 || first3 >= 900 {
		return false
	}

	// Middle 2 digits can't be 00
	middle2, err := strconv.Atoi(cleaned[3:5])
	if err != nil || middle2 == 0 {
		return false
	}

	// Last 4 digits can't be 0000
	last4, err := strconv.Atoi(cleaned[5:])
	if err != nil || last4 == 0 {
		return false
	}

	return true
}

// validateLuhn implements the Luhn algorithm for credit card validation.
//
// This algorithm detects accidental errors in identification numbers:
// 1. Starting from the rightmost digit, double the value of every second digit
// 2. If doubling results in a two-digit number, add those digits together
// 3. Sum all digits in the resulting sequence
// 4. If the sum is divisible by 10, the number is valid
func validateLuhn(number string) bool {
	// Remove spaces and dashes
	number = strings.ReplaceAll(number, " ", "")
	number = strings.ReplaceAll(number, "-", "")

	// Check if all digits and minimum length
	if len(number) < 13 || len(number) > 19 || !regexp.MustCompile(`^\d+$`).MatchString(number) {
		return false
	}

	var sum int
	var alternate bool

	// Process digits from right to left
	for i := len(number) - 1; i >= 0; i-- {
		digit, _ := strconv.Atoi(string(number[i]))

		// Double every second digit
		if alternate {
			digit *= 2
			if digit > 9 {
				digit -= 9 // Same as adding the digits together
			}
		}

		sum += digit
		alternate = !alternate
	}

	// Valid if sum is divisible by 10
	return sum%10 == 0
}

// validateABA checks if a routing number is valid using the checksum algorithm.
//
// ABA routing numbers use a specific checksum algorithm:
// 3(d1+d4+d7) + 7(d2+d5+d8) + (d3+d6+d9) must be divisible by 10
func validateABA(aba string) bool {
	// Must be exactly 9 digits
	if len(aba) != 9 || !regexp.MustCompile(`^\d+$`).MatchString(aba) {
		return false
	}

	// First digit can't be 0
	if aba[0] == '0' {
		return false
	}

	// Calculate checksum using ABA algorithm
	d1, _ := strconv.Atoi(string(aba[0]))
	d2, _ := strconv.Atoi(string(aba[1]))
	d3, _ := strconv.Atoi(string(aba[2]))
	d4, _ := strconv.Atoi(string(aba[3]))
	d5, _ := strconv.Atoi(string(aba[4]))
	d6, _ := strconv.Atoi(string(aba[5]))
	d7, _ := strconv.Atoi(string(aba[6]))
	d8, _ := strconv.Atoi(string(aba[7]))
	d9, _ := strconv.Atoi(string(aba[8]))

	sum := 3*(d1+d4+d7) + 7*(d2+d5+d8) + (d3 + d6 + d9)
	return sum%10 == 0
}
