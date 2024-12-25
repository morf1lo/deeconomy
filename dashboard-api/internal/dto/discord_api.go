package dto

import "github.com/bwmarrin/discordgo"

type DiscordOAuth2Resp struct {
	Message string `json:"message"`
	AccessToken string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type DiscordUserResp struct {
	Message string `json:"message"`
	User *discordgo.User `json:"user"`
}

type DiscordUserGuildsResp struct {
	Guilds []*discordgo.Guild `json:"guilds"`
}
