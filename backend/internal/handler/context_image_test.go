package handler

import (
	"crypto/tls"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func newTestContextWithHost(host string) *gin.Context {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	req, _ := http.NewRequest("GET", "http://"+host+"/api/v1/c/requests/12", nil)
	c.Request = req
	return c
}

func TestImageHost_EmptyAndNull(t *testing.T) {
	cases := []struct {
		name  string
		input string
	}{
		{"empty string", ""},
		{"null string", "null"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			c := newTestContextWithHost("192.168.1.2:5174")
			got := imageHost(c, tc.input)
			if got != tc.input {
				t.Errorf("imageHost() = %q, want %q", got, tc.input)
			}
		})
	}
}

func TestImageHost_RelativePath(t *testing.T) {
	c := newTestContextWithHost("192.168.1.2:5174")
	images := `["/static/c_end/20260418/xxx.jpg"]`
	got := imageHost(c, images)
	var gotPaths []string
	json.Unmarshal([]byte(got), &gotPaths)
	want := "http://192.168.1.2:5174/static/c_end/20260418/xxx.jpg"
	if len(gotPaths) != 1 || gotPaths[0] != want {
		t.Errorf("imageHost() = %v, want [%v]", gotPaths, want)
	}
}

func TestImageHost_OldLocalhostURL(t *testing.T) {
	c := newTestContextWithHost("192.168.1.2:5174")
	images := `["http://localhost:8080/static/c_end/20260418/xxx.jpg"]`
	got := imageHost(c, images)
	var gotPaths []string
	json.Unmarshal([]byte(got), &gotPaths)
	want := "http://192.168.1.2:5174/static/c_end/20260418/xxx.jpg"
	if len(gotPaths) != 1 || gotPaths[0] != want {
		t.Errorf("imageHost() = %v, want [%v]", gotPaths, want)
	}
}

func TestImageHost_HttpsHost(t *testing.T) {
	c := newTestContextWithHost("example.com")
	c.Request.TLS = &tls.ConnectionState{} // 模拟 HTTPS
	images := `["/static/b_end/20260418/banner.png"]`
	got := imageHost(c, images)
	var gotPaths []string
	json.Unmarshal([]byte(got), &gotPaths)
	want := "https://example.com/static/b_end/20260418/banner.png"
	if len(gotPaths) != 1 || gotPaths[0] != want {
		t.Errorf("imageHost() = %v, want [%v]", gotPaths, want)
	}
}

func TestImageHost_InvalidJSON(t *testing.T) {
	c := newTestContextWithHost("192.168.1.2:5174")
	got := imageHost(c, "not valid json")
	if got != "not valid json" {
		t.Errorf("imageHost() = %q, want %q", got, "not valid json")
	}
}

func TestImageHost_MultipleImages(t *testing.T) {
	c := newTestContextWithHost("192.168.1.2:5174")
	images := `["/static/img1.jpg","http://localhost:8080/static/img2.png"]`
	got := imageHost(c, images)
	var gotPaths []string
	json.Unmarshal([]byte(got), &gotPaths)
	if len(gotPaths) != 2 {
		t.Fatalf("len = %d, want 2", len(gotPaths))
	}
	if gotPaths[0] != "http://192.168.1.2:5174/static/img1.jpg" {
		t.Errorf("gotPaths[0] = %q", gotPaths[0])
	}
	if gotPaths[1] != "http://192.168.1.2:5174/static/img2.png" {
		t.Errorf("gotPaths[1] = %q", gotPaths[1])
	}
}
