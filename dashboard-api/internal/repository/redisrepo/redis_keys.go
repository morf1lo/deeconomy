package redisrepo

import "fmt"

var (
	guildKey = "guild:%s" // <guildID>
)

func GuildKey(guildID string) string {
	return fmt.Sprintf(guildKey, guildID)
}
