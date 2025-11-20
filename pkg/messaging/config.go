package messaging

type Config struct {
	LineChannelSecret string
	LineChannelToken  string
}

func (c *Config) MessagingConfig() *Config {
	return c
}
