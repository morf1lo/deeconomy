package main

import (
	"context"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/morf1lo/deeconomy-bot/internal/config"
	"github.com/morf1lo/deeconomy-bot/internal/db"
	"github.com/morf1lo/deeconomy-bot/internal/handler"
	"github.com/morf1lo/deeconomy-bot/internal/repository"
	"github.com/morf1lo/deeconomy-bot/internal/service"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.uber.org/zap"
)

var ctx = context.Background()

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	if err := initEnv(); err != nil {
		logger.Sugar().Fatalf("failed to load environment variables: %s", err.Error())
	}

	if err := initConfig(); err != nil {
		logger.Sugar().Fatalf("failed to initialize yaml config: %s", err.Error())
	}

	mongoConfig := &config.MongoDBConfig{
		URI: os.Getenv("MONGO_URI"),
		DBName: os.Getenv("MONGO_DBNAME"),
	}
	mongoDB, err := db.NewMongoDB(ctx, mongoConfig)
	if err != nil {
		logger.Sugar().Fatalf("failed to connect to mongoDB: %s", err.Error())
	}
	if err := mongoDB.Client().Ping(ctx, &readpref.ReadPref{}); err != nil {
		logger.Sugar().Fatalf("failed to ping mongoDB: %s", err.Error())
	}
	logger.Info("Successfully connected to mongoDB")

	redisOptions := &redis.Options{
		Addr: os.Getenv("REDIS_ADDR"),
		Username: os.Getenv("REDIS_USERNAME"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB: 0,
		Protocol: 3,
		ReadTimeout: time.Second * 5,
		WriteTimeout: time.Second * 5,
	}
	rdb := redis.NewClient(redisOptions)
	pong, err := rdb.Ping(ctx).Result()
	if err != nil {
		logger.Sugar().Fatalf("failed to ping Redis: %s", err.Error())
	}
	logger.Sugar().Infof("Successfully connected to Redis: %s", pong)

	repos := repository.New(logger, mongoDB, rdb)
	services := service.New(logger, repos)
	handlers := handler.New(logger, services)

	botConfig := &config.BotConfig{
		Token: os.Getenv("BOT_TOKEN"),
		Debug: viper.GetBool("bot.debug"),
	}
	bot := handler.NewBot(logger, handlers)
	bot.Start(botConfig)
}

func initEnv() error {
	return godotenv.Load()
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigType("yaml")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
