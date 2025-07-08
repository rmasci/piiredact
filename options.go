package piiredact

type RedactOptions struct {
	ReplaceFunc func(label, match string) string
}

func DefaultOptions() *RedactOptions {
	return &RedactOptions{
		ReplaceFunc: func(label, _ string) string {
			return "[REDACTED:" + label + "]"
		},
	}
}
