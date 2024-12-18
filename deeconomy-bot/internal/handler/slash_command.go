package handler

import (
	"github.com/bwmarrin/discordgo"
	"github.com/spf13/viper"
)

type SlashCommand struct {
	Name string
	Handler func(s *discordgo.Session, i *discordgo.InteractionCreate)
}

func (h *Handler) RegisterDefaultCommands() {
	h.Register(&SlashCommand{
		Name: viper.GetString("bot.commands.wallet.name"),
		Handler: h.WalletCMD,
	})

	h.Register(&SlashCommand{
		Name: viper.GetString("bot.commands.daily.name"),
		Handler: h.DailyCMD,
	})

	h.Register(&SlashCommand{
		Name: viper.GetString("bot.commands.level.name"),
		Handler: h.LevelCMD,
	})

	h.Register(&SlashCommand{
		Name: viper.GetString("bot.commands.transfer.name"),
		Handler: h.TransferCMD,
	})

	h.Register(&SlashCommand{
		Name: viper.GetString("bot.commands.transactions.name"),
		Handler: h.TransactionsCMD,
	})

	h.Register(&SlashCommand{
		Name: viper.GetString("bot.commands.add-xp.name"),
		Handler: h.AddXPCMD,
	})

	h.Register(&SlashCommand{
		Name: viper.GetString("bot.commands.add-money.name"),
		Handler: h.AddMoneyCMD,
	})
}

func (h *Handler) Register(cmd *SlashCommand) {
	h.commands[cmd.Name] = cmd
}

func (h *Handler) InteractionHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if cmd, exists := h.commands[i.ApplicationCommandData().Name]; exists {
		cmd.Handler(s, i)
	} else {
		h.logger.Sugar().Warnf("Unhandled command: %s", i.ApplicationCommandData().Name)
	}
}
