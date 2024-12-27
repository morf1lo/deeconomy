package v1

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

func (h *Handler) oauth2SignIn(c *gin.Context) {
	host := fmt.Sprintf(
		"https://discord.com/oauth2/authorize?client_id=%s&response_type=code&scope=guilds+identify&redirect_uri=%s",
		os.Getenv("CLIENT_ID"),
		os.Getenv("DISCORD_REDIRECT_URI"),
	)
	c.Redirect(http.StatusTemporaryRedirect, host)
}

func (h *Handler) oauth2Authorize(c *gin.Context) {
	code := strings.TrimSpace(c.Query("code"))
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "authorization code is not provided"})
		return
	}

	accessToken, refreshToken, err := h.services.User.Authorize(c.Request.Context(), code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.SetCookie("refreshToken", refreshToken, 3600 * 24 * 7, "/", "localhost", true, true)

	c.JSON(http.StatusOK, gin.H{"ok": true, "accessToken": accessToken})
}

func (h *Handler) oauth2Refresh(c *gin.Context) {
	refreshToken, err := c.Cookie("refreshToken")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	newAccessToken, newRefreshToken, err := h.services.User.RefreshTokens(c.Request.Context(), refreshToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.SetCookie("refreshToken", newRefreshToken, 3600 * 24 * 7, "/", "localhost", true, true)

	c.JSON(http.StatusOK, gin.H{"ok": true, "accessToken": newAccessToken})
}
