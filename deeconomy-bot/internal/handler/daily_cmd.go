package handler

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	"go.mongodb.org/mongo-driver/mongo"
)

func (h *Handler) DailyCMD(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Member.User.Bot {
		return
	}

	userID := i.Member.User.ID
	guildID := i.GuildID

	wallet, err := h.services.Wallet.FindByUserIDAndGuildID(ctx, userID, guildID)
	if err != nil && err != mongo.ErrNoDocuments {
		h.logger.Sugar().Errorf("failed to find wallet by userID(%s) and guildID(%s): %s", userID, guildID, err.Error())
		return
	}

	nextEligibleTime := wallet.LastDailyCollected.Add(time.Hour * 24)
	if time.Now().Before(nextEligibleTime) {
		h.simpleInteractionReplyWithEphemeral(s, i, fmt.Sprintf(
			"You have already collected your dailies. Come back at %s.",
			nextEligibleTime.Format("2006-01-02 15:04"),
		))
		return
	}

	guild, err := h.services.Guild.FindByGuildID(ctx, guildID)
	if err != nil && err != mongo.ErrNoDocuments {
		h.logger.Sugar().Errorf("failed to find guild(%s): %s", guildID, err.Error())
		return
	}

	if err == mongo.ErrNoDocuments {
		wallet.Balance += 100
	} else {
		wallet.Balance += guild.DailyAmount
	}

	wallet.LastDailyCollected = time.Now()

	if err := h.services.Wallet.Update(ctx, wallet); err != nil {
		h.logger.Sugar().Errorf("failed to update user(%s) wallet: %s", userID, err.Error())
		return
	}

	h.simpleInteractionReplyWithEphemeral(s, i, fmt.Sprintf("You have collected your dailies! Now your balance: %d", wallet.Balance))
}
