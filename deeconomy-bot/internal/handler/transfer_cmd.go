package handler

import (
	"context"
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/morf1lo/deeconomy-bot/internal/model"
)

func (h *Handler) TransferCMD(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Member.User.Bot {
		return
	}

	senderID := i.Member.User.ID
	guildID := i.GuildID

	receiverOption := i.ApplicationCommandData().Options[0]
	receiver := receiverOption.UserValue(s)
	if receiver.Bot {
		h.simpleInteractionReplyWithEphemeral(s, i, "You cannot transfer money to a bot.")
		return
	}
	receiverID := receiver.ID
	if receiverID == senderID {
		h.simpleInteractionReplyWithEphemeral(s, i, "You cannot transfer money to yourself.")
		return
	}

	amountOption := i.ApplicationCommandData().Options[1]
	amount := amountOption.IntValue()
	if amount < 100 || amount > 5000 {
		h.simpleInteractionReplyWithEphemeral(s, i, "Minimum money to transfer: 100 and maximum: 5.000")
		return
	}

	// Commission is 5%
	reducedAmount := int64(float64(amount) * 0.95)

	senderWallet, err := h.services.Wallet.FindByUserIDAndGuildID(ctx, senderID, guildID)
	if err != nil {
		h.logger.Sugar().Errorf("failed to get sender(%s) wallet: %s", senderID, err.Error())
		return
	}
	if senderWallet.Balance < amount {
		h.simpleInteractionReplyWithEphemeral(s, i, "You don't have enough money to transfer them.")
		return
	}

	_, err = h.services.Wallet.FindByUserIDAndGuildID(ctx, receiverID, guildID)
	if err != nil {
		h.logger.Sugar().Errorf("failed to get receiver(%s) wallet: %s", receiverID, err.Error())
		return
	}

	h.services.Wallet.DecrBy(ctx, senderID, guildID, amount)
	h.services.Wallet.IncrBy(ctx, receiverID, guildID, reducedAmount)

	newTransaction := &model.Transaction{
		GuildID: guildID,
		SenderID: senderID,
		ReceiverID: receiverID,
		Amount: amount,
		ReducedAmount: reducedAmount,
	}
	go func(newTransaction *model.Transaction)  {
		ctx := context.Background()
		if err := h.services.Transaction.New(ctx, newTransaction); err != nil {
			h.logger.Sugar().Errorf("failed to create new transaction with sender(%s) and receiver(%s): %s", senderID, receiver, err.Error())
			return
		}
	}(newTransaction)

	h.simpleInteractionReplyWithEphemeral(s, i, fmt.Sprintf("You have successfully transfered money to <@%s>", receiverID))
}
