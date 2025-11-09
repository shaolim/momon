package messaging

import (
	"github/shaolim/momon/internal/serverenv"
	"net/http"
)

type messaging struct {
	env    *serverenv.ServerEnv
	config *Config
}

func New(config *Config, env *serverenv.ServerEnv) *messaging {
	return &messaging{
		env:    env,
		config: config,
	}
}

func (m *messaging) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /callback", m.Callback)
	return mux
}
