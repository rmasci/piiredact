// Package piiredact provides enterprise-grade Personally Identifiable Information (PII)
// redaction capabilities for text data. It is designed for high-performance processing
// of transcription chunks from speech-to-text systems or other text sources containing
// sensitive information.
//
// The package offers:
//   - Detection and redaction of multiple PII types (SSN, credit cards, phone numbers, etc.)
//   - Validation algorithms to minimize false positives
//   - Concurrent processing for high-throughput applications
//   - Configurable redaction behavior and pattern matching
//   - Performance metrics and optional logging
//
// Basic usage:
//
//	engine := piiredact.NewRedactionEngine(piiredact.DefaultConfig())
//	redactedChunks, _ := engine.Process(chunks)
//
// For more advanced usage, see the documentation for RedactionEngine and Config.
package piiredact

import (
	"fmt"
	"log"
	"regexp"
	"sync"
	"time"
)

// Chunk represents a single transcription chunk that may contain PII.
//
// UUID is a unique identifier (e.g., uuidv7) assigned to this chunk.
// Speaker is the speaker label (e.g., "A", "B") if available.
// Text is the spoken text, potentially redacted by RedactChunk.
type Chunk struct {
	UUID    string `json:"uuid"`    // Unique identifier for the chunk
	Speaker string `json:"speaker"` // Speaker identifier (e.g., "A", "B")
	Text    string `json:"text"`    // Text content, potentially containing PII
}

// PatternDef defines a pattern for PII detection and its validation function.
//
// Name is used in replacement text (e.g., "[SSN]").
// Regex is the compiled regular expression that matches the pattern.
// Validate is an optional function that confirms matches are valid PII.
type PatternDef struct {
	Name     string            // Name of the PII type (used in redaction)
	Regex    *regexp.Regexp    // Compiled regex pattern for detection
	Validate func(string) bool // Optional validation function to reduce false positives
}

// Config provides configuration options for the redaction engine.
//
// EnabledPatterns controls which patterns are active.
// CustomPatterns allows adding user-defined patterns.
// RedactionFormat defines how redacted text appears.
// MaxConcurrency limits parallel processing.
// Logging enables operational logging.
type Config struct {
	EnabledPatterns map[string]bool // Map of pattern names to enabled status
	CustomPatterns  []PatternDef    // Additional user-defined patterns
	RedactionFormat string          // Format string for redactions (default: "[%s]")
	MaxConcurrency  int             // Maximum number of concurrent goroutines
	Logging         bool            // Whether to log redaction operations
}

// DefaultConfig returns a configuration with sensible defaults.
//
// All built-in patterns are enabled, with standard redaction format "[TYPE]",
// and concurrency set to match available CPU cores.
func DefaultConfig() Config {
	// Create default enabled patterns map with all built-in patterns enabled
	enabled := make(map[string]bool)
	for _, p := range builtinPatterns {
		enabled[p.Name] = true
	}

	return Config{
		EnabledPatterns: enabled,
		CustomPatterns:  []PatternDef{},
		RedactionFormat: "[%s]",
		MaxConcurrency:  8, // Default to 8 concurrent workers
		Logging:         false,
	}
}

// Metrics tracks performance and detection statistics for the redaction engine.
//
// Thread-safe counters for monitoring redaction operations and performance.
type Metrics struct {
	ProcessedChunks  int64            // Total number of chunks processed
	RedactedItems    map[string]int64 // Count of redactions by pattern type
	ProcessingTimeNs int64            // Total processing time in nanoseconds
	mu               sync.Mutex       // Mutex for thread-safe updates
}

// newMetrics initializes a new Metrics instance with zeroed counters.
func newMetrics() *Metrics {
	redactedItems := make(map[string]int64)
	// Initialize counters for all built-in patterns
	for _, p := range builtinPatterns {
		redactedItems[p.Name] = 0
	}

	return &Metrics{
		ProcessedChunks:  0,
		RedactedItems:    redactedItems,
		ProcessingTimeNs: 0,
	}
}

// RedactionEngine provides the main interface for PII redaction operations.
//
// It encapsulates configuration, patterns, and metrics for redaction processing.
type RedactionEngine struct {
	config   Config       // Configuration options
	patterns []PatternDef // Active detection patterns
	logger   *log.Logger  // Optional logger for operations
	metrics  *Metrics     // Performance and detection metrics
}

// NewRedactionEngine creates a new engine with the given configuration.
//
// It initializes the engine with the specified configuration, compiling
// all enabled built-in and custom patterns, and setting up metrics tracking.
func NewRedactionEngine(config Config) *RedactionEngine {
	// Initialize patterns from enabled built-in patterns and custom patterns
	var patterns []PatternDef

	// Add enabled built-in patterns
	for _, p := range builtinPatterns {
		if enabled, exists := config.EnabledPatterns[p.Name]; !exists || enabled {
			// Include pattern if it's not in the map (default) or explicitly enabled
			patterns = append(patterns, p)
		}
	}

	// Add custom patterns
	patterns = append(patterns, config.CustomPatterns...)

	// Create logger if logging is enabled
	var logger *log.Logger
	if config.Logging {
		logger = log.Default()
	}

	return &RedactionEngine{
		config:   config,
		patterns: patterns,
		logger:   logger,
		metrics:  newMetrics(),
	}
}

