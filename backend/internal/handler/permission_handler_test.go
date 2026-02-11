package handler

import (
	"testing"

	"community-elderly-care-platform/internal/middleware"
)

func TestMatchPath(t *testing.T) {
	tests := []struct {
		name     string
		pattern  string
		path     string
		expected bool
	}{
		{
			name:     "exact match",
			pattern:  "/api/v1/b/tasks/pool",
			path:     "/api/v1/b/tasks/pool",
			expected: true,
		},
		{
			name:     "wildcard match",
			pattern:  "/api/v1/b/tasks/*/claim",
			path:     "/api/v1/b/tasks/123/claim",
			expected: true,
		},
		{
			name:     "no match",
			pattern:  "/api/v1/b/tasks/pool",
			path:     "/api/v1/b/tasks/my",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := middleware.MatchPath(tt.pattern, tt.path)
			if result != tt.expected {
				t.Errorf("MatchPath(%q, %q) = %v, want %v", tt.pattern, tt.path, result, tt.expected)
			}
		})
	}
}
