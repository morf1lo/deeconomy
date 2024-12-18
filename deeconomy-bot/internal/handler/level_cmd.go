package handler

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func (h *Handler) LevelCMD(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Member.User.Bot {
		return
	}

	userID := i.Member.User.ID
	guildID := i.GuildID

	level, err := h.services.Level.FindByUserIDAndGuildID(ctx, userID, guildID)
	if err != nil {
		h.logger.Sugar().Errorf("failed to find user(%s) level in guild(%s): %s", userID, guildID, err.Error())
		return
	}

	h.simpleInteractionReplyWithEphemeral(s, i, fmt.Sprintf("LEVEL: %d\nXP: %d", level.Lvl, level.XP))
}
