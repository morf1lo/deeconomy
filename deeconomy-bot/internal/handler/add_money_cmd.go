package handler

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func (h *Handler) AddMoneyCMD(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if !h.IsAdmin(i) {
		return
	}

	adminID := i.Member.User.ID
	guildID := i.GuildID

	receiver := i.ApplicationCommandData().Options[0].UserValue(s)
	amount := i.ApplicationCommandData().Options[1].IntValue()

	if receiver.Bot {
		h.simpleInteractionReplyWithEphemeral(s, i, "You cannot credite money to a bot.")
		return
	} else if amount < 50 || amount > 500 {
		h.simpleInteractionReplyWithEphemeral(s, i, "Minimum amount to credite 50 and maximum 500.")
		return
	} else if receiver.ID == adminID {
		h.simpleInteractionReplyWithEphemeral(s, i, "You cannot credite money to yourself.")
		return
	}

	receiverWallet, err := h.services.Wallet.FindByUserIDAndGuildID(ctx, receiver.ID, guildID)
	if err != nil {
		h.logger.Sugar().Errorf("failed to find receiver(%s) wallet: %s", receiver.ID, err.Error())
		return
	}

	receiverWallet.Balance += amount

	if err := h.services.Wallet.Update(ctx, receiverWallet); err != nil {
		h.logger.Sugar().Errorf("failed to credite money to receiver(%s) wallet: %s", receiver.ID, err.Error())
		return
	}

	h.simpleInteractionReplyWithEphemeral(s, i, fmt.Sprintf("You have successfully credited %d money to <@%s>", amount, receiver.ID))
}
