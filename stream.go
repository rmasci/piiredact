package piiredact

import (
	"bufio"
	"io"
	"strings"
)

// RedactStream reads from r, redacts PII, and writes to w line by line.
func RedactStream(r io.Reader, w io.Writer, opts *RedactOptions) error {
	if opts == nil {
		opts = DefaultOptions()
	}
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		redacted := RedactWithOptions(line, opts)
		_, err := w.Write([]byte(redacted + "\n"))
		if err != nil {
			return err
		}
	}
	return scanner.Err()
}
