package messaging

import "os"

type Config struct {
	LineChannelSecret string
	LineChannelToken  string
}

func NewConfig() *Config {
	return &Config{
		LineChannelSecret: os.Getenv("LINE_CHANNEL_SECRET"),
		LineChannelToken:  os.Getenv("LINE_CHANNEL_TOKEN"),
	}
}
