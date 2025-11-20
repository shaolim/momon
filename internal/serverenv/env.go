package serverenv

import (
	"context"
	"github/shaolim/momon/pkg/database"
	"github/shaolim/momon/pkg/messaging"

	"github.com/openai/openai-go/v3"
)

type Option func(*ServerEnv) *ServerEnv

type ServerEnv struct {
	db               *database.DB
	openaiClient     *openai.Client
	lineMessagingAPI *messaging.LineMessaging
}

func New(opts ...Option) *ServerEnv {
	env := &ServerEnv{}
	for _, f := range opts {
		env = f(env)
	}

	return env
}

func WithOpenAIClient(openaiClient *openai.Client) Option {
	return func(s *ServerEnv) *ServerEnv {
		s.openaiClient = openaiClient
		return s
	}
}

func WithLineMessagingAPI(lineMessagingAPI *messaging.LineMessaging) Option {
	return func(s *ServerEnv) *ServerEnv {
		s.lineMessagingAPI = lineMessagingAPI
		return s
	}
}

func WithDatabase(db *database.DB) Option {
	return func(s *ServerEnv) *ServerEnv {
		s.db = db
		return s
	}
}

func (s *ServerEnv) GetOpenAIClient() *openai.Client {
	return s.openaiClient
}

func (s *ServerEnv) GetDatabase() *database.DB {
	return s.db
}

func (s *ServerEnv) GetLineMessagingAPI() *messaging.LineMessaging {
	return s.lineMessagingAPI
}

func (s *ServerEnv) Close(ctx context.Context) error {
	if s == nil {
		return nil
	}

	if s.db != nil {
		s.db.Close()
	}

	return nil
}
