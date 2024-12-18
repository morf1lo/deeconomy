package handler

import (
	"context"

	"github.com/bwmarrin/discordgo"
	"github.com/morf1lo/deeconomy-bot/internal/lib"
	"github.com/morf1lo/deeconomy-bot/internal/service"
	"go.uber.org/zap"
)

var ctx = context.TODO()

type Handler struct {
	logger *zap.Logger
	services *service.Service
	commands map[string]*SlashCommand
	cooldownManager *lib.CooldownManager
}

func New(logger *zap.Logger, services *service.Service) *Handler {
	return &Handler{
		logger: logger,
		services: services,
		commands: make(map[string]*SlashCommand),
		cooldownManager: lib.NewCooldownManager(0),
	}
}

func (h *Handler) IsAdmin(i *discordgo.InteractionCreate) bool {
	perms := i.Member.Permissions
	return perms&discordgo.PermissionAdministrator != 0
}

func (h *Handler) simpleInteractionReplyWithEphemeral(s *discordgo.Session, i *discordgo.InteractionCreate, content string) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: content,
			Flags: discordgo.MessageFlagsEphemeral,
		},
	})
}

func (h *Handler) giveRoleToUserIfLeveledUp(s *discordgo.Session, userID string, guildID string, level int) error {
	guild, err := h.services.Guild.FindByGuildID(ctx, guildID)
	if err != nil {
		return err
	}

	if guild.LevelRewardingRoles != nil {
		if roleID, exists := guild.LevelRewardingRoles[level]; exists {
			role, err := s.State.Role(guildID, roleID)
			if err != nil {
				h.logger.Sugar().Errorf("failed to fetch role by ID(%s): %s", roleID, err.Error())
				return err
			}

			for lvl, previousRoleID := range guild.LevelRewardingRoles {
				if lvl < level {
					if err := s.GuildMemberRoleRemove(guildID, userID, previousRoleID); err != nil {
						h.logger.Sugar().Error("failed to remove previousRoleID role(%s) from user(%s): %s", previousRoleID, userID, err.Error())
						return err
					}
				}
			}

			if err := s.GuildMemberRoleAdd(guildID, userID, role.ID); err != nil {
				h.logger.Sugar().Errorf("failed to add a role(%s) to member(%s): %s", role.ID, userID, err.Error())
			}
		}
	}

	return nil
}
