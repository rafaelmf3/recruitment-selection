package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	"recruitment-selection/internal/model"
	"recruitment-selection/internal/token"
)

const (
	ContextKeyUserID = "userID"
	ContextKeyRole   = "userRole"
)

// RequireAuth validates the Bearer JWT and stores userID + role in the Gin context.
func RequireAuth(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing or invalid authorization header"})
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		claims := &token.Claims{}

		parsed, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
			return []byte(jwtSecret), nil
		})
		if err != nil || !parsed.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
			return
		}

		userID, err := uuid.Parse(claims.UserID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "malformed token claims"})
			return
		}

		c.Set(ContextKeyUserID, userID)
		c.Set(ContextKeyRole, claims.Role)
		c.Next()
	}
}

// RequireRole aborts with 403 if the authenticated user does not have the required role.
func RequireRole(role model.UserRole) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, _ := c.Get(ContextKeyRole)
		if userRole != role {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "insufficient permissions"})
			return
		}
		c.Next()
	}
}

// InjectTestUser bypasses JWT validation in handler tests by injecting
// a fixed userID and role directly into the Gin context.
func InjectTestUser(userID uuid.UUID, role model.UserRole) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(ContextKeyUserID, userID)
		c.Set(ContextKeyRole, role)
		c.Next()
	}
}
