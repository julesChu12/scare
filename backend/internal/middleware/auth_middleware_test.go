package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"community-elderly-care-platform/pkg/jwt"

	"github.com/gin-gonic/gin"
)

func TestAuthMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)
	jwtManager := jwt.NewManager("test-secret", 1, 2)

	makeRequest := func(authHeader string) *httptest.ResponseRecorder {
		t.Helper()

		r := gin.New()
		r.GET("/secure",
			AuthMiddleware(jwtManager, nil),
			func(c *gin.Context) {
				userID, _ := c.Get("user_id")
				identities, _ := c.Get("user_identities")
				userType, _ := c.Get("user_type")
				stationID, _ := c.Get("station_id")
				tokenID, _ := c.Get("token_id")
				tokenExpiresAt, _ := c.Get("token_expires_at")
				c.JSON(http.StatusOK, gin.H{
					"user_id":           userID,
					"user_type":         userType,
					"station_id":        stationID,
					"identities":        identities,
					"token_id_exists":   tokenID != "",
					"expires_at_is_set": tokenExpiresAt.(time.Time).IsZero() == false,
				})
			},
		)

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/secure", nil)
		if authHeader != "" {
			req.Header.Set("Authorization", authHeader)
		}
		r.ServeHTTP(w, req)
		return w
	}

	if w := makeRequest(""); w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 for missing authorization, got %d", w.Code)
	}

	if w := makeRequest("Token abc"); w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 for invalid authorization format, got %d", w.Code)
	}

	if w := makeRequest("Bearer invalid-token"); w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 for invalid token, got %d", w.Code)
	}

	token, err := jwtManager.GenerateToken(42, "b_end", 9, []string{"staff"}, "staff")
	if err != nil {
		t.Fatalf("failed to generate test token: %v", err)
	}
	if w := makeRequest("Bearer " + token); w.Code != http.StatusOK {
		t.Fatalf("expected 200 for valid token, got %d", w.Code)
	}
}
