package dto

import "github.com/bwmarrin/discordgo"

type DiscordError struct {
	Message string `json:"message"`
	Code int `json:"code"`
}

type DiscordOAuth2Resp struct {
	Message string `json:"message"`
	AccessToken string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type DiscordUserGuildsResp struct {
	Guilds []*discordgo.Guild `json:"guilds"`
}

type DiscordRevokeResp struct {
	Message string `json:"message"`
	AccessToken string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
