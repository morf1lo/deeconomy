package redisrepo

import "fmt"

var (
	guildKey = "guild:%s" // <guildID>
	discordUserKey = "discordUser:%s" // <discordID>
	userGuildsKey = "userGuilds:%s" // <discordID>
	userKey = "user:%s" // <userID>
)

func GuildKey(guildID string) string {
	return fmt.Sprintf(guildKey, guildID)
}

func DiscordUserKey(discordID string) string {
	return fmt.Sprintf(discordUserKey, discordID)
}

func UserGuildsKey(discordID string) string {
	return fmt.Sprintf(userGuildsKey, discordID)
}

func UserKey(id string) string {
	return fmt.Sprintf(userKey, id)
}
