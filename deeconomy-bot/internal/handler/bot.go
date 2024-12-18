package handler

import (
	"github.com/bwmarrin/discordgo"
	"github.com/morf1lo/deeconomy-bot/internal/config"
	// "github.com/spf13/viper"
	"go.uber.org/zap"
)

type Bot struct {
	logger   *zap.Logger
	handlers *Handler
}

func NewBot(logger *zap.Logger, handlers *Handler) *Bot {
	return &Bot{
		logger: logger,
		handlers: handlers,
	}
}

func (b *Bot) Start(cfg *config.BotConfig) {
	s, err := discordgo.New("Bot " + cfg.Token)
	if err != nil {
		b.logger.Sugar().Fatalf("failed to create bot using discordgo: %s", err.Error())
	}

	if err := s.Open(); err != nil {
		b.logger.Sugar().Fatalf("failed to open discord session: %s", err.Error())
	}
	
	b.handlers.RegisterDefaultCommands()

	s.AddHandler(b.handlers.InteractionHandler)
	s.AddHandler(b.handlers.MessageHandler)

	// Delete all previous slash commands
	// allCommands, err := s.ApplicationCommands(s.State.User.ID, "")
	// if err != nil {
	// 	b.logger.Sugar().Fatalf("failed to get all application commands: %s", err.Error())
	// }
	// for _, cmd := range allCommands {
	// 	if err := s.ApplicationCommandDelete(s.State.User.ID, "", cmd.ID); err != nil {
	// 		b.logger.Sugar().Fatalf("failed to delete application command(%s): %s", cmd.Name, err.Error())
	// 	}
	// }

	// Registering all slash commands to bot
	// commands := []*discordgo.ApplicationCommand{
	// 	{
	// 		Name: viper.GetString("bot.commands.wallet.name"),
	// 		Description: viper.GetString("bot.commands.wallet.desc"),
	// 	},
	// 	{
	// 		Name: viper.GetString("bot.commands.daily.name"),
	// 		Description: viper.GetString("bot.commands.daily.desc"),
	// 	},
	// 	{
	// 		Name: viper.GetString("bot.commands.level.name"),
	// 		Description: viper.GetString("bot.commands.level.desc"),
	// 	},
	// 	{
	// 		Name: viper.GetString("bot.commands.transfer.name"),
	// 		Description: viper.GetString("bot.commands.transfer.desc"),
	// 		Options: []*discordgo.ApplicationCommandOption{
	// 			{
	// 				Type: discordgo.ApplicationCommandOptionMentionable,
	// 				Name: viper.GetString("bot.commands.transfer.options.receiver.name"),
	// 				Description: viper.GetString("bot.commands.transfer.options.receiver.desc"),
	// 				Required: viper.GetBool("bot.commands.transfer.options.receiver.required"),
	// 			},
	// 			{
	// 				Type: discordgo.ApplicationCommandOptionInteger,
	// 				Name: viper.GetString("bot.commands.transfer.options.amount.name"),
	// 				Description: viper.GetString("bot.commands.transfer.options.amount.desc"),
	// 				Required: viper.GetBool("bot.commands.transfer.options.amount.required"),
	// 			},
	// 		},
	// 	},
	// 	{
	// 		Name: viper.GetString("bot.commands.transactions.name"),
	// 		Description: viper.GetString("bot.commands.transactions.desc"),
	// 		Options: []*discordgo.ApplicationCommandOption{
	// 			{
	// 				Type: discordgo.ApplicationCommandOptionString,
	// 				Name: viper.GetString("bot.commands.transactions.options.scope.name"),
	// 				Description: viper.GetString("bot.commands.transactions.options.scope.desc"),
	// 				Required: viper.GetBool("bot.commands.transactions.options.scope.required"),
	// 				Choices: []*discordgo.ApplicationCommandOptionChoice{
	// 					{
	// 						Name: viper.GetString("bot.commands.transactions.options.scope.choices.senderOnly.name"),
	// 						Value: viper.GetString("bot.commands.transactions.options.scope.choices.senderOnly.value"),
	// 					},
	// 					{
	// 						Name: viper.GetString("bot.commands.transactions.options.scope.choices.receiverOnly.name"),
	// 						Value: viper.GetString("bot.commands.transactions.options.scope.choices.receiverOnly.value"),
	// 					},
	// 				},
	// 			},
	// 		},
	// 	},
	// 	{
	// 		Name: viper.GetString("bot.commands.add-xp.name"),
	// 		Description: viper.GetString("bot.commands.add-xp.desc"),
	// 		Options: []*discordgo.ApplicationCommandOption{
	// 			{
	// 				Type: discordgo.ApplicationCommandOptionMentionable,
	// 				Name: viper.GetString("bot.commands.add-xp.options.user.name"),
	// 				Description: viper.GetString("bot.commands.add-xp.options.user.desc"),
	// 				Required: viper.GetBool("bot.commands.add-xp.options.user.required"),
	// 			},
	// 			{
	// 				Type: discordgo.ApplicationCommandOptionInteger,
	// 				Name: viper.GetString("bot.commands.add-xp.options.xp.name"),
	// 				Description: viper.GetString("bot.commands.add-xp.options.xp.desc"),
	// 				Required: viper.GetBool("bot.commands.add-xp.options.xp.required"),
	// 			},
	// 		},
	// 	},
	// 	{
	// 		Name: viper.GetString("bot.commands.add-money.name"),
	// 		Description: viper.GetString("bot.commands.add-money.desc"),
	// 		Options: []*discordgo.ApplicationCommandOption{
	// 			{
	// 				Type: discordgo.ApplicationCommandOptionMentionable,
	// 				Name: viper.GetString("bot.commands.add-money.options.user.name"),
	// 				Description: viper.GetString("bot.commands.add-money.options.user.desc"),
	// 				Required: viper.GetBool("bot.commands.add-money.options.user.required"),
	// 			},
	// 			{
	// 				Type: discordgo.ApplicationCommandOptionInteger,
	// 				Name: viper.GetString("bot.commands.add-money.options.amount.name"),
	// 				Description: viper.GetString("bot.commands.add-money.options.amount.desc"),
	// 				Required: viper.GetBool("bot.commands.add-money.options.amount.required"),
	// 			},
	// 		},
	// 	},
	// }
	// for _, cmd := range commands {
	// 	_, err := s.ApplicationCommandCreate(s.State.User.ID, "", cmd)
	// 	if err != nil {
	// 		b.logger.Sugar().Fatalf("failed to register command %s: %s", cmd.Name, err.Error())
	// 	}
	// }

	b.logger.Sugar().Info("Bot is running")

	select {}
}
