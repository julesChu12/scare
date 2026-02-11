package service

import (
	"testing"
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
			name:     "wildcard match single segment",
			pattern:  "/api/v1/b/tasks/*/claim",
			path:     "/api/v1/b/tasks/123/claim",
			expected: true,
		},
		{
			name:     "wildcard match with uuid",
			pattern:  "/api/v1/b/users/*",
			path:     "/api/v1/b/users/abc-123-def",
			expected: true,
		},
		{
			name:     "wildcard in middle",
			pattern:  "/api/v1/*/auth/me",
			path:     "/api/v1/b/auth/me",
			expected: true,
		},
		{
			name:     "multiple wildcards",
			pattern:  "/api/v1/*/notifications/*/read",
			path:     "/api/v1/b/notifications/456/read",
			expected: true,
		},
		{
			name:     "no match - different length",
			pattern:  "/api/v1/b/tasks",
			path:     "/api/v1/b/tasks/pool",
			expected: false,
		},
		{
			name:     "no match - different segment",
			pattern:  "/api/v1/b/tasks/pool",
			path:     "/api/v1/b/tasks/my",
			expected: false,
		},
		{
			name:     "no match - different prefix",
			pattern:  "/api/v1/b/tasks/pool",
			path:     "/api/v2/b/tasks/pool",
			expected: false,
		},
		{
			name:     "empty paths",
			pattern:  "",
			path:     "",
			expected: true,
		},
		{
			name:     "root path",
			pattern:  "/",
			path:     "/",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := matchPath(tt.pattern, tt.path)
			if result != tt.expected {
				t.Errorf("matchPath(%q, %q) = %v, want %v", tt.pattern, tt.path, result, tt.expected)
			}
		})
	}
}
