package main

import (
	"context"
	"os"
	"time"

	"github.com/joho/godotenv"
	v1 "github.com/morf1lo/deeconomy-bot-api/internal/api/v1"
	"github.com/morf1lo/deeconomy-bot-api/internal/config"
	"github.com/morf1lo/deeconomy-bot-api/internal/repository"
	"github.com/morf1lo/deeconomy-bot-api/internal/repository/mongorepo"
	"github.com/morf1lo/deeconomy-bot-api/internal/service"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var ctx = context.Background()

func main() {
	logger, _ := zap.NewProduction()

	if err := initEnv(); err != nil {
		logger.Sugar().Fatalf("failed to load environment variables: %s", err.Error())
	}

	if err := initConfig(); err != nil {
		logger.Sugar().Fatalf("failed to initialize yaml config: %s", err.Error())
	}

	mongoCfg := &config.MongoConfig{
		URI: os.Getenv("MONGO_URI"),
		DBName: os.Getenv("MONGO_DBNAME"),
	}
	db, err := mongorepo.NewMongo(ctx, mongoCfg)
	if err != nil {
		logger.Sugar().Fatalf("failed to connect to MongoDB: %s", err.Error())
	}
	if err := db.Client().Ping(ctx, nil); err != nil {
		logger.Sugar().Fatalf("failed to ping Mongo Client: %s", err.Error())
	}
	logger.Info("Successfully connected to MongoDB")

	redisOpts := &redis.Options{
		Addr: os.Getenv("REDIS_ADDR"),
		Username: os.Getenv("REDIS_USERNAME"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB: 0,
		Protocol: 3,
		ReadTimeout: time.Second * 10,
		WriteTimeout: time.Second * 10,
	}
	rdb := redis.NewClient(redisOpts)
	pong, err := rdb.Ping(ctx).Result()
	if err != nil {
		logger.Sugar().Fatalf("failed to ping Redis: %s", err.Error())
	}
	logger.Sugar().Infof("Successfully connected to Redis: %s", pong)

	repos := repository.New(db, rdb)
	services := service.New(logger, repos)
	handlers := v1.NewHandler(logger, services)

	if err := handlers.InitRoutes().Run(viper.GetString("app.port")); err != nil {
		logger.Sugar().Fatalf("failed to start routing: %s", err.Error())
	}

	logger.Info("Routing started")
}

func initEnv() error {
	return godotenv.Load()
}

func initConfig() error {
	viper.AddConfigPath("./")
	viper.SetConfigType("yaml")
	viper.SetConfigName("app")
	return viper.ReadInConfig()
}
