package messaging

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	messagingapi "github.com/line/line-bot-sdk-go/v8/linebot/messaging_api"
	"github.com/line/line-bot-sdk-go/v8/linebot/webhook"
)

func (m *messaging) Callback(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	lineSignature := r.Header.Get("X-Line-Signature")
	if !webhook.ValidateSignature(m.config.LineChannelSecret, lineSignature, body) {
		slog.Error("invalid line signature", slog.String("signature", lineSignature))
		w.WriteHeader(http.StatusOK)
		return
	}

	var cb webhook.CallbackRequest
	if err := json.Unmarshal(body, &cb); err != nil {
		slog.Error("failed to unmarshal request body", slog.Any("error", err))
		w.WriteHeader(http.StatusOK)
		return
	}

	slog.Info("callback", slog.Any("response", cb))
	go func() {
		err := m.processCallback(cb)
		if err != nil {
			slog.Error("failed to process callback", slog.Any("error", err))
		}
	}()

	w.WriteHeader(http.StatusOK)
}

func (m *messaging) processCallback(callback webhook.CallbackRequest) error {
	for _, event := range callback.Events {
		switch e := event.(type) {
		case webhook.MessageEvent:
			switch message := e.Message.(type) {
			case webhook.TextMessageContent:
				slog.Info("text", slog.String("text", message.Text))
				resp, err := m.env.GetLineMessagingAPI().ReplyMessage(&messagingapi.ReplyMessageRequest{
					ReplyToken: e.ReplyToken,
					Messages: []messagingapi.MessageInterface{
						&messagingapi.TextMessage{
							Text: "Hi i got your message. thank you",
						},
					},
				})
				if err != nil {
					slog.Error("failed to reply message", slog.Any("error", err))
					return err
				}

				slog.Info("reply message", slog.Any("resp", resp))
			default:
				slog.Info("unknown event", slog.Any("event", message))
			}
		case webhook.FollowEvent:
			slog.Info("FollowEvent", slog.Any("event", e))
			if e.Source.GetType() == "user" {
				// TODO: Save User or update the status
				user := e.Source.(webhook.UserSource)
				res, err := m.env.GetLineMessagingAPI().GetProfile(user.UserId)
				if err != nil {
					slog.Error("failed to get profile", slog.String("error", err.Error()))
					return fmt.Errorf("failed to get profile: %w", err)
				}

				slog.Info("user profile", slog.Any("profile", res))
			}

			resp, err := m.env.GetLineMessagingAPI().ReplyMessage(&messagingapi.ReplyMessageRequest{
				ReplyToken: e.ReplyToken,
				Messages: []messagingapi.MessageInterface{
					&messagingapi.TextMessage{
						Text: "Hi thank you for following me!!",
					},
				},
			})
			if err != nil {
				slog.Error("failed to reply message", slog.Any("error", err))
				return err
			}

			slog.Info("reply message follow user", slog.Any("resp", resp))
		case webhook.UnfollowEvent:
			slog.Info("UnfollowEvent", slog.Any("event", e))
			if e.Source.GetType() == "user" {
				// TODO: update user status
			}
		default:
			slog.Info("unknown event", slog.Any("event", e))
		}
	}

	return nil
}
