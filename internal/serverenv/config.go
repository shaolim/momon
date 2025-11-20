package serverenv

import (
	"github/shaolim/momon/pkg/database"
	"github/shaolim/momon/pkg/messaging"
	"os"
)

type LineMessagingConfigProvider interface {
	MessagingConfig() *messaging.Config
}

type DatabaseConfigProvider interface {
	DatabaseConfig() *database.Config
}

type Config struct {
	Messaging LineMessagingConfigProvider
	Database  DatabaseConfigProvider
	Host      string
}

func LoadEnv() *Config {
	messagingConfig := &messaging.Config{
		LineChannelSecret: os.Getenv("LINE_CHANNEL_SECRET"),
		LineChannelToken:  os.Getenv("LINE_CHANNEL_TOKEN"),
	}

	databaseConfig := &database.Config{
		Name:     os.Getenv("DB_NAME"),
		User:     os.Getenv("DB_USER"),
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Password: os.Getenv("DB_PASSWORD"),
	}

	return &Config{
		Messaging: messagingConfig,
		Database:  databaseConfig,
		Host:      os.Getenv("HTTP_PORT"),
	}
}
