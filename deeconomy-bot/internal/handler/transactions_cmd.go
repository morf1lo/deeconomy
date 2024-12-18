package handler

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/morf1lo/deeconomy-bot/internal/service"
)

func (h *Handler) TransactionsCMD(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Member.User.Bot {
		return
	}

	userID := i.Member.User.ID
	guildID := i.GuildID

	options := i.ApplicationCommandData().Options
	scope := service.DefaultScope
	if len(options) > 0 {
		scope = options[0].StringValue()
	}

	transactions, err := h.services.Transaction.FindUserTransactionsInGuild(ctx, userID, guildID, scope)
	if err != nil {
		h.logger.Sugar().Errorf("failed to find user(%s) transactions in guild(%s): %s", userID, guildID, err.Error())
		return
	}

	msg := ""
	for _, transaction := range transactions {
		if transaction.SenderID == userID {
			msg += fmt.Sprintf(
				"--------------------\n**Date**: %s\n**Sender**: `You`\n**Receiver**: <@%s>\n**You**:```diff\n- %d```**Receiver**:```diff\n+ %d```\n--------------------\n",
				transaction.CreatedAt.Format("2006-01-02 15:04"),
				transaction.ReceiverID,
				transaction.Amount,
				transaction.ReducedAmount,
			)
		} else {
			msg += fmt.Sprintf(
				"--------------------\n**Date**: %s\n**Sender**: <@%s>\n**Receiver**: `You`\n**You**:```diff\n+ %d```**Sender**:```diff\n- %d```\n--------------------\n",
				transaction.CreatedAt.Format("2006-01-02 15:04"),
				transaction.SenderID,
				transaction.ReducedAmount,
				transaction.Amount,
			)
		}
	}

	h.simpleInteractionReplyWithEphemeral(s, i, msg)
}
