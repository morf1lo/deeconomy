package v1

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func (h *Handler) AuthMiddleware(c *gin.Context) {
	header := c.GetHeader("Authorization")
	if !strings.HasPrefix(header, "Bearer ") {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "not authorized"})
		c.Abort()
		return
	}

	accessToken := strings.Split(header, " ")[1]
	if accessToken == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "not authorized"})
		c.Abort()
		return
	}

	user, err := h.getUserDataFromTokenClaims(c.Request.Context(), accessToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		c.Abort()
		return
	}

	c.Set("user", *user)

	c.Next()
}
