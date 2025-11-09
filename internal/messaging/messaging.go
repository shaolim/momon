package messaging

import (
	"github/shaolim/momon/internal/serverenv"
	"net/http"
)

type messaging struct {
	env *serverenv.ServerEnv
}

func New(env *serverenv.ServerEnv) *messaging {
	return &messaging{
		env: env,
	}
}

func (m *messaging) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/callback", m.Callback)
	return mux
}
