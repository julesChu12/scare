package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// RequireEndType ensures the authenticated JWT token belongs to the expected end type.
// expectedType: "c_end" or "b_end".
//
// This is intentionally lightweight and independent from Casbin. It prevents tokens
// issued for one end from accessing the other end's routes.
func RequireEndType(expectedType string) gin.HandlerFunc {
	return func(c *gin.Context) {
		value, ok := c.Get("user_type")
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"msg":  "missing user type",
				"data": nil,
			})
			return
		}

		userType, ok := value.(string)
		if !ok || userType == "" {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"msg":  "invalid user type",
				"data": nil,
			})
			return
		}

		if userType != expectedType {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"msg":  "token end type mismatch",
				"data": gin.H{"expected": expectedType, "got": userType},
			})
			return
		}

		c.Next()
	}
}
