package storage

import "testing"

func TestExtractPathPrefix(t *testing.T) {
	cases := []struct {
		name   string
		input  string
		output string
	}{
		{"完整URL带端口", "http://localhost:8080/static", "/static"},
		{"https带端口", "https://example.com:8443/uploads", "/uploads"},
		{"无端口URL", "http://localhost/static", "/static"},
		{"根路径", "http://localhost/", "/"},
		{"多层路径", "https://cdn.example.com/media/images", "/media/images"},
		{"空字符串", "", "/static"},
		{"仅路径无域名", "/static", "/static"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := extractPathPrefix(tc.input)
			if got != tc.output {
				t.Errorf("extractPathPrefix(%q) = %q, want %q", tc.input, got, tc.output)
			}
		})
	}
}
