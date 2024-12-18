package handler

import (
	"errors"
	"fmt"
	"math/rand"

	"github.com/bwmarrin/discordgo"
)

func (h *Handler) MessageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.Bot || h.cooldownManager.IsOnCooldown(m.Author.ID) {
		return
	}

	userID := m.Author.ID
	guildID := m.GuildID

	randomXP, err := getRandomXP(60, 90)
	if err != nil {
		h.logger.Sugar().Errorf(err.Error())
		return
	}

	level, err := h.services.Level.FindByUserIDAndGuildID(ctx, userID, guildID)
	if err != nil {
		h.logger.Sugar().Errorf("failed to find level by userID(%s) and guildID(%s): %s", userID, guildID, err.Error())
		return
	}

	level.XP += randomXP

	isLeveledUp := level.XP >= level.Lvl * 100
	if isLeveledUp {
		level.XP = level.XP - level.Lvl * 100
		level.Lvl += 1
	}

	if err := h.services.Level.Update(ctx, level); err != nil {
		h.logger.Sugar().Errorf("failed to update existing level: %s", err.Error())
		return
	}

	if isLeveledUp {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("<@%s> you have just leveled up to **level %d**", userID, level.Lvl))

		if err := h.giveRoleToUserIfLeveledUp(s, userID, guildID, level.Lvl); err != nil {
			h.logger.Sugar().Errorf("failed to give role to leveled up user(%s): %s", userID, err.Error())
		}
	}

	h.cooldownManager.Add(userID)
}

func getRandomXP(min int, max int) (int, error) {
	if max <= min {
		return 0, errors.New("the max number cannot be equal to the min or be less than the min")
	}

	return min + rand.Intn(max - min), nil
}
