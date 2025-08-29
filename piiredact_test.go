package piiredact

import (
	"fmt"
	"regexp"
	"testing"
	"time"
)

// TestRedactionEngine_Process tests the main processing functionality
func TestRedactionEngine_Process(t *testing.T) {
	// Create test chunks with various PII types
	chunks := []Chunk{
		{"id1", "A", "My SSN is 123-45-6789"},
		{"id2", "B", "My credit card is 4111 1111 1111 1111"},
		{"id3", "A", "Call me at 555-123-4567"},
		{"id4", "B", "My email is user@example.com"},
		{"id5", "A", "No PII in this chunk"},
	}

	// Expected results after redaction
	expected := []Chunk{
		{"id1", "A", "My SSN is [SSN]"},
		{"id2", "B", "My credit card is [CC]"},
		{"id3", "A", "Call me at [PHONE]"},
		{"id4", "B", "My email is [EMAIL]"},
		{"id5", "A", "No PII in this chunk"},
	}

	// Create engine with default config
	engine := NewRedactionEngine(DefaultConfig())

	// Process chunks
	result, err := engine.Process(chunks)

	// Check for errors
	if err != nil {
		t.Fatalf("Process returned error: %v", err)
	}

	// Check results
	if len(result) != len(expected) {
		t.Fatalf("Expected %d chunks, got %d", len(expected), len(result))
	}

	// Compare each chunk
	for i, chunk := range result {
		if chunk.UUID != expected[i].UUID ||
			chunk.Speaker != expected[i].Speaker ||
			chunk.Text != expected[i].Text {
			t.Errorf("Chunk %d mismatch:\nExpected: %+v\nGot: %+v",
				i, expected[i], chunk)
		}
	}
}

// TestRedactionEngine_CustomConfig tests using a custom configuration
func TestRedactionEngine_CustomConfig(t *testing.T) {
	// Create a custom configuration that only redacts SSNs
	config := Config{
		EnabledPatterns: map[string]bool{
			"SSN": true,
			// All others disabled by omission
		},
		RedactionFormat: "***%s***",
		MaxConcurrency:  2,
		Logging:         false,
	}

	// Create test chunks
	chunks := []Chunk{
		{"id1", "A", "My SSN is 123-45-6789 and my card is 4111 1111 1111 1111"},
		{"id2", "B", "Call me at 555-123-4567"},
	}

	// Expected results (only SSN redacted)
	expected := []Chunk{
		{"id1", "A", "My SSN is ***SSN*** and my card is 4111 1111 1111 1111"},
		{"id2", "B", "Call me at 555-123-4567"},
	}

	// Create engine with custom config
	engine := NewRedactionEngine(config)

	// Process chunks
	result, _ := engine.Process(chunks)

	// Compare results
	for i, chunk := range result {
		if chunk.Text != expected[i].Text {
			t.Errorf("Custom config test failed:\nExpected: %s\nGot: %s",
				expected[i].Text, chunk.Text)
		}
	}
}

// TestRedactionEngine_CustomPattern tests adding a custom pattern
func TestRedactionEngine_CustomPattern(t *testing.T) {
	// Create config with a custom pattern for employee IDs
	config := DefaultConfig()
	config.CustomPatterns = []PatternDef{
		{
			Name:     "EMPLOYEE_ID",
			Regex:    regexp.MustCompile(`\bEMP-\d{6}\b`),
			Validate: nil,
		},
	}

	// Create test chunks
	chunks := []Chunk{
		{"id1", "A", "My employee ID is EMP-123456"},
	}

	// Expected results
	expected := []Chunk{
		{"id1", "A", "My employee ID is [EMPLOYEE_ID]"},
	}

	// Create engine with custom pattern
	engine := NewRedactionEngine(config)

	// Process chunks
	result, _ := engine.Process(chunks)

	// Compare results
	if result[0].Text != expected[0].Text {
		t.Errorf("Custom pattern test failed:\nExpected: %s\nGot: %s",
			expected[0].Text, result[0].Text)
	}
}

// TestRedactionEngine_Metrics tests metrics collection
func TestRedactionEngine_Metrics(t *testing.T) {
	// Create engine
	engine := NewRedactionEngine(DefaultConfig())

	// Process chunks with various PII
	chunks := []Chunk{
		{"id1", "A", "SSN: 123-45-6789, Phone: 555-123-4567"},
		{"id2", "B", "Email: user@example.com"},
	}

	// Process chunks
	engine.Process(chunks)

	// Get metrics
	metrics := engine.GetMetrics()

	// Check metrics
	if metrics.ProcessedChunks != 2 {
		t.Errorf("Expected ProcessedChunks=2, got %d", metrics.ProcessedChunks)
	}

	if metrics.RedactedItems["SSN"] != 1 {
		t.Errorf("Expected 1 SSN redaction, got %d", metrics.RedactedItems["SSN"])
	}

	if metrics.RedactedItems["PHONE"] != 1 {
		t.Errorf("Expected 1 PHONE redaction, got %d", metrics.RedactedItems["PHONE"])
	}

	if metrics.RedactedItems["EMAIL"] != 1 {
		t.Errorf("Expected 1 EMAIL redaction, got %d", metrics.RedactedItems["EMAIL"])
	}

	// Test reset
	engine.ResetMetrics()
	metrics = engine.GetMetrics()

	if metrics.ProcessedChunks != 0 {
		t.Errorf("After reset, expected ProcessedChunks=0, got %d", metrics.ProcessedChunks)
	}
}

// TestRedactionEngine_Concurrency tests concurrent processing
func TestRedactionEngine_Concurrency(t *testing.T) {
	// Create a large batch of chunks to test concurrency
	var chunks []Chunk
	for i := 0; i < 1000; i++ {
		chunks = append(chunks, Chunk{
			UUID:    fmt.Sprintf("id%d", i),
			Speaker: "A",
			Text:    "SSN: 123-45-6789, Phone: 555-123-4567, Email: user@example.com",
		})
	}

	// Create configs with different concurrency settings
	config1 := DefaultConfig()
	config1.MaxConcurrency = 1 // Single-threaded

	config2 := DefaultConfig()
	config2.MaxConcurrency = 8 // Multi-threaded

	// Create engines
	engine1 := NewRedactionEngine(config1)
	engine2 := NewRedactionEngine(config2)

	// Process with both engines and measure time
	start1 := time.Now()
	engine1.Process(chunks)
	duration1 := time.Since(start1)

	start2 := time.Now()
	engine2.Process(chunks)
	duration2 := time.Since(start2)

	// Multi-threaded should be faster, but this is not guaranteed
	// so we just log the results rather than making it a test failure
	t.Logf("Single-threaded: %v, Multi-threaded: %v", duration1, duration2)
}
