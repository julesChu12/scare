package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestRequireEndType(t *testing.T) {
	gin.SetMode(gin.TestMode)

	makeRequest := func(setup func(*gin.Context)) *httptest.ResponseRecorder {
		t.Helper()

		r := gin.New()
		r.GET("/check",
			func(c *gin.Context) {
				if setup != nil {
					setup(c)
				}
				c.Next()
			},
			RequireEndType("b_end"),
			func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"ok": true})
			},
		)

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/check", nil)
		r.ServeHTTP(w, req)
		return w
	}

	if w := makeRequest(nil); w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 for missing user_type, got %d", w.Code)
	}

	if w := makeRequest(func(c *gin.Context) { c.Set("user_type", 123) }); w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500 for invalid user_type type, got %d", w.Code)
	}

	if w := makeRequest(func(c *gin.Context) { c.Set("user_type", "c_end") }); w.Code != http.StatusForbidden {
		t.Fatalf("expected 403 for end type mismatch, got %d", w.Code)
	}

	if w := makeRequest(func(c *gin.Context) { c.Set("user_type", "b_end") }); w.Code != http.StatusOK {
		t.Fatalf("expected 200 for matching end type, got %d", w.Code)
	}
}
