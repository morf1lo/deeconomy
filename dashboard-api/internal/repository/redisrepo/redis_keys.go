package redisrepo

import "fmt"

var (
	guildKey = "guild:%s" // <guildID>
	userInfoKey = "userInfo:%s" // <discordID>
	userGuildsKey = "userGuilds:%s" // <discordID>
)

func GuildKey(guildID string) string {
	return fmt.Sprintf(guildKey, guildID)
}

func UserInfoKey(discordID string) string {
	return fmt.Sprintf(userInfoKey, discordID)
}

func UserGuildsKey(discordID string) string {
	return fmt.Sprintf(userGuildsKey, discordID)
}
