package piiredact

import (
	"encoding/json"
	"fmt"
	"regexp"
)

// ExampleBasicUsage demonstrates the simplest way to use the package.
//
// It creates a redaction engine with default settings and processes
// a set of sample chunks containing various types of PII.
func ExampleBasicUsage() {
	// Create sample chunks with PII
	chunks := []Chunk{
		{"018f8f54-5b49-7cc5-9c3f-99b00a5f1cde", "A", "My social security is 401-23-4567"},
		{"018f8f54-5b49-7cc5-9c3f-99b00a5f1cdf", "B", "wait—what?"},
		{"018f8f54-5b49-7cc5-9c3f-99b00a5f1ce0", "A", "my card is 4111 1111 1111 1111"},
		{"018f8f54-5b49-7cc5-9c3f-99b00a5f1ce1", "A", "call me at 404-555-1212"},
	}

	// Create a redaction engine with default configuration
	engine := NewRedactionEngine(DefaultConfig())

	// Process the chunks
	redacted, _ := engine.Process(chunks)

	// Output the results
	jsonBytes, _ := json.MarshalIndent(redacted, "", "  ")
	fmt.Println(string(jsonBytes))

	// Output:
	// [
	//   {
	//     "uuid": "018f8f54-5b49-7cc5-9c3f-99b00a5f1cde",
	//     "speaker": "A",
	//     "text": "My social security is [SSN]"
	//   },
	//   {
	//     "uuid": "018f8f54-5b49-7cc5-9c3f-99b00a5f1cdf",
	//     "speaker": "B",
	//     "text": "wait—what?"
	//   },
	//   {
	//     "uuid": "018f8f54-5b49-7cc5-9c3f-99b00a5f1ce0",
	//     "speaker": "A",
	//     "text": "my card is [CC]"
	//   },
	//   {
	//     "uuid": "018f8f54-5b49-7cc5-9c3f-99b00a5f1ce1",
	//     "speaker": "A",
	//     "text": "call me at [PHONE]"
	//   }
	// ]
}

// ExampleCustomConfiguration demonstrates using a custom configuration.
//
// It shows how to:
// - Enable only specific patterns
// - Add custom patterns
// - Change the redaction format
// - Set concurrency limits
// - Enable logging
func ExampleCustomConfiguration() {
	// Create a custom configuration
	config := Config{
		// Only enable specific patterns
		EnabledPatterns: map[string]bool{
			"SSN":   true,
			"PHONE": true,
			"EMAIL": true,
			// Other patterns disabled by omission
		},

		// Add a custom pattern for detecting internal employee IDs
		CustomPatterns: []PatternDef{
			{
				Name:  "EMPLOYEE_ID",
				Regex: regexp.MustCompile(`\bEMP-\d{6}\b`),
				// No validation needed as pattern is specific
			},
		},

		// Use a different redaction format
		RedactionFormat: "***%s***",

		// Set concurrency to 4 workers
		MaxConcurrency: 4,

		// Enable logging
		Logging: true,
	}

	// Create the engine with custom configuration
	engine := NewRedactionEngine(config)

	// Process sample chunks
	chunks := []Chunk{
		{"id1", "A", "My employee ID is EMP-123456 and my email is user@example.com"},
		{"id2", "B", "My credit card is 4111 1111 1111 1111"}, // CC pattern disabled
	}

	redacted, _ := engine.Process(chunks)

	// Output the results
	jsonBytes, _ := json.MarshalIndent(redacted, "", "  ")
	fmt.Println(string(jsonBytes))

	// Output:
	// [
	//   {
	//     "uuid": "id1",
	//     "speaker": "A",
	//     "text": "My employee ID is ***EMPLOYEE_ID*** and my email is ***EMAIL***"
	//   },
	//   {
	//     "uuid": "id2",
	//     "speaker": "B",
	//     "text": "My credit card is 4111 1111 1111 1111"
	//   }
	// ]
}

// ExampleMetrics demonstrates how to access and use metrics.
//
// It shows processing a batch of chunks and then retrieving
// and displaying the performance metrics.
func ExampleMetrics() {
	engine := NewRedactionEngine(DefaultConfig())

	// Process a batch of chunks with various PII
	chunks := []Chunk{
		{"id1", "A", "SSN: 123-45-6789, Phone: 555-123-4567"},
		{"id2", "B", "Email: user@example.com"},
		{"id3", "A", "Credit card: 4111 1111 1111 1111"},
	}

	engine.Process(chunks)

	// Get the metrics
	metrics := engine.GetMetrics()

	// Display metrics
	fmt.Printf("Processed chunks: %d\n", metrics.ProcessedChunks)
	fmt.Printf("Processing time: %d ns\n", metrics.ProcessingTimeNs)
	fmt.Println("Redacted items:")
	for pattern, count := range metrics.RedactedItems {
		if count > 0 {
			fmt.Printf("  %s: %d\n", pattern, count)
		}
	}

	// Reset metrics for a new measurement period
	engine.ResetMetrics()

	// Output:
	// Processed chunks: 3
	// Processing time: [variable]
	// Redacted items:
	//   SSN: 1
	//   PHONE: 1
	//   EMAIL: 1
	//   CC: 1
}
