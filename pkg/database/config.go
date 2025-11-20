package database

import (
	"net/url"
)

type Config struct {
	Name     string `json:"name"`
	User     string `json:"user"`
	Host     string `json:"host"`
	Port     string `json:"port"`
	Password string `json:"-"`
}

func (c *Config) DatabaseConfig() *Config {
	return c
}

func (c *Config) ConnectionURL() string {
	if c == nil {
		return ""
	}

	host := c.Host
	if host == "" {
		host = "localhost"
	}

	port := c.Port
	if port == "" {
		port = "5432"
	}

	u := url.URL{
		Scheme: "postgres",
		Host:   host + ":" + port,
		Path:   c.Name,
	}

	if c.User != "" || c.Password != "" {
		u.User = url.UserPassword(c.User, c.Password)
	}

	return u.String()
}
