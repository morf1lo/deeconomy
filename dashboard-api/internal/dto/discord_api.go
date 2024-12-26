package dto

import "github.com/bwmarrin/discordgo"

type DiscordOAuth2Resp struct {
	Message string `json:"message"`
	AccessToken string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type DiscordUserGuildsResp struct {
	Guilds []*discordgo.Guild `json:"guilds"`
}

type DiscordError struct {
	Message string `json:"message"`
	Code int `json:"code"`
}
