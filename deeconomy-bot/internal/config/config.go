package config

type BotConfig struct {
	Token string
	Debug bool
}

type MongoDBConfig struct {
	URI string
	DBName string
}
