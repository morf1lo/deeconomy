package redisrepo

import "fmt"

var (
	walletKey = "wallet:%s:%s" // <userID>:<guildID>
	levelKey = "level:%s:%s" // <userID>:<guildID>
	guildKey = "guild:%s" // <guildID>
	transactionsKey = "%s-%s-%s:transactions" // <userID>:<guildID>:<scope>
)

func WalletKey(userID string, guildID string) string {
	return fmt.Sprintf(walletKey, userID, guildID)
}

func LevelKey(userID string, guildID string) string {
	return fmt.Sprintf(levelKey, userID, guildID)
}

func GuildKey(guildID string) string {
	return fmt.Sprintf(guildKey, guildID)
}

func TransactionsKey(userID string, guildID string, scope string) string {
	return fmt.Sprintf(transactionsKey, userID, guildID, scope)
}
