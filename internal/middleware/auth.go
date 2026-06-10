package middleware

import (
	"net/http"
	"strings"

	"task-manager-api/internal/auth"

	"github.com/gin-gonic/gin"
)

const UserIDKey = "userID"

// RequireAuth accepts a bearer token in the Authorization header or the
// httpOnly "token" cookie set at login.
func RequireAuth(jwt *auth.JWTManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := ""

		if header := c.GetHeader("Authorization"); strings.HasPrefix(header, "Bearer ") {
			tokenString = strings.TrimPrefix(header, "Bearer ")
		} else if cookie, err := c.Cookie("token"); err == nil {
			tokenString = cookie
		}

		if tokenString == "" {
			abortUnauthorized(c, "Authentication required")
			return
		}

		userID, err := jwt.Verify(tokenString)
		if err != nil {
			abortUnauthorized(c, "Invalid or expired token")
			return
		}

		c.Set(UserIDKey, userID)
		c.Next()
	}
}

func abortUnauthorized(c *gin.Context, message string) {
	c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
		"error": gin.H{"code": "UNAUTHORIZED", "message": message},
	})
}

func CurrentUserID(c *gin.Context) uint {
	return c.GetUint(UserIDKey)
}
