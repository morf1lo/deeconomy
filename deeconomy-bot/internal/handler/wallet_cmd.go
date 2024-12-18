package handler

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func (h *Handler) WalletCMD(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Member.User.Bot {
		return
	}

	userID := i.Member.User.ID
	guildID := i.GuildID

	wallet, err := h.services.Wallet.FindByUserIDAndGuildID(ctx, userID, guildID)
	if err != nil {
		h.logger.Sugar().Errorf("failed to find wallet by userID(%s) and guildID(%s): %s", userID, guildID, err.Error())
		return
	}

	h.simpleInteractionReplyWithEphemeral(s, i, fmt.Sprintf("Your balance: %d", wallet.Balance))
}
