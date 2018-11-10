package contains

import (
	"bytes"
	"context"
	"testing"
)

func TestContains(t *testing.T) {
	tests := []struct {
		name     string
		text     string
		contains string
		expected bool
	}{
		{
			name:     "success",
			text:     "this is OK",
			contains: "is ",
			expected: true,
		},
		{
			name:     "negative",
			text:     "this is not OK",
			contains: "nope",
			expected: false,
		},
		{
			name:     "multi-line negative",
			text:     "this is not OK\nstill not OK",
			contains: "nope",
			expected: false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			r := bytes.NewBuffer([]byte(test.text))
			ok, _, err := TextInReader(context.Background(), r, test.contains)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if ok != test.expected {
				t.Errorf("expected %v, got %v", test.expected, ok)
			}
		})
	}
}
