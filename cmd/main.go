package main

import (
	"encoding/json"
	"fmt"
	"github.com/rmasci/piiredact"
)

func main() {
	// Create sample data with PII
	chunks := []piiredact.Chunk{
		{"id1", "A", "My SSN is 123-45-6789"},
		{"id2", "B", "Call me at 555-123-4567"},
		{"id3", "A", "My email is user@example.com"},
	}

	// Create a redaction engine with default settings
	engine := piiredact.NewRedactionEngine(piiredact.DefaultConfig())

	// Process the chunks
	redacted, _ := engine.Process(chunks)

	// Output the results
	jsonBytes, _ := json.MarshalIndent(redacted, "", "  ")
	fmt.Println(string(jsonBytes))
}