// Process handles a batch of chunks with metrics and logging.
//
// It processes all chunks according to the engine configuration,
// updates metrics, and returns the redacted chunks.
func (e *RedactionEngine) Process(chunks []Chunk) ([]Chunk, error) {
	startTime := time.Now()

	// Process chunks with configured concurrency
	result := e.processChunks(chunks)

	// Update metrics
	duration := time.Since(startTime)
	e.metrics.mu.Lock()
	e.metrics.ProcessedChunks += int64(len(chunks))
	e.metrics.ProcessingTimeNs += duration.Nanoseconds()
	e.metrics.mu.Unlock()

	// Log summary if enabled
	if e.config.Logging && e.logger != nil {
		e.logger.Printf("Processed %d chunks in %v", len(chunks), duration)
	}

	return result, nil
}

// processChunks handles concurrent processing of multiple chunks.
//
// It uses a worker pool with semaphore to limit concurrency based on
// the engine configuration.
func (e *RedactionEngine) processChunks(chunks []Chunk) []Chunk {
	result := make([]Chunk, len(chunks))

	// If only processing a single chunk or concurrency is set to 1,
	// process sequentially for better efficiency
	if len(chunks) == 1 || e.config.MaxConcurrency == 1 {
		for i, chunk := range chunks {
			result[i] = e.redactChunk(chunk)
		}
		return result
	}

	// Otherwise, process concurrently
	var wg sync.WaitGroup

	// Use a worker pool to limit goroutines
	maxWorkers := e.config.MaxConcurrency
	if maxWorkers <= 0 {
		maxWorkers = 8 // Fallback to default if invalid
	}
	semaphore := make(chan struct{}, maxWorkers)

	// Process each chunk in a separate goroutine
	for i, chunk := range chunks {
		wg.Add(1)
		semaphore <- struct{}{} // Acquire semaphore

		go func(i int, c Chunk) {
			defer wg.Done()
			defer func() { <-semaphore }() // Release semaphore

			// Process the chunk and store in result
			result[i] = e.redactChunk(c)
		}(i, chunk)
	}

	wg.Wait() // Wait for all goroutines to complete
	return result
}

// redactChunk applies PII redaction to a single chunk.
//
// It processes the text with all active patterns, applying validation
// where available, and formats redactions according to configuration.
func (e *RedactionEngine) redactChunk(c Chunk) Chunk {
	redacted := c.Text
	redactionCounts := make(map[string]int)

	// Apply each pattern to the text
	for _, p := range e.patterns {
		// Find all matches for this pattern
		matches := p.Regex.FindAllStringIndex(redacted, -1)

		// Process matches in reverse order to avoid offset issues
		// when replacing text (earlier replacements would change string indices)
		for i := len(matches) - 1; i >= 0; i-- {
			match := matches[i]
			start, end := match[0], match[1]
			potentialPII := redacted[start:end]

			// Skip validation if no validation function or validation passes
			if p.Validate == nil || p.Validate(potentialPII) {
				// Format the redaction according to configuration
				replacement := fmt.Sprintf(e.config.RedactionFormat, p.Name)
				redacted = redacted[:start] + replacement + redacted[end:]
				redactionCounts[p.Name]++
			}
		}
	}

	// Update metrics with redaction counts
	if len(redactionCounts) > 0 {
		e.metrics.mu.Lock()
		for name, count := range redactionCounts {
			e.metrics.RedactedItems[name] += int64(count)
		}
		e.metrics.mu.Unlock()

		// Log redactions if enabled
		if e.config.Logging && e.logger != nil {
			e.logger.Printf("Chunk %s: redacted %v items", c.UUID, redactionCounts)
		}
	}

	// Return the redacted chunk
	c.Text = redacted
	return c
}

// GetMetrics returns a copy of the current metrics.
//
// This provides a thread-safe way to access the engine's performance
// and detection statistics.
func (e *RedactionEngine) GetMetrics() Metrics {
	e.metrics.mu.Lock()
	defer e.metrics.mu.Unlock()

	// Create a deep copy of the metrics
	redactedItems := make(map[string]int64)
	for k, v := range e.metrics.RedactedItems {
		redactedItems[k] = v
	}

	return Metrics{
		ProcessedChunks:  e.metrics.ProcessedChunks,
		RedactedItems:    redactedItems,
		ProcessingTimeNs: e.metrics.ProcessingTimeNs,
	}
}

// ResetMetrics resets all metrics counters to zero.
//
// This is useful for starting a new measurement period.
func (e *RedactionEngine) ResetMetrics() {
	e.metrics.mu.Lock()
	defer e.metrics.mu.Unlock()

	e.metrics.ProcessedChunks = 0
	e.metrics.ProcessingTimeNs = 0
	for k := range e.metrics.RedactedItems {
		e.metrics.RedactedItems[k] = 0
	}
}
