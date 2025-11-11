package messaging

import "os"

type Config struct {
	LineChannelSecret string
}

func NewConfig() *Config {
	return &Config{
		LineChannelSecret: os.Getenv("LINE_CHANNEL_SECRET"),
	}
}
