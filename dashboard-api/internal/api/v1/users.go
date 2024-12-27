package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) usersMe(c *gin.Context) {
	user := h.getUser(c)

	discordUser, err := h.services.User.FindDiscordUser(c.Request.Context(), user.DiscordID, user.DiscordAccessToken, true)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, *discordUser)
}
