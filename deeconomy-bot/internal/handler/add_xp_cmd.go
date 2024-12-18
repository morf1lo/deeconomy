package handler

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func (h *Handler) AddXPCMD(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if !h.IsAdmin(i) {
		return
	}

	adminID := i.Member.User.ID
	guildID := i.GuildID

	receiver := i.ApplicationCommandData().Options[0].UserValue(s)
	xp := i.ApplicationCommandData().Options[1].IntValue()

	if receiver.Bot {
		h.simpleInteractionReplyWithEphemeral(s, i, "You cannot credite XP to a bot.")
		return
	} else if xp < 20 || xp > 1000 {
		h.simpleInteractionReplyWithEphemeral(s, i, "Minimum XP to credite 20 and maximum 1000.")
		return
	} else if receiver.ID == adminID {
		h.simpleInteractionReplyWithEphemeral(s, i, "You cannot credite XP to yourself.")
		return
	}

	receiverLevel, err := h.services.Level.FindByUserIDAndGuildID(ctx, receiver.ID, guildID)
	if err != nil {
		h.logger.Sugar().Errorf("failed to find receiver(%s) level: %s", receiver.ID, err.Error())
		return
	}
	
	receiverLevel.XP += int(xp)
	isLeveledUp := false
	if receiverLevel.XP >= receiverLevel.Lvl * 100 {
		isLeveledUp = true
		for receiverLevel.XP >= receiverLevel.Lvl * 100 {
			receiverLevel.XP = receiverLevel.XP - receiverLevel.Lvl * 100
			receiverLevel.Lvl += 1
		}
	}

	if err := h.services.Level.Update(ctx, receiverLevel); err != nil {
		h.logger.Sugar().Errorf("failed to update receiver(%s) level: %s", receiver.ID, err.Error())
		return
	}

	if isLeveledUp {
		if err := h.giveRoleToUserIfLeveledUp(s, receiver.ID, guildID, receiverLevel.Lvl); err != nil {
			h.logger.Sugar().Errorf("failed to give role to leveled up user(%s): %s", receiver.ID, err.Error())
		}
	}

	h.simpleInteractionReplyWithEphemeral(s, i, fmt.Sprintf("You have successfully credited %d XP to <@%s>", xp, receiver.ID))
}
