package main

import (
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func main() {
	logger, _ := zap.NewProduction()

	if err := initEnv(); err != nil {
		logger.Sugar().Fatalf("failed to load environment variables: %s", err.Error())
	}

	if err := initConfig(); err != nil {
		logger.Sugar().Fatalf("failed to initialize yaml config: %s", err.Error())
	}

	
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
