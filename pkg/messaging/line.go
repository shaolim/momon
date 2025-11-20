package messaging

import (
	"fmt"

	messagingapi "github.com/line/line-bot-sdk-go/v8/linebot/messaging_api"
)

type LineMessaging struct {
	*messagingapi.MessagingApiAPI
}

func NewLineMessaging(config *Config) (*LineMessaging, error) {
	api, err := messagingapi.NewMessagingApiAPI(config.LineChannelToken)
	if err != nil {
		return nil, fmt.Errorf("failed to initiate line messaging API: %w", err)
	}

	return &LineMessaging{
		api,
	}, nil
}
