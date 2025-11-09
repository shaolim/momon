package serverenv

import (
	messagingapi "github.com/line/line-bot-sdk-go/v8/linebot/messaging_api"
	"github.com/openai/openai-go/v3"
)

type Option func(*ServerEnv) *ServerEnv

type ServerEnv struct {
	openaiClient     *openai.Client
	lineMessagingAPI *messagingapi.MessagingApiAPI
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

func WithLineMessagingAPI(lineMessagingAPI *messagingapi.MessagingApiAPI) Option {
	return func(s *ServerEnv) *ServerEnv {
		s.lineMessagingAPI = lineMessagingAPI
		return s
	}
}

func (s *ServerEnv) GetOpenAIClient() *openai.Client {
	return s.openaiClient
}

func (s *ServerEnv) GetLineMessagingAPI() *messagingapi.MessagingApiAPI {
	return s.lineMessagingAPI
}
